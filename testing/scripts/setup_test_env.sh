#!/bin/bash
set -e

echo "🧪 Setting up RWA House Platform Test Environment"
echo "=================================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check prerequisites
check_command() {
    if ! command -v $1 &> /dev/null; then
        echo -e "${RED}Error: $1 is not installed${NC}"
        exit 1
    fi
}

echo "Checking prerequisites..."
check_command node
check_command npm
check_command go
check_command forge
check_command anvil

echo -e "${GREEN}All prerequisites met!${NC}"

# Get script directory
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Function to cleanup on exit
cleanup() {
    echo ""
    echo "Cleaning up..."
    if [ -n "$ANVIL_PID" ]; then
        kill $ANVIL_PID 2>/dev/null || true
    fi
    if [ -n "$BACKEND_PID" ]; then
        kill $BACKEND_PID 2>/dev/null || true
    fi
    if [ -n "$FRONTEND_PID" ]; then
        kill $FRONTEND_PID 2>/dev/null || true
    fi
}
trap cleanup EXIT

# Start Anvil (local Ethereum node)
echo ""
echo "Starting Anvil local blockchain..."
anvil --fork-url "${RPC_URL:-https://eth-mainnet.g.alchemy.com/v2/demo}" \
      --block-time 2 \
      --accounts 10 \
      --balance 10000 &
ANVIL_PID=$!

# Wait for Anvil to start
sleep 5

# Verify Anvil is running
if ! kill -0 $ANVIL_PID 2>/dev/null; then
    echo -e "${RED}Failed to start Anvil${NC}"
    exit 1
fi

echo -e "${GREEN}Anvil running on PID $ANVIL_PID${NC}"

# Deploy contracts
echo ""
echo "Deploying smart contracts..."
cd "$PROJECT_ROOT/contracts/evm"

# Create .env if it doesn't exist
if [ ! -f .env ]; then
    cat > .env << EOF
RPC_URL=http://localhost:8545
PRIVATE_KEY=0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
ETHERSCAN_API_KEY=test
EOF
fi

# Deploy
forge script script/Deploy.s.sol \
    --rpc-url http://localhost:8545 \
    --broadcast \
    --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80

echo -e "${GREEN}Contracts deployed!${NC}"

# Start backend CRE workflow
echo ""
echo "Starting backend CRE workflow..."
cd "$PROJECT_ROOT/backend/cre"

# Create config if it doesn't exist
if [ ! -f config.local.json ]; then
    cat > config.local.json << EOF
{
  "ethereum": {
    "rpc_url": "http://localhost:8545",
    "chain_id": 31337,
    "contract_address": "0x...",
    "private_key": "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
  },
  "encryption": {
    "threshold": 3,
    "total_shares": 5
  },
  "server": {
    "port": 8080,
    "host": "localhost"
  }
}
EOF
fi

go build -o cre-workflow .
./cre-workflow &
BACKEND_PID=$!

sleep 3

echo -e "${GREEN}Backend running on PID $BACKEND_PID${NC}"

# Start frontend web
echo ""
echo "Starting frontend web..."
cd "$PROJECT_ROOT/RWA-House-UI/web"

# Create .env.local if it doesn't exist
if [ ! -f .env.local ]; then
    cat > .env.local << EOF
VITE_API_URL=http://localhost:8080
VITE_RPC_URL=http://localhost:8545
VITE_HOUSE_RWA_ADDRESS=0x...
VITE_CHAIN_ID=31337
VITE_EXPECTED_CHAIN_ID=31337
EOF
fi

npm install
npm run dev &
FRONTEND_PID=$!

sleep 5

echo -e "${GREEN}Frontend running on PID $FRONTEND_PID${NC}"

# Print summary
echo ""
echo "=================================================="
echo -e "${GREEN}Test Environment Ready!${NC}"
echo "=================================================="
echo ""
echo "Services:"
echo "  - Anvil (Blockchain): http://localhost:8545"
echo "  - Backend API:        http://localhost:8080"
echo "  - Frontend Web:       http://localhost:5173"
echo ""
echo "Test Accounts (Anvil defaults):"
echo "  - Account 0: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266 (10,000 ETH)"
echo "  - Account 1: 0x70997970C51812dc3A010C7d01b50e0d17dc79C8 (10,000 ETH)"
echo "  - Account 2: 0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC (10,000 ETH)"
echo ""
echo "To run tests:"
echo "  npm run test:contracts  # Smart contract tests"
echo "  npm run test:backend    # Backend tests"
echo "  npm run test:e2e        # E2E tests"
echo "  npm run test:security   # Security tests"
echo ""
echo "Press Ctrl+C to stop all services"
echo ""

# Wait for interrupt
wait
