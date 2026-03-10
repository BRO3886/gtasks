# Advanced Usage and Integration

This reference covers advanced usage patterns, scripting, and integration scenarios for gtasks CLI.

## Scripting and Automation

### Batch Task Creation from File

Create multiple tasks from a text file:

```bash
# tasks.txt format:
# Task title | Notes | Due date

while IFS='|' read -r title notes due; do
  gtasks tasks add -l "Work" -t "$title" -n "$notes" -d "$due"
done < tasks.txt
```

### Task Completion Script

Mark all tasks containing a keyword as complete:

```bash
#!/bin/bash
KEYWORD="meeting"
LIST="Work"

# Get tasks as JSON
TASKS=$(gtasks tasks view -l "$LIST" --format=json)

# Find matching task numbers and mark complete
echo "$TASKS" | jq -r ".[] | select(.title | contains(\"$KEYWORD\")) | .number" | while read num; do
  gtasks tasks done "$num" -l "$LIST"
done
```

## Output Parsing

### JSON Processing with jq

Extract specific information:

```bash
# Get all overdue tasks
gtasks tasks view --format=json -l "Work" | jq '[.[] | select(.due != "" and .status == "pending") | select(.due < (now | strftime("%Y-%m-%d")))]'

# Count pending vs completed
gtasks tasks view -i --format=json -l "Work" | jq 'group_by(.status) | map({status: .[0].status, count: length})'

# Get tasks due this week
gtasks tasks view --format=json -l "Work" --sort=due | jq '[.[] | select(.due != "") | select(.due <= (now + 604800 | strftime("%Y-%m-%d")))]'
```

### CSV Processing

Import to spreadsheet tools:

```bash
# Export all lists to separate CSV files
for list in $(gtasks tasklists view | grep -oP '\[\d+\] \K.*'); do
  filename="${list// /_}.csv"
  gtasks tasks view -l "$list" --format=csv > "$filename"
done
```

## Integration Patterns

### Calendar Integration

Create calendar events for tasks with due dates:

```bash
# Pseudocode - integrate with calendar tool
gtasks tasks view --format=json -l "Work" | jq -r '.[] | select(.due != "") | "\(.title),\(.due),\(.description)"' | while IFS=',' read -r title date desc; do
  # Add to calendar using your calendar CLI tool
  # cal add "$title" "$date" --description "$desc"
done
```

### Notification System

Set up notifications for tasks due soon:

```bash
#!/bin/bash
# Add to crontab: 0 9 * * * /path/to/task-reminder.sh

TOMORROW=$(date -d '+1 day' +%Y-%m-%d)

gtasks tasks view --format=json -l "Work" | jq -r ".[] | select(.due == \"$TOMORROW\") | .title" | while read task; do
  # Send notification
  notify-send "Task Due Tomorrow" "$task"
done
```

### Sync with Other Task Managers

Export and import pattern:

```bash
# Export from gtasks
gtasks tasks view --format=json -l "Work" > gtasks_export.json

# Transform data (example for generic task manager)
jq '[.[] | {
  title: .title,
  description: .description,
  due_date: .due,
  completed: (.status == "completed")
}]' gtasks_export.json > transformed.json
```

## Performance Optimization

### Caching Task Lists

For scripts that run frequently, cache task list names:

```bash
# Cache task lists for 1 hour
CACHE_FILE="/tmp/gtasks_lists_cache"
CACHE_DURATION=3600

if [ ! -f "$CACHE_FILE" ] || [ $(( $(date +%s) - $(stat -f %m "$CACHE_FILE") )) -gt $CACHE_DURATION ]; then
  gtasks tasklists view > "$CACHE_FILE"
fi

cat "$CACHE_FILE"
```

### Parallel Operations

Process multiple lists concurrently:

```bash
#!/bin/bash
# Process each list in parallel
gtasks tasklists view | grep -oP '\[\d+\] \K.*' | xargs -P 4 -I {} bash -c 'gtasks tasks view -l "{}" --format=json > "$(echo {} | tr " " "_").json"'
```

## API Rate Limiting

Google Tasks API has rate limits. For bulk operations:

1. **Add delays between requests:**
```bash
for task in task1 task2 task3; do
  gtasks tasks add -l "Work" -t "$task"
  sleep 1  # 1 second delay
done
```

2. **Batch operations where possible:**
```bash
# Instead of multiple view commands, store the result
TASKS=$(gtasks tasks view -l "Work" --format=json)
# Then process locally using jq
```

## Error Handling in Scripts

Robust error handling:

```bash
#!/bin/bash
set -e  # Exit on error

# Check if authenticated
if ! gtasks tasklists view &> /dev/null; then
  echo "Error: Not authenticated. Run 'gtasks login'" >&2
  exit 1
fi

# Check if task list exists
if ! gtasks tasks view -l "Work" &> /dev/null; then
  echo "Error: Task list 'Work' not found" >&2
  exit 1
fi

# Proceed with operations
gtasks tasks add -l "Work" -t "New Task" || {
  echo "Error: Failed to create task" >&2
  exit 1
}
```

## Cross-Platform Considerations

### Date Commands

Different platforms have different `date` command syntax:

```bash
# Linux
TOMORROW=$(date -d '+1 day' +%Y-%m-%d)

# macOS
TOMORROW=$(date -v+1d +%Y-%m-%d)

# Cross-platform using gtasks native date parsing
gtasks tasks add -l "Work" -t "Task" -d "tomorrow"  # Works everywhere
```

### File Paths

Use HOME instead of ~ in scripts:

```bash
# Good
TOKEN_FILE="$HOME/.gtasks/token.json"

# Avoid
TOKEN_FILE="~/.gtasks/token.json"  # May not expand in all contexts
```

## Debugging

Enable verbose output for troubleshooting:

```bash
# Check if token file exists
ls -la ~/.gtasks/token.json

# Verify token is valid (login again if expired)
gtasks login

# Test API connectivity
gtasks tasklists view

# Check command output
gtasks tasks view -l "Work" --format=json | jq '.'
```

## Security Considerations

1. **Protect token file:**
```bash
chmod 600 ~/.gtasks/token.json
```

2. **Don't commit tokens to version control:**
```bash
echo ".gtasks/" >> .gitignore
```

3. **Use environment variables for sensitive data:**
```bash
# If scripting login automation (not recommended)
export GTASKS_CLIENT_ID="..."
export GTASKS_CLIENT_SECRET="..."
```

4. **Regularly rotate credentials:**
```bash
gtasks logout
gtasks login
```

## Monitoring and Logging

Add logging to automation scripts:

```bash
#!/bin/bash
LOG_FILE="$HOME/.gtasks/automation.log"

log() {
  echo "[$(date '+%Y-%m-%d %H:%M:%S')] $*" | tee -a "$LOG_FILE"
}

log "Starting task sync"
gtasks tasks view -l "Work" > /dev/null && log "Success" || log "Failed"
```

## Testing

Test script template:

```bash
#!/bin/bash
# test_gtasks.sh

test_auth() {
  gtasks tasklists view &> /dev/null
  if [ $? -eq 0 ]; then
    echo "✓ Authentication working"
    return 0
  else
    echo "✗ Authentication failed"
    return 1
  fi
}

test_list_creation() {
  TEST_LIST="Test_$(date +%s)"
  gtasks tasklists add -t "$TEST_LIST"

  if gtasks tasks view -l "$TEST_LIST" &> /dev/null; then
    echo "✓ List creation working"
    gtasks tasklists rm <<< "1"  # Cleanup
    return 0
  else
    echo "✗ List creation failed"
    return 1
  fi
}

# Run tests
test_auth
test_list_creation
```
