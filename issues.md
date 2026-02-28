# Issues Log (Sovereign-Properties)

Updated: 2026-02-27

## Patch status (latest)

- ✅ Fixed: Auth now required for all privileged `/workflow/trigger` actions.
- ✅ Fixed: Request `actorAddress` must match authenticated bearer wallet.
- ✅ Fixed: Owner/renter/buyer/creator role checks added for critical workflow actions.
- ✅ Fixed: `WORKFLOW_AUTH_SECRET` is now mandatory (insecure static fallback removed).
- ✅ Fixed: Frontend SIWE no longer falls back to local demo auth on backend auth failure.
- ✅ Fixed: API 401 now clears auth token and emits a session-expired event.
- ✅ Fixed: AuthProvider now listens for auth-expired and resets session/UI state.
- ✅ Fixed: Added Vercel CSP header config (`RWA-House-UI/web/vercel.json`) to align deployed CSP behavior.
- ✅ Improved: Removed nonce-time dependency from CRE bill payment reference/hash fallback paths in Go handlers.
- ⚠️ Still requires follow-up: full CRE node-mode aggregation refactor for external KYC verifier path.

## Critical

### 1) Unauthenticated access to privileged CRE workflow actions
- **Severity:** Critical
- **What:** `POST /workflow/trigger` is publicly reachable and only enforces auth for `claim_key`. Other privileged actions (`mint`, `set_kyc`, `create_listing`, `sell`, `rent`, `create_bill`, `pay_bill`) can be called without a bearer token.
- **Evidence:**
  - `backend/zkpassport-session-service/server.cjs:1875-1877` (route exposed)
  - `backend/zkpassport-session-service/workflow-trigger.cjs:2832-2850` (all actions)
  - `backend/zkpassport-session-service/workflow-trigger.cjs:2881-2915` (auth guard only for `claim_key`)
- **Why this is dangerous:** these actions submit onchain writes through workflow-authorized execution (contracts enforce `onlyCRE`), so backend auth is the security boundary.
- **Related contract context:** `contracts/evm/src/HouseRWA.sol:268-275,316-325,337-341,401` (`onlyCRE`-gated operations).
- **Fix:** require authenticated wallet for every state-changing action, and bind request actor to token owner/role server-side.

### 2) Auth token secret has insecure fallback
- **Severity:** Critical
- **What:** token-signing secret falls back to a static development string if env vars are missing.
- **Evidence:** `backend/zkpassport-session-service/workflow-trigger.cjs:209-213`
- **Why this is dangerous:** predictable fallback enables token forgery risk if deployment misses `WORKFLOW_AUTH_SECRET`.
- **Fix:** hard-fail startup unless `WORKFLOW_AUTH_SECRET` is present and sufficiently strong.

## High

### 3) Caller-controlled identity fields on privileged actions
- **Severity:** High
- **What:** action payloads trust caller-supplied addresses (`ownerAddress`, `buyerAddress`, `provider`, etc.) without mandatory auth binding (except `claim_key`).
- **Evidence:**
  - `handleMint`: `workflow-trigger.cjs:2220-2223`
  - `handleSell`: `workflow-trigger.cjs:2463-2471`
  - `handleCreateBill`: `workflow-trigger.cjs:2628-2631`
- **Impact:** unauthorized callers can initiate workflows on behalf of others when auth is absent.
- **Fix:** enforce bearer auth and derive actor server-side; reject mismatches.

### 4) Anonymous KYC bypass is enabled by payload
- **Severity:** High (production), Low/Expected (demo)
- **What:** `kycProvider=none` short-circuits KYC writes/verification.
- **Evidence:** `backend/zkpassport-session-service/workflow-trigger.cjs:2178-2182`
- **Impact:** if left enabled in production, can bypass intended KYC gating flows.
- **Fix:** environment-gate anonymous mode (`ALLOW_ANON_KYC=false` in prod), reject `none` when disabled.

### 5) CRE determinism risks in workflow handler code
- **Severity:** High
- **What:** CRE handlers use non-deterministic local time/random-like values in workflow logic.
- **Evidence:**
  - `backend/cre/handlers/http.go:718` (`time.Now().UnixNano()` in payment reference)
  - `backend/cre/handlers/http.go:949-955` (`time.Now()` in KYC fallback hash/expiry)
  - `backend/cre/handlers/kyc.go:68,77-81,171` (`time.Now()` for expiry/hash material)
- **Impact:** can break multi-node consensus behavior or create inconsistent outputs across nodes.
- **Fix:** use CRE deterministic primitives/runtime time where required; avoid entropy/time in consensus-critical paths.

### 6) External KYC HTTP call path is not CRE node-mode aggregated
- **Severity:** High
- **What:** KYC verifier uses plain `net/http` call inside CRE handler.
- **Evidence:** `backend/cre/handlers/kyc.go:102-137`
- **Impact:** no explicit node-level aggregation strategy; centralized dependency can undermine Byzantine assumptions and determinism.
- **Fix:** move external data fetch to CRE capability flow with explicit consensus/aggregation policy.

## Medium

### 7) EVM reads do not specify explicit finality mode
- **Severity:** Medium
- **What:** `CallContract` requests do not set finalized/safe read constraints.
- **Evidence:** `backend/cre/handlers/http.go:1098-1104`
- **Impact:** potential reorg sensitivity for read-before-write logic.
- **Fix:** set explicit finality/read confidence where available; add post-write confirmation policy for critical flows.

### 8) Frontend auth can enter local demo mode and drift from backend auth
- **Severity:** Medium
- **What:** on SIWE failure, UI creates demo user state even when backend token is missing/invalid.
- **Evidence:** `RWA-House-UI/web/src/components/AuthProvider.tsx:621-636,640`
- **Impact:** repeated 401s (`/auth/me`, `/notifications`, `/houses`) and broken KYC/session UX.
- **Fix:** if SIWE fails, block protected routes and prompt re-auth; do not proceed as authenticated demo user on production endpoints.

### 9) Stale/invalid auth tokens are not auto-cleared on 401
- **Severity:** Medium
- **What:** API client throws on 4xx but does not clear auth token or trigger re-login.
- **Evidence:** `RWA-House-UI/shared/src/utils/api.ts:555-577`
- **Impact:** polling endpoints continue failing until manual logout/refresh.
- **Fix:** intercept 401 specifically, clear token/session, and route user to re-auth.

### 10) CSP mismatch is still a blocker for ZKPassport verification flows
- **Severity:** Medium
- **What:** runtime reports show CSP blocking required connections during proof verification (bridge/RPC paths), even when scan/QR starts.
- **Evidence:** observed browser errors during testing (`connect-src`/`Fetch API cannot load ... violates CSP`).
- **Code reference:** CSP currently defined in meta tag at `RWA-House-UI/web/index.html:15`.
- **Impact:** QR scan may work, but app-side proof verification/handshake fails or loops.
- **Fix:** align deployed CSP headers (not only meta tag) with exact zkPassport + RPC hosts used in runtime.

## Operational blockers observed in testing

1. **KYC completion loop:** QR scan succeeds, but app-side confirm/verify gets stuck or returns failed verification.
2. **401 storm on protected APIs:** `/notifications`, `/auth/me`, and owner-scoped `/houses` repeatedly fail when token/session is stale.
3. **Workflow auth errors:** `claim_key` correctly requires auth, but other state-changing workflow actions remain too open.

## Recommended remediation order

1. Lock down `/workflow/trigger` auth for all write actions.
2. Enforce required `WORKFLOW_AUTH_SECRET` at startup (remove static fallback).
3. Gate `kycProvider=none` behind explicit non-production flag.
4. Fix CSP at deployment-header level for zkPassport and required RPC hosts.
5. Fix frontend auth-state drift (no demo-auth fallback on protected production paths).
6. Refactor CRE handler non-deterministic/external-call paths to CRE-safe deterministic patterns.
