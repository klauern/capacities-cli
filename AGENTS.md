# AGENTS.md

## Project Overview
This project is a CLI tool for interacting with the [Capacities API](https://api.capacities.io/docs). The current focus is on implementing the `save-to-daily-note` functionality, allowing users to quickly append text to their daily notes from the command line.

## Technology Stack
- **Language**: Go (Golang)
- **CLI Framework**: [`urfave/cli/v3`](https://github.com/urfave/cli/v3) (Alpha/Beta)
- **Configuration**: `gopkg.in/yaml.v3`
- **HTTP Client**: Standard `net/http`

## Configuration
The CLI requires a configuration file located at `~/.config/capacities/config.yaml`.
```yaml
token: YOUR_API_TOKEN
default_space_id: YOUR_SPACE_ID
```

## Development Workflow
- **Build**: `go build -o capacities cmd/capacities/main.go`
- **Run**: `./capacities daily save "Note text"`
- **Test**: `go test ./...`

---

## Issue Management (Beads)

This project uses the **Beads** framework for issue and task management. As an AI agent, you should use `bd` to track your work, manage dependencies, and find new tasks.

### Core Workflow

1.  **Find Work**: Run `bd ready` to see issues that are unblocked and ready to be worked on.
2.  **Create Issues**: When you discover new work or bugs, create an issue using `bd create "Title"`.
    *   Use `-p` to set priority (0=highest).
    *   Use `-t` to set type (e.g., `bug`, `feature`, `task`).
3.  **Manage Dependencies**: If a task depends on another, use `bd dep add <child> <parent>`.
    *   Example: `bd dep add bd-5 bd-2` means `bd-2` blocks `bd-5`.
4.  **Update Status**:
    *   Start working: `bd update <id> --status in_progress`
    *   Complete work: `bd close <id>`

### Commands Reference

*   `bd list`: List all issues.
*   `bd show <id>`: Show details of an issue.
*   `bd ready`: Show actionable issues.
*   `bd create "<title>"`: Create a new issue.
*   `bd update <id> ...`: Update issue fields.
*   `bd close <id>`: Close an issue.

### Tips for Agents

*   Always check `bd ready` before starting a new task if you don't have one assigned.
*   Break down large tasks into smaller issues and link them with dependencies.
*   Use the `--json` flag if you need to parse output programmatically (though standard output is designed to be readable).
