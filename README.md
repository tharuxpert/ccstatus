# ccstatus

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
