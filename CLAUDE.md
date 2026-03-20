# CLAUDE.md

## Development workflow

- **Always follow TDD**: Write or update tests first, verify they fail, then implement the change, then verify tests pass.

## Build & test commands

- Build: `go build -o zo-cli .`
- Test: `go test ./...`
- Vet: `go vet ./...`

## Project structure

- `main.go` — entry point
- `cmd/` — cobra commands (ask, chat, models, personas, config)
- `internal/api/` — HTTP client, request/response types
- `internal/config/` — config file load/save
