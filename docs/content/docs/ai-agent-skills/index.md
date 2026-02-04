---
title: 'AI Agent Skills'
weight: 5
---

## Overview

GTasks includes [Agent Skills](https://agentskills.io) support for AI agents like Claude Code. Agent Skills provide a standardized format for teaching AI agents how to effectively use command-line tools.

The gtasks-cli skill contains comprehensive instructions, examples, and helper scripts that enable AI agents to intelligently assist with Google Tasks management.

## Location

Skills are located in the [`skills/gtasks-cli/`](https://github.com/BRO3886/gtasks/tree/master/skills/gtasks-cli) directory of the repository.

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
- **TASK-IMPORT-TEMPLATE.txt**: Template for bulk task imports

### Helper Scripts

Ready-to-run bash scripts:

- **daily-report.sh**: Generate comprehensive daily task reports
- **backup-tasks.sh**: Backup all tasks and task lists to JSON
- **import-tasks.sh**: Bulk import tasks from text files

## Using with AI Agents

### Claude Code

If you're using [Claude Code](https://claude.com/claude-code), the skill may be automatically loaded when relevant, or you can explicitly load it:

```bash
# Let Claude discover and load skills automatically
claude "show me my Google Tasks"

# Or explicitly reference the skill
claude /plugin skills/gtasks-cli
```

### Other AI Agents

Any AI agent that supports the [Agent Skills specification](https://agentskills.io/spec) can use this skill. Refer to your AI agent's documentation for loading skills.

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
- [GTasks Skills README](https://github.com/BRO3886/gtasks/blob/master/skills/README.md)

## Contributing

To improve the AI skills:

1. Edit files in `skills/gtasks-cli/`
2. Test with actual gtasks commands
3. Update version number in SKILL.md frontmatter
4. Submit a pull request

See the [skills README](https://github.com/BRO3886/gtasks/blob/master/skills/README.md) for detailed contribution guidelines.
