// Package keeper implements the state machine logic for x/mempoolshield.
package keeper

import (
	"fmt"
	"time"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"rampage/x/mempoolshield/types"
)

// OracleStatusKey is the KV store key for the oracle availability status.
var (
	ParamsKey            = []byte{0x01}
	OracleStatusKey      = []byte{0x02}
	ProhibitedEntityKey  = []byte{0x03}
	ThreatLevelKey       = []byte{0x04}
	OracleLastPingKey    = []byte{0x05}
)

// OracleAvailabilityTimeout is the maximum age of an oracle ping before it is
// considered unavailable. After this duration, capital routing is fail-closed.
const OracleAvailabilityTimeout = 5 * time.Minute

// Keeper maintains the link to data storage and exposes getter/setter methods
// for the mempoolshield module's state.
type Keeper struct {
	cdc              codec.BinaryCodec
	storeService     store.KVStoreService
	logger           log.Logger
	authorityAddress string
}

// NewKeeper creates a new mempoolshield Keeper.
func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	logger log.Logger,
	authorityAddress string,
) Keeper {
	if _, err := sdk.AccAddressFromBech32(authorityAddress); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s", err))
	}
	return Keeper{
		cdc:              cdc,
		storeService:     storeService,
		logger:           logger.With("module", fmt.Sprintf("x/%s", types.ModuleName)),
		authorityAddress: authorityAddress,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger() log.Logger { return k.logger }

// GetParams returns the current mempoolshield parameters from the KV store.
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	storeAdapter := k.storeService.OpenKVStore(ctx)
	bz, err := storeAdapter.Get(ParamsKey)
	if err != nil || bz == nil {
		return types.DefaultParams()
	}
	var params types.Params
	if err := k.cdc.Unmarshal(bz, &params); err != nil {
		return types.DefaultParams()
	}
	return params
}

// SetParams persists mempoolshield parameters to the KV store.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	storeAdapter := k.storeService.OpenKVStore(ctx)
	bz, err := k.cdc.Marshal(&params)
	if err != nil {
		return err
	}
	return storeAdapter.Set(ParamsKey, bz)
}

// GetThreatLevel returns the current operational threat level (1-4).
func (k Keeper) GetThreatLevel(ctx sdk.Context) int32 {
	return k.GetParams(ctx).ThreatLevel
}

// SetThreatLevel updates the threat level. Only callable via governance.
// This does NOT halt the chain — it changes the capital-routing filter behavior.
func (k Keeper) SetThreatLevel(ctx sdk.Context, level int32) error {
	params := k.GetParams(ctx)
	params.ThreatLevel = level
	return k.SetParams(ctx, params)
}

// IsOracleAvailable returns true if the oracle committee has checked in recently.
// "Available" means at least one oracle ping has been received within OracleAvailabilityTimeout.
// If unavailable, the fail-closed rule applies to capital-routing txs ONLY.
func (k Keeper) IsOracleAvailable(ctx sdk.Context) bool {
	storeAdapter := k.storeService.OpenKVStore(ctx)
	bz, err := storeAdapter.Get(OracleLastPingKey)
	if err != nil || bz == nil {
		// No oracle ping ever received. Testnet: return true to allow pass-through.
		// Production: change this to return false (fail-closed by default).
		return true // TODO: change to false for mainnet
	}
	var lastPing time.Time
	if err := lastPing.UnmarshalBinary(bz); err != nil {
		return false
	}
	return time.Since(lastPing) < OracleAvailabilityTimeout
}

// RecordOraclePing updates the oracle's last-seen timestamp.
func (k Keeper) RecordOraclePing(ctx sdk.Context) error {
	storeAdapter := k.storeService.OpenKVStore(ctx)
	now := ctx.BlockTime()
	bz, err := now.MarshalBinary()
	if err != nil {
		return err
	}
	return storeAdapter.Set(OracleLastPingKey, bz)
}

// OracleApproves checks whether the given transaction is approved by the oracle committee.
// For testnet phase, this is a stub that checks the local prohibited entity list.
// Production: replace with 5-of-7 oracle committee ZK-proof verification.
func (k Keeper) OracleApproves(ctx sdk.Context, txBytes []byte) bool {
	params := k.GetParams(ctx)
	// Testnet stub: check local prohibited entity list.
	// In production, this calls out to the 5-of-7 oracle signer committee.
	// The local list is managed via governance MsgUpdateProhibitedEntities.
	if len(params.ProhibitedEntities) == 0 {
		return true // empty list — approve all (testnet default)
	}
	// TODO: decode txBytes, extract recipient addresses, check against prohibited list.
	_ = txBytes
	return true
}

// AddProhibitedEntity adds an entity to the local prohibited list.
// For testnet use only; production list comes from oracle committee.
func (k Keeper) AddProhibitedEntity(ctx sdk.Context, entity string) error {
	params := k.GetParams(ctx)
	params.ProhibitedEntities = append(params.ProhibitedEntities, entity)
	return k.SetParams(ctx, params)
}
