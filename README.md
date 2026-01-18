# ccstatus

<img width="1024" height="500" alt="12" src="https://github.com/user-attachments/assets/d762d9e4-2c5f-47d5-b44b-a8c2f80288bc" />

A CLI tool that shows Claude Code session and weekly usage in the statusline.

## Overview

ccstatus is designed specifically for Claude Code. It uses Anthropic's API to display real usage data directly in your Claude Code statusline, giving you visibility into your session and weekly consumption without leaving your workflow.

## Installation

### Homebrew (recommended)

```bash
brew tap tharuxpert/ccstatus
brew install ccstatus
```

### Manual install from source


Requires Go 1.21 or later.

```bash
go install github.com/tharuxpert/ccstatus@latest
```

## Setup

Run the install command to configure ccstatus:

```bash
ccstatus install
```

This command safely updates `~/.claude/settings.json` to add the statusline configuration. No configuration is modified without user consent—you will be prompted before any changes are made.

## Usage

### Commands

| Command | Description |
|---------|-------------|
| `ccstatus install` | Configure ccstatus in Claude Code settings |
| `ccstatus uninstall` | Remove ccstatus from Claude Code settings |
| `ccstatus config` | Configure statusline display options |
| `ccstatus doctor` | Run diagnostic checks on your configuration |
| `ccstatus version` | Print the version number |

Note: Running `ccstatus` without arguments outputs statusline data. This is intended to be called by Claude Code and will not produce meaningful output in a normal terminal session.

## Configuration

After running `ccstatus install`, the following entry is added to your Claude Code settings:

```json
{
  "statusline": {
    "command": "ccstatus"
  }
}
```

This tells Claude Code to execute ccstatus and display its output in the statusline.

### Statusline Options

Use `ccstatus config` to customize what information is displayed in the statusline. This interactive command allows you to toggle the following options:

- **Session Usage**: Show current session usage percentage
- **Weekly Usage**: Show weekly usage percentage
- **Reset Times**: Show when usage limits reset
- **Git Branch**: Show current git branch name

Configuration is saved to `~/.claude/ccstatus.json` and takes effect immediately.

## Compatibility

- macOS
- Claude Code CLI

## Security

- OAuth tokens are retrieved from the system keychain (stored by Claude Code)
- No telemetry or tracking
- No data is sent to third parties

## Feedback & Ideas

💬 Feedback & ideas: https://github.com/tharuxpert/ccstatus/discussions

## License

MIT
