# Rampage Public Testnet — Node Operations Kit

> **Status:** Pre-launch hand-off. Copy-paste configs for standing up the seed
> and public RPC (state-sync) nodes, plus a faucet stub. Use AFTER genesis is
> generated per [`01-GENESIS-BUILD.md`](./01-GENESIS-BUILD.md). Run on dedicated
> hosts — never on the live continuity/dev chain.

## Identity recap

| Field | Value |
|-------|-------|
| Chain ID | `rampage-testnet-1` |
| Binary | `rampaged` |
| Home dir | `~/.rampage` (service user: `rampage`) |
| P2P port | `26656` |
| RPC port | `26657` |
| Denom | `urpm` |

## 1. Common node setup (run on every host)

```bash
# Create a dedicated, non-login service user
sudo useradd -m -s /usr/sbin/nologin rampage

# Install the binary (built per 01-GENESIS-BUILD.md)
sudo install -m 0755 rampaged /usr/local/bin/rampaged

# Initialize and drop in the published genesis
sudo -u rampage rampaged init "$MONIKER" --chain-id rampage-testnet-1 --home /home/rampage/.rampage
sudo -u rampage curl -sL \
  https://raw.githubusercontent.com/RampageRPM1/rampage/main/docs/static/genesis.json \
  -o /home/rampage/.rampage/config/genesis.json

# Verify the genesis hash matches the published checksum
sha256sum /home/rampage/.rampage/config/genesis.json
# compare against docs/static/genesis.sha256
```

## 2. systemd service (`/etc/systemd/system/rampaged.service`)

```ini
[Unit]
Description=Rampage L1 node (rampage-testnet-1)
After=network-online.target
Wants=network-online.target

[Service]
User=rampage
Group=rampage
ExecStart=/usr/local/bin/rampaged start --home /home/rampage/.rampage
Restart=on-failure
RestartSec=3
LimitNOFILE=65535

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now rampaged
sudo journalctl -u rampaged -f
```

## 3. Seed node — `config.toml` overrides

A seed node only does peer discovery (PEX); it does not stay connected.

```toml
# [p2p]
seed_mode = true
pex = true
laddr = "tcp://0.0.0.0:26656"
# Set this to the node's public address so peers can dial back:
external_address = "tcp://<SEED_PUBLIC_IP>:26656"
max_num_inbound_peers = 240
max_num_outbound_peers = 40
```

Get the seed's node-id (this is what validators put in their `seeds` list):

```bash
sudo -u rampage rampaged tendermint show-node-id --home /home/rampage/.rampage
# Final seed string: <node-id>@<SEED_PUBLIC_IP>:26656
```

## 4. Public RPC + state-sync provider — `config.toml` / `app.toml`

Run >= 2 of these so validators can state-sync from a quorum.

```toml
# config.toml -> [rpc]
laddr = "tcp://0.0.0.0:26657"
# Allow public read access (front with TLS/nginx in production):
cors_allowed_origins = ["*"]

# config.toml -> base / [statesync]
# Keep recent snapshots available for peers to state-sync FROM this node
# (snapshot settings live in app.toml below)
```

```toml
# app.toml -> [state-sync]
snapshot-interval = 1000
snapshot-keep-recent = 5

# app.toml -> base
minimum-gas-prices = "0.001urpm"
```

## 5. Validators: state-sync config block to publish in VALIDATOR-GUIDE.md

Once two RPCs are live, capture a recent trust height/hash and publish this
ready-to-paste block for new validators:

```bash
RPC1="https://rpc-1.rampage.example:443"
LATEST=$(curl -s "$RPC1/block" | jq -r '.result.block.header.height')
TRUST_HEIGHT=$((LATEST - 1000))
TRUST_HASH=$(curl -s "$RPC1/block?height=$TRUST_HEIGHT" | jq -r '.result.block_id.hash')
echo "height=$TRUST_HEIGHT hash=$TRUST_HASH"
```

```toml
# config.toml -> [statesync]
enable = true
rpc_servers = "https://rpc-1.rampage.example:443,https://rpc-2.rampage.example:443"
trust_height = <TRUST_HEIGHT>
trust_hash = "<TRUST_HASH>"
trust_period = "168h0m0s"
```

## 6. Faucet stub (testnet only)

Dispenses 1 RPM (`1000000urpm`) per request from the `truth_reserve` account,
matching `config.yml`'s faucet block. This is a minimal reference — add rate
limiting / captcha before exposing publicly.

```bash
#!/usr/bin/env bash
# faucet.sh — usage: ./faucet.sh <recipient-rampage-address>
set -euo pipefail

CHAIN_ID="rampage-testnet-1"
FAUCET_KEY="truth_reserve"      # key must exist in the faucet host keyring
AMOUNT="1000000urpm"           # 1 RPM
NODE="https://rpc-1.rampage.example:443"
FEES="2000urpm"

RECIPIENT="${1:?usage: faucet.sh <rampage-address>}"

rampaged tx bank send "$FAUCET_KEY" "$RECIPIENT" "$AMOUNT" \
  --chain-id "$CHAIN_ID" \
  --node "$NODE" \
  --fees "$FEES" \
  --keyring-backend file \
  --yes
```

Optional tiny HTTP wrapper (Flask) for a web faucet:

```python
# faucet_server.py
import subprocess, re
from flask import Flask, request, jsonify

app = Flask(__name__)
ADDR_RE = re.compile(r"^rampage1[0-9a-z]{38}$")

@app.post("/request")
def request_funds():
    addr = (request.json or {}).get("address", "")
    if not ADDR_RE.match(addr):
        return jsonify(error="invalid rampage address"), 400
    r = subprocess.run(["./faucet.sh", addr], capture_output=True, text=True)
    if r.returncode != 0:
        return jsonify(error=r.stderr.strip()), 500
    return jsonify(status="sent", amount="1000000urpm", to=addr)
```

## 7. Smoke test before announcing

```bash
# Node is producing/syncing blocks?
curl -s localhost:26657/status | jq -r '.result.sync_info.latest_block_height'

# Reachable from OUTSIDE the host network (run from a different machine):
curl -s https://rpc-1.rampage.example:443/status | jq -r '.result.node_info.network'
# expected: rampage-testnet-1

# Faucet works end-to-end:
./faucet.sh rampage1<some-test-address>
```

When seeds + >= 2 state-sync RPCs are reachable externally and the faucet
dispenses, return to [`01-GENESIS-BUILD.md`](./01-GENESIS-BUILD.md) Step 6 and
remove the "TESTNET NOT YET LIVE" banner from `docs/VALIDATOR-GUIDE.md`.
