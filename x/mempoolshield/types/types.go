// Package types defines the data structures for the x/mempoolshield module.
// Implements the Compliance Membrane per CHAIN-SPEC-v1.5.1 Section 6.2.
package types

const (
	ModuleName = "mempoolshield"
	StoreKey   = ModuleName

	// Threat level constants per CHAIN-SPEC-v1.5.1 Section 8.
	ThreatLevelStandard  = 1 // Default: oracle-screened capital routing
	ThreatLevelElevated  = 2 // Enhanced screening
	ThreatLevelHigh      = 3 // Dual oracle confirmation
	ThreatLevelConflict  = 4 // Level 4 / Art. VIII: capital routing SUSPENDED
)

// FailsafeMode defines the behavior when the oracle is unavailable.
type FailsafeMode string

const (
	// FailsafeLockdownCapitalOnly rejects capital-routing txs only; blocks continue.
	// This is the REQUIRED default per CHAIN-SPEC-v1.5.1 Liveness Invariant #3.
	FailsafeLockdownCapitalOnly FailsafeMode = "LOCKDOWN_CAPITAL_ONLY"
	// FailsafePassThrough allows all txs when oracle is unavailable (NOT RECOMMENDED).
	FailsafePassThrough FailsafeMode = "PASS_THROUGH"
)

// Params defines the configurable parameters for x/mempoolshield.
type Params struct {
	// Enabled controls whether the Mempool Shield is active.
	Enabled bool `json:"enabled"`
	// ThreatLevel is the current operational threat level (1-4).
	ThreatLevel int32 `json:"threat_level"`
	// OracleThreshold is the minimum oracle signatures required (default: 5).
	OracleThreshold int32 `json:"oracle_threshold"`
	// OracleSigners is the total number of oracle nodes (default: 7).
	OracleSigners int32 `json:"oracle_signers"`
	// FailsafeDefault defines behavior when the oracle is unavailable.
	// MUST be LOCKDOWN_CAPITAL_ONLY to satisfy liveness invariants.
	FailsafeDefault FailsafeMode `json:"failsafe_default"`
	// ProhibitedEntities is the local testnet allow/deny list.
	// Populated from the 5-of-7 oracle committee in production.
	ProhibitedEntities []string `json:"prohibited_entities"`
}

// DefaultParams returns the default Mempool Shield parameters.
func DefaultParams() Params {
	return Params{
		Enabled:            true,
		ThreatLevel:        ThreatLevelStandard,
		OracleThreshold:    5,
		OracleSigners:      7,
		FailsafeDefault:    FailsafeLockdownCapitalOnly,
		ProhibitedEntities: []string{},
	}
}

// OracleNode represents a registered Mempool Shield oracle signer.
type OracleNode struct {
	// Address is the bech32 validator address.
	Address string `json:"address"`
	// FeedType identifies the compliance feed (e.g., "OFAC", "UN", "EU", "FATF").
	FeedType string `json:"feed_type"`
	// Active indicates whether this node is currently participating.
	Active bool `json:"active"`
	// LastUpdateBlock is the last block at which this node submitted an oracle update.
	LastUpdateBlock int64 `json:"last_update_block"`
}

// ScreeningResult is the oracle committee's response for a capital-routing transaction.
type ScreeningResult struct {
	// TxHash identifies the transaction.
	TxHash string `json:"tx_hash"`
	// Approved is true if the transaction passed all compliance checks.
	Approved bool `json:"approved"`
	// MatchedEntities lists any prohibited entity identifiers found.
	MatchedEntities []string `json:"matched_entities"`
	// OracleSignatures contains the approving/rejecting oracle signatures.
	OracleSignatures int32 `json:"oracle_signatures"`
}
