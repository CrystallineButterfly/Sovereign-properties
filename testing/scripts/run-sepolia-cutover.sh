#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
DEPLOY_DIR="$ROOT_DIR/testing/deployment"
TIMESTAMP="$(date -u +"%Y%m%dT%H%M%SZ")"
LOG_FILE="$DEPLOY_DIR/cutover_${TIMESTAMP}.log"
SUMMARY_FILE="$DEPLOY_DIR/cutover_${TIMESTAMP}.md"

mkdir -p "$DEPLOY_DIR"
exec > >(tee -a "$LOG_FILE") 2>&1

require_env() {
  local key="$1"
  if [[ -z "${!key:-}" ]]; then
    echo "ERROR: ${key} is required."
    exit 1
  fi
}

validate_private_key() {
  if [[ ! "$PRIVATE_KEY" =~ ^0x[0-9a-fA-F]{64}$ ]]; then
    echo "ERROR: PRIVATE_KEY must be 0x + 64 hex characters."
    exit 1
  fi
}

echo "=== RWA Sepolia + CRE cutover started @ ${TIMESTAMP} ==="

require_env PRIVATE_KEY
require_env SEPOLIA_RPC
require_env CRE_WORKFLOW_ADDRESS
require_env CRE_FORWARDER_ADDRESS
validate_private_key

echo ""
echo "[1/10] Sepolia RPC check"
CHAIN_ID="$(cast chain-id --rpc-url "$SEPOLIA_RPC")"
if [[ "$CHAIN_ID" != "11155111" ]]; then
  echo "ERROR: Expected Sepolia chain id 11155111, got ${CHAIN_ID}."
  exit 1
fi
echo "Sepolia RPC OK (chain id ${CHAIN_ID})"

echo ""
echo "[2/10] Contract tests"
cd "$ROOT_DIR/contracts/evm"
forge test -q

echo ""
echo "[3/10] Contract deployment"
./deploy-sepolia.sh

LATEST_DEPLOYMENT="$(ls -1t deployments/houserwa_11155111_*.json 2>/dev/null | head -n 1 || true)"
if [[ -z "$LATEST_DEPLOYMENT" ]]; then
  echo "ERROR: No Sepolia deployment artifact found in contracts/evm/deployments."
  exit 1
fi

readarray -t DEPLOY_VALUES < <(
  python - "$LATEST_DEPLOYMENT" <<'PY'
import json
import sys

with open(sys.argv[1], encoding="utf-8") as f:
    data = json.load(f)

print(data.get("proxyAddress", ""))
print(data.get("receiverAddress", ""))
print(data.get("implementationAddress", ""))
print(data.get("network", ""))
PY
)

PROXY_ADDRESS="${DEPLOY_VALUES[0]}"
RECEIVER_ADDRESS="${DEPLOY_VALUES[1]}"
IMPLEMENTATION_ADDRESS="${DEPLOY_VALUES[2]}"
NETWORK_NAME="${DEPLOY_VALUES[3]}"

if [[ -z "$PROXY_ADDRESS" || -z "$RECEIVER_ADDRESS" ]]; then
  echo "ERROR: Could not parse proxy/receiver from $LATEST_DEPLOYMENT."
  exit 1
fi

echo ""
echo "[4/10] Updating CRE config files"
python - "$ROOT_DIR" "$PROXY_ADDRESS" "$RECEIVER_ADDRESS" "$SEPOLIA_RPC" <<'PY'
import json
import os
import sys

root, proxy, receiver, rpc = sys.argv[1:]
targets = [
    os.path.join(root, "RWA-House-CRE", "config.staging.json"),
    os.path.join(root, "backend", "cre", "config", "config.staging.json"),
]

for path in targets:
    with open(path, encoding="utf-8") as f:
        data = json.load(f)
    data["houseRWAContractAddr"] = proxy
    data["houseRWAReceiverAddr"] = receiver
    if "rpcURL" in data:
        data["rpcURL"] = rpc
    with open(path, "w", encoding="utf-8") as f:
        json.dump(data, f, indent=2)
        f.write("\n")
    print(f"updated {path}")
PY

echo ""
echo "[5/10] Updating root and web env files"
ROOT_ENV="$ROOT_DIR/.env"
WEB_ENV="$ROOT_DIR/RWA-House-UI/web/.env"

if [[ ! -f "$WEB_ENV" ]]; then
  cp "$ROOT_DIR/RWA-House-UI/web/.env.example" "$WEB_ENV"
fi

for key in \
  HOUSE_RWA_CONTRACT_ADDRESS \
  CRE_CONTRACT_ADDRESS \
  CRE_RECEIVER_ADDRESS \
  VITE_HOUSE_RWA_ADDRESS \
  VITE_CRE_RECEIVER_ADDRESS \
  VITE_RPC_URL \
  VITE_EXPECTED_CHAIN_ID
do
  sed -i "/^${key}=/d" "$ROOT_ENV" "$WEB_ENV" 2>/dev/null || true
done

{
  echo "HOUSE_RWA_CONTRACT_ADDRESS=${PROXY_ADDRESS}"
  echo "CRE_CONTRACT_ADDRESS=${PROXY_ADDRESS}"
  echo "CRE_RECEIVER_ADDRESS=${RECEIVER_ADDRESS}"
} >> "$ROOT_ENV"

{
  echo "VITE_HOUSE_RWA_ADDRESS=${PROXY_ADDRESS}"
  echo "VITE_CRE_RECEIVER_ADDRESS=${RECEIVER_ADDRESS}"
  echo "VITE_RPC_URL=${SEPOLIA_RPC}"
  echo "VITE_EXPECTED_CHAIN_ID=11155111"
} >> "$WEB_ENV"

echo ""
echo "[6/10] CRE wasm compile checks"
cd "$ROOT_DIR/RWA-House-CRE"
GOCACHE=/tmp/go-build-cache GOOS=wasip1 GOARCH=wasm go test ./...
cd "$ROOT_DIR/backend/cre"
GOCACHE=/tmp/go-build-cache GOOS=wasip1 GOARCH=wasm go test ./...

echo ""
echo "[7/10] CRE auth check"
cd "$ROOT_DIR/RWA-House-CRE"
cre whoami

echo ""
echo "[8/10] CRE simulate (mint fixture)"
cre workflow simulate . --target staging-settings --non-interactive --trigger-index 0 \
  --http-payload @../backend/cre/simulations/mint.json

echo ""
echo "[9/10] CRE deploy"
cre workflow deploy . --target staging-settings

echo ""
echo "[10/10] Web build"
cd "$ROOT_DIR/RWA-House-UI/web"
npm run build

cat > "$SUMMARY_FILE" <<EOF
# Cutover Summary

- Timestamp: ${TIMESTAMP}
- Network: ${NETWORK_NAME}
- Chain ID: 11155111
- Proxy: ${PROXY_ADDRESS}
- Receiver: ${RECEIVER_ADDRESS}
- Implementation: ${IMPLEMENTATION_ADDRESS}
- Deployment artifact: ${LATEST_DEPLOYMENT}
- Full log: ${LOG_FILE}
EOF

echo ""
echo "=== Cutover complete ==="
echo "Summary: ${SUMMARY_FILE}"
echo "Log: ${LOG_FILE}"
