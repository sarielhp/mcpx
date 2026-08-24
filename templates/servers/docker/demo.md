# Docker MCP Server — Token-Savings Demo

This demo shows how MCP servers reduce token usage by delegating to specialized tools instead of requiring AI to explain complex APIs or commands.

## How MCP Saves Tokens

**Without MCP**: The AI must explain tool syntax, API calls, and usage patterns
- Example: "I'll use the GitHub API to list your repositories. This requires authenticating with a token, making a GET request to https://api.github.com/user/repos, and parsing the JSON response..."
- Estimated tokens: ~100

**With MCP**: The AI calls the tool directly with parameters
- Example: `docker_ps()`
- Estimated tokens: ~15

**Savings**: ~85 tokens (6.7× reduction)

## Docker MCP Server Tools

- **docker_ps**
- **docker_images**
- **docker_build**
- **docker_run**
- **docker_logs**

## Try it yourself

1. Configure the server: `mcpx add docker`
2. Ask opencode: "list running containers"
3. Remove the server: `mcpx rm docker`
4. Ask the same question again
5. Compare the response lengths

## Reference

This demo was generated using `tools/mcpbench` to illustrate the concept of MCP token savings.

