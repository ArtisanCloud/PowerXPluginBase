# Repository Guidelines

## Project Structure & Module Organization
The backend Go service lives in `backend/`, with entrypoints at `backend/cmd/plugin` for runtime and `backend/cmd/database` for migrations. Domain logic is grouped under `backend/internal/` by concern (e.g. `services/`, `router/`). The Nuxt 4 admin UI is in `web-admin/`, and shared Make targets sit in `make-files/`. Configuration samples reside in `backend/etc/`, while release artefacts are emitted into `dist/` or `target/`.

## Build, Test, and Development Commands
Use `make dev-setup` once to download Go modules and install `golangci-lint`. `make run` starts the backend with dev-friendly defaults; add `make migrate` first if schema changes. Build Go binaries via `make build`, and bundle the admin UI for the proxied host path with `make frontend-build`. End-to-end release bundles are prepared by `make dist` (PowerX install) or `make release` (full package).

## Coding Style & Naming Conventions
Go sources must stay `gofmt`-clean; run `make fmt` before pushing. Lint with `make lint`, resolving `golangci-lint` issues promptly. Keep package names lower_snake, exported APIs in UpperCamel, and file names matching the package purpose. TypeScript/Vue code follows the Nuxt ESLint presets (two-space indent, kebab-case component files, `useFoo` composables). Prefer explicit PowerX domain names such as `powerxTenantService`.

## Testing Guidelines
Backend unit tests belong alongside code as `_test.go` files; run all suites with `make test`. For coverage snapshots use `make test-coverage`, which writes `backend/coverage.html` for review. Name tests `Test<Thing>_<Expectation>` and seed fixtures through helpers under `backend/internal/testdata`. Frontend stories or component tests should use `@nuxt/test-utils` under `web-admin/tests/` when UI logic warrants validation.

## Commit & Pull Request Guidelines
Adopt Conventional Commits: `feat(domain): add workspace sync` or `fix(router): guard nil payload`. Group related changes per commit and keep messages under 72 characters after the scope. Pull requests should reference the tracking issue, describe schema or configuration impacts, and include a concise test plan (`make test`, `npm run build`, etc.). Attach screenshots or API samples when UI or contract behaviour shifts.

## Configuration Tips
Copy `backend/etc/config.example.yaml` to `config.yaml` for local overrides and avoid committing secrets. Environment-sensitive values (DB DSNs, PowerX bindings) are read from env vars in `make run`; document any new keys in `backend/etc/README.md`. Keep generated binaries out of version control—use the `clean` targets before publishing branches.
