# OASC - OpenAPI Specification Combiner

A CLI tool for combining multiple OpenAPI specifications into a single file.

## Features

- Combines multiple OpenAPI 3.x specifications
- Supports both YAML and JSON formats
- Cross-platform support (macOS, Linux, Windows)
- Debug mode for detailed operation logs

## Installation

### Go install

```bash
# Install specific version
go install github.com/haryoiro/oasc/cmd/oasc@v0.1.35

# Install latest version
go install github.com/haryoiro/oasc/cmd/oasc@latest
```

### Build from Source

```bash
# Install Go 1.24 or later
go install github.com/haryoiro/oasc/cmd/oasc@latest
```

## Usage

```bash
# Basic usage
oasc -f spec1.yaml -f spec2.yaml -o merged.yaml

# Specify output format
oasc -f spec1.yaml -f spec2.yaml -o merged.json --format json

# Enable debug mode
oasc -f spec1.yaml -f spec2.yaml -o merged.yaml --debug
```

### Options

- `-f, --file`: Input OpenAPI file paths (can be specified multiple times)
- `-o, --output`: Output file path (default: merged.yaml)
- `-F, --format`: Output format (json or yaml)
- `--debug`: Enable debug logging
- `--version`: Show version information

## Development

### Prerequisites

- Go 1.24 or later
- [mise](https://github.com/jdx/mise) for task management
- [goreleaser](https://goreleaser.com/) for releases

### Development Setup

```bash
# Install mise
curl https://mise.jdx.dev/install.sh | sh

# Install development tools
mise install

# Run tests
mise run test
```

### Release Process

```bash
# Create release tag (this will push a tag and trigger GitHub Actions)
mise run release-tag 0.1.35
```

GitHub Actionsが自動でビルド・リリースを行います。

> バージョン管理はGitのタグのみで行われます。

## License

MIT License
