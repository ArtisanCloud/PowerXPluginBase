#!/usr/bin/env bash
# shellcheck disable=SC2086
set -euo pipefail

if ! command -v python3 >/dev/null 2>&1; then
  echo "python3 is required to run the mock webhook target" >&2
  exit 1
fi

PORT="${1:-8089}"

echo "Starting mock webhook target on http://0.0.0.0:${PORT}"
echo "Incoming requests will be logged to stdout. Press Ctrl+C to stop."

python3 -u - "${PORT}" <<'PY'
import json
import sys
import threading
from http.server import BaseHTTPRequestHandler, HTTPServer
from time import time

PORT = int(sys.argv[1]) if len(sys.argv) > 1 else 8089


class MockWebhookHandler(BaseHTTPRequestHandler):
    def do_POST(self):
        content_length = int(self.headers.get("Content-Length", 0))
        body = self.rfile.read(content_length or 0)
        payload = body.decode("utf-8", errors="replace")

        print("=" * 80)
        print(f"[mock-webhook] {self.command} {self.path}")
        for header, value in self.headers.items():
            print(f"{header}: {value}")
        if payload:
            print("\nPayload:")
            print(payload)
        else:
            print("\n(no payload)")
        print("=" * 80)

        self.send_response(200)
        self.send_header("Content-Type", "application/json")
        self.end_headers()
        self.wfile.write(
            json.dumps(
                {
                    "status": "ok",
                    "received_at": time(),
                }
            ).encode("utf-8")
        )

    def do_GET(self):
        self.send_response(200)
        self.send_header("Content-Type", "application/json")
        self.end_headers()
        self.wfile.write(
            json.dumps(
                {
                    "status": "ready",
                    "path": self.path,
                }
            ).encode("utf-8")
        )

    def log_message(self, fmt, *args):
        # Silence default logging; we already print a structured summary.
        return


def serve():
    server = HTTPServer(("0.0.0.0", PORT), MockWebhookHandler)
    try:
        server.serve_forever()
    except KeyboardInterrupt:
        pass
    finally:
        server.server_close()


if __name__ == "__main__":
    thread = threading.Thread(target=serve, daemon=False)
    thread.start()
    try:
        thread.join()
    except KeyboardInterrupt:
        pass
PY
