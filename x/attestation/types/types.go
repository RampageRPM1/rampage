// Package types defines the core data structures for the x/attestation module.
// Implements the Seven-Seal Truth Attestation system per CHAIN-SPEC-v1.5.1 Section 6.2.
package types

import (
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName = "attestation"
	StoreKey   = ModuleName
	RouterKey  = ModuleName

	// Oracle thresholds per CHAIN-SPEC-v1.5.1 Section 9
	MinOracleSignatures = 5
	MaxOracleNodes      = 7

	// Verification levels
	VerificationLevelUnverified   = 1
	VerificationLevelCorroborated = 2
	VerificationLevelCrossVerified = 3
	VerificationLevelOnTheGround  = 4
)

// Verdict represents the outcome of a Seven-Seal attestation consensus.
type Verdict string

const (
	VerdictTrue           Verdict = "TRUE"
	VerdictFalse          Verdict = "FALSE"
	VerdictUnverifiable   Verdict = "UNVERIFIABLE"
	VerdictPartiallyTrue  Verdict = "PARTIALLY_TRUE"
	VerdictMisleading     Verdict = "MISLEADING"
	VerdictContextMissing Verdict = "CONTEXT_MISSING"
)

// JurisdictionFlag identifies special jurisdiction handling per Constitution Art. VIII.
type JurisdictionFlag string

const (
	JurisdictionStandard        JurisdictionFlag = "STANDARD"
	JurisdictionElevated        JurisdictionFlag = "ELEVATED"
	JurisdictionHighRisk        JurisdictionFlag = "HIGH_RISK"
	JurisdictionConflictZone    JurisdictionFlag = "CONFLICT_ZONE" // Level 4 / Art. VIII
)

// ValidatorSignature represents a single oracle node's signed attestation.
type ValidatorSignature struct {
	// ValidatorAddress is the bech32 address of the signing oracle node.
	ValidatorAddress string `json:"validator_address"`
	// Signature is the ed25519 signature over the attestation hash.
	Signature []byte `json:"signature"`
	// ModelID identifies the AI model node (e.g., "gpt4o", "claude", "gemini").
	ModelID string `json:"model_id"`
	// Timestamp when this oracle node submitted its signature.
	Timestamp time.Time `json:"timestamp"`
}

// CorrectionRecord captures an amendment to a prior attestation.
type CorrectionRecord struct {
	// CorrectedAt is the block height at which the correction was applied.
	CorrectedAt int64 `json:"corrected_at"`
	// Reason describes why the correction was made.
	Reason string `json:"reason"`
	// PreviousVerdict is the verdict before correction.
	PreviousVerdict Verdict `json:"previous_verdict"`
	// NewVerdict is the corrected verdict.
	NewVerdict Verdict `json:"new_verdict"`
	// InitiatedBy is the validator address that submitted MsgCorrectAttestation.
	InitiatedBy string `json:"initiated_by"`
	// OracleSignatures are the new 5-of-7 signatures approving the correction.
	OracleSignatures []ValidatorSignature `json:"oracle_signatures"`
}

// Attestation is the on-chain record of a Seven-Seal truth verification query.
// Schema per CHAIN-SPEC-v1.5.1 Section 6.2.
type Attestation struct {
	// QueryID is the unique identifier for this truth query (UUID v4).
	QueryID string `json:"query_id"`
	// Claim is the human-readable statement being verified.
	Claim string `json:"claim"`
	// ConsensusPct is the percentage of oracle nodes that agreed on the verdict.
	ConsensusPct sdkmath.LegacyDec `json:"consensus_pct"`
	// Verdict is the outcome determined by the Seven-Seal oracle network.
	Verdict Verdict `json:"verdict"`
	// VerificationLevel indicates the depth of verification (1-4).
	VerificationLevel int32 `json:"verification_level"`
	// AuditHash is the SHA-256 hash of the full Seven-Seal audit bundle.
	AuditHash string `json:"audit_hash"`
	// AuditURL points to the off-chain audit bundle (truthoracle.ai).
	AuditURL string `json:"audit_url"`
	// ValidatorSigs contains the oracle node signatures (requires >= 5-of-7).
	ValidatorSigs []ValidatorSignature `json:"validator_sigs"`
	// CorrectionHistory is the immutable audit trail of all corrections.
	CorrectionHistory []CorrectionRecord `json:"correction_history"`
	// SubmittedBy is the validator address that submitted this attestation.
	SubmittedBy string `json:"submitted_by"`
	// BlockHeight is the block at which this attestation was finalized.
	BlockHeight int64 `json:"block_height"`
	// Timestamp is the wall-clock time of finalization.
	Timestamp time.Time `json:"timestamp"`
	// JurisdictionFlags identifies any special jurisdiction handling.
	JurisdictionFlags []JurisdictionFlag `json:"jurisdiction_flags"`
	// Active indicates whether this attestation is the current canonical record.
	Active bool `json:"active"`
}

// MsgSubmitAttestation is the message type for submitting a new attestation.
// Only credentialed validators (registered oracle nodes) may submit.
type MsgSubmitAttestation struct {
	Sender            string               `json:"sender"`
	QueryID           string               `json:"query_id"`
	Claim             string               `json:"claim"`
	Verdict           Verdict              `json:"verdict"`
	ConsensusPct      sdkmath.LegacyDec    `json:"consensus_pct"`
	VerificationLevel int32                `json:"verification_level"`
	AuditHash         string               `json:"audit_hash"`
	AuditURL          string               `json:"audit_url"`
	OracleSignatures  []ValidatorSignature `json:"oracle_signatures"`
	JurisdictionFlags []JurisdictionFlag   `json:"jurisdiction_flags"`
}

// GetSigners returns the signer addresses for MsgSubmitAttestation.
func (msg MsgSubmitAttestation) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{signer}
}

// ValidateBasic performs stateless validation on MsgSubmitAttestation.
func (msg MsgSubmitAttestation) ValidateBasic() error {
	if msg.QueryID == "" {
		return ErrInvalidQueryID
	}
	if msg.Claim == "" {
		return ErrInvalidClaim
	}
	if len(msg.OracleSignatures) < MinOracleSignatures {
		return ErrInsufficientOracleSignatures
	}
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return ErrInvalidSender
	}
	return nil
}

// MsgCorrectAttestation allows a quorum of validators to correct a prior attestation.
// Full audit trail is preserved per Constitution Art. IV.
type MsgCorrectAttestation struct {
	Sender           string               `json:"sender"`
	QueryID          string               `json:"query_id"`
	NewVerdict       Verdict              `json:"new_verdict"`
	Reason           string               `json:"reason"`
	OracleSignatures []ValidatorSignature `json:"oracle_signatures"`
}

// GetSigners returns the signer addresses for MsgCorrectAttestation.
func (msg MsgCorrectAttestation) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{signer}
}

// ValidateBasic performs stateless validation on MsgCorrectAttestation.
func (msg MsgCorrectAttestation) ValidateBasic() error {
	if msg.QueryID == "" {
		return ErrInvalidQueryID
	}
	if msg.Reason == "" {
		return ErrCorrectionReasonRequired
	}
	if len(msg.OracleSignatures) < MinOracleSignatures {
		return ErrInsufficientOracleSignatures
	}
	return nil
}

// Params defines the configurable parameters for x/attestation.
type Params struct {
	// MinOracleSignatures is the minimum 5-of-7 oracle signatures required.
	MinOracleSignatures int32 `json:"min_oracle_signatures"`
	// MaxOracleNodes is the maximum number of registered oracle nodes (7).
	MaxOracleNodes int32 `json:"max_oracle_nodes"`
	// RequireJurisdictionFlag enforces jurisdiction tagging on all attestations.
	RequireJurisdictionFlag bool `json:"require_jurisdiction_flag"`
}

// DefaultParams returns the default parameters for x/attestation.
func DefaultParams() Params {
	return Params{
		MinOracleSignatures:     MinOracleSignatures,
		MaxOracleNodes:          MaxOracleNodes,
		RequireJurisdictionFlag: false,
	}
}
