#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"
CONFIG_FILE="${REPO_ROOT}/backend/etc/config.yaml"
HOST_VALUES="${REPO_ROOT}/backend/etc/host-values.yaml"

log() {
  echo "[runtime-ops-debug] $*"
}

ensure_config() {
  if [[ ! -f "${CONFIG_FILE}" ]]; then
    log "未找到 backend/etc/config.yaml，请先复制 config.example.yaml"
    exit 1
  fi
}

print_runtime_config() {
  ensure_config
  log "Runtime Ops 配置摘要（来自 ${CONFIG_FILE}）:"
  awk '
    /^[[:space:]]*runtime_ops:/ {print; in=1; next}
    in && /^[[:alnum:]_+-]+:/ {in=0}
    in {print}
  ' "${CONFIG_FILE}"
  echo
  log "监控配置:"
  awk '
    /^[[:space:]]*monitoring:/ {print; in=1; next}
    in && /^[[:alnum:]_+-]+:/ {in=0}
    in {print}
  ' "${CONFIG_FILE}"

  if [[ -f "${HOST_VALUES}" ]]; then
    log "检测到 host-values.yaml，将覆盖宿主注入参数: ${HOST_VALUES}"
  fi

  if grep -q 'path:[[:space:]]*"/api/v1/admin/runtime/metrics"' "${CONFIG_FILE}"; then
    log "Prometheus 抓取端点已配置为 /api/v1/admin/runtime/metrics"
  else
    log "警告: 未在 config.yaml 中发现 runtime metrics 路径，默认 /api/v1/admin/runtime/metrics"
  fi
}

server_bind_addr() {
  ensure_config
  awk '
    /^[[:space:]]*bind_addr:/ {
      gsub(/"/, "", $2);
      print $2;
      exit
    }
  ' "${CONFIG_FILE}"
}

metrics_path() {
  ensure_config
  local path
  path=$(awk '
    /^[[:space:]]*monitoring:/ {inMon=1; next}
    inMon && /^[[:space:]]*metrics:/ {inMetrics=1; next}
    inMetrics && /^[[:space:]]*path:/ {
      gsub(/"/, "", $2);
      print $2;
      exit
    }
    inMetrics && /^[[:space:]]*[[:alnum:]_+-]+:/ {exit}
  ' "${CONFIG_FILE}")
  if [[ -z "${path}" ]]; then
    path="/api/v1/admin/runtime/metrics"
  fi
  echo "${path}"
}

invoke_metrics() {
  ensure_config
  local bind_addr path url
  bind_addr="$(server_bind_addr)"
  path="$(metrics_path)"
  bind_addr="${bind_addr:-:8086}"
  path="${path:-/api/v1/admin/runtime/metrics}"
  if [[ "${bind_addr}" == :* ]]; then
    url="http://127.0.0.1${bind_addr}${path}"
  else
    url="http://${bind_addr}${path}"
  fi
  log "尝试访问 Prometheus 指标端点: ${url}"
  curl -fsS "${url}" | head -n 20
}

bootstrap_pipeline() {
  print_runtime_config
  log "模拟执行 bootstrap pipeline: unpack -> port reservation -> env injection -> process launch -> health registration"
  if [[ -f "${HOST_VALUES}" ]]; then
    log "宿主注入配置已发现，模拟合并 host-values.yaml"
  fi
}

mcp_sequence() {
  log "模拟 MCP REGISTER -> HEARTBEAT -> STALE 流程"
  log "发送 REGISTER payload (示意)"
  cat <<'JSON'
{
  "plugin_id": "com.powerx.plugins.base",
  "tenant_id": "demo",
  "capabilities_hash": "sha256:demo",
  "heartbeat_interval": 15
}
JSON
  log "心跳超时后应进入 stale 状态，并触发 quota/metrics 更新"
}

usage() {
  cat <<'EOF'
Usage: runtime_ops_debug.sh <command>
Commands:
  bootstrap       打印配置摘要并演示 bootstrap 流程
  mcp-heartbeat   模拟 MCP 注册与心跳流程
  inspect         仅打印 runtime_ops 与监控配置
  metrics         尝试请求 Prometheus 指标端点
EOF
}

command="${1:-bootstrap}"
case "${command}" in
  bootstrap) bootstrap_pipeline ;;
  mcp-heartbeat) mcp_sequence ;;
  inspect) print_runtime_config ;;
  metrics) invoke_metrics ;;
  *)
    usage >&2
    exit 1
    ;;
esac
