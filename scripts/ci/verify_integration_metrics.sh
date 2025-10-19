#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")"/.. && pwd)"
METRICS_FILE="$ROOT_DIR/backend/internal/observability/integration/metrics.go"
DASHBOARD_FILE="$ROOT_DIR/docs/observability/integration-dashboard.json"
declare -A REQUIRED=(
  ["powerx_integration_envelopes_total"]="SC-001"
  ["powerx_integration_webhook_attempts_total"]="SC-002"
  ["powerx_integration_webhook_delivery_seconds"]="SC-002"
  ["powerx_integration_secrets_rotations_due"]="SC-003"
  ["powerx_integration_idempotency_events_total"]="SC-005"
)

for metric in "${!REQUIRED[@]}"; do
  if ! grep -q "$metric" "$METRICS_FILE"; then
    echo "[ERROR] required metric '$metric' not found in metrics implementation" >&2
    exit 1
  fi
  if ! grep -q "$metric" "$DASHBOARD_FILE"; then
    echo "[ERROR] dashboard missing reference to '$metric' (covers ${REQUIRED[$metric]})" >&2
    exit 1
  fi
  echo "[OK] metric '$metric' present (covers ${REQUIRED[$metric]})"
done

# Ensure Secrets rotation metric captures due windows
if ! grep -q 'powerx_integration_secrets_rotations_due{window="due_24h"}' "$DASHBOARD_FILE"; then
  echo "[ERROR] dashboard missing due_24h window view for secrets rotation" >&2
  exit 1
fi

# Ensure DLQ backlog panel present for SC-002/SC-005 alerting
if ! grep -q 'status="dlq"' "$DASHBOARD_FILE"; then
  echo "[ERROR] dashboard missing webhook DLQ backlog panel" >&2
  exit 1
fi

echo "[OK] integration metrics verified"
