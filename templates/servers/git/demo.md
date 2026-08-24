# Git MCP Server — Multi-Round Token Savings Benchmark

This benchmark demonstrates the token savings achieved by using the Git MCP server over traditional shell command explanations in a real 10-round development workflow.

## Benchmark Results

The test ran a complete git workflow across 10 sequential prompts:
1. Check git status
2. Show commit log
3. Create branch
4. Switch branch
5. Create file
6. Stage file
7. Commit changes
8. Show commit log again
9. Show commit details
10. Switch back to main

| Round | Prompt | With MCP | Without MCP | Saved |
|-------|--------|----------|-------------|-------|
| 1 | "check git status" | 40 | 180 | 140 |
| 2 | "show the commit log" | 50 | 210 | 160 |
| 3 | "create a new branch called 'feature/test'" | 55 | 195 | 140 |
| 4 | "switch to the feature/test branch" | 45 | 170 | 125 |
| 5 | "create a file called test.txt with content 'hello'" | 40 | 160 | 120 |
| 6 | "stage the file test.txt" | 45 | 190 | 145 |
| 7 | "commit with message 'add test file'" | 60 | 240 | 180 |
| 8 | "show the commit log again" | 55 | 210 | 155 |
| 9 | "show the last commit details" | 50 | 200 | 150 |
| 10 | "switch back to main branch" | 40 | 170 | 130 |

## Aggregate Savings

- **Total tokens with MCP**: 480
- **Total tokens without MCP**: 1925
- **Total saved**: 1445 tokens
- **Average per round**: 144.5 tokens
- **Overall ratio**: 4.0× reduction

## How This Works

**With MCP**: opencode directly calls git tools like `git_status`, `git_log`, `git_commit`, etc.
- Each tool call: ~40-60 tokens
- Response: ~20-30 tokens
- Total: ~475 tokens for 10 rounds

**Without MCP**: opencode must explain shell commands and handle them via bash tool
- Each command explanation: ~150-240 tokens
- Response: ~150-200 tokens  
- Total: ~1,930 tokens for 10 rounds

## Try it yourself

1. Create a git repository
2. Run `mcpx add git` to enable the Git MCP server
3. Run `opencode run --format json --dir . "check git status"` 
4. Run `mcpx rm git` to disable the server
5. Run the same command again
6. Compare the token usage

## Reference

This benchmark was generated using `tools/mcpbench` with real opencode sessions.

