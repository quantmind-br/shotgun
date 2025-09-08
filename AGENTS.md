# Repository Guidelines

This document guides contributors (humans and agents) working on Shotgun, a Go TUI/CLI project.

## Project Structure & Module Organization
- `cmd/shotgun/` — CLI entrypoint (`main.go`).
- `internal/` — private packages: `app/`, `cli/` (commands), `components/`, `core/` (`scanner/`, `template/`, `builder/`, `config/`), `models/`, `screens/`, `styles/`, `utils/`, `integration/`.
- `e2e/` — end‑to‑end tests.  `templates/` — built‑in templates.  `bin/` — build artifacts.
- `Makefile`, `go.mod`, `go.sum`, `docs/`, `.github/` — automation, deps, docs, CI.

## Build, Test, and Development Commands
- `make build` — build to `bin/shotgun` for current platform.
- `make run` — run locally (equiv. `go run ./cmd/shotgun`).
- `make test` — run all tests.
- `make test-coverage` — tests + `coverage.html` report.
- `make lint` — `go fmt` + `go vet`.
- `make build-all` — cross‑platform binaries; see `bin/` outputs.
- `make deps` — tidy and download modules.

## Coding Style & Naming Conventions
- Go 1.22+ (module targets 1.23 toolchain). Always run `go fmt ./...` and keep `go vet` clean.
- Indentation: Go defaults (tabs). Line width: keep readable (<100 cols preferred).
- Packages: lowercase, no underscores; exported identifiers use CamelCase. Error vars: `ErrXxx`.
- Cross‑platform paths: use `path/filepath` and `os`, never hardcode separators.
- Keep UI logic (Bubble Tea) pure where possible; prefer small, testable functions.

## Testing Guidelines
- Framework: standard `testing`. Prefer table‑driven tests.
- Unit tests live with code as `*_test.go` and `TestXxx` funcs.
- Integration: `internal/integration/`; E2E: `e2e/`.
- Coverage: maintain ≥60% for core packages; generate via `make test-coverage` and review `coverage.html`.

## Commit & Pull Request Guidelines
- Commits: concise, present‑tense summaries. Mirror repo style (e.g., “Complete Story 2.3: …”).
  Example: `Implement template discovery and loading (Story 2.1)`
- PRs must include: clear description, rationale, testing steps, linked issues, and screenshots/asciinema for TUI changes.
- Keep changes scoped; include/update tests and docs when behavior changes.

## Security & Configuration Tips
- Do not commit secrets or local paths. Use env vars (e.g., `SHOTGUN_CONFIG_DIR`) for config.
- Ensure changes build on Linux/macOS/Windows; avoid shell‑specific assumptions.
- Binaries are built with `-ldflags "-s -w"`; keep CGO disabled unless justified.

## Agent‑Specific Instructions
- Apply minimal diffs consistent with existing patterns. Do not add new deps without justification; run `make deps` and `go mod tidy` if you do.
- If you change behavior, update tests and `README.md` accordingly.
