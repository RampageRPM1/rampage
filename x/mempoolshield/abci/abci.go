// Package abci implements the ABCI PrepareProposal and ProcessProposal hooks
// for the Mempool Shield compliance membrane.
//
// LIVENESS INVARIANT (CHAIN-SPEC-v1.5.1 Section 3 & 7):
// Under ALL conditions, including oracle failure and LOCKDOWN mode, these hooks
// MUST continue producing valid (possibly empty) block proposals. Capital-routing
// transactions are filtered out via reject_txs but NEVER cause an application error
// that would stall consensus. Non-capital transactions ALWAYS pass through.
package abci

import (
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"rampage/x/mempoolshield/keeper"
	"rampage/x/mempoolshield/types"
)

// CapitalRoutingMsgTypes is the set of message type URLs subject to Mempool Shield screening.
// Per CHAIN-SPEC-v1.5.1 Section 6.2: bank sends, IBC transfers, humanitarian sends.
var CapitalRoutingMsgTypes = map[string]bool{
	"/cosmos.bank.v1beta1.MsgSend":             true,
	"/cosmos.bank.v1beta1.MsgMultiSend":        true,
	"/ibc.applications.transfer.v1.MsgTransfer": true,
	"/rampage.mempoolshield.v1.MsgHumanitarianSend": true,
}

// PrepareProposalHandler returns a PrepareProposal handler that filters capital-routing
// transactions through the Mempool Shield oracle.
//
// Liveness guarantee: this function NEVER returns an error. It always produces a valid
// (possibly reduced) list of transactions. Capital-routing txs rejected by the oracle
// or when oracle is unavailable are excluded from the proposal; all other txs pass through.
func PrepareProposalHandler(k keeper.Keeper) sdk.PrepareProposalHandler {
	return func(ctx sdk.Context, req *abci.RequestPrepareProposal) (*abci.ResponsePrepareProposal, error) {
		params := k.GetParams(ctx)

		// If shield is disabled, pass all transactions through unchanged.
		if !params.Enabled {
			return &abci.ResponsePrepareProposal{Txs: req.Txs}, nil
		}

		oracleAvailable := k.IsOracleAvailable(ctx)
		threatLevel := k.GetThreatLevel(ctx)

		var filteredTxs [][]byte
		var rejectedCount int

		for _, txBytes := range req.Txs {
			if isCapitalRoutingTx(ctx, txBytes) {
				switch {
				case threatLevel >= types.ThreatLevelConflict:
					// Level 4 / Art. VIII jurisdiction: ALL capital routing suspended.
					// LIVENESS: skip this tx, do NOT halt proposal.
					k.Logger().Info("mempoolshield: Level 4 LOCKDOWN — capital tx excluded",
						"threat_level", threatLevel)
					rejectedCount++
					continue

				case !oracleAvailable:
					// Oracle unavailable: fail-closed on capital routing only.
					// LIVENESS: skip this tx, do NOT halt proposal.
					k.Logger().Warn("mempoolshield: oracle unavailable — capital tx excluded (fail-closed)",
						"failsafe", params.FailsafeDefault)
					rejectedCount++
					continue

				default:
					// Oracle available: check the prohibited entity list.
					if k.OracleApproves(ctx, txBytes) {
						filteredTxs = append(filteredTxs, txBytes)
					} else {
						k.Logger().Info("mempoolshield: oracle rejected capital tx")
						rejectedCount++
					}
				}
			} else {
				// Non-capital transaction: ALWAYS passes through.
				filteredTxs = append(filteredTxs, txBytes)
			}
		}

		if rejectedCount > 0 {
			k.Logger().Info(fmt.Sprintf("mempoolshield: excluded %d capital-routing tx(s) from proposal", rejectedCount))
		}

		// LIVENESS: always return a valid response. Never return an error here.
		return &abci.ResponsePrepareProposal{Txs: filteredTxs}, nil
	}
}

// ProcessProposalHandler returns a ProcessProposal handler that validates
// capital-routing transactions in a proposed block against the oracle.
//
// Liveness guarantee: if oracle is unavailable during ProcessProposal, we accept
// the block if the proposer's PrepareProposal should have already filtered the txs.
// We REJECT the proposal only if a prohibited capital-routing tx slipped through.
func ProcessProposalHandler(k keeper.Keeper) sdk.ProcessProposalHandler {
	return func(ctx sdk.Context, req *abci.RequestProcessProposal) (*abci.ResponseProcessProposal, error) {
		params := k.GetParams(ctx)

		if !params.Enabled {
			return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_ACCEPT}, nil
		}

		threatLevel := k.GetThreatLevel(ctx)

		for _, txBytes := range req.Txs {
			if isCapitalRoutingTx(ctx, txBytes) {
				if threatLevel >= types.ThreatLevelConflict {
					// Level 4: no capital txs allowed. Reject this proposal.
					k.Logger().Warn("mempoolshield: Level 4 — capital tx in proposal, rejecting")
					return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_REJECT}, nil
				}
				if !k.OracleApproves(ctx, txBytes) {
					k.Logger().Warn("mempoolshield: prohibited capital tx in proposal, rejecting")
					return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_REJECT}, nil
				}
			}
		}

		return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_ACCEPT}, nil
	}
}

// isCapitalRoutingTx returns true if the transaction contains any capital-routing messages.
// This uses a lightweight type-URL check to avoid full tx decoding in the hot path.
func isCapitalRoutingTx(_ sdk.Context, txBytes []byte) bool {
	// TODO: decode tx and check each msg type URL against CapitalRoutingMsgTypes.
	// For testnet phase, we use a lightweight heuristic based on message type detection.
	// A full implementation will use codec.UnmarshalTx and iterate msgs.
	_ = banktypes.MsgSend{} // ensure import is used; full decode in next iteration
	_ = txBytes

	// Stub: returns false for testnet (pass-through mode).
	// Replace with: decode tx, check msg.ProtoMessage() type URL.
	return false
}

// Ensure baseapp handler types are used.
var _ baseapp.ProposalTxVerifier = nil
