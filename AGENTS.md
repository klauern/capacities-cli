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

<!-- BEGIN BEADS INTEGRATION -->
## Issue Tracking with bd (beads)

**IMPORTANT**: This project uses **bd (beads)** for ALL issue tracking. Do NOT use markdown TODOs, task lists, or other tracking methods.

### Why bd?

- Dependency-aware: Track blockers and relationships between issues
- Git-friendly: Dolt-powered version control with native sync
- Agent-optimized: JSON output, ready work detection, discovered-from links
- Prevents duplicate tracking systems and confusion

### Quick Start

**Check for ready work:**

```bash
bd ready --json
```

**Create new issues:**

```bash
bd create "Issue title" --description="Detailed context" -t bug|feature|task -p 0-4 --json
bd create "Issue title" --description="What this issue is about" -p 1 --deps discovered-from:bd-123 --json
```

**Claim and update:**

```bash
bd update <id> --claim --json
bd update bd-42 --priority 1 --json
```

**Complete work:**

```bash
bd close bd-42 --reason "Completed" --json
```

### Issue Types

- `bug` - Something broken
- `feature` - New functionality
- `task` - Work item (tests, docs, refactoring)
- `epic` - Large feature with subtasks
- `chore` - Maintenance (dependencies, tooling)

### Priorities

- `0` - Critical (security, data loss, broken builds)
- `1` - High (major features, important bugs)
- `2` - Medium (default, nice-to-have)
- `3` - Low (polish, optimization)
- `4` - Backlog (future ideas)

### Workflow for AI Agents

1. **Check ready work**: `bd ready` shows unblocked issues
2. **Claim your task atomically**: `bd update <id> --claim`
3. **Work on it**: Implement, test, document
4. **Discover new work?** Create linked issue:
   - `bd create "Found bug" --description="Details about what was found" -p 1 --deps discovered-from:<parent-id>`
5. **Complete**: `bd close <id> --reason "Done"`

### Auto-Sync

bd automatically syncs via Dolt:

- Each write auto-commits to Dolt history
- Use `bd dolt push`/`bd dolt pull` for remote sync
- No manual export/import needed!

### Important Rules

- ✅ Use bd for ALL task tracking
- ✅ Always use `--json` flag for programmatic use
- ✅ Link discovered work with `discovered-from` dependencies
- ✅ Check `bd ready` before asking "what should I work on?"
- ❌ Do NOT create markdown TODO lists
- ❌ Do NOT use external issue trackers
- ❌ Do NOT duplicate tracking systems

For more details, see README.md and docs/QUICKSTART.md.

## Landing the Plane (Session Completion)

**When ending a work session**, you MUST complete ALL steps below. Work is NOT complete until `git push` succeeds.

**MANDATORY WORKFLOW:**

1. **File issues for remaining work** - Create issues for anything that needs follow-up
2. **Run quality gates** (if code changed) - Tests, linters, builds
3. **Update issue status** - Close finished work, update in-progress items
4. **PUSH TO REMOTE** - This is MANDATORY:
   ```bash
   git pull --rebase
   bd sync
   git push
   git status  # MUST show "up to date with origin"
   ```
5. **Clean up** - Clear stashes, prune remote branches
6. **Verify** - All changes committed AND pushed
7. **Hand off** - Provide context for next session

**CRITICAL RULES:**
- Work is NOT complete until `git push` succeeds
- NEVER stop before pushing - that leaves work stranded locally
- NEVER say "ready to push when you are" - YOU must push
- If push fails, resolve and retry until it succeeds

<!-- END BEADS INTEGRATION -->
