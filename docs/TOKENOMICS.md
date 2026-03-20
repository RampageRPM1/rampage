# Rampage L1 Tokenomics

> **RPM Token Economics and Distribution Model**

---

## Token Overview

| Parameter | Value |
|-----------|-------|
| **Token Name** | RPM |
| **Base Denomination** | urpm (1 RPM = 1,000,000 urpm) |
| **Chain** | Rampage L1 (Sovereign Cosmos SDK Chain) |
| **Consensus** | CometBFT (Delegated Proof-of-Stake) |
| **Initial Supply** | Defined at genesis |
| **Inflation Model** | Dynamic, bounded |

---

## Token Utility

The RPM token serves multiple critical functions within the Rampage L1 ecosystem:

### 1. Network Security (Staking)
- Validators and delegators stake RPM to secure the network
- Staked RPM earns block rewards and a share of transaction fees
- Staking creates economic alignment between validators and network health

### 2. Governance
- RPM holders participate in on-chain governance through the `x/governor` module
- Voting power is proportional to staked RPM
- Governance proposals include parameter changes, treasury allocation, and protocol upgrades
- Constitutional constraints prevent governance attacks on core protocol parameters

### 3. Transaction Fees
- All transactions on Rampage L1 require RPM for gas fees
- Minimum gas price is enforced by validators
- Fee revenue is distributed to validators and delegators

### 4. Truth Attestation
- RPM is required to submit and stake on truth attestations
- Economic incentives align attestors with accuracy
- Slashing penalties for false attestation protect data integrity

---

## Inflation Model

Rampage L1 uses a dynamic inflation model similar to the Cosmos Hub, with bounded parameters:

| Parameter | Value |
|-----------|-------|
| **Inflation Min** | 7% |
| **Inflation Max** | 20% |
| **Inflation Rate Change** | 13% per year |
| **Bonded Target** | 67% |
| **Block Reward Distribution** | Proportional to stake |

### How It Works

- When the bonded ratio (staked tokens / total supply) falls below 67%, inflation increases to incentivize staking
- When the bonded ratio exceeds 67%, inflation decreases
- This creates a self-regulating economic mechanism that targets network security

---

## Staking Economics

### Validator Rewards

Validators earn rewards from two sources:

1. **Block Rewards:** New RPM minted per block, distributed proportionally to stake weight
2. **Transaction Fees:** All gas fees collected in each block

### Commission

Validators set a commission rate on rewards earned by their delegators:

| Parameter | Constraint |
|-----------|------------|
| Initial Commission | Set by validator (0-100%) |
| Max Commission | Set at validator creation (immutable ceiling) |
| Max Commission Change Rate | Maximum daily commission increase |

### Delegator Returns

Delegators receive staking rewards minus the validator's commission:

```
Delegator Reward = (Delegator Stake / Validator Total Stake) * Block Rewards * (1 - Commission Rate)
```

---

## Slashing

Slashing protects the network by penalizing malicious or negligent validators:

| Violation | Penalty | Jail Duration |
|-----------|---------|---------------|
| **Double Signing** | 5% of staked RPM | Permanent (tombstoned) |
| **Downtime** (missing >95% of blocks in window) | 0.01% of staked RPM | 10 minutes |

### Slashing Impact

- Slashed tokens are burned (removed from supply)
- Delegators staked with a slashed validator also lose proportional stake
- Double signing results in permanent removal from the validator set

---

## Unbonding

| Parameter | Value |
|-----------|-------|
| **Unbonding Period** | 21 days |
| **Redelegation** | Instant (with cooldown) |

During the unbonding period:
- Tokens do not earn rewards
- Tokens are still subject to slashing
- After 21 days, tokens become liquid and transferable

---

## Governance Economics

### Proposal Deposits

- Governance proposals require a minimum RPM deposit
- Deposits are returned if the proposal passes or fails to reach quorum
- Deposits are burned if the proposal is vetoed (>33% NoWithVeto)

### Voting Parameters

| Parameter | Value |
|-----------|-------|
| Voting Period | 14 days |
| Quorum | 33.4% of bonded stake |
| Threshold | 50% Yes votes |
| Veto Threshold | 33.4% NoWithVeto |

### Constitutional Constraints

The `x/governor` module enforces constitutional constraints on governance:
- Certain parameters cannot be modified via governance (protocol constants)
- Emergency proposals have shorter voting periods but higher quorum requirements
- Validator weight decay for non-participating validators

---

## Fee Structure

### Gas Pricing

| Parameter | Value |
|-----------|-------|
| Minimum Gas Price | 0.001 urpm |
| Gas Limit Per Block | Configurable |
| Fee Burn Rate | 0% (all fees to validators) |

### Fee Distribution

- Block proposer receives a bonus share of fees
- Remaining fees are distributed proportionally to all bonded validators
- Community pool receives a configurable percentage

---

## Community Pool

A percentage of all block rewards and fees is directed to the community pool:

- Controlled by governance proposals
- Used for ecosystem development, grants, and infrastructure
- Cannot be spent without an approved governance proposal

---

## Token Distribution (Genesis)

The genesis token allocation will be published prior to mainnet launch. The distribution is designed to:

- Ensure sufficient decentralization of the validator set
- Fund core development and ecosystem growth
- Incentivize early validators and testnet participants
- Establish a community pool for ongoing development

Detailed allocation percentages will be included in the genesis ceremony documentation.

---

## Economic Security Model

Rampage L1's economic security is derived from:

1. **Stake-weighted consensus:** The cost of attacking the network is proportional to the value of staked RPM
2. **Slashing penalties:** Malicious behavior results in permanent economic loss
3. **Unbonding period:** Prevents rapid destabilization of the validator set
4. **Constitutional governance:** Prevents economic parameter manipulation through governance attacks
5. **Dynamic inflation:** Self-regulating mechanism to maintain target security ratio

---

> **Disclaimer:** Tokenomics parameters are subject to change through governance. Values listed reflect the initial configuration and may be modified via approved governance proposals, subject to constitutional constraints.

---

*Copyright 2025 Shea Patrick Kastl. Rampage L1 tokenomics are proprietary design elements protected under the NOTICE file. See [NOTICE](../NOTICE) for details.*
