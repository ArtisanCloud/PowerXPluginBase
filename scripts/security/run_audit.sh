#!/usr/bin/env bash
set -euo pipefail

TIMESTAMP=$(date -u +"%Y%m%dT%H%M%SZ")
ROOT_DIR=$(cd "$(dirname "$0")/../.." && pwd)
OUTPUT_DIR="$ROOT_DIR/build/security/$TIMESTAMP"
mkdir -p "$OUTPUT_DIR"

run_step() {
  local name="$1"; shift
  local cmd=($@)
  echo " - ${name}"
  if command -v "${cmd[0]}" >/dev/null 2>&1; then
    if ! "${cmd[@]}" &>"$OUTPUT_DIR/${name}.log"; then
      echo "   ! ${name} failed (see ${OUTPUT_DIR}/${name}.log)"
      return 1
    fi
  else
    echo "   ! ${cmd[0]} not available, skipping" | tee -a "$OUTPUT_DIR/${name}.log"
  fi
  return 0
}

STATUS=PASSED

if ! run_step golangci-lint golangci-lint run ./...; then
  STATUS=FAILED
fi

if ! run_step govulncheck govulncheck ./...; then
  STATUS=FAILED
fi

if ! run_step gosec gosec ./...; then
  STATUS=FAILED
fi

if command -v npm >/dev/null 2>&1; then
  echo " - npm audit"
  if ! (cd "$ROOT_DIR/web-admin" && npm audit --production &>"$OUTPUT_DIR/npm-audit.log"); then
    STATUS=FAILED
  fi
else
  echo " - npm audit (skipped)" | tee -a "$OUTPUT_DIR/npm-audit.log"
fi

if ! run_step trivy trivy fs --exit-code 0 --no-progress "$ROOT_DIR"; then
  STATUS=FAILED
fi

if command -v cosign >/dev/null 2>&1; then
  if ! cosign version &>"$OUTPUT_DIR/cosign.log"; then
    STATUS=FAILED
  fi
fi

cat >"$OUTPUT_DIR/report.json" <<JSON
{
  "timestamp": "$TIMESTAMP",
  "status": "$STATUS",
  "artifacts": {
    "directory": "build/security/$TIMESTAMP"
  }
}
JSON

echo "security audit status: $STATUS"
[[ "$STATUS" == "PASSED" ]]
