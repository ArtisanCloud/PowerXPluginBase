#!/usr/bin/env bash
set -euo pipefail

# Export audit logs into a timestamped bundle for compliance review.
# Usage:
#   scripts/security/audit_export.sh [OUTPUT_DIR] [AUDIT_LOG_PATH]
#
# OUTPUT_DIR defaults to dist/security. AUDIT_LOG_PATH defaults to logs/audit.log
# but can be overridden via arguments or the AUDIT_LOG_PATH environment variable.

OUTPUT_DIR=${1:-dist/security}
AUDIT_SOURCE=${2:-${AUDIT_LOG_PATH:-logs/audit.log}}
TIMESTAMP=$(date -u +"%Y%m%dT%H%M%SZ")

if [[ ! -f "${AUDIT_SOURCE}" ]]; then
  echo "audit_export: audit log not found at ${AUDIT_SOURCE}" >&2
  exit 1
fi

mkdir -p "${OUTPUT_DIR}"

ARCHIVE_NAME="audit-${TIMESTAMP}.tar.gz"
ARCHIVE_PATH="${OUTPUT_DIR}/${ARCHIVE_NAME}"

tar -czf "${ARCHIVE_PATH}" -C "$(dirname "${AUDIT_SOURCE}")" "$(basename "${AUDIT_SOURCE}")"

cat <<EOF
Audit log exported:
  Source : ${AUDIT_SOURCE}
  Output : ${ARCHIVE_PATH}
EOF
