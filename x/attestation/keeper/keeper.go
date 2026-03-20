// Package keeper implements the state machine logic for x/attestation.
// Handles storage and retrieval of attestation records per CHAIN-SPEC-v1.5.1.
package keeper

import (
	"fmt"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"rampage/x/attestation/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods
// for the attestation module's state.
type Keeper struct {
	cdc          codec.BinaryCodec
	storeService store.KVStoreService
	logger       log.Logger

	// authorityAddress is the bech32 address authorized for governance operations.
	authorityAddress string
}

// NewKeeper creates a new attestation Keeper.
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
func (k Keeper) Logger() log.Logger {
	return k.logger
}

// GetAuthorityAddress returns the governance authority address.
func (k Keeper) GetAuthorityAddress() string {
	return k.authorityAddress
}

// Store key prefixes
var (
	AttestationKeyPrefix       = []byte{0x01} // prefix for attestation records
	OracleNodeKeyPrefix        = []byte{0x02} // prefix for registered oracle nodes
	ParamsKeyPrefix            = []byte{0x03} // prefix for module params
)

// AttestationKey returns the store key for a given query ID.
func AttestationKey(queryID string) []byte {
	return append(AttestationKeyPrefix, []byte(queryID)...)
}

// OracleNodeKey returns the store key for a given oracle node address.
func OracleNodeKey(addr string) []byte {
	return append(OracleNodeKeyPrefix, []byte(addr)...)
}

// SetAttestation persists an attestation record to the KV store.
func (k Keeper) SetAttestation(ctx sdk.Context, attestation types.Attestation) error {
	storeAdapter := k.storeService.OpenKVStore(ctx)
	bz, err := k.cdc.Marshal(&attestation)
	if err != nil {
		return fmt.Errorf("failed to marshal attestation: %w", err)
	}
	return storeAdapter.Set(AttestationKey(attestation.QueryID), bz)
}

// GetAttestation retrieves an attestation by query ID.
// Returns ErrAttestationNotFound if the record does not exist.
func (k Keeper) GetAttestation(ctx sdk.Context, queryID string) (types.Attestation, error) {
	storeAdapter := k.storeService.OpenKVStore(ctx)
	bz, err := storeAdapter.Get(AttestationKey(queryID))
	if err != nil {
		return types.Attestation{}, err
	}
	if bz == nil {
		return types.Attestation{}, types.ErrAttestationNotFound
	}
	var attestation types.Attestation
	if err := k.cdc.Unmarshal(bz, &attestation); err != nil {
		return types.Attestation{}, fmt.Errorf("failed to unmarshal attestation: %w", err)
	}
	return attestation, nil
}

// HasAttestation returns true if an attestation exists for the given query ID.
func (k Keeper) HasAttestation(ctx sdk.Context, queryID string) bool {
	_, err := k.GetAttestation(ctx, queryID)
	return err == nil
}

// IsRegisteredOracle checks whether an address is a registered oracle node.
func (k Keeper) IsRegisteredOracle(ctx sdk.Context, addr string) bool {
	storeAdapter := k.storeService.OpenKVStore(ctx)
	bz, err := storeAdapter.Get(OracleNodeKey(addr))
	if err != nil || bz == nil {
		return false
	}
	return true
}

// RegisterOracleNode adds an oracle node address to the registered set.
// This can only be called via governance (Constitutional Amendment class).
func (k Keeper) RegisterOracleNode(ctx sdk.Context, addr string, modelID string) error {
	storeAdapter := k.storeService.OpenKVStore(ctx)
	return storeAdapter.Set(OracleNodeKey(addr), []byte(modelID))
}

// AppendCorrectionRecord appends a correction to an existing attestation's history.
func (k Keeper) AppendCorrectionRecord(
	ctx sdk.Context,
	queryID string,
	correction types.CorrectionRecord,
) error {
	attestation, err := k.GetAttestation(ctx, queryID)
	if err != nil {
		return err
	}
	// Update the verdict and append the correction record.
	attestation.CorrectionHistory = append(attestation.CorrectionHistory, correction)
	attestation.Verdict = correction.NewVerdict
	return k.SetAttestation(ctx, attestation)
}
