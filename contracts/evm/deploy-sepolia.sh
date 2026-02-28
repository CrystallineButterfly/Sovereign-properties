#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo "╔══════════════════════════════════════════════════════════════╗"
echo "║                     HouseRWA Sepolia Deploy                 ║"
echo "╚══════════════════════════════════════════════════════════════╝"
echo ""

RPC_URL="${SEPOLIA_RPC:-${RPC_URL:-}}"
if [[ -z "$RPC_URL" ]]; then
  echo -e "${RED}ERROR:${NC} Set SEPOLIA_RPC (or RPC_URL) before running deploy-sepolia.sh."
  exit 1
fi

if [[ -z "${PRIVATE_KEY:-}" ]]; then
  echo -e "${RED}ERROR:${NC} PRIVATE_KEY is required."
  exit 1
fi

if [[ ! "$PRIVATE_KEY" =~ ^0x[0-9a-fA-F]{64}$ ]]; then
  echo -e "${RED}ERROR:${NC} PRIVATE_KEY must be 0x + 64 hex chars."
  exit 1
fi

validate_address() {
  local value="$1"
  local label="$2"
  if [[ ! "$value" =~ ^0x[0-9a-fA-F]{40}$ ]]; then
    echo -e "${RED}ERROR:${NC} ${label} must be a valid 0x-prefixed 20-byte address."
    exit 1
  fi
  if [[ "$value" == "0x0000000000000000000000000000000000000000" ]]; then
    echo -e "${RED}ERROR:${NC} ${label} cannot be the zero address."
    exit 1
  fi
}

DEPLOYER_ADDRESS="$(cast wallet address --private-key "$PRIVATE_KEY")"
FEE_RECIPIENT="${FEE_RECIPIENT:-$DEPLOYER_ADDRESS}"
CRE_WORKFLOW_ADDRESS="${CRE_WORKFLOW_ADDRESS:-}"
CRE_FORWARDER_ADDRESS="${CRE_FORWARDER_ADDRESS:-}"
CRE_STRICT_MODE="${CRE_STRICT_MODE:-false}"

if [[ -z "$CRE_WORKFLOW_ADDRESS" ]]; then
  echo -e "${RED}ERROR:${NC} CRE_WORKFLOW_ADDRESS is required for Sepolia deployments."
  exit 1
fi

if [[ -z "$CRE_FORWARDER_ADDRESS" ]]; then
  echo -e "${RED}ERROR:${NC} CRE_FORWARDER_ADDRESS is required for Sepolia deployments."
  exit 1
fi

validate_address "$FEE_RECIPIENT" "FEE_RECIPIENT"
validate_address "$CRE_WORKFLOW_ADDRESS" "CRE_WORKFLOW_ADDRESS"
validate_address "$CRE_FORWARDER_ADDRESS" "CRE_FORWARDER_ADDRESS"

echo -e "${BLUE}Configuration${NC}"
echo "  Network: Sepolia (11155111)"
echo "  Deployer: $DEPLOYER_ADDRESS"
echo "  Fee recipient: $FEE_RECIPIENT"
echo "  CRE workflow: $CRE_WORKFLOW_ADDRESS"
echo "  CRE forwarder: $CRE_FORWARDER_ADDRESS"
echo "  Strict mode: $CRE_STRICT_MODE"
echo ""

echo -e "${BLUE}Checking HouseRWA runtime size (EIP-170)...${NC}"
BYTECODE_HEX="$(forge inspect src/HouseRWA.sol:HouseRWA deployedBytecode)"
RUNTIME_SIZE_BYTES="$(( (${#BYTECODE_HEX} - 2) / 2 ))"
if (( RUNTIME_SIZE_BYTES > 24576 )); then
  echo -e "${RED}ERROR:${NC} HouseRWA runtime is ${RUNTIME_SIZE_BYTES} bytes (> 24,576)."
  echo -e "${YELLOW}This contract will revert on Sepolia/Mainnet due to EIP-170.${NC}"
  echo "For local demo only, use: ./testing/scripts/run-anvil-cutover.sh"
  exit 1
fi
echo -e "${GREEN}✓ Runtime size OK (${RUNTIME_SIZE_BYTES} bytes)${NC}"
echo ""

echo -e "${BLUE}Checking RPC connectivity...${NC}"
CHAIN_ID="$(cast chain-id --rpc-url "$RPC_URL")"
if [[ "$CHAIN_ID" != "11155111" ]]; then
  echo -e "${RED}ERROR:${NC} RPC returned chain id $CHAIN_ID (expected 11155111 for Sepolia)."
  exit 1
fi
echo -e "${GREEN}✓ Sepolia RPC is reachable${NC}"
echo ""

echo -e "${BLUE}Running contract tests...${NC}"
forge test -q
echo -e "${GREEN}✓ Tests passed${NC}"
echo ""

VERIFY_ARGS=()
if [[ -n "${ETHERSCAN_API_KEY:-}" ]]; then
  VERIFY_ARGS=(--verify)
  echo -e "${GREEN}✓ Etherscan verification enabled${NC}"
else
  echo -e "${YELLOW}⚠ ETHERSCAN_API_KEY not set; skipping verification${NC}"
fi
echo ""

echo -e "${BLUE}Broadcasting deployment...${NC}"
PRIVATE_KEY="$PRIVATE_KEY" \
CRE_WORKFLOW_ADDRESS="$CRE_WORKFLOW_ADDRESS" \
CRE_FORWARDER_ADDRESS="$CRE_FORWARDER_ADDRESS" \
FEE_RECIPIENT="$FEE_RECIPIENT" \
CRE_STRICT_MODE="$CRE_STRICT_MODE" \
forge script script/Deploy.s.sol:Deploy \
  --rpc-url "$RPC_URL" \
  --broadcast \
  -vvvv \
  "${VERIFY_ARGS[@]}"

LATEST_FILE="$(ls -1t deployments/houserwa_11155111_*.json 2>/dev/null | head -n 1 || true)"
if [[ -z "$LATEST_FILE" ]]; then
  echo -e "${YELLOW}⚠ Deployment completed but no Sepolia deployment file was found.${NC}"
  exit 0
fi

echo ""
echo -e "${GREEN}Deployment artifact:${NC} $LATEST_FILE"
echo -e "${BLUE}Update CRE + UI with:${NC}"
echo "  - proxyAddress  -> houseRWAContractAddr / VITE_HOUSE_RWA_ADDRESS"
echo "  - receiverAddress -> houseRWAReceiverAddr / VITE_CRE_RECEIVER_ADDRESS"
