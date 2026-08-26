# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v0.2.1] - 2026-08-24

### Added
- **Init Command**: Added `mcpx init` command to bootstrap workspaces from presets
- **Multi-Config Removal**: `mcpx rm` now removes servers from all config file types:
  - `opencode.json` (OpenCode format)
  - `.agents/plugins/mcp/mcp_config.json` (Antigravity format)
  - `.cursor/mcp.json` (Cursor format)
  - `.claude/mcp.json` (Claude format)
  - `.vscode/mcp.json` (VS Code format)
  - `mcp.json` (Standard format)
- **Subdirectory Config Support**: Configs now written to proper subdirectories:
  - OpenCode: `.opencode/opencode.json`
  - Antigravity: `.agents/plugins/mcp/mcp_config.json`
  - Cursor: `.cursor/mcp.json`
  - Claude: `.claude/mcp.json`
  - VS Code: `.vscode/mcp.json`
- **Recursive Recommendations**: `--with-recommended` flag now works recursively using BFS algorithm
- **Comprehensive Unit Tests**: Added Go unit tests for all packages:
  - `internal/types/types_test.go`
  - `internal/envvault/envvault_test.go`
  - `internal/config/config_test.go`
  - `internal/runner/runner_test.go`
  - `main_test.go`
- **Documentation**: Added README.md and LICENSE files

### Changed
- **Recommendation Relationships**: All 23 server templates now have meaningful `recommends` relationships
- **Config Writers**: Updated to support subdirectory paths and new config formats
- **Command Registration**: Improved command structure for better CLI handling

### Fixed
- **Multi-Config Removal**: Fixed `rm` command to properly handle all config file types
- **Init Command**: Implemented proper workspace initialization with preset support

## [v0.2.0] - Initial Release

### Added
- Core MCP server management functionality
- Template registry with 23 server templates
- Config writers for OpenCode, Antigravity, Cursor, Claude, VS Code
- Parallel health checking with `mcpx validate`
- Environment variable management with `mcpx auth`
- Template management with `mcpx template`
- Preset support with `mcpx add` and `mcpx list`

[Unreleased]: https://github.com/sarielhp/mcpx/compare/v0.2.1...HEAD
[v0.2.1]: https://github.com/sarielhp/mcpx/compare/v0.2.0...v0.2.1