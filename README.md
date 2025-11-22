# Capacities CLI

A command-line interface for interacting with the [Capacities.io](https://capacities.io) API. This tool allows you to manage your Capacities workspace directly from your terminal.

## Features

- **Daily Notes**: Save content to your daily notes
- **Search**: Search across your Capacities workspace
- **Spaces**: List and view information about your spaces
- **Web Links**: Save web links to your Capacities workspace
- **Configuration**: Easy setup and configuration management

## Installation

### Prerequisites

- Go 1.25.4 or later

### Build from Source

```bash
# Clone the repository
git clone https://github.com/klauern/capacities-cli.git
cd capacities-cli

# Build using Task
task build

# Or build directly with Go
go build -o bin/capacities ./cmd/capacities
```

The binary will be available at `bin/capacities`.

## Configuration

Before using the CLI, you need to configure it with your Capacities API token:

1. Get your API token from [Capacities Settings](https://app.capacities.io/settings/security)
2. Run the configuration command:

```bash
./bin/capacities configure
```

This will create a configuration file at `~/.config/capacities/config.yaml` with the following structure:

```yaml
token: YOUR_API_TOKEN
default_space_id: YOUR_SPACE_ID
```

## Usage

### Daily Notes

Save content to your daily note:

```bash
capacities daily save "Your note content here"
```

### Search

Search your Capacities workspace:

```bash
capacities search "search query"
```

### Spaces

List all spaces:

```bash
capacities spaces
```

Get information about a specific space:

```bash
capacities space-info --space-id YOUR_SPACE_ID
```

### Web Links

Save a web link:

```bash
capacities save-weblink --url "https://example.com" --title "Example Site"
```

## Development

This project uses [Task](https://taskfile.dev/) for common development tasks.

### Available Tasks

- `task build` - Build the binary
- `task fmt` - Format code using gofumpt
- `task deps` - Download dependencies and install development tools
- `task test` - Run tests

### Project Structure

```
capacities-cli/
├── cmd/
│   └── capacities/      # Main application entry point
├── internal/
│   ├── api/            # API client implementation
│   ├── cli/            # CLI command implementations
│   └── config/         # Configuration management
├── Taskfile.yml        # Task definitions
└── README.md          # This file
```

## Technology Stack

- **Language**: Go (Golang)
- **CLI Framework**: [urfave/cli/v3](https://github.com/urfave/cli/v3)
- **Configuration**: YAML via `gopkg.in/yaml.v3`
- **HTTP Client**: Standard `net/http`

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Issue Management

This project uses [Beads](https://github.com/jonny-so/beads) for issue tracking. See [AGENTS.md](AGENTS.md) for details on the development workflow and issue management.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Resources

- [Capacities.io](https://capacities.io)
- [Capacities API Documentation](https://api.capacities.io/docs)
