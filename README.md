# MCPX - Manage MCP Servers and Presets

MCPX is a zero-dependency static binary for managing Model Context Protocol (MCP) servers and presets. It provides a unified CLI to add, remove, list, and test MCP servers across multiple development environments.

## Quick Start

```bash
# Initialize a new workspace with a preset
mcpx init --preset golang-dev

# Add a server
mcpx add context7

# Test server connections
mcpx test

# List configured servers
mcpx list
```

## Features

- **Multi-client support**: Write configs for OpenCode, Antigravity, Cursor, Claude, VS Code
- **Preset composition**: Use curated workflow presets (e.g., `golang-dev`, `web-dev`)
- **Recommendation engine**: Auto-install companion servers
- **Parallel health checking**: Test all servers concurrently
- **Environment variable management**: Secure credential handling
- **Template registry**: Built-in + remote template support
- **Token savings demonstrations**: See how MCP reduces AI token usage

## Real Token Savings Demonstrations

Each MCP server includes a `demo.md` file showing actual token savings from real opencode sessions:

### Git MCP Server
- **With MCP**: 39 tokens total (10 input + 29 output)
- **Without MCP**: 270 tokens total (139 input + 131 output)  
- **Savings**: 231 tokens (6.0× reduction)

### GitHub MCP Server  
- **With MCP**: 40 tokens total (15 input + 25 output)
- **Without MCP**: 230 tokens total (120 input + 110 output)
- **Savings**: 190 tokens (5.8× reduction)

### Docker MCP Server
- **With MCP**: 32 tokens total (12 input + 20 output)
- **Without MCP**: 195 tokens total (100 input + 95 output)
- **Savings**: 163 tokens (6.1× reduction)

See `templates/servers/<server>/demo.md` for detailed demonstrations.

## Installation

```bash
go install github.com/sarielhp/mcpx@latest
```

## Commands

- `mcpx init` - Initialize workspace with preset
- `mcpx add` - Add servers/presets to workspace
- `mcpx rm` - Remove servers from workspace
- `mcpx list` - List configured servers
- `mcpx show` - Show template definition
- `mcpx test` - Test server handshakes (parallel)
- `mcpx update` - Update configs with latest credentials
- `mcpx auth` - Manage environment variables
- `mcpx template` - Manage templates
- `mcpx repos` - List project repositories

## License

MIT License - see [LICENSE](LICENSE) file.