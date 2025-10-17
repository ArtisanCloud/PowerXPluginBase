#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"
CONFIG_FILE="${REPO_ROOT}/backend/etc/config.yaml"

tenant="demo"
plugin="com.powerx.plugins.base"
capability="bootstrap"
qps=10
duration=10

usage() {
  cat <<'EOF'
Usage: quota_burst.sh [--tenant <id>] [--plugin <id>] [--capability <name>] [--qps <int>] [--duration <seconds>]
Simulates sustained runtime traffic to validate quota enforcement and metrics.
EOF
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --tenant) tenant="$2"; shift 2 ;;
    --plugin) plugin="$2"; shift 2 ;;
    --capability) capability="$2"; shift 2 ;;
    --qps) qps="$2"; shift 2 ;;
    --duration) duration="$2"; shift 2 ;;
    --help|-h) usage; exit 0 ;;
    *)
      echo "Unknown argument: $1" >&2
      usage >&2
      exit 1
      ;;
  esac
done

if [[ ! -f "${CONFIG_FILE}" ]]; then
  echo "[quota-burst] config.yaml not found, please copy backend/etc/config.example.yaml first" >&2
  exit 1
fi

bind_addr="$(awk '/^[[:space:]]*bind_addr:/ {gsub(/"/, "", $2); print $2; exit}' "${CONFIG_FILE}")"
if [[ -z "${bind_addr}" ]]; then
  bind_addr=":8086"
fi
if [[ "${bind_addr}" == :* ]]; then
  base_url="http://127.0.0.1${bind_addr}/api/v1/admin/runtime"
else
  base_url="http://${bind_addr}/api/v1/admin/runtime"
fi

total_calls=$(( qps * duration ))
interval="$(awk -v q="${qps}" 'BEGIN { if (q <= 0) { print 0; exit } printf("%.4f", 1.0 / q) }')"

echo "[quota-burst] targeting ${base_url}/quota/overrides"
echo "[quota-burst] tenant=${tenant} plugin=${plugin} capability=${capability} qps=${qps} duration=${duration}s total_calls=${total_calls}"

success=0
for ((i=1; i<=total_calls; i++)); do
  payload=$(cat <<JSON
{
  "plugin_id": "${plugin}",
  "tenant_id": "${tenant}",
  "capability": "${capability}",
  "action": "throttle"
}
JSON
)
  if curl -fsS \
      -X POST \
      -H "Content-Type: application/json" \
      -d "${payload}" \
      "${base_url}/quota/overrides" >/dev/null; then
    success=$((success + 1))
  fi
  if (( i < total_calls )) && [[ "${interval}" != "0" ]]; then
    sleep "${interval}"
  fi
done

echo "[quota-burst] Completed ${success}/${total_calls} override requests"
echo "[quota-burst] Fetching latest metrics snapshot:"
"${SCRIPT_DIR}/runtime_ops_debug.sh" metrics || true
