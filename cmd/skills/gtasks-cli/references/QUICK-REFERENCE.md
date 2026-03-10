# Quick Reference Card

Fast lookup for common gtasks commands. For detailed explanations, see the main SKILL.md file.

## System Checks (Run First!)

```bash
# Check if gtasks is installed (cross-platform)
gtasks --version 2>/dev/null || gtasks.exe --version 2>/dev/null

# macOS/Linux: which command
which gtasks

# Windows: where command
where gtasks

# Check environment variables (macOS/Linux)
echo $GTASKS_CLIENT_ID
echo $GTASKS_CLIENT_SECRET

# Check environment variables (Windows PowerShell)
echo $env:GTASKS_CLIENT_ID
echo $env:GTASKS_CLIENT_SECRET

# Check authentication status
gtasks tasklists view &>/dev/null && echo "✓ Authenticated" || echo "✗ Not authenticated"
```

## Authentication

```bash
gtasks login                    # Authenticate with Google
gtasks logout                   # Remove credentials
```

## Task Lists

```bash
gtasks tasklists view                    # List all task lists
gtasks tasklists add -t "List Name"      # Create task list
gtasks tasklists rm                      # Delete task list (interactive)
gtasks tasklists update -t "New Name"    # Rename task list (interactive)
```

## View Tasks

```bash
gtasks tasks view                        # View tasks (interactive list selection)
gtasks tasks view -l "Work"              # View tasks in specific list
gtasks tasks view -i                     # Include completed tasks
gtasks tasks view --completed            # Show only completed tasks
gtasks tasks view --sort=due             # Sort by due date
gtasks tasks view --sort=title           # Sort by title
gtasks tasks view --format=json          # JSON output
gtasks tasks view --format=csv           # CSV output
```

## Create Tasks

```bash
gtasks tasks add                                          # Interactive mode
gtasks tasks add -t "Title"                               # With title only
gtasks tasks add -t "Title" -n "Notes"                    # With notes
gtasks tasks add -t "Title" -d "2024-12-25"              # With due date
gtasks tasks add -t "Title" -n "Notes" -d "tomorrow"     # All fields
gtasks tasks add -l "Work" -t "Title"                    # Specify list
```

## Complete Tasks

```bash
gtasks tasks done                   # Interactive selection
gtasks tasks done 1                 # Complete task #1
gtasks tasks done 3 -l "Work"       # Complete task #3 in Work list
```

## Delete Tasks

```bash
gtasks tasks rm                     # Interactive selection
gtasks tasks rm 2                   # Delete task #2
gtasks tasks rm 1 -l "Personal"     # Delete task #1 in Personal list
```

## Task Details

```bash
gtasks tasks info                   # Interactive selection
gtasks tasks info 1                 # Show details for task #1
gtasks tasks info 2 -l "Work"       # Show details for task #2 in Work list
```

## Date Format Examples

All these work with the `-d` flag:

```
2024-12-25          # ISO format
Dec 25, 2024        # Month day, year
December 25         # Month day (current year)
12/25/2024          # US format
tomorrow            # Relative day
next Friday         # Relative named day
in 3 days           # Relative duration
```

## Common Workflows

### Add task with deadline
```bash
gtasks tasks add -l "Work" -t "Submit proposal" -d "next Friday"
```

### Check today's tasks
```bash
gtasks tasks view -l "Work" --sort=due
```

### Complete multiple tasks
```bash
gtasks tasks done -l "Work"    # Select first task
gtasks tasks done -l "Work"    # Select next task
```

### Export tasks
```bash
gtasks tasks view --format=json > tasks.json
gtasks tasks view --format=csv > tasks.csv
```

### Create shopping list
```bash
gtasks tasklists add -t "Shopping"
gtasks tasks add -l "Shopping" -t "Milk"
gtasks tasks add -l "Shopping" -t "Bread"
gtasks tasks add -l "Shopping" -t "Eggs"
```

## Flags Reference

### Global Flags

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--tasklist` | `-l` | Specify task list by name |

### View Tasks Flags

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--include-completed` | `-i` | Include completed tasks |
| `--completed` | | Show only completed tasks |
| `--sort` | | Sort by: due, title, position |
| `--format` | | Output format: table, json, csv |

### Add Task Flags

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--title` | `-t` | Task title (required in flag mode) |
| `--note` | `-n` | Task notes/description |
| `--due` | `-d` | Due date (flexible format) |

### Add Task List Flags

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--title` | `-t` | Task list title (required) |

## Output Formats

### Table (Default)
```
No  Title              Description         Status     Due
1   Finish report      Q4 analysis         pending    25 December 2024
2   Team meeting       Weekly sync         pending    -
```

### JSON
```json
[
  {
    "number": 1,
    "title": "Finish report",
    "description": "Q4 analysis",
    "status": "pending",
    "due": "2024-12-25"
  }
]
```

### CSV
```
No,Title,Description,Status,Due
1,Finish report,Q4 analysis,pending,25 December 2024
2,Team meeting,Weekly sync,pending,-
```

## Task Numbering

- Tasks are numbered starting from 1
- Task numbers are shown in view/list output
- Task numbers can change when tasks are added/deleted/sorted
- Always view the list first to get current numbers

## Status Values

- `pending` - Task needs action (not completed)
- `completed` - Task is marked as done

## Error Messages

| Error | Cause | Solution |
|-------|-------|----------|
| "command not found: gtasks" | GTasks not installed | Download from [releases](https://github.com/BRO3886/gtasks/releases) and add to PATH |
| "Failed to get service" | Not authenticated or missing env vars | Check env vars, then run `gtasks login` |
| Missing GTASKS_CLIENT_ID/SECRET | Environment variables not set | Export GTASKS_CLIENT_ID and GTASKS_CLIENT_SECRET |
| "incorrect task-list name" | List doesn't exist | Check with `gtasks tasklists view` |
| "Incorrect task number" | Invalid task number | Run `gtasks tasks view` to see valid numbers |
| "Date format incorrect" | Unparseable date | Use format like "2024-12-25" or "tomorrow" |

## Tips

1. **Task numbers change** - Always view list before using numbers
2. **Use -l flag** - Avoid interactive prompts in scripts
3. **Flexible dates** - Natural language works: "tomorrow", "next week"
4. **JSON for parsing** - Use `--format=json` when processing with jq
5. **Include completed** - Use `-i` to see full task history

## Getting Help

```bash
gtasks --help                  # General help
gtasks tasks --help            # Tasks command help
gtasks tasklists --help        # Task lists command help
```
