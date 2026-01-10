# Limelight

A macOS CLI tool for proactive automation of Philips Hue lights and scenes.

## Features

### Phase 1 (Completed)
- Hue Bridge V2 API client with authentication
- 1Password CLI integration for secure credential storage
- CLI commands for:
  - Initial setup and bridge pairing (`setup`)
  - Listing and controlling lights (`lights list`, `lights set`)
  - Listing and activating scenes (`scenes list`, `scenes activate`)

## Installation

### Prerequisites
- Go 1.21 or later
- 1Password CLI (optional, for secure credential storage)

### Build
```bash
go build -o limelight ./cmd/limelight
```

## Usage

### Initial Setup
Run the setup wizard to pair with your Hue Bridge:
```bash
./limelight setup
```

This will:
1. Prompt for your Hue Bridge IP address
2. Configure 1Password integration (if available)
3. Guide you through the button press authentication flow
4. Store credentials securely

### List Lights
```bash
./limelight lights list
```

### Control a Light
```bash
# Turn on a light
./limelight lights set <light-id> --on

# Turn off a light
./limelight lights set <light-id> --off

# Set brightness
./limelight lights set <light-id> --on --brightness 75
```

### List Scenes
```bash
./limelight scenes list
```

### Activate a Scene
```bash
./limelight scenes activate <scene-id>
```

## Configuration

Configuration is stored in `~/.config/limelight/config.json` and includes:
- Bridge IP address
- 1Password item name (for API key storage)
- Location coordinates (for future sunrise/sunset automation)

API keys are stored securely in 1Password when the CLI is available.

## Project Structure

```
limelight/
├── cmd/
│   └── limelight/          # CLI entry point and commands
│       └── commands/       # Command implementations
├── internal/
│   ├── bridge/             # Hue V2 API client
│   ├── credentials/        # Config and 1Password integration
│   ├── db/                 # Database layer (future)
│   ├── automation/         # Automation engine (future)
│   ├── presence/           # macOS presence detection (future)
│   ├── astro/              # Sunrise/sunset calculations (future)
│   ├── nlp/                # Natural language parser (future)
│   └── daemon/             # Background service (future)
└── README.md
```

## Development Status

See `HUE_AUTOMATION_LLM_PLAN.org` for the complete implementation plan.

**Current Phase:** Phase 1 (Foundation) - Completed

**Next Phase:** Phase 2 (Database & Core Logic)
- SQLite schema design
- Automation rules storage
- Sunrise/sunset calculator
