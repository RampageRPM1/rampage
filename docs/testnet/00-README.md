# Rampage L1 — Testnet Deployment Kit

This folder contains the documentation, configuration, and operational
playbook for bringing Rampage L1 to a public testnet. It is documentation
only. Nothing here modifies the running chain or any node configuration
automatically.

## Status: PRE-LAUNCH — NOT YET JOINABLE

The public testnet is **not live yet**. Three launch-blocking issues must be
resolved before any external validator can join. These were verified directly
against this repository on 2026-06-09.

### Blocker 1 — Chain identity is contradictory

The authoritative spec and the live chain disagree on core identity fields:

| Field | CHAIN-SPEC-v1.5.1.md (target) | Live chain / committed genesis |
|-------|-------------------------------|--------------------------------|
| Chain ID | rampage-testnet-1 | rampage-1 |
| Bech32 prefix | rampage | ramp |
| Base denom | urpm | rpm |

These cannot both be canonical. Decision required: launch the public testnet
on the spec values (rampage-testnet-1 / rampage / urpm) and keep the existing
rampage-1 / ramp / rpm node running untouched as the dev/continuity chain.

### Blocker 2 — Published genesis file does not exist

VALIDATOR-GUIDE.md instructs validators to download genesis from
docs/static/genesis.json, but docs/static/ currently contains only
openapi.json. That download URL returns 404. A real genesis.json must be
generated on a clean host (no Ignite alice/bob accounts, capped faucet float,
urpm consistent across staking/mint) and committed with its SHA256.

### Blocker 3 — Validator guide has placeholder coordinates

The genesis SHA256, seed node IDs, and persistent peers are still TBD
placeholders and must be filled with real values before launch.

## Execution order

1. Ratify the canonical identity (Blocker 1).
2. Generate and publish docs/static/genesis.json + SHA256 (Blocker 2).
3. Fill validator-guide network coordinates (Blocker 3).
4. Provision VPS, harden, build rampaged.
5. Stand up nginx + public RPC/API/faucet endpoints.
6. Deploy explorer.
7. Announce.

## Safety constraints honored

- The live chain is never stopped or mutated.
- No chain commands were run against the home directory.
- No code under app/, cmd/, or x/ is modified by this kit.
