#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")"/.. && pwd)"

pushd "$ROOT_DIR" >/dev/null

echo "[CI] Formatting & Linting"
make fmt lint lint-admin

echo "[CI] Running full test matrix"
make test-all

echo "[CI] Running webhook replay drill"
make integration-smoke

echo "[CI] Verifying integration metrics definitions"
scripts/ci/verify_integration_metrics.sh

echo "[CI] Building release artefacts"
make build

popd >/dev/null
