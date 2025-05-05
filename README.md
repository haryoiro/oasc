# OASC - OpenAPI Specification Combiner

A CLI tool for combining multiple OpenAPI specifications into a single file.

## Features

- Combines multiple OpenAPI 3.x specifications
- Supports both YAML and JSON formats
- Cross-platform support (macOS, Linux, Windows)
- Debug mode for detailed operation logs

## Installation

### Build from Source

```bash
# Install Go 1.24 or later
go install github.com/haryoiro/oasc@latest
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

### Creating a Release

1. Update the version in `version.go`
2. Commit your changes
3. Create and push a tag:
   ```bash
   git tag -a v0.1.0 -m "Release v0.1.0"
   git push origin v0.1.0
   ```
4. GitHub Actions will automatically create a release with binaries for all platforms

## License

MIT License
