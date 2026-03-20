// Package keeper implements the 955 Operational Governor.
// Enforces the 95% truth verification / 5% humanitarian operational split
// per CHAIN-SPEC-v1.5.1 Section 6.2 and Rampage Constitution Art. VI.
package keeper

import (
	"encoding/binary"
	"fmt"

	"cosmossdk.io/core/store"
	sdkmath "cosmossdk.io/math"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"rampage/x/governor955/types"
)

const (
	// RatioPublishInterval is the number of blocks between on-chain ratio publications.
	// Approximately quarterly at 6s/block and 10,000 TPS target.
	RatioPublishInterval = int64(1000)
)

var (
	ParamsKey              = []byte{0x01}
	TruthSpendKey          = []byte{0x02} // cumulative urpm spent on truth verification
	HumanitarianSpendKey   = []byte{0x03} // cumulative urpm spent on humanitarian
	LastRatioPublishKey    = []byte{0x04} // last block height at which ratio was published
	RatioHistoryKeyPrefix  = []byte{0x05} // prefix for ratio history records
	LDFBalanceKey          = []byte{0x06} // Legal Defense Fund current balance
)

// SpendCategory classifies treasury outflows.
type SpendCategory string

const (
	// SpendTruthVerification is for attestation, journalist support, oracle ops.
	SpendTruthVerification SpendCategory = "TRUTH_VERIFICATION"
	// SpendHumanitarian is for Art. VIII humanitarian access operations.
	SpendHumanitarian SpendCategory = "HUMANITARIAN"
)

// Keeper maintains the link to data storage and exposes getter/setter methods
// for the governor955 module's state.
type Keeper struct {
	cdc              codec.BinaryCodec
	storeService     store.KVStoreService
	logger           log.Logger
	authorityAddress string
	bankKeeper       types.BankKeeper
}

// NewKeeper creates a new governor955 Keeper.
func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	logger log.Logger,
	authorityAddress string,
	bankKeeper types.BankKeeper,
) Keeper {
	if _, err := sdk.AccAddressFromBech32(authorityAddress); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s", err))
	}
	return Keeper{
		cdc:              cdc,
		storeService:     storeService,
		logger:           logger.With("module", fmt.Sprintf("x/%s", types.ModuleName)),
		authorityAddress: authorityAddress,
		bankKeeper:       bankKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger() log.Logger { return k.logger }

// RecordSpend records a treasury outflow under a given category.
// Called by any module that spends from the Rampage treasury.
func (k Keeper) RecordSpend(ctx sdk.Context, amount sdkmath.Int, category SpendCategory) error {
	switch category {
	case SpendTruthVerification:
		current := k.getTruthSpend(ctx)
		return k.setTruthSpend(ctx, current.Add(amount))
	case SpendHumanitarian:
		current := k.getHumanitarianSpend(ctx)
		return k.setHumanitarianSpend(ctx, current.Add(amount))
	default:
		return fmt.Errorf("unknown spend category: %s", category)
	}
}

// GetCurrentRatio returns the current truth verification / humanitarian ratio.
// Returns (truthPct, humanitarianPct) as sdk.Dec values (0-1 range).
func (k Keeper) GetCurrentRatio(ctx sdk.Context) (sdkmath.LegacyDec, sdkmath.LegacyDec) {
	truth := k.getTruthSpend(ctx)
	humanitarian := k.getHumanitarianSpend(ctx)
	total := truth.Add(humanitarian)

	if total.IsZero() {
		// No spend yet — return target ratio.
		return sdkmath.LegacyNewDecWithPrec(95, 2), sdkmath.LegacyNewDecWithPrec(5, 2)
	}

	truthPct := sdkmath.LegacyNewDecFromInt(truth).Quo(sdkmath.LegacyNewDecFromInt(total))
	humanitarianPct := sdkmath.LegacyNewDecFromInt(humanitarian).Quo(sdkmath.LegacyNewDecFromInt(total))
	return truthPct, humanitarianPct
}

// IsRatioDrifting returns true if the ratio has deviated from the 95/5 target.
// "Drifting" is defined as the truth verification allocation falling below 90%.
func (k Keeper) IsRatioDrifting(ctx sdk.Context) bool {
	truthPct, _ := k.GetCurrentRatio(ctx)
	threshold := sdkmath.LegacyNewDecWithPrec(90, 2) // 90% minimum
	return truthPct.LT(threshold)
}

// MaybePublishRatio publishes the current ratio to the chain log every RatioPublishInterval blocks.
// Should be called from EndBlock.
func (k Keeper) MaybePublishRatio(ctx sdk.Context) {
	currentBlock := ctx.BlockHeight()
	lastPublish := k.getLastRatioPublish(ctx)

	if currentBlock-lastPublish < RatioPublishInterval {
		return
	}

	truthPct, humanitarianPct := k.GetCurrentRatio(ctx)
	k.logger.Info("governor955: periodic ratio publication",
		"block", currentBlock,
		"truth_verification_pct", truthPct.String(),
		"humanitarian_pct", humanitarianPct.String(),
		"drifting", k.IsRatioDrifting(ctx),
	)

	// Emit an SDK event so indexers and the TruthOracle.ai dashboard can track it.
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"governor955_ratio_published",
			sdk.NewAttribute("block_height", fmt.Sprintf("%d", currentBlock)),
			sdk.NewAttribute("truth_verification_pct", truthPct.String()),
			sdk.NewAttribute("humanitarian_pct", humanitarianPct.String()),
			sdk.NewAttribute("is_drifting", fmt.Sprintf("%t", k.IsRatioDrifting(ctx))),
		),
	)

	k.setLastRatioPublish(ctx, currentBlock)
}

// --- Internal helpers ---

func (k Keeper) getTruthSpend(ctx sdk.Context) sdkmath.Int {
	storeAdapter := k.storeService.OpenKVStore(ctx)
	bz, _ := storeAdapter.Get(TruthSpendKey)
	if bz == nil {
		return sdkmath.ZeroInt()
	}
	var i sdkmath.Int
	k.cdc.MustUnmarshal(bz, &i)
	return i
}

func (k Keeper) setTruthSpend(ctx sdk.Context, amount sdkmath.Int) error {
	storeAdapter := k.storeService.OpenKVStore(ctx)
	bz := k.cdc.MustMarshal(&amount)
	return storeAdapter.Set(TruthSpendKey, bz)
}

func (k Keeper) getHumanitarianSpend(ctx sdk.Context) sdkmath.Int {
	storeAdapter := k.storeService.OpenKVStore(ctx)
	bz, _ := storeAdapter.Get(HumanitarianSpendKey)
	if bz == nil {
		return sdkmath.ZeroInt()
	}
	var i sdkmath.Int
	k.cdc.MustUnmarshal(bz, &i)
	return i
}

func (k Keeper) setHumanitarianSpend(ctx sdk.Context, amount sdkmath.Int) error {
	storeAdapter := k.storeService.OpenKVStore(ctx)
	bz := k.cdc.MustMarshal(&amount)
	return storeAdapter.Set(HumanitarianSpendKey, bz)
}

func (k Keeper) getLastRatioPublish(ctx sdk.Context) int64 {
	storeAdapter := k.storeService.OpenKVStore(ctx)
	bz, _ := storeAdapter.Get(LastRatioPublishKey)
	if bz == nil {
		return 0
	}
	return int64(binary.BigEndian.Uint64(bz))
}

func (k Keeper) setLastRatioPublish(ctx sdk.Context, height int64) {
	storeAdapter := k.storeService.OpenKVStore(ctx)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, uint64(height))
	_ = storeAdapter.Set(LastRatioPublishKey, bz)
}
