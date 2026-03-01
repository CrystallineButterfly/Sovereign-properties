#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd -- "$SCRIPT_DIR/../.." && pwd)"
WEB_DIR="$REPO_ROOT/RWA-House-UI/web"
LOG_DIR="$REPO_ROOT/testing/deployment"
TIMESTAMP="$(date -u +"%Y%m%dT%H%M%SZ")"
LOG_FILE="$LOG_DIR/vercel_deploy_${TIMESTAMP}.log"
XDG_CACHE_HOME="${XDG_CACHE_HOME:-/tmp/vercel-cache}"

HOUSE_RWA_DEFAULT="0x990e1EB2Dd8fA8007533Ab50bE262A44EEF172ee"
API_DEFAULT="https://zkpassport-api-production.up.railway.app"

log() {
  printf '[deploy-vercel-demo] %s\n' "$*"
}

die() {
  printf '[deploy-vercel-demo] ERROR: %s\n' "$*" >&2
  exit 1
}

load_env_file() {
  local env_file="$1"
  if [[ -f "$env_file" ]]; then
    log "Loading env from ${env_file#$REPO_ROOT/}"
    set -a
    # shellcheck disable=SC1090
    source "$env_file"
    set +a
  fi
}

require_command() {
  command -v "$1" >/dev/null 2>&1 || die "Required command not found: $1"
}

require_value() {
  local name="$1"
  local value="$2"
  [[ -n "$value" ]] || die "$name must be set before deployment."
}

upsert_vercel_env() {
  local name="$1"
  local value="$2"
  local target="$3"
  printf '%s' "$value" | vercel env add "$name" "$target" --force "${VERCEL_ARGS[@]}" \
    >/dev/null
}

main() {
  mkdir -p "$LOG_DIR"
  mkdir -p "$XDG_CACHE_HOME"
  export XDG_CACHE_HOME

  require_command vercel
  require_command npm

  load_env_file "$REPO_ROOT/.env"
  load_env_file "$WEB_DIR/.env"
  load_env_file "$WEB_DIR/.env.local"
  if [[ -n "${DEPLOY_ENV_FILE:-}" ]]; then
    load_env_file "$DEPLOY_ENV_FILE"
  fi

  export VITE_API_URL="${VITE_API_URL:-$API_DEFAULT}"
  export VITE_ZKPASSPORT_API_URL="${VITE_ZKPASSPORT_API_URL:-$VITE_API_URL}"
  export VITE_RPC_URL="${VITE_RPC_URL:-${VITE_API_URL%/}/rpc}"
  export VITE_HOUSE_RWA_ADDRESS="${VITE_HOUSE_RWA_ADDRESS:-${HOUSE_RWA_CONTRACT_ADDRESS:-$HOUSE_RWA_DEFAULT}}"
  export VITE_EXPECTED_CHAIN_ID="${VITE_EXPECTED_CHAIN_ID:-11155111}"
  export VITE_SUPPORTED_CHAIN_IDS="${VITE_SUPPORTED_CHAIN_IDS:-11155111}"
  export VITE_ENABLE_CHAIN_RPC_FALLBACK="${VITE_ENABLE_CHAIN_RPC_FALLBACK:-false}"
  export VITE_ENABLE_PUBLIC_RPC_CANDIDATES="${VITE_ENABLE_PUBLIC_RPC_CANDIDATES:-false}"
  export VITE_XMTP_ENV="${VITE_XMTP_ENV:-production}"
  export VITE_MAX_HOUSE_SCAN="${VITE_MAX_HOUSE_SCAN:-500}"
  export VERCEL_PROJECT_NAME="${VERCEL_PROJECT_NAME:-sovereign-properties}"
  export VERCEL_TARGETS="${VERCEL_TARGETS:-production}"

  require_value "VITE_PRIVY_APP_ID" "${VITE_PRIVY_APP_ID:-}"

  VERCEL_ARGS=()
  if [[ -n "${VERCEL_TOKEN:-}" ]]; then
    VERCEL_ARGS+=(--token "$VERCEL_TOKEN")
  fi
  if [[ -n "${VERCEL_SCOPE:-}" ]]; then
    VERCEL_ARGS+=(--scope "$VERCEL_SCOPE")
  fi

  if ! vercel whoami "${VERCEL_ARGS[@]}" >/dev/null 2>&1; then
    die "No Vercel credentials. Run 'vercel login' once or export VERCEL_TOKEN and rerun."
  fi

  cd "$WEB_DIR"

  if [[ ! -f "$WEB_DIR/.vercel/project.json" ]]; then
    log "Linking Vercel project: $VERCEL_PROJECT_NAME"
    vercel link --yes --project "$VERCEL_PROJECT_NAME" "${VERCEL_ARGS[@]}" \
      >/dev/null
  fi

  local targets_csv="$VERCEL_TARGETS"
  IFS=',' read -r -a targets <<< "$targets_csv"
  for raw_target in "${targets[@]}"; do
    local target
    target="$(printf '%s' "$raw_target" | xargs)"
    [[ -n "$target" ]] || continue
    log "Syncing Vercel env vars for target=$target"
    upsert_vercel_env "VITE_PRIVY_APP_ID" "$VITE_PRIVY_APP_ID" "$target"
    upsert_vercel_env "VITE_API_URL" "$VITE_API_URL" "$target"
    upsert_vercel_env "VITE_ZKPASSPORT_API_URL" "$VITE_ZKPASSPORT_API_URL" "$target"
    upsert_vercel_env "VITE_RPC_URL" "$VITE_RPC_URL" "$target"
    upsert_vercel_env "VITE_HOUSE_RWA_ADDRESS" "$VITE_HOUSE_RWA_ADDRESS" "$target"
    upsert_vercel_env "VITE_EXPECTED_CHAIN_ID" "$VITE_EXPECTED_CHAIN_ID" "$target"
    upsert_vercel_env "VITE_SUPPORTED_CHAIN_IDS" "$VITE_SUPPORTED_CHAIN_IDS" "$target"
    upsert_vercel_env \
      "VITE_ENABLE_CHAIN_RPC_FALLBACK" \
      "$VITE_ENABLE_CHAIN_RPC_FALLBACK" \
      "$target"
    upsert_vercel_env \
      "VITE_ENABLE_PUBLIC_RPC_CANDIDATES" \
      "$VITE_ENABLE_PUBLIC_RPC_CANDIDATES" \
      "$target"
    upsert_vercel_env "VITE_XMTP_ENV" "$VITE_XMTP_ENV" "$target"
    upsert_vercel_env "VITE_MAX_HOUSE_SCAN" "$VITE_MAX_HOUSE_SCAN" "$target"
  done

  log "Installing web dependencies"
  npm install >/dev/null

  log "Building web app"
  npm run build >/dev/null

  log "Deploying to Vercel"
  local deploy_output deploy_url
  deploy_output="$(
    vercel --prod --yes "${VERCEL_ARGS[@]}" 2>&1 | tee "$LOG_FILE"
  )"
  deploy_url="$(printf '%s\n' "$deploy_output" | grep -Eo 'https://[^[:space:]]+' | tail -n1)"
  require_value "deployment URL" "$deploy_url"

  printf '\nDeployed URL: %s\n' "$deploy_url"
  printf 'Deploy log: %s\n' "$LOG_FILE"
}

main "$@"
