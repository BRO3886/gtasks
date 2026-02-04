# GTasks Skills Directory

This directory contains [Agent Skills](https://agentskills.io) that enable AI agents to effectively use the gtasks CLI tool.

## What are Agent Skills?

Agent Skills are a standardized format for providing AI agents (like Claude) with instructions on how to use tools and perform tasks. They follow a specific directory structure and format that allows agents to:

1. Discover available skills through metadata
2. Load instructions when needed
3. Access reference documentation progressively
4. Execute helper scripts

## Available Skills

### gtasks-cli

A comprehensive skill for managing Google Tasks from the command line.

**Location:** `skills/gtasks-cli/`

**What it includes:**
- Complete command reference for all gtasks operations
- Task and task list management workflows
- Common usage patterns and examples
- Advanced scripting and automation guides
- Helper scripts for daily reports, backups, and imports

**When to use:**
- User wants to view, create, or manage Google Tasks
- User mentions to-do lists, task lists, or Google Tasks
- User needs to check what tasks they have
- User wants to add tasks with deadlines
- User needs to organize tasks into lists

## Directory Structure

```
skills/
├── README.md                    # This file
└── gtasks-cli/                  # Google Tasks CLI skill
    ├── SKILL.md                 # Main skill instructions (required)
    ├── references/              # Additional documentation
    │   └── ADVANCED.md         # Advanced usage patterns
    └── scripts/                 # Helper scripts
        ├── daily-report.sh     # Generate daily task report
        ├── backup-tasks.sh     # Backup all tasks
        └── import-tasks.sh     # Import tasks from file
```

## Using Skills with Claude Code

When using Claude Code CLI, skills in this directory can be loaded using the `/plugin` command or automatically detected based on the skill descriptions.

### Manual Loading

```bash
claude /plugin skills/gtasks-cli
```

### Automatic Detection

Skills with good descriptions will be automatically suggested when relevant. For example, if you say "show me my Google Tasks", Claude may automatically load the gtasks-cli skill.

## Creating New Skills

To create a new skill for this project:

1. Create a new directory under `skills/` with a descriptive name
2. Add a `SKILL.md` file with YAML frontmatter and instructions
3. Optionally add `references/` and `scripts/` directories
4. Follow the [Agent Skills specification](https://agentskills.io/spec)

### Minimum SKILL.md Format

```yaml
---
name: skill-name
description: What the skill does and when to use it
---

# Skill Instructions

Write clear, actionable instructions here...
```

## Validation

Validate skills using the Agent Skills reference implementation:

```bash
# Install the validator (if available)
npm install -g @agentskills/validator

# Validate a skill
agentskills validate skills/gtasks-cli
```

## Best Practices

1. **Clear descriptions**: Make the `description` field specific about when to use the skill
2. **Progressive disclosure**: Keep SKILL.md concise, move detailed info to references
3. **Executable scripts**: Include ready-to-run helper scripts
4. **Examples**: Provide real-world examples throughout
5. **Error handling**: Document common errors and solutions

## Contributing

When adding new commands or features to gtasks, update the corresponding skill:

1. Add new commands to `SKILL.md`
2. Update examples and workflows
3. Add advanced patterns to `references/ADVANCED.md`
4. Create helper scripts if applicable

## Resources

- [Agent Skills Documentation](https://agentskills.io)
- [Agent Skills Specification](https://agentskills.io/spec)
- [Example Skills](https://github.com/agentskills/agentskills)
- [GTasks Repository](https://github.com/BRO3886/gtasks)
