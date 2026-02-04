#!/bin/bash
# import-tasks.sh - Import tasks from a text file
#
# Usage: ./import-tasks.sh <tasklist_name> <input_file>
#
# Input file format (pipe-separated):
# Task Title|Notes|Due Date
#
# Example:
# Buy groceries|Milk, bread, eggs|2024-12-25
# Call dentist||tomorrow
# Finish report|Q4 analysis|next Friday

set -e

if [ $# -lt 2 ]; then
  echo "Usage: $0 <tasklist_name> <input_file>"
  echo
  echo "Input file format (pipe-separated):"
  echo "Task Title|Notes|Due Date"
  echo
  echo "Example:"
  echo "Buy groceries|Milk, bread, eggs|2024-12-25"
  echo "Call dentist||tomorrow"
  echo "Finish report|Q4 analysis|next Friday"
  exit 1
fi

TASKLIST="$1"
INPUT_FILE="$2"

if [ ! -f "$INPUT_FILE" ]; then
  echo "Error: File not found: $INPUT_FILE" >&2
  exit 1
fi

# Check authentication
if ! gtasks tasklists view &> /dev/null; then
  echo "Error: Not authenticated. Please run 'gtasks login'" >&2
  exit 1
fi

# Check if task list exists
if ! gtasks tasks view -l "$TASKLIST" &> /dev/null; then
  echo "Error: Task list '$TASKLIST' not found" >&2
  echo
  echo "Available task lists:"
  gtasks tasklists view
  exit 1
fi

echo "========================================"
echo "  Importing Tasks to: $TASKLIST"
echo "========================================"
echo

TOTAL=0
SUCCESS=0
FAILED=0

while IFS='|' read -r title notes due; do
  # Skip empty lines and comments
  if [ -z "$title" ] || [[ "$title" =~ ^[[:space:]]*# ]]; then
    continue
  fi

  TOTAL=$((TOTAL + 1))

  echo "[$TOTAL] Importing: $title"

  # Build command
  CMD="gtasks tasks add -l \"$TASKLIST\" -t \"$title\""

  if [ -n "$notes" ]; then
    CMD="$CMD -n \"$notes\""
  fi

  if [ -n "$due" ]; then
    CMD="$CMD -d \"$due\""
  fi

  # Execute command
  if eval $CMD 2>/dev/null; then
    SUCCESS=$((SUCCESS + 1))
    echo "   ✓ Success"
  else
    FAILED=$((FAILED + 1))
    echo "   ✗ Failed"
  fi

  # Small delay to respect API rate limits
  sleep 0.5

done < "$INPUT_FILE"

echo
echo "========================================"
echo "  Import Complete"
echo "========================================"
echo "Total: $TOTAL"
echo "Successful: $SUCCESS"
echo "Failed: $FAILED"
echo

if [ $FAILED -gt 0 ]; then
  echo "⚠️  Some imports failed. Check the output above for details."
  exit 1
fi

echo "✅ All tasks imported successfully!"
