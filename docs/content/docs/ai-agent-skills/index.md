---
title: 'AI Agent Skills'
weight: 5
---

## Overview

GTasks includes an embedded [Agent Skill](https://agentskills.io) for teaching compatible AI agents how to use `gtasks` effectively.

The `gtasks-cli` skill contains command references, examples, and workflow guidance for Google Tasks management. You can install it directly with `gtasks skills install`.

## Install with gtasks

Use the built-in skill management commands:

```bash
gtasks skills status
gtasks skills install
gtasks skills install --agent codex
gtasks skills uninstall --agent codex
```

Supported install targets:

- `claude` -> `~/.claude/skills/gtasks-cli/`
- `codex` -> `~/.agents/skills/gtasks-cli/`
- `openclaw` -> `~/.openclaw/skills/gtasks-cli/`

`gtasks skills install` copies the embedded skill files into the selected agent's skill directory so they are available in future sessions.

## Repository Location

For contributors, the canonical skill files are located in the [`internal/skills/assets/gtasks-cli/`](https://github.com/BRO3886/gtasks/tree/master/internal/skills/assets/gtasks-cli) directory of the repository.

## What's Included

### Main Skill File (SKILL.md)

The primary instruction file contains:
- Complete command reference for all gtasks operations
- Cross-platform installation checks (macOS, Linux, Windows)
- Environment variable setup instructions
- Authentication workflows
- Task and task list management examples
- Common usage patterns and best practices
- Error handling guides

### References

Additional documentation for advanced usage:

- **QUICK-REFERENCE.md**: Fast lookup table for common commands
- **ADVANCED.md**: Scripting, automation, and integration patterns

## Using with AI Agents

### Claude Code

If you're using [Claude Code](https://claude.com/claude-code), install the skill first:

```bash
gtasks skills install --agent claude
```

After installation, Claude can discover and load the skill automatically when relevant.

### Other AI Agents

Other agents can use the installed skill if they support the [Agent Skills specification](https://agentskills.io/spec) or scan compatible skill directories.

## Prerequisites for AI Usage

Before AI agents can use gtasks, ensure:

1. **GTasks is installed** and available in PATH
2. **Environment variables are set**:
   ```bash
   export GTASKS_CLIENT_ID="your-client-id.apps.googleusercontent.com"
   export GTASKS_CLIENT_SECRET="your-client-secret"
   ```
3. **You're authenticated**: Run `gtasks login`

The AI agent will check these prerequisites and guide you through setup if needed.

## Example AI Interactions

Once the skill is loaded, you can ask your AI agent:

- "Show me my tasks for today"
- "Add a task to my Work list: Finish the report by Friday"
- "Mark task #3 as complete"
- "Create a new task list called Shopping"
- "What tasks are overdue?"
- "Generate a daily task report"
- "Export my tasks to JSON"

## Benefits

**For Users:**
- Natural language task management
- Intelligent assistance with complex operations
- Automated workflows and reports
- Cross-platform compatibility

**For AI Agents:**
- Comprehensive command documentation
- Progressive loading for efficiency
- Built-in error handling patterns
- Platform-specific guidance

## Learn More

- [Agent Skills Website](https://agentskills.io)
- [Agent Skills Specification](https://agentskills.io/spec)
- [Canonical gtasks skill files](https://github.com/BRO3886/gtasks/tree/master/internal/skills/assets/gtasks-cli)

## Contributing

To improve the AI skills:

1. Edit files in `internal/skills/assets/gtasks-cli/`
2. Test with actual `gtasks` commands
3. Test `gtasks skills status` and at least one install flow
4. Update version number in SKILL.md frontmatter if needed
5. Submit a pull request

See the canonical skill directory for contribution guidance and examples.
