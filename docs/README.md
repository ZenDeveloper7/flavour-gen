# Documentation Index

Welcome to flavour-gen documentation.

## Getting Started

- [Installation](./installation.md)
- [Quick Start](./quickstart.md)
- [Input Format](./input-format.md)

## Architecture

### Core Files

- [main.go](./main.md) - Entry point
- [cmd/root.go](./cmd/root.md) - Root command
- [cmd/create.go](./cmd/create.md) - Create subcommand

### Packages

- [pkg/config](./pkg/config.md) - Configuration parsing
- [pkg/flavor](./pkg/flavour.md) - Theme duplication
- [pkg/icon](./pkg/icon.md) - Icon generation
- [pkg/keystore](./pkg/keystore.md) - Keystore generation

## Project Structure

```
flavour-gen/
├── main.go              # Entry point
├── cmd/
│   ├── root.go        # Root command
│   └── create.go      # Create subcommand
├── pkg/
│   ├── config/        # Configuration
│   ├── flavor/        # Theme processing
│   ├── icon/          # Icon generation
│   └── keystore/      # Keystore generation
├── docs/              # This documentation
└── templates/         # Template files
```

## Additional Resources

- [GitHub Repository](https://github.com/ZenDeveloper7/flavour-gen)
- [Releases](https://github.com/ZenDeveloper7/flavour-gen/releases)
- [Homebrew Tap](https://github.com/ZenDeveloper7/homebrew-tap)
