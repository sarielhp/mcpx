# Github MCP Server — Token-Savings Demo

This demo shows how MCP servers reduce token usage by delegating to specialized tools instead of requiring AI to explain complex APIs or commands.

## How MCP Saves Tokens

**Without MCP**: The AI must explain tool syntax, API calls, and usage patterns
- Example: "I'll use the GitHub API to list your repositories. This requires authenticating with a token, making a GET request to https://api.github.com/user/repos, and parsing the JSON response..."
- Estimated tokens: ~100

**With MCP**: The AI calls the tool directly with parameters
- Example: `github_repo_list()`
- Estimated tokens: ~15

**Savings**: ~85 tokens (6.7× reduction)

## Github MCP Server Tools

- **github_repo_list**
- **github_issue_list**
- **github_pr_list**
- **github_file_read**
- **github_file_write**

## Try it yourself

1. Configure the server: `mcpx add github`
2. Ask opencode: "list my repositories"
3. Remove the server: `mcpx rm github`
4. Ask the same question again
5. Compare the response lengths

## Reference

This demo was generated using `tools/mcpbench` to illustrate the concept of MCP token savings.

