# MCP Benchmark Tools

This directory contains tools for benchmarking MCP (Model Context Protocol) servers.

## `mcpbench` - Token Savings Benchmark

The `mcpbench` script benchmarks how much token usage is reduced when using MCP servers versus traditional shell command explanations.

### Usage

```bash
# Benchmark a single server
tools/mcpbench --server git

# Generate demo files for all servers
tools/mcpbench --all --write

# Benchmark and generate demo for a specific server
tools/mcpbench --server git --write
```

### What it does

1. Creates a demo.md file in each server's template directory
2. Shows the token savings concept for that server
3. Lists the available tools
4. Provides instructions for users to try it themselves

### Example Output

The script generates a markdown file that shows:
- How much token usage is reduced with MCP
- Example prompts and expected token counts
- The list of tools provided by the server
- Instructions for users to test it themselves

### Supported Servers

Currently supports:
- git (with 12 tools)
- github (with 5 tools)
- docker (with 5 tools)

More servers can be added by extending the PROMPTS hash in the script.