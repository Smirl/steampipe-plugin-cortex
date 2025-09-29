# Repository Guidelines

## Project Structure & Module Organization
- `main.go`: Steampipe plugin entrypoint.
- `cortex/`: Plugin implementation (tables, models, utils, tests). New tables follow `table_cortex_<name>.go` and `_test.go`.
- `docs/`: User docs. Update `docs/tables/*.md` when adding/changing tables. Keep `README.md` and `docs/index.md` in sync.
- `config/`: Steampipe plugin spec and example config (`cortex.spc`).

## Build, Test, and Development Commands
- `make install`: Build and install the plugin into `~/.steampipe/plugins/hub.steampipe.io/plugins/smirl/cortex@latest/`.
- `make install-local`: Build to local plugin paths for development.
- `make test`: Installs locally then runs `steampipe query ./test.sql` for an end‑to‑end smoke check.
- `go test ./...`: Run Go unit tests (e.g., `go test -v ./cortex`).
- `go fmt ./...`: Format Go code before committing.

## Coding Style & Naming Conventions
- Language: Go modules. Use standard `gofmt` (tabs, gofmt import ordering).
- Files: Tables as `cortex/table_cortex_<resource>.go`; tests as `<file>_test.go`.
- Functions/vars: Go idioms (`CamelCase` exported, `camelCase` internal). Keep table/column names consistent with docs.
- Linting: Prefer idiomatic Go; avoid global state; keep functions small.

## Testing Guidelines
- Framework: Go `testing` in `cortex/*_test.go` plus a SQL smoke test (`test.sql`).
- Run: `go test -cover ./cortex` and `make test` before PRs.
- Add tests when changing table logic, hydration, transforms, and error paths.
- Use realistic fixtures; avoid network calls in unit tests.

## Commit & Pull Request Guidelines
- Commits: Concise, imperative. Conventional style encouraged (e.g., `feat(entity): add group filter`).
- PRs: Describe motivation, approach, and impact. Link related issues.
- Include: updated docs in `docs/tables/`, examples if SQL shape changes, and any config notes.
- CI: Ensure build, tests, and formatting pass locally before requesting review.

## Security & Configuration Tips
- Never commit secrets. Use env vars: `CORTEX_API_KEY`, `CORTEX_BASE_URL`.
- Test against non‑prod Cortex where possible. Document config changes in README/docs.
