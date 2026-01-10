# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Limelight is a macOS CLI tool for proactive automation of Philips Hue lights and scenes. It integrates with the Hue Bridge V2 API and uses 1Password CLI for secure credential storage.

**Current Status:** Phase 1 (Foundation) completed. See `HUE_AUTOMATION_LLM_PLAN.org` for the complete implementation roadmap.

## Build and Development Commands

### Build
```bash
make build       # Format, test, and build (recommended)
make all         # Same as make build
```

The build process runs:
1. `gofmt -w .` - Format all Go files
2. `go test ./...` - Run all tests
3. `go build -o limelight ./cmd/limelight` - Build the binary

### Other Commands
```bash
make fmt         # Format code only
make test        # Run tests only (no verbose output)
make clean       # Remove binary and clean artifacts

# Run tests with verbose output
go test -v ./...

# Run tests for a specific package
go test -v ./internal/bridge

# Run a specific test
go test -v ./internal/bridge -run TestFunctionName
```

### Running the CLI
```bash
./limelight setup              # Initial setup and bridge pairing
./limelight lights list        # List all lights
./limelight lights set <id>    # Control a light
./limelight scenes list        # List all scenes
./limelight scenes activate <id>  # Activate a scene
```

## Architecture

### Core Components

1. **Hue Bridge Client (`internal/bridge/`)**
   - `client.go`: Base HTTP client with V2 API support, uses self-signed cert (InsecureSkipVerify)
   - `auth.go`: Button press authentication flow with retry logic (30 attempts, 2s interval)
   - `lights.go`: Light resource operations
   - `scenes.go`: Scene resource operations
   - `groups.go`: Group/room resource operations
   - All API calls go through `doRequest()` which handles auth headers and error responses

2. **Credential Management (`internal/credentials/`)**
   - `config.go`: Manages `~/.config/limelight/config.json` with bridge IP and 1Password item name
   - `onepassword.go`: Integrates with `op` CLI for secure API key storage
   - Config file uses atomic writes (write to .tmp then rename)
   - API keys are stored in 1Password items with `api_key[concealed]` field

3. **CLI Commands (`cmd/limelight/commands/`)**
   - `setup.go`: Interactive setup wizard for bridge pairing
   - `lights.go`: Light control commands
   - `scenes.go`: Scene management commands
   - All commands receive a `*zap.Logger` instance

### Key Design Patterns

- **Error Handling**: Uses `github.com/cockroachdb/errors` for wrapping with context
- **Logging**: Uses `go.uber.org/zap` with development config in main
- **HTTP Client**: 10-second timeout, TLS verification disabled for Hue Bridge self-signed certs
- **CLI Framework**: Built with `github.com/spf13/cobra`

### Configuration

- **Config Location**: `~/.config/limelight/config.json`
- **Config Fields**:
  - `bridge_ip`: Hue Bridge IP address
  - `onepassword_item_name`: 1Password item storing API key
  - `latitude`, `longitude`: For future sunrise/sunset automation

### Authentication Flow

The setup process (`internal/bridge/auth.go`):
1. User provides bridge IP address
2. User is prompted to press the physical button on the bridge
3. Polls `/api` endpoint every 2 seconds for up to 60 seconds
4. Waits for error code 101 (link button not pressed) to clear
5. Returns username (API key) on success
6. Stores API key in 1Password if available

### Hue V2 API Details

- **Base URL**: `https://{bridge_ip}/clip/v2`
- **Authentication**: `hue-application-key` header (not URL parameter)
- **Content-Type**: `application/json`
- **TLS**: Self-signed certificate (verification disabled)

## Future Development Areas

Based on `HUE_AUTOMATION_LLM_PLAN.org`, upcoming phases include:

- **Phase 2**: SQLite schema for automation rules, sunrise/sunset calculator
- **Phase 3**: Automation engine with triggers/conditions/actions
- **Phase 4**: macOS presence detection
- **Phase 5**: Natural language parser for automation commands
- **Phase 6**: Testing and documentation

Empty package directories already exist for:
- `internal/db/` - Database layer
- `internal/automation/` - Rule engine
- `internal/presence/` - macOS presence detection
- `internal/astro/` - Sunrise/sunset calculations
- `internal/nlp/` - Natural language parser
- `internal/daemon/` - Background service

## Dependencies

- `github.com/cockroachdb/errors` - Error handling with wrapping
- `go.uber.org/zap` - Structured logging
- `github.com/spf13/cobra` - CLI framework
- `github.com/stretchr/testify` - Testing (imported but no tests written yet)

## Important Notes

- No tests exist yet (Phase 1 focused on functionality)
- 1Password CLI (`op`) integration is optional but recommended for security
- The tool is macOS-specific due to planned presence detection features
- Bridge IP must be manually configured (no mDNS discovery implemented)
