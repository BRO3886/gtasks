#!/bin/bash
# backup-tasks.sh - Backup all Google Tasks to JSON files
#
# Usage: ./backup-tasks.sh [backup_directory]
#
# Creates a timestamped backup of all task lists and their tasks
# in the specified directory (default: ~/gtasks_backups)

set -e

BACKUP_DIR="${1:-$HOME/gtasks_backups}"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_PATH="$BACKUP_DIR/$TIMESTAMP"

echo "========================================"
echo "  Google Tasks Backup"
echo "========================================"
echo

# Check authentication
if ! gtasks tasklists view &> /dev/null; then
  echo "Error: Not authenticated. Please run 'gtasks login'" >&2
  exit 1
fi

# Create backup directory
mkdir -p "$BACKUP_PATH"
echo "📁 Creating backup in: $BACKUP_PATH"
echo

# Get all task lists
LISTS=$(gtasks tasklists view 2>/dev/null)

if [ -z "$LISTS" ]; then
  echo "No task lists found."
  exit 0
fi

# Backup task lists metadata
echo "$LISTS" > "$BACKUP_PATH/tasklists.txt"
echo "✓ Saved task lists metadata"

# Counter for statistics
TOTAL_LISTS=0
TOTAL_TASKS=0

# Export each task list
echo "$LISTS" | sed -n 's/^\[\([0-9]*\)\] \(.*\)/\2/p' | while IFS= read -r list; do
  TOTAL_LISTS=$((TOTAL_LISTS + 1))

  # Sanitize filename
  filename=$(echo "$list" | tr ' /' '_' | tr -cd '[:alnum:]_-')

  echo "  📋 Backing up: $list"

  # Export tasks including completed ones
  TASKS=$(gtasks tasks view -l "$list" -i --format=json 2>/dev/null || echo "[]")

  # Save to file
  echo "$TASKS" > "$BACKUP_PATH/${filename}.json"

  # Count tasks
  TASK_COUNT=$(echo "$TASKS" | jq '. | length')
  TOTAL_TASKS=$((TOTAL_TASKS + TASK_COUNT))

  echo "     → $TASK_COUNT tasks saved"
done

echo
echo "========================================"
echo "  Backup Complete"
echo "========================================"
echo "Location: $BACKUP_PATH"
echo "Lists backed up: $TOTAL_LISTS"
echo "Total tasks: $TOTAL_TASKS"
echo

# Create a manifest file
cat > "$BACKUP_PATH/manifest.json" <<EOF
{
  "backup_date": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
  "backup_timestamp": "$TIMESTAMP",
  "total_lists": $TOTAL_LISTS,
  "total_tasks": $TOTAL_TASKS
}
EOF

echo "✓ Manifest created: $BACKUP_PATH/manifest.json"

# Create latest symlink
ln -sf "$TIMESTAMP" "$BACKUP_DIR/latest"
echo "✓ Latest backup linked"

echo
echo "To restore from this backup, use the restore-tasks.sh script"
echo "or manually import the JSON files."
