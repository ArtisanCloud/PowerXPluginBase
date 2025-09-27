# Repository Guidelines

## Project Structure & Modules
- `backend/`: Go plugin backend (cmd, internal, etc, bin). Example: `backend/cmd/plugin/main.go`.
- `web-admin/`: Nuxt 4 admin UI (TypeScript, Pinia, i18n). Start from `web-admin/app/`.
- `docs/`: Documentation; `plugin.yaml`: plugin manifest; `Makefile`: task runner.


## Build, Test, and Development
- Backend
  - `make dev`: Run DB migrate then start backend (`POWERX_*` envs set).
  - `make run`: Start backend locally on `:8091`.
  - `make build` / `make build-linux`: Compile backend binary to `backend/bin/plugin`.
  - `make migrate` / `make migrate-cmd`: Apply migrations.
  - `make test`: Run Go tests. `make test-coverage`: generate coverage report.
  - `make package`: Zip plugin with binary, `plugin.yaml`, and `web-admin/` (excludes `node_modules`).
- Frontend
  - From `web-admin/`: `npm i`, `npm run dev -- --port 3036` to start UI. `npm run build` to build; `npm run preview` to serve build.

## Coding Style & Naming
- Go
  - Use `make fmt` (go fmt) and `make lint` (golangci-lint). Prefer tabs, `gofmt`-style imports.
  - Packages: lower_snake; exported types/functions use PascalCase; private use camelCase.
  - Keep layers clean: `internal/domain`, `services`, `transport/http`, `router`.
- Frontend (Nuxt 4 + TS)
  - ESLint via Nuxt preset. Prefer 2-space indent, PascalCase components, camelCase variables.
  - Pages, components, and composables live under `web-admin/app/`.

## Testing Guidelines
- Go: standard `testing` with `_test.go`. Place tests alongside packages, name functions `TestXxx`.
- Coverage: aim ≥70% for changed code. Run `make test-coverage` before PRs.
- Frontend: Nuxt test utils available; add tests where changes are significant.

## Commit & Pull Requests
- Commits: small, descriptive, imperative (“Add team routes”). Group logical changes; avoid mixed backend/frontend in one commit when possible.
- PRs: include summary, screenshots for UI, steps to validate, and any schema or config changes (`POWERX_*`). Link issues. Passing CI, lint, and coverage required.

## Security & Config
- Sensitive config via env: `POWERX_BIND_ADDR`, `POWERX_DB_SCHEMA`, `POWERX_LOG_LEVEL`, `POWERX_DEV_MODE`. Do not hardcode secrets.
- Validate inputs at transport layer; use centralized error/response helpers for consistency.

## Agent-Specific Notes
- Do not alter public plugin contracts without discussion. Keep router registration patterns and middleware layering consistent. Update `Makefile` targets if build or layout changes.


# Communication
- 默认使用简体中文与我交流和解释。
- 生成代码时：注释、README、commit message 优先中文，必要时双语（中 + 英术语）。

# Style
- 终端命令给出可复制的一行版本；危险操作前先解释风险并征询确认。