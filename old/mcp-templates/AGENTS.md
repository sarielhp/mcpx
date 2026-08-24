# Agent Instructions: MCP Templates & Server Management

This directory (`~/.config/mcp-templates/`) is the central repository for MCP server definitions and project templates.

## Structure
- `servers/`: Individual server definitions (JSON fragments).
- `presets/`: Combined bundles for specific workflows (e.g. `golang-dev`, `full-stack`, `docker-stacks`, `paper-writing`).
- `bin/mcp-init`: Ruby CLI tool to initialize, manage, and verify MCP configurations across projects.

## Instantiation & Usage

By default, `mcp-init` configures **both Antigravity-CLI and OpenCode** for the current directory:

```bash
# Initialize presets for both Antigravity-CLI and OpenCode (default)
mcp-init golang-dev
mcp-init julia-dev
mcp-init latex-academic
mcp-init latex-starter
mcp-init latex-overleaf
mcp-init latex-full
mcp-init full-stack
mcp-init docker-stacks

# Merge multiple presets / individual servers
mcp-init golang-dev docker context7

# Target specific environments
mcp-init --opencode golang-dev
mcp-init --antigravity golang-dev
mcp-init --cursor golang-dev
mcp-init --claude golang-dev
mcp-init --vscode golang-dev
mcp-init --all full-stack

# Verification & Testing
mcp-init test                  # Test all MCP server handshakes in servers/
mcp-init test golang-dev       # Test servers defined in a preset
mcp-init test gopls context7   # Test specific MCP servers

# Utility commands
mcp-init list                  # List available presets and individual servers
mcp-init show <name>           # Display definition of a preset or server
mcp-init status                # Check active MCP configuration in target workspace
mcp-init overleaf-auth         # Extract Overleaf session cookie & token into .env
mcp-init remove <server...>    # Remove servers from workspace config
mcp-init --dry-run <name>      # Preview changes without modifying files
```

## New Hierarchical CLI Structure

The `mcp-init` tool now supports a hierarchical command structure with contextual help:

```bash
# Local Workspace Management (workspace-specific operations)
mcp-init local list                     # List active MCP servers in workspace
mcp-init local add <items...>          # Add/merge preset(s) or server(s) into workspace
mcp-init local rm <names...>            # Remove server(s) from workspace config
mcp-init local test [names...]          # Test MCP server handshakes in workspace
mcp-init local update [names...]       # Update local MCP configs with latest credentials

# Template Registry Commands (template registry operations)
mcp-init template list                  # List all available presets & individual servers
mcp-init template show <name>          # Display JSON template definition

# Auth Commands
mcp-init auth set <server|VAR> [VALUE]  # Set/update an API key in template .env
mcp-init auth list                      # Show status of all MCP keys
mcp-init auth get <server|VAR>          # Show currently configured value
mcp-init auth rm <server|VAR>           # Remove a key from .env
mcp-init auth overleaf [options]       # Extract Overleaf session cookie into .env

# All commands support contextual --help
mcp-init --help                         # Main help
mcp-init local --help                   # Local commands help
mcp-init template --help                # Template commands help
mcp-init auth --help                    # Auth commands help
mcp-init local add --help               # Specific command help
```
