# Rampage Public Testnet — Genesis Build Runbook

> **Status:** Pre-launch hand-off. Run this on a CLEAN host (not the live
> continuity/dev chain). Produces `docs/static/genesis.json` and its SHA256.
> Target launch window: within 1-2 weeks of this document's merge.

## Canonical identity (do not deviate)

| Field | Value |
|-------|-------|
| Chain ID | `rampage-testnet-1` |
| Bech32 prefix | `rampage` |
| Base denom | `urpm` (1 RPM = 1_000_000 urpm) |
| Binary | `rampaged` |
| Home dir | `~/.rampage` |
| Source of truth | `config.yml` + `docs/CHAIN-SPEC-v1.5.1.md` |

The legacy live Mac chain (`rampage-1` / `ramp` / `rpm`) is the continuity/dev
chain ONLY. Never export its genesis to the public testnet.

## Prerequisites (clean host)

- Ubuntu 22.04 LTS (or equivalent), 4 vCPU / 16 GB RAM / 200 GB SSD
- Go 1.21+ and `make` installed
- `jq` and `sha256sum` installed
- Fresh clone of this repo; NO existing `~/.rampage` directory

```bash
# Wipe any prior state so nothing leaks from a dev box
rm -rf ~/.rampage
```

## Step 1 — Build the binary from source

```bash
git clone https://github.com/RampageRPM1/rampage.git
cd rampage
make install            # builds and installs `rampaged`
rampaged version        # confirm it runs
```

## Step 2 — Generate genesis from config.yml via Ignite

The authoritative parameters live in `config.yml` (urpm denom, 21 max
validators, 21-day unbonding, gov/staking/slashing params, and the custom
`mempoolshield`, `governor955`, and `attestation` modules).

```bash
# Ignite reads config.yml and scaffolds a chain-id-correct genesis
ignite chain init --home ~/.rampage

# Verify the chain id was applied
jq -r '.chain_id' ~/.rampage/config/genesis.json
# expected: rampage-testnet-1
```

If you prefer a manual cosmos-sdk flow instead of Ignite:

```bash
rampaged init <moniker> --chain-id rampage-testnet-1 --home ~/.rampage
# then add genesis accounts, gentx, and collect-gentxs per config.yml amounts
```

## Step 3 — Validate the genesis

```bash
rampaged genesis validate --home ~/.rampage

# Spot-check the critical invariants:
jq -r '.chain_id' ~/.rampage/config/genesis.json
jq -r '.app_state.staking.params.bond_denom' ~/.rampage/config/genesis.json   # urpm
jq -r '.app_state.staking.params.max_validators' ~/.rampage/config/genesis.json # 21
jq -r '.app_state.staking.params.unbonding_time' ~/.rampage/config/genesis.json # 1814400s
jq -e '.app_state.mempoolshield' ~/.rampage/config/genesis.json > /dev/null && echo mempoolshield OK
jq -e '.app_state.governor955' ~/.rampage/config/genesis.json > /dev/null && echo governor955 OK
jq -e '.app_state.attestation' ~/.rampage/config/genesis.json > /dev/null && echo attestation OK
```

All checks must pass before publishing.

## Step 4 — Publish genesis to the repo (via PR, never direct to main)

```bash
cp ~/.rampage/config/genesis.json docs/static/genesis.json
sha256sum docs/static/genesis.json | tee docs/static/genesis.sha256

git checkout -b publish-testnet-genesis
git add docs/static/genesis.json docs/static/genesis.sha256
git commit -m "feat(testnet): publish rampage-testnet-1 genesis + sha256"
git push origin publish-testnet-genesis
# open a PR into main (CODEOWNERS will request your review automatically)
```

The published genesis will then be reachable at:
`https://raw.githubusercontent.com/RampageRPM1/rampage/main/docs/static/genesis.json`

## Step 5 — Stand up seed / RPC nodes and record coordinates

Before announcing the testnet, run at least one seed node and one public RPC
with state-sync enabled, then fill these into `docs/VALIDATOR-GUIDE.md`:

- **Seeds:** `<node-id>@<host>:26656` (run `rampaged tendermint show-node-id`)
- **State-sync RPC servers:** two public RPC endpoints, e.g.
  `https://rpc-1.rampage.example:443,https://rpc-2.rampage.example:443`
- **Trust height / trust hash:** from a recent block on the seed RPC:
  ```bash
  curl -s <rpc>/block | jq -r '.result.block.header.height'
  curl -s <rpc>/block?height=<height> | jq -r '.result.block_id.hash'
  ```

## Step 6 — Final go-live checklist

- [ ] `docs/static/genesis.json` published and reachable (no 404)
- [ ] `docs/static/genesis.sha256` published; hash matches
- [ ] Seeds populated in VALIDATOR-GUIDE.md (real node IDs, not placeholders)
- [ ] State-sync RPC + trust height/hash populated
- [ ] Faucet endpoint live (truth_reserve, 1 RPM/request per config.yml)
- [ ] "TESTNET NOT YET LIVE" banner removed from VALIDATOR-GUIDE.md
- [ ] At least 1 seed + 1 RPC reachable from outside the host network

Once all boxes are checked, the network is live and external validators can
follow the Validator Guide end-to-end.
