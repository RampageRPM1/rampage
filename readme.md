# Rampage — Censorship-Resistant Truth Engine

> **Chain:** `rampage-testnet-1` | **Version:** v1.5.1 | **Status:** Active Testnet Development  
> **Token:** RPM (`urpm` on-chain) | **Supply:** 21,000,000 RPM  
> **Constitution:** v1.5 | **Whitepaper:** v1.5.1

Rampage is a sovereign Cosmos SDK L1 blockchain purpose-built as a censorship-resistant truth verification and humanitarian reporting infrastructure. It is governed by the Rampage Constitution v1.5 and implements the Seven-Seal oracle verification system.

---

## Architecture Overview

### Core Principles
- **Human rights primacy** — Articles I–IV of the Constitution are immutable on-chain
- **Nonviolence** — civilian journalism and humanitarian access only
- **Truth verification** — Seven-Seal AI oracle consensus with 5-of-7 signature threshold
- **Liveness guarantee** — the chain NEVER stops; capital restrictions affect only capital-routing txs

### Chain Identity

| Parameter | Value |
|---|---|
| Chain ID | `rampage-testnet-1` |
| Binary | `rampaged` |
| Bech32 Prefix | `rampage` |
| Base Denom | `urpm` (1,000,000 urpm = 1 RPM) |
| Total Supply | 21,000,000 RPM |
| Consensus | CometBFT (1/3 Byzantine fault tolerance) |
| Block Target | < 7 second finality |

---

## Custom Modules

### `x/attestation` — Seven-Seal Truth Attestation
Stores on-chain consensus records from the Seven-Seal off-chain AI oracle network. Requires ≥ 5-of-7 oracle node signatures to finalize an attestation. Supports corrections with full immutable audit trail. Verification levels 1–4 (single source through on-the-ground verified).

### `x/mempoolshield` — Compliance Membrane
ABCI `PrepareProposal`/`ProcessProposal` hooks that screen all capital-routing transactions (bank sends, IBC transfers) through a 5-of-7 oracle signer committee. Fail-closed on oracle unavailability — capital txs are rejected but **block production is never halted**. Implements 4-level threat matrix including Level 4 / Article VIII conflict-zone suspension.

### `x/governor955` — Operational Split Enforcer
Tracks all treasury outflows and enforces the 95% / 5% truth verification / humanitarian access split mandated by Constitution Art. VI. Publishes the current ratio on-chain every 1,000 blocks. Triggers mandatory governance review if the ratio drifts for 18+ months. The Legal Defense Fund (5% of total supply) is ring-fenced.

---

## Governance (Constitution v1.5)

| Class | Quorum | Approval | Voting Period |
|---|---|---|---|
| Standard | 51% | 67% | 7 days |
| Constitutional Amendment | 60% | 80% | 30 days |
| Emergency Legal Compliance | 40% | 75% | 24 hours |
| Humanitarian Escalation (Art. VIII) | 50% | 75% | 30 days |
| Level 4 / Conflict Mode | 75% | 80% | 30 days |

**Articles I–IV are immutable** and cannot be altered by any governance proposal.

---

## Token Allocation

| Allocation | RPM | % |
|---|---|---|
| Truth Verification Reserve | 10,500,000 | 50% |
| Legal Defense Fund | 1,050,000 | 5% |
| Ecosystem / Developer Fund | 4,200,000 | 20% |
| Community Airdrop | 2,100,000 | 10% |
| Humanitarian Access Reserve | 1,050,000 | 5% |
| Founding Contributors (2yr vest) | 2,100,000 | 10% |

---

## Staking

- **Truth Bearer minimum:** 100,000 RPM
- **Unbonding period:** 21 days
- **Max validators (testnet):** 21
- **Slashing:** 5% verification misconduct, 10% capital routing violation, 5% double-sign

---

## Development

### Prerequisites
- Go 1.24+
- [Ignite CLI](https://ignite.com/cli)

### Run Testnet

```bash
git clone https://github.com/RampageRPM1/rampage
cd rampage
ignite chain serve
```

### Build Binary

```bash
ignite chain build
```

---

## Repository Structure

```
rampage/
├── app/                    # App wiring (register new modules here)
├── x/
│   ├── attestation/        # Seven-Seal truth attestation module
│   ├── mempoolshield/      # Compliance membrane + ABCI hooks
│   ├── governor955/        # 95/5 operational split enforcer
│   └── rampage/            # Base module
├── docs/
│   └── CHAIN-SPEC-v1.5.1.md  # Authoritative implementation contract
├── config.yml              # Testnet config (urpm, rampage-testnet-1)
└── readme.md               # This file
```

---

## Links

- [TruthOracle.ai](https://truthoracle.ai) — Seven-Seal verification dashboard
- [Rampage Constitution v1.5](docs/CHAIN-SPEC-v1.5.1.md)
- [Whitepaper v1.5.1](docs/) — *uploading*

---

*Rampage is an early-stage prototype testnet. All parameters are subject to change prior to mainnet launch.*
