# Rampage Chain Specification v1.5.1

> **Status:** Active — `rampage-testnet-1`  
> **Date:** 2026-03-20  
> **Supersedes:** All prior Ignite scaffold defaults

---

## 1. Chain Identity

| Parameter | Value |
|---|---|
| Chain ID | `rampage-testnet-1` |
| Binary | `rampaged` |
| Bech32 Prefix | `rampage` |
| Base Denom | `urpm` |
| Display Denom | `rpm` |
| Denom Factor | 1,000,000 urpm = 1 RPM |
| Total Supply | 21,000,000 RPM (21,000,000,000,000 urpm) |
| Consensus Engine | CometBFT (1/3 Byzantine fault tolerance) |
| Target TPS | 10,000 |
| Block Finality | < 7 seconds |

---

## 2. Initial Supply Allocation

| Allocation | RPM | urpm | Notes |
|---|---|---|---|
| Truth Verification Reserve | 10,500,000 | 10,500,000,000,000 | 50% — validator rewards + attestation incentives |
| Legal Defense Fund | 1,050,000 | 1,050,000,000,000 | 5% — Constitution Art. IX, immutable |
| Ecosystem / Developer Fund | 4,200,000 | 4,200,000,000,000 | 20% — governed by 955 Governor |
| Community Airdrop | 2,100,000 | 2,100,000,000,000 | 10% — Truth Bearer onboarding |
| Humanitarian Access Reserve | 1,050,000 | 1,050,000,000,000 | 5% — 955 Governor humanitarian split |
| Founding Contributors | 2,100,000 | 2,100,000,000,000 | 10% — 2-year vesting, 6-month cliff |
| **Total** | **21,000,000** | **21,000,000,000,000** | |

---

## 3. Liveness Invariants

The following invariants are **non-negotiable** and MUST be preserved by all module implementations:

1. Consensus remains live under up to 1/3 byzantine voting power.
2. No governance action, oracle failure, or module error may halt block production.
3. Only explicitly defined, temporary LOCKDOWN modes may restrict specific transaction
   classes (e.g., capital routing). LOCKDOWN MUST NOT prevent block proposal or finalization
   for non-capital state transitions.
4. In any fail-safe mode, `PrepareProposal` and `ProcessProposal` MUST continue producing
   and accepting non-capital blocks. Capital-routing txs are returned in `reject_txs` only —
   they never cause an application error that stalls consensus.

---

## 4. Governance Parameters (Constitution v1.5)

All thresholds derive from Rampage Constitution Articles V–XI. Articles I–IV are **immutable**
and MUST NOT be alterable by any on-chain governance proposal.

| Governance Class | Min Deposit | Quorum | Approval Threshold | Voting Period | Enactment Delay |
|---|---|---|---|---|---|
| Standard (Art. V–VII, X–XI) | 1,000 RPM | 51% | 67% | 7 days | 48 hours |
| Constitutional Amendment (Art. V–XI) | 10,000 RPM | 60% | 80% | 30 days | 90 days |
| Emergency Legal Compliance | 500 RPM | 40% | 75% | 24 hours | Immediate |
| Humanitarian Escalation (Art. VIII) | 5,000 RPM | 50% | 75% | 30 days | 90 days |
| Level 4 / Conflict Mode | 5,000 RPM | 75% | 80% | 30 days | 90 days |

**Immutable parameters (Articles I–IV — human rights primacy, nonviolence,
civilians-only, truth verification standards):**  
These MUST be enforced at the application level by rejecting any `MsgSubmitProposal`
that would modify them, regardless of governance outcome.

---

## 5. Staking / SVA Parameters

| Parameter | Value |
|---|---|
| Standard Truth Bearer Minimum | 100,000 RPM |
| Unbonding Period | 21 days |
| Max Validators (testnet) | 21 |
| Max Validators (mainnet target) | 100 |
| SVA Pre-Mainnet Carve-Out | Enabled — governance-defined per-validator override |
| SVA Post-Mainnet Convergence | 6 months after mainnet launch |
| Slash: Verification Misconduct | 5% stake |
| Slash: Capital Routing Violation | 10% stake + jailing |
| Slash: Double Sign | 5% stake + tombstone |
| Slash: Downtime | 1% stake + jailing (unjail after 10 min) |

---

## 6. Module Inventory

### 6.1 Standard Cosmos SDK Modules (Inherited)

- `x/auth`, `x/bank`, `x/staking`, `x/slashing`, `x/distribution`
- `x/gov` (parameterized per Section 4 above)
- `x/params`, `x/upgrade`, `x/evidence`
- `x/ibc` (IBC transfers subject to Mempool Shield screening)

### 6.2 Rampage Custom Modules

#### `x/attestation` — Seven-Seal Truth Attestation
- Stores on-chain consensus records from the Seven-Seal off-chain oracle network.
- Schema: `queryID`, `consensusPct`, `verdict`, `verificationLevel`, `auditHash`,
  `validatorSigs[]`, `correctionHistory[]`, `timestamp`, `jurisdictionFlags`.
- `MsgSubmitAttestation`: callable only by credentialed validators.
- `MsgCorrectAttestation`: correction with full audit trail, quorum-gated.
- Query endpoints: `GetAttestation`, `ListAttestations`, `GetAuditTrail`.
- Off-chain Seven-Seal engine signs and pushes via authenticated RPC.

#### `x/mempoolshield` — Compliance Membrane
- ABCI `PrepareProposal`/`ProcessProposal` hooks inspect all capital-routing messages:
  `MsgSend` (bank), `MsgTransfer` (IBC), `MsgHumanitarianSend` (custom).
- Oracle interface: 5-of-7 signer committee (OFAC / UN / EU / FATF feeds).
- Zero-knowledge proof layer preserves recipient privacy during screening.
- **Fail-closed rule:** On oracle unavailability, capital-routing txs are rejected.
  Non-capital txs and block production are NEVER affected.
- Threat levels: 1 (standard) → 4 (Level 4 / conflict mode, e.g., Art. VIII jurisdictions).
- Prohibited entity list sourced from 5-of-7 oracle committee; local allow/deny list
  used in testnet phase pending live oracle integration.

#### `x/governor955` — Operational Expenditure Enforcer
- Tracks all treasury outflows and classifies them: `TRUTH_VERIFICATION` or `HUMANITARIAN`.
- Enforces 95% / 5% split (truth verification / humanitarian access).
- Publishes ratio on-chain each 1,000 blocks (quarterly equivalent).
- If ratio drifts for 18+ months, triggers mandatory governance review proposal.
- Legal Defense Fund (5% of total supply) is ring-fenced and non-redistributable without
  Constitutional Amendment governance class approval.

---

## 7. ABCI Hook Specification (Mempool Shield)

```
PrepareProposal(txs):
  filtered = []
  for tx in txs:
    if isCapitalRoutingTx(tx):
      if oracleAvailable() and oracleApproves(tx):
        filtered.append(tx)
      else:
        // oracle down or tx prohibited — reject this tx only
        rejectTx(tx)
    else:
      filtered.append(tx)  // non-capital tx always passes
  return filtered  // NEVER return error; always return a valid (possibly empty) block

ProcessProposal(txs):
  for tx in txs:
    if isCapitalRoutingTx(tx) and not oracleApproves(tx):
      return REJECT  // reject the proposal, not the chain
  return ACCEPT
```

---

## 8. Jurisdiction & Threat Level Matrix

| Level | Trigger | Capital Routing | Attestation | Governance |
|---|---|---|---|---|
| 1 — Standard | Default | Oracle-screened | Normal | Standard thresholds |
| 2 — Elevated | Regional conflict indicator | Enhanced screening | Priority queue | Standard thresholds |
| 3 — High | Active conflict zone | Dual oracle confirmation | Expedited | Emergency class eligible |
| 4 — Conflict Mode | Art. VIII jurisdictions (e.g., Iran) | SUSPENDED | Humanitarian-only | 75% quorum / 80% approval |

---

## 9. Seven-Seal Oracle Architecture (On-Chain Interface)

- 7 independent AI model nodes (off-chain) submit signed attestation bundles.
- On-chain module requires ≥ 5-of-7 signatures to accept an attestation as finalized.
- Each model node identified by its validator-controlled public key registered in `x/attestation`.
- Model keys rotatable via Constitutional Amendment governance proposal only.
- Attestation verification levels:
  - Level 1: Single-source, unverified claim
  - Level 2: Corroborated (2+ independent sources)
  - Level 3: Cross-verified (5-of-7 oracle consensus, documented sources)
  - Level 4: Independently verified on-the-ground (physical presence + oracle consensus)

---

## 10. Genesis Configuration Target

```yaml
chain_id: rampage-testnet-1
denom: urpm
initial_height: 1
gov:
  voting_params:
    voting_period: "604800s"  # 7 days standard
  deposit_params:
    min_deposit: [{amount: "1000000000", denom: "urpm"}]  # 1000 RPM
  tally_params:
    quorum: "0.51"
    threshold: "0.67"
    veto_threshold: "0.334"
staking:
  params:
    bond_denom: urpm
    max_validators: 21
    unbonding_time: "1814400s"  # 21 days
slashing:
  params:
    slash_fraction_double_sign: "0.05"
    slash_fraction_downtime: "0.01"
mempoolshield:
  enabled: true
  threat_level: 1
  oracle_threshold: 5
  oracle_signers: 7
  failsafe_default: LOCKDOWN_CAPITAL_ONLY  # reject capital txs only, never halt blocks
governor955:
  enabled: true
  truth_verification_target: "0.95"
  humanitarian_target: "0.05"
  review_trigger_months: 18
attendance:
  min_oracle_signatures: 5
  max_oracle_nodes: 7
```

---

## 11. File & Module Layout

```
rampage/
├── app/
│   ├── app.go               # Register x/attestation, x/mempoolshield, x/governor955
│   ├── app_config.go        # depinject wiring for new modules
│   └── genesis.go           # Updated genesis with urpm + constitutional params
├── x/
│   ├── rampage/             # Existing base module (keep, rename to x/core if needed)
│   ├── attestation/         # NEW — Seven-Seal on-chain interface
│   │   ├── keeper/
│   │   ├── module/
│   │   └── types/
│   ├── mempoolshield/       # NEW — Compliance membrane + ABCI hooks
│   │   ├── abci/
│   │   ├── keeper/
│   │   ├── module/
│   │   └── types/
│   └── governor955/         # NEW — 955 operational split enforcer
│       ├── keeper/
│       ├── module/
│       └── types/
├── docs/
│   ├── CHAIN-SPEC-v1.5.1.md # THIS FILE
│   └── CONSTITUTION-v1.5.md # Constitutional law reference
├── config.yml               # Updated: urpm denom, testnet accounts, validators
└── readme.md                # Updated: v1.5.1 architecture overview
```

---

## 12. Implementation Checklist

- [ ] `config.yml` — update denom to `urpm`, chain identity, testnet accounts
- [ ] `app/genesis.go` — set constitutional genesis params
- [ ] `x/attestation` — full module skeleton with keeper, types, MsgServer
- [ ] `x/mempoolshield` — ABCI hooks, oracle stub, fail-closed logic
- [ ] `x/governor955` — treasury tracker, 95/5 enforcer, governance trigger
- [ ] `app/app.go` — register all three new modules
- [ ] `x/gov` params — wire v1.5 governance thresholds
- [ ] `docs/CONSTITUTION-v1.5.md` — pin constitutional law on-chain
- [ ] Tag `v0.1.0-testnet` once chain builds and produces blocks

---

*This document is the authoritative implementation contract between the Rampage
Constitution v1.5 / Whitepaper v1.5.1 and the Go codebase. All module authors
MUST implement to this spec. Deviations require a CHAIN-SPEC amendment commit.*
