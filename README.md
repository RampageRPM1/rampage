# Rampage L1

> **Censorship-resistant truth engine built on a sovereign Layer 1 blockchain.**

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](go.mod)
[![Status](https://img.shields.io/badge/Status-Prototype%20v1.5-orange)](CHANGELOG.md)
[![Network](https://img.shields.io/badge/Network-Testnet%20Preparation-yellow)](docs/)

---

## Overview

Rampage L1 is a sovereign Layer 1 blockchain protocol designed for censorship-resistant data verification, truth attestation, and constitutional on-chain governance. Built on the Cosmos SDK and CometBFT consensus engine, Rampage L1 provides a permissionless, decentralized foundation for applications requiring verifiable, immutable truth records.

**Key Properties:**
- **Censorship Resistance** – No single actor can suppress or alter submitted data
- **Constitutional Governance** – On-chain governance bound by an immutable constitutional framework
- **Sovereign Chain** – Independent L1, not a rollup or sidechain
- **IBC Ready** – Native Inter-Blockchain Communication for cross-chain interoperability
- **Decentralized Validation** – Byzantine fault-tolerant consensus via CometBFT

---

## Architecture

```
Rampage L1
├── Consensus Layer      CometBFT (Byzantine Fault Tolerant)
├── Application Layer    Cosmos SDK v0.50.x
├── Core Modules
│   ├── x/governor       Custom constitutional governance module
│   ├── x/bank           Token transfers
│   ├── x/staking        Validator staking and delegation
│   ├── x/slashing       Validator accountability
│   ├── x/mint           Token emission
│   ├── x/distribution   Reward distribution
│   └── x/ibc            Inter-Blockchain Communication
├── cmd/rampaged         Chain daemon binary
├── app/                 Application wiring
├── proto/               Protobuf definitions
└── docs/                Chain specification and documentation
```

### x/governor Module

The `x/governor` module is Rampage's proprietary governance innovation. Unlike standard Cosmos SDK governance, the governor module implements:

- **Constitutional Constraints** – Proposals that violate constitutional parameters are automatically rejected
- **Weighted Voting** – Validator vote weight accounts for stake, participation history, and constitutional standing
- **Quorum Enforcement** – Configurable quorum requirements per proposal type
- **Proposal Lifecycle Management** – Full lifecycle from submission through execution with veto and emergency provisions

---

## Current Status

| Component | Status |
|-----------|--------|
| Core chain scaffold | Complete |
| x/governor module | Complete |
| Protobuf definitions | Complete |
| Chain specification v1.5 | Complete |
| IBC integration | Complete |
| Testnet genesis configuration | In Progress |
| Block explorer integration | Planned |
| Public testnet launch | Planned |
| Security audit | Pre-mainnet |
| Mainnet | Future |

**Current Version:** v1.5.0-prototype

---

## Getting Started

### Prerequisites

- Go 1.21 or higher
- `make`
- `ignite` CLI (optional, for scaffolding)
- `buf` (for protobuf compilation)

### Building from Source

```bash
# Clone the repository
git clone https://github.com/RampageRPM1/rampage.git
cd rampage

# Build the chain binary
make build

# Install the binary
make install
```

### Running a Local Development Node

```bash
# Initialize a local chain
rampaged init mynode --chain-id rampage-devnet-1

# Start the node
rampaged start
```

### Running Tests

```bash
make test
```

---

## Documentation

| Document | Description |
|----------|-------------|
| [Whitepaper-v1.5.2.md](docs/Whitepaper-v1.5.2.md) | Rampage L1 Technical Whitepaper — authoritative architecture and governance specification |
 | [docs/CHAIN-SPEC-v1.5.md](docs/) | Full chain specification including architecture, consensus parameters, validator economics, and governance constitution |
| [CONTRIBUTING.md](CONTRIBUTING.md) | How to contribute to Rampage L1 |
| [SECURITY.md](SECURITY.md) | Security policy and vulnerability reporting |
| [CHANGELOG.md](CHANGELOG.md) | Version history and roadmap |
| [NOTICE](NOTICE) | Copyright and trademark information |

---

## Roadmap

### Testnet Phase (Current)
- [ ] Genesis configuration finalization
- [ ] Validator onboarding process
- [ ] Faucet deployment
- [ ] Block explorer integration
- [ ] IBC channel configuration

### Pre-Mainnet Phase
- [ ] External security audit
- [ ] Governance stress testing
- [ ] Economic parameter validation
- [ ] Community validator set expansion

### Mainnet
- [ ] Genesis ceremony
- [ ] Mainnet launch
- [ ] Ecosystem development

---

## Contributing

Contributions are welcome under the terms described in [CONTRIBUTING.md](CONTRIBUTING.md). All contributors must agree to the Rampage Contributor License Agreement.

**Please note:** Rampage L1 is the proprietary intellectual property of Shea Patrick Kastl. The source code is made available under the Apache License 2.0, but the Rampage name, brand, architecture concepts, and governance model are protected. See [NOTICE](NOTICE) and [CONTRIBUTING.md](CONTRIBUTING.md) for details.

---

## Security

For security vulnerabilities, **do not open a public issue**. See [SECURITY.md](SECURITY.md) for responsible disclosure procedures.

---

## License

Copyright 2025 Shea Patrick Kastl

Licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE) for full terms.

This project incorporates components from the Cosmos SDK, CometBFT, IBC-Go, and Ignite CLI, all under Apache License 2.0. See [NOTICE](NOTICE) for full attribution.

---

## Contact

Project maintained by **Shea Patrick Kastl**

For partnership inquiries, licensing questions, or permissions beyond the scope of this license, contact through the GitHub repository.

---

*Rampage L1 – Building the infrastructure for verifiable truth.*
