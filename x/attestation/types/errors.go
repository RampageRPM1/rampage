package types

import "cosmossdk.io/errors"

var (
	// ErrInvalidQueryID is returned when the query ID is empty or malformed.
	ErrInvalidQueryID = errors.Register(ModuleName, 1, "invalid query ID")

	// ErrInvalidClaim is returned when the claim text is empty.
	ErrInvalidClaim = errors.Register(ModuleName, 2, "claim text cannot be empty")

	// ErrInsufficientOracleSignatures is returned when fewer than 5-of-7 oracle signatures are provided.
	ErrInsufficientOracleSignatures = errors.Register(ModuleName, 3, "insufficient oracle signatures: minimum 5-of-7 required")

	// ErrInvalidSender is returned when the sender address is not a valid bech32 address.
	ErrInvalidSender = errors.Register(ModuleName, 4, "invalid sender address")

	// ErrAttestationNotFound is returned when a queried attestation does not exist.
	ErrAttestationNotFound = errors.Register(ModuleName, 5, "attestation not found")

	// ErrAttestationAlreadyExists is returned when a duplicate query ID is submitted.
	ErrAttestationAlreadyExists = errors.Register(ModuleName, 6, "attestation already exists for this query ID")

	// ErrUnauthorizedOracle is returned when the sender is not a registered oracle node.
	ErrUnauthorizedOracle = errors.Register(ModuleName, 7, "sender is not a registered oracle node")

	// ErrCorrectionReasonRequired is returned when no reason is given for a correction.
	ErrCorrectionReasonRequired = errors.Register(ModuleName, 8, "correction reason is required")

	// ErrInvalidVerificationLevel is returned when the verification level is out of range (1-4).
	ErrInvalidVerificationLevel = errors.Register(ModuleName, 9, "verification level must be between 1 and 4")

	// ErrInvalidAuditHash is returned when the audit hash is empty or malformed.
	ErrInvalidAuditHash = errors.Register(ModuleName, 10, "audit hash is required and must be a valid SHA-256 hex string")
)
