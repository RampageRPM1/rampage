# Rampage L1 — Testnet Deployment Kit

This folder contains the documentation, configuration, and operational
playbook for bringing Rampage L1 to a public testnet. It is documentation
only. Nothing here modifies the running chain or any node configuration
automatically.

## Canonical chain identity (RATIFIED 2026-06-09)

The project is standardized on the repository/spec identity, matching
config.yml and CHAIN-SPEC-v1.5.1.md:

| Field | Canonical value |
|-------|-----------------|
| Chain ID | rampage-testnet-1 |
| Bech32 prefix | rampage (rampagevaloper / rampagevalcons) |
| Base denom | urpm |
| Display denom | rpm (1 RPM = 1,000,000 urpm) |

config.yml, CHAIN-SPEC-v1.5.1.md, and VALIDATOR-GUIDE.md already use these
values, so no source changes are required to adopt them.

### Legacy dev/continuity chain

The node that has been running continuously on the founder's Mac uses the
older Ignite identity (chain-id rampage-1, prefix ramp, denom rpm). It is
NOT the public testnet. It is preserved, untouched and never stopped, as the
dev/continuity chain. The public testnet rampage-testnet-1 is launched
separately from config.yml, which already produces the canonical identity
above. This honors the standing requirement that the live chain is never
halted.

## Status: PRE-LAUNCH — NOT YET JOINABLE

The public testnet is not live yet. Two launch-blocking issues remain. These
were verified directly against this repository on 2026-06-09.

### Blocker 1 — Published genesis file does not exist

VALIDATOR-GUIDE.md instructs validators to download genesis from
docs/static/genesis.json, but docs/static/ currently contains only
openapi.json. That download URL returns 404. A real genesis.json must be
generated on a clean host from config.yml (no Ignite alice/bob accounts,
capped faucet float, urpm consistent across staking/mint) and committed with
its SHA256.

### Blocker 2 — Validator guide has placeholder coordinates

The genesis SHA256, seed node IDs, and persistent peers are still TBD
placeholders and must be filled with real values before launch.

## Execution order

1. (DONE) Ratify canonical identity: rampage-testnet-1 / rampage / urpm.
2. Generate and publish docs/static/genesis.json + SHA256 (Blocker 1).
3. Fill validator-guide network coordinates (Blocker 2).
4. Provision VPS, harden, build rampaged.
5. Stand up nginx + public RPC/API/faucet endpoints.
6. Deploy explorer.
7. Announce.

## Safety constraints honored

- The live chain is never stopped or mutated.
- No chain commands were run against the home directory.
- No code under app/, cmd/, or x/ is modified by this kit.
