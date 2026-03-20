# Changelog

All notable changes to Rampage L1 will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

> **Note:** Rampage L1 is currently in pre-release development. Versions below 1.0.0 are prototype milestones. Stability and API compatibility are not guaranteed until v1.0.0.

---

## [Unreleased]

### In Progress
- Testnet infrastructure and validator onboarding
- Full genesis configuration for public testnet
- Block explorer integration
- IBC channel configuration for cross-chain connectivity
- Governance module stress testing
- Faucet service for testnet tokens

---

## [1.5.0-prototype] - 2025

### Added
- **x/governor module:** Custom governance module implementing Rampage's constitutional governance model with weighted voting, quorum enforcement, and proposal lifecycle management
- **Chain specification v1.5:** CHAIN-SPEC document detailing architecture, consensus parameters, validator economics, and governance constitution
- **IBC integration:** Inter-Blockchain Communication support via ibc-go for cross-chain asset transfers
- **CometBFT consensus:** Byzantine fault-tolerant consensus engine with configurable block time and validator set management
- **Cosmos SDK foundation:** Full Cosmos SDK v0.50.x application framework
- **Protobuf definitions:** Complete proto definitions for all custom modules
- **Ignite CLI scaffolding:** Initial chain scaffold with standard modules (auth, bank, staking, slashing, gov, mint, distribution, crisis)
- **LICENSE:** Apache License 2.0
- **NOTICE:** Copyright and intellectual property notice for Shea Patrick Kastl / Rampage L1
- **CONTRIBUTING.md:** Contributor guidelines including IP notice and CLA
- **SECURITY.md:** Security vulnerability reporting policy
- **CODE_OF_CONDUCT.md:** Community standards and enforcement guidelines

### Architecture Highlights
- L1 blockchain built on Cosmos SDK / CometBFT stack
- Custom governor module for on-chain constitutional governance
- Designed for censorship-resistant data verification and truth attestation
- Native token economics with staking and slashing
- Modular design for future extension

---

## [1.0.0-alpha] - 2025 (Initial Scaffold)

### Added
- Initial repository creation via Ignite CLI
- Base chain configuration
- Standard Cosmos SDK modules initialized
- Go module setup
- GitHub Actions CI workflow
- Makefile with standard build targets
- buf.yaml / buf.lock for protobuf management

---

## Version Roadmap

| Version | Status | Description |
|---------|--------|-------------|
| 1.5.0-prototype | Current | Architecture complete, testnet preparation |
| 1.6.0-testnet | Planned | Public testnet launch |
| 1.7.0-testnet | Planned | Governance testing, IBC channels live |
| 1.8.0-testnet | Planned | Security audit preparation |
| 2.0.0-rc | Planned | Release candidate, mainnet preparation |
| 2.0.0 | Future | Mainnet launch |

---

*Maintained by Shea Patrick Kastl. Copyright 2025 Shea Patrick Kastl. All rights reserved under Apache License 2.0.*
