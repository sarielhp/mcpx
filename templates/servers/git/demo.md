# Git MCP Server — Token-Savings Demo

This demo shows how MCP servers reduce token usage by delegating to specialized tools instead of requiring AI to explain shell commands.

## How MCP Saves Tokens

**Without MCP**: The AI must explain shell command syntax, flags, and usage
- Example: "I'll run `git log --oneline --graph -5` to show the last 5 commits. The `--oneline` flag condenses each commit to a single line, and `--graph` shows the branching structure..."
- Estimated tokens: ~150

**With MCP**: The AI calls the tool directly with parameters
- Example: `git_log({"max_count": 5, "oneline": true, "graph": true})`
- Estimated tokens: ~20

**Savings**: ~130 tokens (6.5× reduction)

## Git MCP Server Tools

The Git MCP server provides 12 tools for repository management:

- **git_status** - Check repository status
- **git_log** - Show commit history  
- **git_commit** - Create commits
- **git_diff** - Show differences
- **git_add** - Stage files
- **git_reset** - Unstage files
- **git_branch** - List branches
- **git_checkout** - Switch branches
- **git_create_branch** - Create new branches
- **git_show** - Show commit details
- **git_diff_unstaged** - Show unstaged changes
- **git_diff_staged** - Show staged changes

## Try it yourself

1. Configure the server: `mcpx add git`
2. Ask opencode: "show git log with --oneline --graph for the last 5 commits"
3. Remove the server: `mcpx rm git`
4. Ask the same question again
5. Compare the response lengths

## Reference

This demo was generated using `tools/mcpbench` to illustrate the concept of MCP token savings.

