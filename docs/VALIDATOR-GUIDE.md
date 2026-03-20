# Rampage L1 Validator Guide

> **A comprehensive guide for running a Rampage L1 validator node on the testnet.**

---

## Table of Contents

1. [Overview](#overview)
2. [Hardware Requirements](#hardware-requirements)
3. [Software Prerequisites](#software-prerequisites)
4. [Building from Source](#building-from-source)
5. [Node Initialization](#node-initialization)
6. [Genesis Configuration](#genesis-configuration)
7. [Joining the Testnet](#joining-the-testnet)
8. [Creating a Validator](#creating-a-validator)
9. [Validator Operations](#validator-operations)
10. [Monitoring](#monitoring)
11. [Security Best Practices](#security-best-practices)
12. [Troubleshooting](#troubleshooting)

---

## Overview

Rampage L1 validators are responsible for proposing and validating blocks on the network. The chain runs CometBFT (Tendermint) Byzantine fault-tolerant consensus, where validators stake RPM tokens and participate in block production proportional to their stake weight.

Validators who act honestly earn block rewards and transaction fee shares. Validators who double-sign, go offline, or otherwise violate protocol rules are subject to slashing penalties.

---

## Hardware Requirements

### Minimum (Testnet)

| Resource | Specification |
|----------|---------------|
| CPU | 4 cores / 8 threads (x86_64) |
| RAM | 16 GB |
| Storage | 500 GB SSD (NVMe recommended) |
| Network | 100 Mbps symmetric |
| OS | Ubuntu 22.04 LTS or Debian 12 |

### Recommended (Production / Mainnet)

| Resource | Specification |
|----------|---------------|
| CPU | 8+ cores / 16 threads |
| RAM | 32 GB |
| Storage | 1 TB NVMe SSD |
| Network | 1 Gbps symmetric |
| OS | Ubuntu 22.04 LTS |

> **Note:** Storage requirements will grow over time as the chain produces blocks. Plan for expansion.

---

## Software Prerequisites

- **Go 1.21+**: [Download](https://go.dev/dl/)
- **make**: Build toolchain
- **git**: Version control
- **jq**: JSON processing (for script helpers)
- **curl**: HTTP client

### Install Go (Ubuntu/Debian)

```bash
wget https://go.dev/dl/go1.21.6.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
echo 'export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin' >> ~/.bashrc
source ~/.bashrc
go version
```

---

## Building from Source

```bash
# Clone the repository
git clone https://github.com/RampageRPM1/rampage.git
cd rampage

# Build the chain binary
make build

# Verify the build
./build/rampaged version

# Install to $GOPATH/bin
make install
rampaged version
```

---

## Node Initialization

```bash
# Set your moniker (validator display name)
export MONIKER="your-validator-name"
export CHAIN_ID="rampage-testnet-1"

# Initialize the node
rampaged init $MONIKER --chain-id $CHAIN_ID

# This creates:
# ~/.rampage/config/config.toml     - CometBFT configuration
# ~/.rampage/config/app.toml        - Application configuration
# ~/.rampage/config/genesis.json    - Genesis file (will be replaced)
# ~/.rampage/data/                  - Chain data directory
```

### Key Generation

```bash
# Create a new validator key
rampaged keys add validator --keyring-backend test

# IMPORTANT: Save the mnemonic phrase securely!
# This key controls your validator and staked funds.

# View your validator address
rampaged keys show validator -a --keyring-backend test

# View your validator operator address
rampaged keys show validator --bech val -a --keyring-backend test
```

> **Security Warning:** On mainnet, use `--keyring-backend os` or a hardware security module (HSM). The `test` backend is for testnet only.

---

## Genesis Configuration

```bash
# Download the official testnet genesis file
curl -o ~/.rampage/config/genesis.json \
  https://raw.githubusercontent.com/RampageRPM1/rampage/main/docs/static/genesis.json

# Verify the genesis file hash
sha256sum ~/.rampage/config/genesis.json
# Expected: [will be published with testnet launch]
```

---

## Joining the Testnet

### Configure Peers

```bash
# Set persistent peers in config.toml
# Seed nodes will be published at testnet launch
SEEDS="[seed-node-id]@[seed-ip]:26656"
sed -i "s/seeds = \"\"/seeds = \"$SEEDS\"/" ~/.rampage/config/config.toml

# Set minimum gas prices
sed -i 's/minimum-gas-prices = ""/minimum-gas-prices = "0.001urpm"/' ~/.rampage/config/app.toml
```

### Start the Node

```bash
# Start syncing
rampaged start

# Or run as a systemd service (recommended)
sudo tee /etc/systemd/system/rampaged.service > /dev/null <<EOF
[Unit]
Description=Rampage L1 Node
After=network-online.target

[Service]
User=$USER
ExecStart=$(which rampaged) start
Restart=always
RestartSec=3
LimitNOFILE=65535

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable rampaged
sudo systemctl start rampaged

# Monitor logs
journalctl -u rampaged -f
```

### Verify Sync Status

```bash
# Check if the node is catching up
rampaged status | jq '.SyncInfo.catching_up'
# Should return `false` when fully synced

# Check latest block height
rampaged status | jq '.SyncInfo.latest_block_height'
```

---

## Creating a Validator

> **Wait until your node is fully synced before creating a validator.**

```bash
rampaged tx staking create-validator \
  --amount=1000000urpm \
  --pubkey=$(rampaged tendermint show-validator) \
  --moniker=$MONIKER \
  --chain-id=$CHAIN_ID \
  --commission-rate="0.10" \
  --commission-max-rate="0.20" \
  --commission-max-change-rate="0.01" \
  --min-self-delegation="1" \
  --gas="auto" \
  --gas-adjustment=1.5 \
  --gas-prices="0.001urpm" \
  --from=validator \
  --keyring-backend=test
```

### Verify Your Validator

```bash
# Check if your validator is in the active set
rampaged query staking validators --output json | \
  jq '.validators[] | select(.description.moniker=="'$MONIKER'")'

# Check validator signing info
rampaged query slashing signing-info $(rampaged tendermint show-validator)
```

---

## Validator Operations

### Delegate Additional Tokens

```bash
rampaged tx staking delegate \
  $(rampaged keys show validator --bech val -a --keyring-backend test) \
  500000urpm \
  --from=validator \
  --chain-id=$CHAIN_ID \
  --gas=auto \
  --keyring-backend=test
```

### Unjail a Validator

If your validator is jailed for downtime:

```bash
rampaged tx slashing unjail \
  --from=validator \
  --chain-id=$CHAIN_ID \
  --gas=auto \
  --keyring-backend=test
```

### Withdraw Rewards

```bash
rampaged tx distribution withdraw-rewards \
  $(rampaged keys show validator --bech val -a --keyring-backend test) \
  --commission \
  --from=validator \
  --chain-id=$CHAIN_ID \
  --gas=auto \
  --keyring-backend=test
```

---

## Monitoring

### Prometheus Metrics

Enable Prometheus in `config.toml`:

```toml
[instrumentation]
prometheus = true
prometheus_listen_addr = ":26660"
```

### Key Metrics to Monitor

- **Block height** – Ensure the node is keeping up with the network
- **Missed blocks** – Excessive missed blocks lead to jailing
- **Peer count** – Maintain adequate peer connections
- **Memory and disk usage** – Plan capacity accordingly
- **Validator uptime** – Target 99.9%+ uptime

### Health Check

```bash
# Quick status check
curl -s localhost:26657/status | jq '.result.sync_info'

# Check number of peers
curl -s localhost:26657/net_info | jq '.result.n_peers'
```

---

## Security Best Practices

1. **Key Management:** Use a separate machine for key storage. Consider HSM for mainnet.
2. **Firewall:** Only expose port 26656 (P2P). Block RPC (26657) and API (1317) from public access.
3. **Sentry Nodes:** Run sentry nodes in front of your validator to mitigate DDoS.
4. **Updates:** Monitor the repository for security patches and upgrade promptly.
5. **Backups:** Regularly back up `~/.rampage/config/priv_validator_key.json` and `~/.rampage/data/priv_validator_state.json`.
6. **SSH Hardening:** Use key-based authentication only. Disable root login.
7. **Monitoring Alerts:** Set up alerts for missed blocks, low disk space, and node crashes.

---

## Troubleshooting

### Node Won't Start

```bash
# Check logs for errors
journalctl -u rampaged --no-pager -n 50

# Reset chain data (if corrupt)
rampaged tendermint unsafe-reset-all
# Then re-download genesis and restart
```

### Node Stuck Syncing

- Verify persistent peers are correct and reachable
- Check firewall allows outbound connections on port 26656
- Ensure system clock is synchronized (use NTP)

### Validator Jailed

- Wait for the jailing period to expire
- Ensure the node is fully synced and producing blocks
- Run the unjail command above

---

## Testnet Faucet

Testnet RPM tokens for staking are available through the faucet (URL to be published at testnet launch).

---

*Maintained by Shea Patrick Kastl. Copyright 2025 Shea Patrick Kastl. For questions, open a GitHub Discussion or Issue.*
