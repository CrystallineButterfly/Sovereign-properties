#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
CONTRACTS_DIR="$ROOT_DIR/contracts/evm"
WEB_ENV="$ROOT_DIR/RWA-House-UI/web/.env"
ROOT_ENV="$ROOT_DIR/.env"
DEPLOYMENT_DIR="$ROOT_DIR/testing/deployment"
TIMESTAMP="$(date -u +"%Y%m%dT%H%M%SZ")"
LOG_FILE="$DEPLOYMENT_DIR/anvil_cutover_${TIMESTAMP}.log"
SUMMARY_FILE="$DEPLOYMENT_DIR/anvil_cutover_${TIMESTAMP}.md"

ANVIL_PORT="${ANVIL_PORT:-8545}"
ANVIL_RPC="${ANVIL_RPC:-http://127.0.0.1:${ANVIL_PORT}}"
ANVIL_CHAIN_ID="${ANVIL_CHAIN_ID:-31337}"
ANVIL_CODE_SIZE_LIMIT="${ANVIL_CODE_SIZE_LIMIT:-200000}"
DEFAULT_ANVIL_PRIVATE_KEY="0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

mkdir -p "$DEPLOYMENT_DIR"
exec > >(tee -a "$LOG_FILE") 2>&1

start_anvil_if_needed() {
  if cast chain-id --rpc-url "$ANVIL_RPC" >/dev/null 2>&1; then
    echo "Anvil already running at $ANVIL_RPC"
    return
  fi

  echo "Starting Anvil on $ANVIL_RPC"
  echo "Chain=${ANVIL_CHAIN_ID}, code-size-limit=${ANVIL_CODE_SIZE_LIMIT}"
  nohup anvil \
    --host 127.0.0.1 \
    --port "$ANVIL_PORT" \
    --chain-id "$ANVIL_CHAIN_ID" \
    --code-size-limit "$ANVIL_CODE_SIZE_LIMIT" \
    --silent \
    > "$DEPLOYMENT_DIR/anvil_${TIMESTAMP}.log" 2>&1 &

  sleep 2
  cast chain-id --rpc-url "$ANVIL_RPC" >/dev/null
}

update_kv() {
  local file="$1"
  local key="$2"
  local value="$3"
  if [[ ! -f "$file" ]]; then
    return
  fi
  sed -i "/^${key}=/d" "$file"
  printf "%s=%s\n" "$key" "$value" >> "$file"
}

echo "=== Anvil cutover started @ ${TIMESTAMP} ==="
start_anvil_if_needed

CHAIN_ID="$(cast chain-id --rpc-url "$ANVIL_RPC")"
if [[ "$CHAIN_ID" != "$ANVIL_CHAIN_ID" ]]; then
  echo "ERROR: Expected chain id ${ANVIL_CHAIN_ID}, got ${CHAIN_ID}."
  exit 1
fi
echo "Anvil RPC check: chain id ${CHAIN_ID}"

PRIVATE_KEY="${PRIVATE_KEY:-$DEFAULT_ANVIL_PRIVATE_KEY}"
if [[ ! "$PRIVATE_KEY" =~ ^0x[0-9a-fA-F]{64}$ ]]; then
  echo "ERROR: PRIVATE_KEY must be 0x + 64 hex chars."
  exit 1
fi

DEPLOYER_ADDRESS="$(cast wallet address --private-key "$PRIVATE_KEY")"
echo "Using deployer: $DEPLOYER_ADDRESS"

echo "[1/6] Running contract tests"
cd "$CONTRACTS_DIR"
forge test -q

echo "[2/6] Deploying contracts to Anvil"
PRIVATE_KEY="$PRIVATE_KEY" \
CRE_WORKFLOW_ADDRESS="${CRE_WORKFLOW_ADDRESS:-$DEPLOYER_ADDRESS}" \
CRE_FORWARDER_ADDRESS="${CRE_FORWARDER_ADDRESS:-$DEPLOYER_ADDRESS}" \
FEE_RECIPIENT="${FEE_RECIPIENT:-$DEPLOYER_ADDRESS}" \
forge script script/Deploy.s.sol:Deploy \
  --rpc-url "$ANVIL_RPC" \
  --disable-code-size-limit \
  --broadcast \
  -vvvv

LATEST_DEPLOYMENT="$(ls -1t "$CONTRACTS_DIR"/deployments/houserwa_31337_*.json | head -n 1)"
readarray -t DEPLOY_VALUES < <(
  python - "$LATEST_DEPLOYMENT" <<'PY'
import json
import sys
with open(sys.argv[1], encoding="utf-8") as f:
    data = json.load(f)
print(data["proxyAddress"])
print(data["receiverAddress"])
print(data["implementationAddress"])
PY
)
PROXY_ADDRESS="${DEPLOY_VALUES[0]}"
RECEIVER_ADDRESS="${DEPLOY_VALUES[1]}"
IMPLEMENTATION_ADDRESS="${DEPLOY_VALUES[2]}"

echo "[3/6] Updating local env/config wiring"
update_kv "$ROOT_ENV" "HOUSE_RWA_CONTRACT_ADDRESS" "$PROXY_ADDRESS"
update_kv "$ROOT_ENV" "CRE_CONTRACT_ADDRESS" "$PROXY_ADDRESS"
update_kv "$ROOT_ENV" "CRE_RECEIVER_ADDRESS" "$RECEIVER_ADDRESS"
update_kv "$WEB_ENV" "VITE_HOUSE_RWA_ADDRESS" "$PROXY_ADDRESS"
update_kv "$WEB_ENV" "VITE_CRE_RECEIVER_ADDRESS" "$RECEIVER_ADDRESS"
update_kv "$WEB_ENV" "VITE_RPC_URL" "$ANVIL_RPC"
update_kv "$WEB_ENV" "VITE_EXPECTED_CHAIN_ID" "$ANVIL_CHAIN_ID"

python - "$ROOT_DIR" "$PROXY_ADDRESS" "$RECEIVER_ADDRESS" <<'PY'
import json
import os
import sys

root, proxy, receiver = sys.argv[1:]
targets = [
    os.path.join(root, "RWA-House-CRE", "config.anvil.json"),
    os.path.join(root, "backend", "cre", "config", "config.anvil.json"),
]

for path in targets:
    with open(path, encoding="utf-8") as f:
        data = json.load(f)
    data["houseRWAContractAddr"] = proxy
    data["houseRWAReceiverAddr"] = receiver
    data["rpcURL"] = "http://127.0.0.1:8545"
    with open(path, "w", encoding="utf-8") as f:
        json.dump(data, f, indent=2)
        f.write("\n")
    print(f"updated {path}")
PY

echo "[4/6] WASM compile checks"
cd "$ROOT_DIR/RWA-House-CRE"
GOCACHE=/tmp/go-build-cache GOOS=wasip1 GOARCH=wasm go test ./...
cd "$ROOT_DIR/backend/cre"
GOCACHE=/tmp/go-build-cache GOOS=wasip1 GOARCH=wasm go test ./...

echo "[5/6] CRE simulate (anvil-settings)"
cd "$ROOT_DIR"
SECRET_VALUE="${SECRET_VALUE:-local-dev-secret}" \
cre workflow simulate RWA-House-CRE --target anvil-settings --non-interactive --trigger-index 0 \
  --http-payload @"$ROOT_DIR/backend/cre/simulations/mint.json"

echo "[6/6] Web build"
cd "$ROOT_DIR/RWA-House-UI/web"
npm run build

cat > "$SUMMARY_FILE" <<EOF
# Anvil Cutover Summary

- Timestamp: ${TIMESTAMP}
- RPC: ${ANVIL_RPC}
- Chain ID: ${ANVIL_CHAIN_ID}
- Proxy: ${PROXY_ADDRESS}
- Receiver: ${RECEIVER_ADDRESS}
- Implementation: ${IMPLEMENTATION_ADDRESS}
- Deployment artifact: ${LATEST_DEPLOYMENT}
- Full log: ${LOG_FILE}
EOF

echo "=== Anvil cutover complete ==="
echo "Summary: ${SUMMARY_FILE}"
echo "Log: ${LOG_FILE}"
