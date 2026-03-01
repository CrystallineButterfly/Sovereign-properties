#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
CRE_DIR="$ROOT_DIR/RWA-House-CRE"
SERVICE_DIR="$ROOT_DIR/backend/zkpassport-session-service"
SIM_DIR="$ROOT_DIR/backend/cre/simulations"
DEPLOYMENT_DIR="$ROOT_DIR/testing/deployment"
TIMESTAMP="$(date -u +"%Y%m%dT%H%M%SZ")"
LOG_FILE="$DEPLOYMENT_DIR/cre_auth_smoke_${TIMESTAMP}.log"
SUMMARY_FILE="$DEPLOYMENT_DIR/cre_auth_smoke_${TIMESTAMP}.md"

ANVIL_PORT="${ANVIL_PORT:-8545}"
ANVIL_RPC="${ANVIL_RPC:-http://127.0.0.1:${ANVIL_PORT}}"
SERVICE_PORT="${SERVICE_PORT:-8787}"
SERVICE_URL="${SERVICE_URL:-http://127.0.0.1:${SERVICE_PORT}}"

SIM_ACTIONS=(mint create_listing sell rent create_bill pay_bill claim_key)
OPTIONAL_SIM_ACTIONS=(claim_key)

PASS_COUNT=0
FAIL_COUNT=0
WARN_COUNT=0
declare -a FAIL_MESSAGES
declare -a WARN_MESSAGES

ANVIL_STARTED=0
ANVIL_PID=""
SERVICE_STARTED=0
SERVICE_PID=""

mkdir -p "$DEPLOYMENT_DIR"
exec > >(tee -a "$LOG_FILE") 2>&1

print_pass() {
  local message="$1"
  PASS_COUNT=$((PASS_COUNT + 1))
  echo "✅ PASS: $message"
}

print_fail() {
  local message="$1"
  FAIL_COUNT=$((FAIL_COUNT + 1))
  FAIL_MESSAGES+=("$message")
  echo "❌ FAIL: $message"
}

print_warn() {
  local message="$1"
  WARN_COUNT=$((WARN_COUNT + 1))
  WARN_MESSAGES+=("$message")
  echo "⚠️  WARN: $message"
}

require_cmd() {
  local cmd="$1"
  if ! command -v "$cmd" >/dev/null 2>&1; then
    print_fail "Missing required command: $cmd"
    return 1
  fi
  return 0
}

wait_for_http_ok() {
  local url="$1"
  local attempts="${2:-30}"
  local sleep_seconds="${3:-1}"

  for _ in $(seq 1 "$attempts"); do
    if curl -fsS "$url" >/dev/null 2>&1; then
      return 0
    fi
    sleep "$sleep_seconds"
  done
  return 1
}

is_optional_action() {
  local action="$1"
  for opt in "${OPTIONAL_SIM_ACTIONS[@]}"; do
    if [[ "$opt" == "$action" ]]; then
      return 0
    fi
  done
  return 1
}

cleanup() {
  if [[ "$SERVICE_STARTED" -eq 1 && -n "$SERVICE_PID" ]]; then
    kill "$SERVICE_PID" >/dev/null 2>&1 || true
    wait "$SERVICE_PID" 2>/dev/null || true
  fi

  if [[ "$ANVIL_STARTED" -eq 1 && -n "$ANVIL_PID" ]]; then
    kill "$ANVIL_PID" >/dev/null 2>&1 || true
    wait "$ANVIL_PID" 2>/dev/null || true
  fi
}

trap cleanup EXIT

start_anvil_if_needed() {
  if cast chain-id --rpc-url "$ANVIL_RPC" >/dev/null 2>&1; then
    print_pass "Anvil RPC already reachable at $ANVIL_RPC"
    return
  fi

  echo "Starting Anvil at $ANVIL_RPC ..."
  nohup anvil \
    --host 127.0.0.1 \
    --port "$ANVIL_PORT" \
    --chain-id 31337 \
    --silent \
    > "$DEPLOYMENT_DIR/anvil_smoke_${TIMESTAMP}.log" 2>&1 &
  ANVIL_PID=$!
  ANVIL_STARTED=1

  local started=0
  for _ in $(seq 1 20); do
    if cast chain-id --rpc-url "$ANVIL_RPC" >/dev/null 2>&1; then
      started=1
      break
    fi
    sleep 1
  done

  if [[ "$started" -eq 1 ]]; then
    print_pass "Anvil started successfully"
  else
    print_fail "Failed to start Anvil at $ANVIL_RPC"
  fi
}

start_service_if_needed() {
  if wait_for_http_ok "$SERVICE_URL/healthz" 3 1; then
    print_pass "Workflow adapter already reachable at $SERVICE_URL"
    return
  fi

  echo "Starting workflow adapter at $SERVICE_URL ..."
  (
    cd "$SERVICE_DIR"
    node server.cjs
  ) > "$DEPLOYMENT_DIR/workflow_adapter_smoke_${TIMESTAMP}.log" 2>&1 &
  SERVICE_PID=$!
  SERVICE_STARTED=1

  if wait_for_http_ok "$SERVICE_URL/healthz" 30 1; then
    print_pass "Workflow adapter started successfully"
  else
    print_fail "Workflow adapter failed to become ready at $SERVICE_URL"
  fi
}

load_env() {
  if [[ ! -f "$ROOT_DIR/.env" ]]; then
    print_fail "Missing required env file: $ROOT_DIR/.env"
    return 1
  fi

  set -a
  # shellcheck disable=SC1090
  source "$ROOT_DIR/.env"
  set +a

  export SECRET_VALUE="${SECRET_VALUE:-local-dev-secret}"
  print_pass "Loaded .env and exported SECRET_VALUE for CRE simulation"
}

resolve_private_key() {
  local candidate=""
  for key_name in WORKFLOW_PRIVATE_KEY PRIVATE_KEY CRE_ETH_PRIVATE_KEY; do
    candidate="${!key_name:-}"
    if [[ "$candidate" =~ ^0x[0-9a-fA-F]{64}$ ]]; then
      printf "%s" "$candidate"
      return 0
    fi
  done
  return 1
}

run_auth_smoke() {
  local private_key wallet_address msg signature token
  local unauth_code mismatch_code valid_code valid_success

  if ! private_key="$(resolve_private_key)"; then
    print_fail "No valid private key found in WORKFLOW_PRIVATE_KEY/PRIVATE_KEY/CRE_ETH_PRIVATE_KEY"
    return
  fi
  print_pass "Resolved signing private key for auth smoke"

  wallet_address="$(cast wallet address --private-key "$private_key")"
  msg="workflow-auth-smoke-${TIMESTAMP}"
  signature="$(cast wallet sign --private-key "$private_key" "$msg")"

  curl -sS \
    -X POST "$SERVICE_URL/auth/verify-wallet" \
    -H "Content-Type: application/json" \
    -d "{\"address\":\"$wallet_address\",\"signature\":\"$signature\",\"message\":\"$msg\"}" \
    > "$DEPLOYMENT_DIR/auth_verify_${TIMESTAMP}.json"

  token="$(
    python3 - "$DEPLOYMENT_DIR/auth_verify_${TIMESTAMP}.json" <<'PY'
import json
import sys
payload = json.load(open(sys.argv[1], encoding="utf-8"))
print(((payload.get("data") or {}).get("token")) or "")
PY
  )"

  if [[ -z "$token" ]]; then
    print_fail "Auth verify did not return a bearer token"
    return
  fi
  print_pass "Auth verify returned bearer token"

  unauth_code="$(
    curl -sS -o "$DEPLOYMENT_DIR/unauth_trigger_${TIMESTAMP}.json" -w "%{http_code}" \
      -X POST "$SERVICE_URL/workflow/trigger" \
      -H "Content-Type: application/json" \
      -d '{"action":"set_kyc","walletAddress":"0x1111111111111111111111111111111111111111","kycProvider":"none"}'
  )"
  if [[ "$unauth_code" == "401" ]]; then
    print_pass "Unauthenticated workflow action is blocked (401)"
  else
    print_fail "Expected 401 for unauthenticated workflow action, got $unauth_code"
  fi

  mismatch_code="$(
    curl -sS -o "$DEPLOYMENT_DIR/mismatch_trigger_${TIMESTAMP}.json" -w "%{http_code}" \
      -X POST "$SERVICE_URL/workflow/trigger" \
      -H "Authorization: Bearer $token" \
      -H "Content-Type: application/json" \
      -d '{"action":"set_kyc","actorAddress":"0x1111111111111111111111111111111111111111","walletAddress":"0x1111111111111111111111111111111111111111","kycProvider":"none"}'
  )"
  if [[ "$mismatch_code" == "403" ]]; then
    print_pass "Actor mismatch is blocked (403)"
  else
    print_fail "Expected 403 for actor mismatch, got $mismatch_code"
  fi

  valid_code="$(
    curl -sS -o "$DEPLOYMENT_DIR/valid_trigger_${TIMESTAMP}.json" -w "%{http_code}" \
      -X POST "$SERVICE_URL/workflow/trigger" \
      -H "Authorization: Bearer $token" \
      -H "Content-Type: application/json" \
      -d "{\"action\":\"set_kyc\",\"actorAddress\":\"$wallet_address\",\"walletAddress\":\"$wallet_address\",\"kycProvider\":\"none\"}"
  )"
  valid_success="$(
    python3 - "$DEPLOYMENT_DIR/valid_trigger_${TIMESTAMP}.json" <<'PY'
import json
import sys
payload = json.load(open(sys.argv[1], encoding="utf-8"))
print(str(bool(payload.get("success"))).lower())
PY
  )"

  if [[ "$valid_code" == "200" && "$valid_success" == "true" ]]; then
    print_pass "Authenticated actor can execute allowed action (200 + success=true)"
  else
    print_fail "Expected 200+success=true for authenticated action, got code=$valid_code success=$valid_success"
  fi
}

run_cre_simulations() {
  local action
  for action in "${SIM_ACTIONS[@]}"; do
    echo "Running CRE simulation: $action"
    if (
      cd "$CRE_DIR"
      cre -R .. -e ../.env -T anvil-settings workflow simulate RWA-House-CRE \
        --non-interactive --trigger-index 0 \
        --http-payload @"$SIM_DIR/${action}.json"
    ) > "$DEPLOYMENT_DIR/sim_${action}_${TIMESTAMP}.log" 2>&1; then
      print_pass "CRE simulation passed: $action"
    else
      if is_optional_action "$action"; then
        print_warn "CRE simulation non-zero for optional action '$action' (see sim_${action}_${TIMESTAMP}.log)"
      else
        print_fail "CRE simulation failed: $action (see sim_${action}_${TIMESTAMP}.log)"
      fi
    fi
  done
}

write_summary() {
  {
    echo "# CRE + Auth Smoke Summary"
    echo
    echo "- Timestamp: $TIMESTAMP"
    echo "- Service URL: $SERVICE_URL"
    echo "- Anvil RPC: $ANVIL_RPC"
    echo "- Pass: $PASS_COUNT"
    echo "- Warn: $WARN_COUNT"
    echo "- Fail: $FAIL_COUNT"
    echo "- Log: $LOG_FILE"
    echo
    if (( WARN_COUNT > 0 )); then
      echo "## Warnings"
      for msg in "${WARN_MESSAGES[@]}"; do
        echo "- $msg"
      done
      echo
    fi
    if (( FAIL_COUNT > 0 )); then
      echo "## Failures"
      for msg in "${FAIL_MESSAGES[@]}"; do
        echo "- $msg"
      done
      echo
    fi
  } > "$SUMMARY_FILE"
}

main() {
  echo "=== CRE + Auth smoke start ($TIMESTAMP) ==="

  require_cmd cre || true
  require_cmd cast || true
  require_cmd anvil || true
  require_cmd curl || true
  require_cmd python3 || true
  require_cmd node || true

  if (( FAIL_COUNT > 0 )); then
    write_summary
    echo "Aborting due to missing prerequisites."
    echo "Summary: $SUMMARY_FILE"
    exit 1
  fi

  load_env || true
  if (( FAIL_COUNT > 0 )); then
    write_summary
    echo "Aborting due to env load failure."
    echo "Summary: $SUMMARY_FILE"
    exit 1
  fi

  if [[ -z "${WORKFLOW_AUTH_SECRET:-}" ]]; then
    print_fail "WORKFLOW_AUTH_SECRET is required for hardened auth path."
    write_summary
    echo "Summary: $SUMMARY_FILE"
    exit 1
  fi

  start_anvil_if_needed
  start_service_if_needed

  if (( FAIL_COUNT > 0 )); then
    write_summary
    echo "Aborting due to startup failures."
    echo "Summary: $SUMMARY_FILE"
    exit 1
  fi

  run_auth_smoke
  run_cre_simulations
  write_summary

  echo
  echo "=== Smoke complete ==="
  echo "Pass=$PASS_COUNT Warn=$WARN_COUNT Fail=$FAIL_COUNT"
  echo "Summary: $SUMMARY_FILE"
  echo "Log: $LOG_FILE"

  if (( FAIL_COUNT > 0 )); then
    exit 1
  fi
}

main "$@"
