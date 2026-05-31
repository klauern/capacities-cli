# CLAUDE.md

## Commands
- **Build**: `task build` (output: `bin/capacities`) or `go build -o bin/capacities ./cmd/capacities`
- **Test**: `task test` (runs with `-race`) or `go test -race ./...`
- **Lint**: `task lint` (golangci-lint)
- **Format**: `task fmt` (uses `gofumpt`, not `gofmt`)
- **All CI checks**: `task ci`

> Go caches are set to `.cache/` in the repo root. If running `go` directly outside `task`, set:
> `GOCACHE=$PWD/.cache/go-build GOTMPDIR=$PWD/.cache/go-tmp go test ./...`

## Architecture
- `cmd/capacities/main.go` — entry point, registers all commands
- `internal/cli/` — one file per command + `format.go` for shared output helpers
- `internal/api/client.go` — API types and HTTP client
- `internal/config/config.go` — config loading from `~/.config/capacities/config.yaml`

## Code Conventions
- CLI framework: `urfave/cli/v3` (not Cobra) — access flags via `cmd.String("name")`, `cmd.Bool("name")`
- All table-printing commands must include `--format (table|json)` flag defaulting to `table`
- Use shared formatters from `internal/cli/format.go`: `printSpaces`, `printStructures`, `printLookupResults`, `printTabTable`, `printJSON`
- New commands go in `internal/cli/<name>.go` and are registered in `main.go`
