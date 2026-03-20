// Package types defines the data structures for the x/governor955 module.
// Implements the 955 Operational Split Enforcer per CHAIN-SPEC-v1.5.1 Section 6.2.
package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName = "governor955"
	StoreKey   = ModuleName

	// TruthVerificationTarget is the target allocation for truth verification operations (95%).
	TruthVerificationTarget = "0.95"
	// HumanitarianTarget is the target allocation for humanitarian access operations (5%).
	HumanitarianTarget = "0.05"
	// DriftThreshold is the minimum truth verification allocation before review is triggered (90%).
	DriftThreshold = "0.90"
	// ReviewTriggerMonths is the number of months of drift before mandatory governance review.
	ReviewTriggerMonths = 18
)

// BankKeeper defines the expected bank module keeper interface for governor955.
type BankKeeper interface {
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
}

// Params defines the configurable parameters for x/governor955.
type Params struct {
	// TruthVerificationTarget is the target ratio for truth verification (default: 0.95).
	TruthVerificationTarget sdkmath.LegacyDec `json:"truth_verification_target"`
	// HumanitarianTarget is the target ratio for humanitarian access (default: 0.05).
	HumanitarianTarget sdkmath.LegacyDec `json:"humanitarian_target"`
	// ReviewTriggerMonths is how many months of drift triggers mandatory governance review.
	ReviewTriggerMonths int32 `json:"review_trigger_months"`
	// LegalDefenseFundAddress is the bech32 address of the ring-fenced Legal Defense Fund.
	LegalDefenseFundAddress string `json:"legal_defense_fund_address"`
}

// DefaultParams returns the default 955 governor parameters per Constitution Art. VI.
func DefaultParams() Params {
	return Params{
		TruthVerificationTarget: sdkmath.LegacyMustNewDecFromStr(TruthVerificationTarget),
		HumanitarianTarget:      sdkmath.LegacyMustNewDecFromStr(HumanitarianTarget),
		ReviewTriggerMonths:     ReviewTriggerMonths,
		LegalDefenseFundAddress: "", // Set at genesis
	}
}

// RatioRecord is an on-chain snapshot of the 955 ratio at a given block.
type RatioRecord struct {
	// BlockHeight is when the snapshot was taken.
	BlockHeight int64 `json:"block_height"`
	// TruthVerificationPct is the truth verification spend ratio at this block.
	TruthVerificationPct sdkmath.LegacyDec `json:"truth_verification_pct"`
	// HumanitarianPct is the humanitarian spend ratio at this block.
	HumanitarianPct sdkmath.LegacyDec `json:"humanitarian_pct"`
	// IsDrifting indicates whether the ratio is below the 90% drift threshold.
	IsDrifting bool `json:"is_drifting"`
}

// SpendRecord tracks an individual treasury outflow.
type SpendRecord struct {
	// Amount is the quantity of urpm spent.
	Amount sdkmath.Int `json:"amount"`
	// Category is the classification of the spend (TRUTH_VERIFICATION or HUMANITARIAN).
	Category string `json:"category"`
	// BlockHeight is the block at which the spend occurred.
	BlockHeight int64 `json:"block_height"`
	// Memo is an optional human-readable description of the expenditure.
	Memo string `json:"memo"`
}
