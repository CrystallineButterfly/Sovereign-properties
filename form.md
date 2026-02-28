# Submission Form Draft

## Project name
Sovereign-Properties

## 1 line project description (under ~80-100 characters)
Private RWA housing platform using Chainlink CRE for private sale, rent, billing, and key handoff.

## Full project description
Sovereign-Properties is a private real-estate RWA platform that tokenizes houses and automates the full lifecycle of ownership and tenancy. Users mint homes as onchain assets, create sale/rental listings, execute buy/rent flows, and manage bills, while sensitive metadata and document access remain private.

The core problem it solves is that traditional onchain real-estate workflows expose too much information publicly. Sovereign-Properties uses a commitment + encrypted key exchange model so private property data is not written in plaintext onchain. Only authorized participants (owner, buyer, renter, allowed counterparties) can access private document data through authenticated, role-gated app/API flows.

Chainlink CRE is used to orchestrate offchain workflow logic and secure EVM actions for minting, listing, sale, rental, billing, and key-claim flows.

## How is it built?
- Smart contracts: Solidity (Foundry), `HouseRWA` + receiver architecture
- Workflow engine: Chainlink CRE (Go workflows + HTTP trigger routing)
- Backend: Node/Go services for auth/session, role-gated private data, workflow trigger mediation
- Frontend: React + TypeScript web app with wallet auth (SIWE/Privy) and protected user flows
- Privacy design: onchain commitments + offchain encrypted payloads + recipient-bound encrypted key exchange
- Infra/testing: Sepolia + Anvil simulation, CRE workflow simulations, auth + workflow smoke tests

## What challenges did you run into?
- Hardening workflow auth so privileged CRE actions require authenticated actor binding
- Removing insecure fallback auth behavior and handling session expiry cleanly on frontend
- Keeping privacy guarantees while preserving usable buy/rent/bill UX
- Environment consistency issues (RPC/contract alignment, stale key hashes across runs)
- Determinism and reliability concerns in workflow paths (time/random and external dependency handling)
- CSP/network integration friction around third-party wallet/KYC/browser extension environments

## Link to project repo
REPO_URL_HERE

## Chainlink Usage
Primary Chainlink CRE usage is in:
- `RWA-House-CRE/workflow.yaml`
- `RWA-House-CRE/main.go`
- `backend/cre/workflow.yaml`
- `backend/cre/handlers/http.go`
- `backend/cre/workflows/mint.go`
- `backend/cre/workflows/sale.go`
- `backend/cre/workflows/rental.go`
- `backend/cre/workflows/payment.go`
- `contracts/evm/src/HouseRWAReceiver.sol`

Repo-link form (replace base URL once public):
- `REPO_URL_HERE/tree/main/RWA-House-CRE/workflow.yaml`
- `REPO_URL_HERE/tree/main/backend/cre/handlers/http.go`

## Project Demo
VIDEO_URL_HERE (must be public and under 5 minutes)

## Which Chainlink prize track(s) are you applying to?
- DeFi and Tokenization
- Privacy

## Which sponsor track(s) are you applying to?
TBD

## Submitter name
K42 (Kell)

## Submitter email
k42.radicle.eth@ethermail.io

## Are you participating in a team or individually?
Individual
