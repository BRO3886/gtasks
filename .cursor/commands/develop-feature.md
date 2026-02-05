# Develop Feature

Implement a new feature for the gtasks CLI tool.

## Instructions

Before implementing, explore the codebase to understand:
- How similar features are implemented in `cmd/tasks.go` and `cmd/tasklists.go`
- How API calls are structured in `api/tasks.go` and `api/tasklists.go`
- Existing patterns for error handling, prompts, and output formatting

Then follow the dev pipeline:

1. **Implement the feature** following existing code patterns
2. **Build and test**:
   - Credentials are in `.env` (`GTASKS_CLIENT_ID` and `GTASKS_CLIENT_SECRET`)
   - Run `make dev EMBED_CREDS=1` to build with embedded credentials from `.env`
   - Test with `./gtasks <command>`
3. **Update docs/** if the feature adds/changes commands (Hugo site in `docs/content/`)
4. **Update README.md** if needed for user-facing changes
5. **Commit** with a clear message
6. **Do NOT push or release** unless explicitly told

## Feature to Implement

Will be provided by the user.
