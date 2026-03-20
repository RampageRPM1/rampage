// Package module implements the Cosmos SDK AppModule interface for x/attestation.
package module

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	"cosmossdk.io/core/appmodule"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"rampage/x/attestation/keeper"
	"rampage/x/attestation/types"
)

const ConsensusVersion = 1

var (
	_ module.AppModuleBasic = AppModule{}
	_ appmodule.AppModule   = AppModule{}
)

// AppModule implements the attestation module.
type AppModule struct {
	cdc    codec.Codec
	keeper keeper.Keeper
}

// NewAppModule creates a new AppModule for x/attestation.
func NewAppModule(cdc codec.Codec, keeper keeper.Keeper) AppModule {
	return AppModule{cdc: cdc, keeper: keeper}
}

// Name returns the module name.
func (AppModule) Name() string { return types.ModuleName }

// RegisterLegacyAminoCodec registers amino codec (required by SDK).
func (AppModule) RegisterLegacyAminoCodec(_ *codec.LegacyAmino) {}

// RegisterInterfaces registers the module's interface types.
func (AppModule) RegisterInterfaces(_ codectypes.InterfaceRegistry) {}

// RegisterGRPCGatewayRoutes registers gRPC gateway routes.
func (AppModule) RegisterGRPCGatewayRoutes(_ client.Context, _ *runtime.ServeMux) {}

// GetTxCmd returns the root tx command for x/attestation.
func (AppModule) GetTxCmd() *cobra.Command { return nil }

// GetQueryCmd returns the root query command for x/attestation.
func (AppModule) GetQueryCmd() *cobra.Command { return nil }

// DefaultGenesis returns the default genesis state as raw JSON.
func (am AppModule) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesis())
}

// ValidateGenesis validates the genesis state.
func (am AppModule) ValidateGenesis(cdc codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
	var gs types.GenesisState
	if err := cdc.UnmarshalJSON(bz, &gs); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
	}
	return gs.Validate()
}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (AppModule) IsOnePerModuleType() {}

// IsAppModule implements the appmodule.AppModule interface.
func (AppModule) IsAppModule() {}

// ConsensusVersion implements AppModule/HasConsensusVersion.
func (AppModule) ConsensusVersion() uint64 { return ConsensusVersion }

// InitGenesis initializes the module's state from a genesis state.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, bz json.RawMessage) {
	var gs types.GenesisState
	cdc.MustUnmarshalJSON(bz, &gs)
	// TODO: apply genesis state to keeper (oracle node registrations, initial params)
	_ = gs
}

// ExportGenesis returns the module's genesis state as raw JSON.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	gs := types.GenesisState{
		Params: am.keeper.GetParams(ctx),
	}
	return cdc.MustMarshalJSON(&gs)
}

// RegisterServices registers gRPC services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	// TODO: register MsgServer and QueryServer once proto-generated stubs exist
	_ = cfg
}

// EndBlock executes all ABCI EndBlock logic for the attestation module.
func (am AppModule) EndBlock(_ context.Context) error {
	// No periodic EndBlock logic required for attestation at this time.
	return nil
}
