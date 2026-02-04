---
name: gtasks-cli
description: Manage Google Tasks from the command line - view, create, update, delete tasks and task lists. Use when the user asks to interact with Google Tasks, manage to-do items, create task lists, mark tasks complete, or check their Google Tasks.
license: MIT
compatibility: Requires gtasks CLI tool to be installed and authenticated
metadata:
  author: BRO3886
  version: "1.0"
allowed-tools: Bash(gtasks:*)
---

# Google Tasks CLI Skill

This skill enables you to manage Google Tasks directly from the command line using the `gtasks` CLI tool.

## Prerequisites

Before using any commands, ensure the following requirements are met:

### 1. GTasks Installation

Check if gtasks is installed on the system:

```bash
# Cross-platform check (works on macOS, Linux, Windows Git Bash)
gtasks --version 2>/dev/null || gtasks.exe --version 2>/dev/null || echo "gtasks not found"

# Or use which/where commands
# macOS/Linux:
which gtasks

# Windows (Command Prompt):
where gtasks

# Windows (PowerShell):
Get-Command gtasks
```

**If gtasks is not installed:**

1. Download the binary for your system from [GitHub Releases](https://github.com/BRO3886/gtasks/releases)
2. Install it:
   - **macOS/Linux**: Move to `/usr/local/bin` or add to PATH
   - **Windows**: Add to a folder in your PATH environment variable
3. Verify installation: `gtasks --version`

**IMPORTANT for Agents:** Always check if gtasks is installed before attempting to use it. If the command is not found, inform the user and provide installation instructions.

### 2. Environment Variables

Set up Google OAuth2 credentials as environment variables:

```bash
export GTASKS_CLIENT_ID="your-client-id.apps.googleusercontent.com"
export GTASKS_CLIENT_SECRET="your-client-secret"
```

**How to get credentials:**
1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select an existing one
3. Enable the Google Tasks API
4. Create OAuth2 credentials (Application type: "Web application")
5. Add authorized redirect URIs:
   - `http://localhost:8080/callback`
   - `http://localhost:8081/callback`
   - `http://localhost:8082/callback`
   - `http://localhost:9090/callback`
   - `http://localhost:9091/callback`

**For persistent setup**, add these to your shell profile (`~/.bashrc`, `~/.zshrc`, etc.):

```bash
echo 'export GTASKS_CLIENT_ID="your-client-id"' >> ~/.bashrc
echo 'export GTASKS_CLIENT_SECRET="your-client-secret"' >> ~/.bashrc
source ~/.bashrc
```

### 2. Authentication

Once environment variables are set, authenticate with Google:

```bash
gtasks login
```

This will open a browser for OAuth2 authentication. The token is stored in `~/.gtasks/token.json`.

## Core Concepts

- **Task Lists**: Containers that hold tasks (like "Work", "Personal", "Shopping")
- **Tasks**: Individual to-do items within a task list
- **Task Properties**: Title (required), notes/description (optional), due date (optional), status (pending/completed)

## Command Structure

All commands follow this pattern:
```
gtasks [command] [subcommand] [flags] [arguments]
```

## Authentication

### Login
```bash
gtasks login
```
Opens browser for Google OAuth2 authentication. Required before using any other commands.

### Logout
```bash
gtasks logout
```
Removes stored credentials from `~/.gtasks/token.json`.

## Task List Management

### View All Task Lists
```bash
gtasks tasklists view
```
Displays all task lists with numbered indices.

**Output Example:**
```
[1] My Tasks
[2] Work
[3] Personal
```

### Create a Task List
```bash
gtasks tasklists add -t "Work Projects"
gtasks tasklists add --title "Shopping List"
```
Creates a new task list with the specified title.

**Flags:**
- `-t, --title`: Task list title (required)

### Delete a Task List
```bash
gtasks tasklists rm
```
Interactive prompt to select and delete a task list.

### Update Task List Title
```bash
gtasks tasklists update -t "New Title"
```
Interactive prompt to select a task list and update its title.

**Flags:**
- `-t, --title`: New title for the task list (required)

## Task Management

All task commands can optionally specify a task list using the `-l` flag. If omitted, you'll be prompted to select one interactively.

### View Tasks

**Basic view:**
```bash
gtasks tasks view
gtasks tasks view -l "Work"
```

**Include completed tasks:**
```bash
gtasks tasks view --include-completed
gtasks tasks view -i
```

**Show only completed tasks:**
```bash
gtasks tasks view --completed
```

**Sort tasks:**
```bash
gtasks tasks view --sort=due        # Sort by due date
gtasks tasks view --sort=title      # Sort by title
gtasks tasks view --sort=position   # Sort by position (default)
```

**Output formats:**
```bash
gtasks tasks view --format=table    # Table format (default)
gtasks tasks view --format=json     # JSON output
gtasks tasks view --format=csv      # CSV output
```

**Table Output Example:**
```
Tasks in Work:
No  Title              Description         Status     Due
1   Finish report      Q4 analysis         pending    25 December 2024
2   Team meeting       Weekly sync         pending    -
3   Code review        PR #123            completed  20 December 2024
```

**JSON Output Example:**
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

### Create a Task

**Interactive mode:**
```bash
gtasks tasks add
gtasks tasks add -l "Work"
```
Prompts for title, notes, and due date.

**Flag mode:**
```bash
gtasks tasks add -t "Buy groceries"
gtasks tasks add -t "Finish report" -n "Q4 analysis" -d "2024-12-25"
gtasks tasks add -t "Call dentist" -d "tomorrow"
gtasks tasks add -t "Team meeting" -d "Dec 25"
```

**Flags:**
- `-t, --title`: Task title (required for non-interactive mode)
- `-n, --note`: Task notes/description (optional)
- `-d, --due`: Due date (optional, flexible format)

**Date Format Examples:**
The date parser supports many formats:
- `2024-12-25` (ISO format)
- `Dec 25, 2024`
- `December 25`
- `tomorrow`
- `next Friday`
- `12/25/2024`

See [dateparse examples](https://github.com/araddon/dateparse#extended-example) for all supported formats.

### Mark Task as Complete

**With task number:**
```bash
gtasks tasks done 1
gtasks tasks done 3 -l "Work"
```

**Interactive mode:**
```bash
gtasks tasks done
gtasks tasks done -l "Personal"
```
Prompts to select a task from the list.

### Delete a Task

**With task number:**
```bash
gtasks tasks rm 2
gtasks tasks rm 1 -l "Shopping"
```

**Interactive mode:**
```bash
gtasks tasks rm
gtasks tasks rm -l "Work"
```
Prompts to select a task to delete.

### View Task Details

**With task number:**
```bash
gtasks tasks info 1
gtasks tasks info 3 -l "Work"
```

**Interactive mode:**
```bash
gtasks tasks info
gtasks tasks info -l "Personal"
```

**Output Example:**
```
Task: Finish report
Status: Needs action
Due: 25 December 2024
Notes: Complete Q4 analysis and submit to manager

Links:
  - https://docs.google.com/document/d/...

View in Google Tasks: https://tasks.google.com/...
```

## Common Workflows

### Quick Task Creation
When a user says "add a task to my work list":
```bash
gtasks tasks add -l "Work" -t "Task title"
```

### Check Today's Tasks
```bash
gtasks tasks view --sort=due
```

### Complete Multiple Tasks
```bash
gtasks tasks done -l "Work"
# Interactive prompt appears, select task
gtasks tasks done -l "Work"
# Repeat as needed
```

### View All Tasks Across Lists
Run view command multiple times for each list, or first list all task lists:
```bash
gtasks tasklists view
gtasks tasks view -l "Work"
gtasks tasks view -l "Personal"
```

### Export Tasks
```bash
gtasks tasks view --format=json > tasks.json
gtasks tasks view --format=csv > tasks.csv
```

## Best Practices

1. **Always check authentication first**: If commands fail with authentication errors, run `gtasks login`

2. **Use task list flag for automation**: When scripting or when the user specifies a list name, use `-l` flag to avoid interactive prompts

3. **Leverage flexible date parsing**: The `--due` flag accepts natural language dates like "tomorrow", "next week", etc.

4. **Use appropriate output format**:
   - Table format for human-readable output
   - JSON for parsing/integration with other tools
   - CSV for spreadsheet import

5. **Task numbers are ephemeral**: Task numbers change when tasks are added, completed, or deleted. Always view the list first to get current numbers.

6. **Handle missing lists gracefully**: If a user specifies a non-existent list name, the command will error. Always verify list names first with `gtasks tasklists view`.

## Error Handling

Common errors and solutions:

- **"Failed to get service"** or **Authentication errors**:
  - First, ensure environment variables are set: `echo $GTASKS_CLIENT_ID`
  - If variables are not set, export them (see Prerequisites section)
  - Then run `gtasks login` to authenticate
- **"incorrect task-list name"**: The specified list name doesn't exist. Use `gtasks tasklists view` to see available lists
- **"Incorrect task number"**: The task number is invalid. Use `gtasks tasks view` to see current task numbers
- **"Date format incorrect"**: The date string couldn't be parsed. Use formats like "2024-12-25", "tomorrow", or "Dec 25"

## Examples

### Example 1: Create a shopping list and add items
```bash
gtasks tasklists add -t "Shopping"
gtasks tasks add -l "Shopping" -t "Milk"
gtasks tasks add -l "Shopping" -t "Bread"
gtasks tasks add -l "Shopping" -t "Eggs"
```

### Example 2: Review and complete work tasks
```bash
gtasks tasks view -l "Work" --sort=due
gtasks tasks done 1 -l "Work"
```

### Example 3: Add task with deadline
```bash
gtasks tasks add -l "Work" -t "Submit proposal" -n "Include budget and timeline" -d "next Friday"
```

### Example 4: Export completed tasks
```bash
gtasks tasks view --completed --format=json -l "Work" > completed_work.json
```

## Tips for Agents

### Before Running Any Commands

1. **Check gtasks installation first**:
   ```bash
   # Try to run gtasks version check
   gtasks --version 2>/dev/null || gtasks.exe --version 2>/dev/null
   ```
   If this fails, inform the user that gtasks is not installed and provide installation instructions from the Prerequisites section.

2. **Verify environment variables are set**:
   ```bash
   # Check if variables exist (macOS/Linux)
   [ -n "$GTASKS_CLIENT_ID" ] && echo "GTASKS_CLIENT_ID is set" || echo "GTASKS_CLIENT_ID is not set"
   [ -n "$GTASKS_CLIENT_SECRET" ] && echo "GTASKS_CLIENT_SECRET is set" || echo "GTASKS_CLIENT_SECRET is not set"

   # Windows PowerShell
   if ($env:GTASKS_CLIENT_ID) { "GTASKS_CLIENT_ID is set" } else { "GTASKS_CLIENT_ID is not set" }
   if ($env:GTASKS_CLIENT_SECRET) { "GTASKS_CLIENT_SECRET is set" } else { "GTASKS_CLIENT_SECRET is not set" }
   ```

3. **Check authentication status**:
   ```bash
   gtasks tasklists view &>/dev/null && echo "Authenticated" || echo "Not authenticated - run 'gtasks login'"
   ```

### General Tips

- When the user mentions "tasks" without specifying a tool, ask if they want to use Google Tasks
- If the user asks about their tasks, first run `gtasks tasklists view` to see available lists
- Always confirm which task list to use if not specified by the user
- When creating tasks with dates, prefer explicit date formats (YYYY-MM-DD) over relative terms for clarity
- Remember that task numbers are 1-indexed and change after modifications
- If a command requires interaction but you're running non-interactively, use flags to provide all required information
