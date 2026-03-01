# Sovereign-Properties (Chainlink CRE Hackathon)

Private real-estate RWA platform for tokenization, private document handoff, and CRE-driven sale, rental, and billing workflows.

## Use Case

This app supports private real-estate operations:

- Tokenize a property as an onchain RWA asset
- Create listings (set sale/rent price) through CRE workflow actions
- Run private sales and rentals
- Exchange encrypted access/document keys securely
- Create/pay recurring bills
- Apply optional KYC-gated workflow rules
- Expose private property data only to authorized participants in the app/API

---

## Stack / Architecture

### Onchain
- **Solidity + Foundry**
- `HouseRWA` (UUPS proxy) + `HouseRWAReceiver`
- Target network: **Sepolia** (live target), **Anvil** (local simulation)

### Offchain workflow
- **Chainlink CRE workflow** in Go
- HTTP trigger action routing (`mint`, `create_listing`, `sell`, `rent`, `create_bill`, `pay_bill`, `claim_key`)
- CRE EVM writes via receiver contract

### App
- **Web:** React + TypeScript + Privy wallet auth + external wallet connectivity
- **Backend services:** Go workflow runtime + optional ZKPassport verifier + `/workflow/trigger` API adapter

### External integrations
- Stripe API path for fiat/billing flow
- KYC verifier endpoint support (mock or ZKPassport-style proof flow)

### Privacy model (seller↔buyer, landlord↔renter)

The project uses a **commitment + encrypted key exchange** model:

1. **Mint writes commitments onchain, not full private metadata**
   - Mint computes a salted `documentHash` and a `metadataCommitment`.
   - Onchain mint stores commitment-style values and an opaque pointer like
     `cre://private/<nonce>`.
   - Private metadata is not published in plaintext onchain.

2. **Private metadata is stored offchain in encrypted form**
   - The workflow adapter stores private mint payloads in
     `backend/zkpassport-session-service/.workflow-private-store.json`.
   - Payloads are encrypted at rest (AES-GCM) and linked to token IDs.

3. **Frontend/API visibility is role-gated**
   - Authorized viewers include:
     - current owner (buyer after sale),
     - original owner/minter (seller),
     - active renter,
     - `allowedBuyer` for private listings.
   - Unauthorized `GET /houses/:id` is redacted.
   - Unauthorized `GET /houses/:id/documents` and `/houses/:id/bills` return `403`.

4. **Seller→buyer and landlord→renter key delivery**
   - Sale and rental paths create encrypted key exchanges bound to an
     intended recipient.
   - `claim_key` enforces claimant/recipient matching in workflow/API.
   - The key material exposed by claim paths is ciphertext only.

5. **Session auth for private routes**
   - `/auth/verify-wallet` verifies wallet signatures and issues signed bearer
     tokens using `WORKFLOW_AUTH_SECRET`.
   - Private house routes and `claim_key` require authenticated callers.
   - Private document content can be configured to require KYC via
     `REQUIRE_KYC_FOR_PRIVATE_DOCUMENTS` (default: `false`).

6. **Optional KYC mode for CRE actions**
   - CRE action payloads now accept `kycProvider=none`.
   - In `none` mode, the workflow skips KYC proof verification and skips
     `setKYCVerification` writes.
   - This enables fully anonymous demo flows while keeping auth + role checks.

### Privacy guarantees checklist

| Guarantee | How it is enforced | Evidence pointers |
|---|---|---|
| House private metadata is not written in plaintext onchain (CRE mint path) | Mint stores `metadataCommitment` + salted `documentHash` + opaque `cre://private/...` pointer | `backend/cre/handlers/http.go` (`handleMint`), `backend/zkpassport-session-service/workflow-trigger.cjs` (`handleMint`) |
| Offchain private metadata is encrypted at rest | Private store payload is encrypted (AES-GCM) before persistence | `backend/zkpassport-session-service/workflow-trigger.cjs` (`encryptPrivatePayload`, `decryptPrivatePayload`) |
| Unauthorized users cannot read private house details from app API | House reads are redacted; private docs/bills endpoints return `403` for non-authorized viewers | `backend/zkpassport-session-service/server.cjs` (`projectHouseForViewer`, `/houses/:id`, `/houses/:id/documents`, `/houses/:id/bills`) |
| Seller→buyer / landlord→renter key handoff is recipient-bound | `claim_key` validates claimant against intended recipient before returning ciphertext | `backend/zkpassport-session-service/workflow-trigger.cjs` (`handleClaimKey`, claim auth gate), `backend/cre/handlers/http.go` (`handleClaimKey`) |
| Private route access requires authenticated wallet sessions | Wallet signature verification issues signed bearer tokens; private routes parse/verify bearer wallet | `backend/zkpassport-session-service/workflow-trigger.cjs` (`handleVerifyWallet`, token signing/parsing), `backend/zkpassport-session-service/server.cjs` (`extractViewerWalletAddress`) |
| Anonymous CRE execution is supported | `kycProvider=none` short-circuits KYC verification + onchain KYC writes | `backend/zkpassport-session-service/workflow-trigger.cjs` (`ensureKYCFromPayload`), `backend/cre/handlers/kyc.go`, `backend/cre/handlers/http.go` (`writeKYCVerification`) |
| CRE mint privacy flow validated in simulation | Manual mint simulate run returned `success: true` with `private onchain commitment` message and CRE write success statuses | Evidence section: **Manual CRE mint privacy simulation (February 23, 2026, 00:28 UTC)** |

**Scope note:** ownership, listing terms, rental state, and billing state are still onchain-visible by design.
CRE does not hide already-written onchain data; it prevents plaintext private data from being written onchain.
Private property metadata and encrypted document/access keys are protected by the workflow/API controls above.

---

## Chainlink File Index

These are the project files that implement Chainlink CRE / Chainlink-connected behavior.

### CRE workflow definition + runtime
- [`RWA-House-CRE/workflow.yaml`](RWA-House-CRE/workflow.yaml)
- [`RWA-House-CRE/main.go`](RWA-House-CRE/main.go)
- [`project.yaml`](project.yaml)

### Backend CRE implementation
- [`backend/cre/workflow.yaml`](backend/cre/workflow.yaml)
- [`backend/cre/main.go`](backend/cre/main.go)
- [`backend/cre/handlers/http.go`](backend/cre/handlers/http.go)
- [`backend/cre/handlers/kyc.go`](backend/cre/handlers/kyc.go)
- [`backend/cre/workflows/mint.go`](backend/cre/workflows/mint.go)
- [`backend/cre/workflows/sale.go`](backend/cre/workflows/sale.go)
- [`backend/cre/workflows/rental.go`](backend/cre/workflows/rental.go)
- [`backend/cre/workflows/payment.go`](backend/cre/workflows/payment.go)
- [`backend/cre/pkg/evm/client.go`](backend/cre/pkg/evm/client.go)
- [`backend/cre/config/config.go`](backend/cre/config/config.go)

### EVM receiver + deployment wiring for CRE writes
- [`contracts/evm/src/HouseRWAReceiver.sol`](contracts/evm/src/HouseRWAReceiver.sol)
- [`contracts/evm/src/keystone/IReceiver.sol`](contracts/evm/src/keystone/IReceiver.sol)
- [`contracts/evm/src/HouseRWA.sol`](contracts/evm/src/HouseRWA.sol)
- [`contracts/evm/script/Deploy.s.sol`](contracts/evm/script/Deploy.s.sol)
- [`contracts/evm/deploy-sepolia.sh`](contracts/evm/deploy-sepolia.sh)

### CRE simulation payloads
- [`backend/cre/simulations/mint.json`](backend/cre/simulations/mint.json)
- [`backend/cre/simulations/create_listing.json`](backend/cre/simulations/create_listing.json)
- [`backend/cre/simulations/sell.json`](backend/cre/simulations/sell.json)
- [`backend/cre/simulations/rent.json`](backend/cre/simulations/rent.json)
- [`backend/cre/simulations/create_bill.json`](backend/cre/simulations/create_bill.json)
- [`backend/cre/simulations/pay_bill.json`](backend/cre/simulations/pay_bill.json)
- [`backend/cre/simulations/claim_key.json`](backend/cre/simulations/claim_key.json)

### Local validation automation used to generate evidence

The evidence below was produced with local project automation during testing.
Those scripts and generated artifacts are environment-specific and may not be
included in the public submission tree.

### App integration surfaces (calls into CRE-backed flows)
- [`RWA-House-UI/web/src/components/MintHouseForm.tsx`](RWA-House-UI/web/src/components/MintHouseForm.tsx)
- [`RWA-House-UI/web/src/components/ListingForm.tsx`](RWA-House-UI/web/src/components/ListingForm.tsx)
- [`RWA-House-UI/web/src/pages/MarketplacePage.tsx`](RWA-House-UI/web/src/pages/MarketplacePage.tsx)
- [`RWA-House-UI/web/src/pages/CreateBillPage.tsx`](RWA-House-UI/web/src/pages/CreateBillPage.tsx)
- [`RWA-House-UI/mobile/src/screens/MintScreen.tsx`](RWA-House-UI/mobile/src/screens/MintScreen.tsx)
- [`RWA-House-UI/mobile/src/screens/MarketplaceScreen.tsx`](RWA-House-UI/mobile/src/screens/MarketplaceScreen.tsx)

---

## Evidence

### CRE + auth smoke (automated, February 27, 2026, 01:34 UTC)

Generated in a local validation environment from repo root:

```bash
./testing/scripts/run-cre-auth-smoke.sh
```

Run artifacts (local validation outputs):
- `testing/deployment/cre_auth_smoke_20260227T013400Z.md`
- `testing/deployment/cre_auth_smoke_20260227T013400Z.log`

Observed results:
- Auth hardening checks:
  - unauthenticated `/workflow/trigger` write action blocked with `401` ✅
  - bearer token + mismatched `actorAddress` blocked with `403` ✅
  - bearer token + matching actor accepted (`200`, `success=true`) ✅
- CRE simulations in one pass:
  - `mint` ✅
  - `create_listing` ✅
  - `sell` ✅
  - `rent` ✅
  - `create_bill` ✅
  - `pay_bill` ✅
  - `claim_key` ✅
- Final summary: `Pass=15 Warn=0 Fail=0`

### Local CRE workflow simulation

- Run artifact summary (local output): `testing/deployment/anvil_cutover_20260219T045410Z.log`
- Includes successful CRE report-based writes, e.g.:
  - `txStatus=TX_STATUS_SUCCESS`
  - `receiverStatus=RECEIVER_CONTRACT_EXECUTION_STATUS_SUCCESS`

### Successful CRE CLI simulations (manual run on February 21, 2026)

Executed from `RWA-House-CRE/`:

```bash
export SECRET_VALUE=local-dev-secret
for p in mint create_listing sell rent create_bill pay_bill claim_key; do
  cre -R .. -e ../.env -T anvil-settings workflow simulate RWA-House-CRE \
    --non-interactive --trigger-index 0 \
    --http-payload @../backend/cre/simulations/${p}.json
done
```

Observed results:
- `mint` ✅ success (`TX_STATUS_SUCCESS`, `RECEIVER_CONTRACT_EXECUTION_STATUS_SUCCESS`)
- `create_listing` ✅ success (rerun on **February 21, 2026, 21:15 UTC**:
  `TX_STATUS_SUCCESS`, `RECEIVER_CONTRACT_EXECUTION_STATUS_SUCCESS`,
  message: `listing created successfully`)
- `sell` ✅ success (`TX_STATUS_SUCCESS`, `RECEIVER_CONTRACT_EXECUTION_STATUS_SUCCESS`)
- `rent` ✅ success (`TX_STATUS_SUCCESS`, `RECEIVER_CONTRACT_EXECUTION_STATUS_SUCCESS`)
- `create_bill` ✅ success (`TX_STATUS_SUCCESS`, `RECEIVER_CONTRACT_EXECUTION_STATUS_SUCCESS`)
- `pay_bill` ✅ success (`TX_STATUS_SUCCESS`, `RECEIVER_CONTRACT_EXECUTION_STATUS_SUCCESS`)
- `claim_key` ⚠ returned `key exchange not found` with the simulation payload; this is expected when the key hash
  does not exist in the currently simulated chain state.

### Manual CRE mint privacy simulation (February 23, 2026, 00:28 UTC)

Executed from `RWA-House-CRE/`:

```bash
cre -R .. -e ../.env -T anvil-settings workflow simulate RWA-House-CRE \
  --non-interactive --trigger-index 0 \
  --http-payload @../backend/cre/simulations/mint.json
```

Observed evidence:
- Workflow ran `action=mint` and submitted CRE EVM reports successfully.
- Write reply showed:
  - `txStatus=TX_STATUS_SUCCESS`
  - `receiverStatus=RECEIVER_CONTRACT_EXECUTION_STATUS_SUCCESS`
- Simulation response decoded to:
  - `success: true`
  - `message: "house minted successfully with private onchain commitment"`
  - `metadataCommitment: 0x2762e34d57be3ed1b67c6ae7edfdfb29bc3b2e69965ecd4013255d803f2ae9f8`
  - `documentURI: cre://private/4919bbe2d899f344d14647a19d1ff15c765f`
  - `sharesCount: 5`, `threshold: 3`
- `txHash` and `tokenId` were empty in this local simulation run, which is
  consistent with simulator behavior noted in logs.

### Sepolia deployment + verification evidence (February 21, 2026)

Deployed with `contracts/evm/deploy-sepolia.sh` and verified on Etherscan:

- **HouseRWA implementation:** `0x6d43697D2308b67784927e4E4387465429Ba47da`
  - https://sepolia.etherscan.io/address/0x6d43697D2308b67784927e4E4387465429Ba47da#code
- **HouseRWA proxy (use for app interactions):** `0x990e1EB2Dd8fA8007533Ab50bE262A44EEF172ee`
  - https://sepolia.etherscan.io/address/0x990e1EB2Dd8fA8007533Ab50bE262A44EEF172ee#code
- **HouseRWA receiver:** `0x65844014526C32Ef8e68a80CD99a01aA4588D5BA`
  - https://sepolia.etherscan.io/address/0x65844014526C32Ef8e68a80CD99a01aA4588D5BA#code

Deployment artifact:
- `contracts/evm/deployments/houserwa_11155111_1771707912.json`

### CRE simulation ↔ Sepolia interaction model

- The CRE workflow logic is executed and validated via `cre workflow simulate`
  (actions: mint/create_listing/sell/rent/create_bill/pay_bill).
- The same HouseRWA/receiver call path used in simulation is what the app uses against the Sepolia proxy/receiver.
- Because live CRE workflow onboarding is still early-access for this account (`cre account list-key` returned
  `No linked owners found`), current Sepolia wiring uses temporary EOA values for:
  - `CRE_WORKFLOW_ADDRESS`
  - `CRE_FORWARDER_ADDRESS`
  - `FEE_RECIPIENT`
- This preserves end-to-end testnet functionality now; once CRE owner/forwarder are issued, replace those env values,
  re-deploy (or update receiver forwarder), and enable strict mode if desired.

### Live `/workflow/trigger` smoke on Sepolia (February 21, 2026)

Against `http://localhost:8787` with current env:

- `mint` ✅ success  
  `txHash=0x077e99dc5092a066fe7384045c9b4f0b521214b132889768e3bdd60bb1fe77fb`  
  minted token `2`
- `create_listing` ✅ success  
  `txHash=0x3fcde3d43c2495e5844a13aba64eafaa1af46eded768564e45515ff25af6d17f`
- `create_bill` ✅ success  
  `txHash=0x1dab74e47602685ab7d61464ec7aef6d9f1e0ec3d93a48bfda3cee0a13da46a1`
- `pay_bill` ✅ success  
  `txHash=0x7c4489b6ac6da2c8ab843ffcd62a2aa192b8418f8e2ff04eff8b97fe328a03ad`

`/kyc/zkpassport/verify` and `/kyc/verify` routes also respond correctly with validation errors on malformed proofs,
confirming the verifier routes are wired for frontend + CRE usage.

### Live end-to-end smoke on Railway + Sepolia (February 22, 2026)

Against `https://zkpassport-api-production.up.railway.app`:

- `mint` ✅ success  
  `txHash=0x31017f9653f0d3c99941648bc69bf949da08d54a3b896086311c5ccf77772122`
- `create_listing` (for sale) ✅ success  
  `txHash=0xe61d6a9643f5154b0b2a00954852e16f21c04e6b15fd904dba049104c8bbd8b8`
- `sell` ✅ success  
  `txHash=0x9c09d6a2cf0236a62a59bdb36fc4875c85038c776884f9c89668639bae69375f`
- `claim_key` ✅ success (same encrypted key returned for intended recipient)
- `set_kyc` preflight ✅ success (added for rental UX reliability)
- `create_listing` (for rent) ✅ success  
  `txHash=0x671b9dbf480d59238c3ad41924c5f84e008a17d47c8b17f21367dc7646d1efc7`
- onchain `depositForRental` ✅ success  
  `txHash=0x811c549fc934c3c0bd1c2d11cf423818f64e5211881d12a14e2fb8bb07c767b0`
- `rent` ✅ success  
  `txHash=0x4658f30d421a92df5a07f06448dbdd5d1de6d822d27dd5802ceec37643551adf`
- `create_bill` ✅ success  
  `txHash=0x18b0ec8e881f6c8613b2638b9f7b69492919fe68b8cac295f765459cfb6e85da`
- `pay_bill` ✅ success  
  `txHash=0x8febc091e2a078f7f9d6d660968621731026f8c21f1812b31dde454124399163`

This confirms that the frontend-connected adapter and Sepolia contract path are
operating for all hackathon demo flows.

### Contract + workflow checks run

```bash
cd contracts/evm
forge build
forge test -q

cd /home/k42/Auditz/CLH/RWA-Houses
testing/scripts/run-anvil-cutover.sh
```

---

## Point-in-time Sepolia Status (February 21, 2026)

Sepolia deployment is complete and verified:

- `HouseRWA` runtime: **20,991 bytes**
- Sepolia EIP-170 limit: **24,576 bytes**
- Proxy address in use: `0x990e1EB2Dd8fA8007533Ab50bE262A44EEF172ee`
- Receiver address in use: `0x65844014526C32Ef8e68a80CD99a01aA4588D5BA`

The deployment script (`contracts/evm/deploy-sepolia.sh`) enforces this size check before broadcast.

To use CRE-driven listing/price actions (`create_listing`), deploy a build that includes
`HouseRWA.createListingFromWorkflow` and the receiver selector update in `Deploy.s.sol`.

---

## Local Run (fast path, local validation environment)

```bash
cd /home/k42/Auditz/CLH/RWA-Houses
testing/scripts/run-anvil-cutover.sh
```

This local validation flow performs:
- local Anvil deployment
- CRE workflow simulation
- env/config wiring
- web build

---

## Web API Adapter (`/workflow/trigger` + ZKPassport)

For hackathon demos, the web app can use `backend/zkpassport-session-service` as a lightweight API layer.

Start it with:

```bash
cd backend/zkpassport-session-service
HOST=0.0.0.0 PORT=8787 CORS_ORIGIN="http://localhost:3000,http://127.0.0.1:3000,http://localhost:5173,http://127.0.0.1:5173" node server.cjs
```

Implemented endpoints:
- `POST /workflow/trigger` (`mint`, `set_kyc`, `create_listing`, `sell`, `rent`, `create_bill`, `pay_bill`, `claim_key`)
- `POST /rpc` (JSON-RPC proxy to configured Sepolia RPC; recommended for browser-safe onchain reads)
- `POST /auth/verify-wallet`
- `POST /auth/refresh`
- `POST /auth/logout`
- `POST /kyc/zkpassport/verify` (frontend proof verification)
- `POST /kyc/verify` (CRE-style verifier response)
- `POST /kyc/zkpassport/session`
- `GET /kyc/zkpassport/session/:sessionId`
- `GET /healthz`

`/workflow/trigger` KYC behavior by payload:
- `kycProvider=none` → skip KYC proof + skip onchain KYC write
- `kycProvider=mock` → mock KYC path writes `setKYCVerification`
- `kycProvider=zkpassport` → requires `kycProof`, verifies proof, then writes KYC

When running locally, point the UI to this adapter:

- `VITE_API_URL=http://localhost:8787`
- `VITE_ZKPASSPORT_API_URL=http://localhost:8787`
- `VITE_RPC_URL=http://localhost:8787/rpc`
- `VITE_ENABLE_CHAIN_RPC_FALLBACK=false` (recommended for hosted/browser deployments)
- `VITE_ENABLE_PUBLIC_RPC_CANDIDATES=false` (prevents probing public RPC URLs that may be blocked by CSP)

If you intentionally enable browser direct-RPC fallback, set explicit allowed hosts:

- `VITE_ENABLE_CHAIN_RPC_FALLBACK=true`
- `VITE_FALLBACK_RPC_URLS=https://your-allowed-rpc.example,https://another-rpc.example`

If your browser still cannot reach `localhost:8787` (common with WSL/VM setups),
use your machine IP/hostname in `VITE_ZKPASSPORT_API_URL` instead.

If the UI shows `Unable to start ZKPassport flow`, verify:
- adapter health: `curl http://localhost:8787/healthz`
- session creation:  
  `curl -X POST http://localhost:8787/kyc/zkpassport/session -H 'Content-Type: application/json' -d '{"walletAddress":"0x1111111111111111111111111111111111111111"}'`

If the UI shows `ZKPassport proof verification failed`, verify:
- proof verify endpoint:
  `curl -X POST http://localhost:8787/kyc/zkpassport/verify -H 'Content-Type: application/json' -d '{"walletAddress":"0x1111111111111111111111111111111111111111","proof":{"proofs":[],"queryResult":{}}}'`

Dev fallback: the Vite app now proxies `/kyc/zkpassport/*` and `/healthz` to
`http://127.0.0.1:8787` (override with `ZKPASSPORT_PROXY_TARGET`) so local
button clicks can still work even if direct cross-origin calls are blocked.

### Required env for `/workflow/trigger`

- `WORKFLOW_RPC_URL` (or `CRE_RPC_URL` / `SEPOLIA_RPC`)
- `WORKFLOW_PRIVATE_KEY` (or `CRE_ETH_PRIVATE_KEY` / `PRIVATE_KEY`)
- `HOUSE_RWA_CONTRACT_ADDRESS` (Sepolia proxy)
- `WORKFLOW_AUTH_SECRET` (required for signed session tokens used by private house data routes)
- `WORKFLOW_MAX_BODY_BYTES` (optional, defaults to `125829120`) for larger mint
  document payloads from the web form
- `REQUIRE_KYC_FOR_PRIVATE_DOCUMENTS` (optional, defaults to `false`; set `true`
  to enforce KYC before private document content can be viewed)

Optional auth hardening flags:
- `WORKFLOW_AUTH_TOKEN_TTL_SECONDS` (defaults to `43200`, i.e. 12 hours)
- `WORKFLOW_ALLOW_INSECURE_BEARER=false` (default). Set `true` only for short-lived local demos.

The adapter auto-loads `backend/zkpassport-session-service/.env` and root `.env`
when those files exist.

Note: the adapter loads `ethers` from `RWA-House-UI/web/node_modules`, so run
`cd RWA-House-UI/web && npm install` before starting the adapter service.

### Recommended env for ZKPassport

- `ZKPASSPORT_DOMAIN` (defaults to `demo.zkpassport.id`)
- `ZKPASSPORT_SCOPE`
- `ZKPASSPORT_DEV_MODE=true` for hackathon/demo mode
- `ZKPASSPORT_VERIFY_WRITING_DIRECTORY=/tmp` (recommended in hosted/serverless environments)
- Optional:
  - `ZKPASSPORT_BRIDGE_URL`
  - `ZKPASSPORT_CLOUD_PROVER_URL`
  - `ZKPASSPORT_VALIDITY_SECONDS`

### Frontend ZKPassport env (Vercel/web)

- `VITE_ZKPASSPORT_API_URL` (backend verifier base URL)
- `VITE_ZKPASSPORT_DOMAIN`
- `VITE_ZKPASSPORT_SCOPE`
- `VITE_ZKPASSPORT_DEV_MODE` (`true` for demo)
- Optional branding:
  - `VITE_ZKPASSPORT_APP_NAME`
  - `VITE_ZKPASSPORT_APP_LOGO`
  - `VITE_ZKPASSPORT_PURPOSE`

Vercel-safe flow used by this repo:
1. Browser starts request with `@zkpassport/sdk`.
2. Browser receives proof callback.
3. Browser sends proof to backend `POST /kyc/zkpassport/verify`.
4. Backend verifies proof and the frontend reuses that proof in CRE action payloads (`mint/create_listing/sell/rent`).

This keeps CRE + smart-contract actions unchanged while making ZKPassport compatible with hosted frontends.

### Optional anonymous mode (web UI + CRE)

The web dashboard now exposes:
- **Choose to KYC**
- **Choose to be anon**

When users select **Choose to be anon**:
- UI stores `RWA_KYC_PROVIDER=none`
- CRE action payloads include `kycProvider=none`
- Backend workflow paths skip KYC verification/writes
- Authenticated, role-gated privacy controls still apply

When users select **Choose to KYC**:
- UI re-enables `mock` / `zkpassport` routes and optional proof JSON flow
- Existing KYC verification + write behavior is preserved

Reference docs:
- https://docs.zkpassport.id/api
- https://docs.zkpassport.id/getting-started

---

## Sepolia Deploy (Production Path)

From `contracts/evm/`:

```bash
export PRIVATE_KEY=0x...
export SEPOLIA_RPC=https://<your-sepolia-rpc>
export CRE_WORKFLOW_ADDRESS=0x...   # required
export CRE_FORWARDER_ADDRESS=0x...  # required
export FEE_RECIPIENT=0x...
export ETHERSCAN_API_KEY=...   # optional, enables verification

./deploy-sepolia.sh
```

What `deploy-sepolia.sh` does:
- validates private key format
- validates Sepolia chain ID (`11155111`)
- enforces EIP-170 runtime size guard
- runs `forge test -q`
- broadcasts proxy + implementation + receiver deployment
- writes deployment artifact to `contracts/evm/deployments/houserwa_11155111_*.json`

### Temporary wallet flow (Foundry)

```bash
cast wallet new
```

Use the generated private key as `PRIVATE_KEY`, fund that address with Sepolia ETH from a faucet, then run `./deploy-sepolia.sh`.

---

## Repository Structure

- `contracts/evm/` — Solidity contracts + deployment scripts
- `RWA-House-CRE/` — primary CRE workflow package
- `backend/cre/` — backend CRE handlers/workflows
- `RWA-House-UI/web/` — web app
- `testing/` — integration/security tests and cutover scripts
