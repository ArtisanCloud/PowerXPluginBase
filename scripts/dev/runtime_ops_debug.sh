#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"
CONFIG_FILE="${REPO_ROOT}/backend/etc/config.yaml"
HOST_VALUES="${REPO_ROOT}/backend/etc/host-values.yaml"

if [[ ! -f "${CONFIG_FILE}" ]]; then
  echo "[runtime-ops-debug] 未找到 backend/etc/config.yaml，请先复制 config.example.yaml" >&2
  exit 1
fi

if [[ -f "${HOST_VALUES}" ]]; then
  echo "[runtime-ops-debug] 使用宿主提供的 host-values.yaml 覆盖关键参数"
fi

case "${1:-bootstrap}" in
  bootstrap)
    echo "[runtime-ops-debug] 模拟执行 bootstrap pipeline ..."
    echo "读取配置: ${CONFIG_FILE}"
    [[ -f "${HOST_VALUES}" ]] && echo "读取宿主注入配置: ${HOST_VALUES}"
    ;;
  mcp-heartbeat)
    echo "[runtime-ops-debug] 模拟 MCP REGISTER -> HEARTBEAT -> STALE 流程"
    ;;
  *)
    echo "Usage: $0 [bootstrap|mcp-heartbeat]" >&2
    exit 1
    ;;
esac
