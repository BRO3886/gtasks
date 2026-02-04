#!/bin/bash
# daily-report.sh - Generate a daily task report across all task lists
#
# Usage: ./daily-report.sh [--format=table|json|csv]
#
# Generates a comprehensive report of all tasks across all lists,
# highlighting tasks due today and overdue tasks.

set -e

FORMAT="${1:-table}"
FORMAT="${FORMAT#--format=}"

echo "==================================="
echo "  Google Tasks Daily Report"
echo "  $(date '+%A, %B %d, %Y')"
echo "==================================="
echo

# Check authentication
if ! gtasks tasklists view &> /dev/null; then
  echo "Error: Not authenticated. Please run 'gtasks login'" >&2
  exit 1
fi

# Get all task lists
LISTS=$(gtasks tasklists view | sed -n 's/^\[\([0-9]*\)\] \(.*\)/\2/p')

if [ -z "$LISTS" ]; then
  echo "No task lists found."
  exit 0
fi

TODAY=$(date +%Y-%m-%d)

# Process each list
while IFS= read -r list; do
  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
  echo "📋 $list"
  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

  # Get tasks for this list
  TASKS=$(gtasks tasks view -l "$list" --format=json 2>/dev/null || echo "[]")

  # Count tasks
  TOTAL=$(echo "$TASKS" | jq '. | length')
  PENDING=$(echo "$TASKS" | jq '[.[] | select(.status == "pending")] | length')
  COMPLETED=$(echo "$TASKS" | jq '[.[] | select(.status == "completed")] | length')

  if [ "$TOTAL" -eq 0 ]; then
    echo "   No tasks"
    echo
    continue
  fi

  echo "   Total: $TOTAL | Pending: $PENDING | Completed: $COMPLETED"
  echo

  # Show pending tasks sorted by due date
  if [ "$PENDING" -gt 0 ]; then
    echo "   📌 Pending Tasks:"
    echo "$TASKS" | jq -r '.[] | select(.status == "pending") | "      • \(.title)\(.due | if . != "" then " (Due: \(.))" else "" end)"'
    echo
  fi

  # Highlight overdue tasks
  if [ "$FORMAT" = "table" ]; then
    OVERDUE=$(echo "$TASKS" | jq --arg today "$TODAY" '[.[] | select(.status == "pending" and .due != "" and .due < $today)] | length')
    if [ "$OVERDUE" -gt 0 ]; then
      echo "   ⚠️  OVERDUE TASKS: $OVERDUE"
      echo "$TASKS" | jq -r --arg today "$TODAY" '.[] | select(.status == "pending" and .due != "" and .due < $today) | "      • \(.title) (Due: \(.due))"'
      echo
    fi

    # Highlight tasks due today
    DUE_TODAY=$(echo "$TASKS" | jq --arg today "$TODAY" '[.[] | select(.status == "pending" and .due == $today)] | length')
    if [ "$DUE_TODAY" -gt 0 ]; then
      echo "   🎯 DUE TODAY: $DUE_TODAY"
      echo "$TASKS" | jq -r --arg today "$TODAY" '.[] | select(.status == "pending" and .due == $today) | "      • \(.title)"'
      echo
    fi
  fi

done <<< "$LISTS"

echo "==================================="
echo "  Summary"
echo "==================================="

# Calculate totals across all lists
ALL_TASKS=$(gtasks tasklists view | sed -n 's/^\[\([0-9]*\)\] \(.*\)/\2/p' | while IFS= read -r list; do
  gtasks tasks view -l "$list" --format=json 2>/dev/null || echo "[]"
done | jq -s 'add')

TOTAL_ALL=$(echo "$ALL_TASKS" | jq '. | length')
PENDING_ALL=$(echo "$ALL_TASKS" | jq '[.[] | select(.status == "pending")] | length')
COMPLETED_ALL=$(echo "$ALL_TASKS" | jq '[.[] | select(.status == "completed")] | length')
OVERDUE_ALL=$(echo "$ALL_TASKS" | jq --arg today "$TODAY" '[.[] | select(.status == "pending" and .due != "" and .due < $today)] | length')
DUE_TODAY_ALL=$(echo "$ALL_TASKS" | jq --arg today "$TODAY" '[.[] | select(.status == "pending" and .due == $today)] | length')

echo "Total Tasks: $TOTAL_ALL"
echo "Pending: $PENDING_ALL"
echo "Completed: $COMPLETED_ALL"
echo "Overdue: $OVERDUE_ALL"
echo "Due Today: $DUE_TODAY_ALL"
echo

if [ "$DUE_TODAY_ALL" -gt 0 ] || [ "$OVERDUE_ALL" -gt 0 ]; then
  echo "⚡ Action needed! You have tasks that require attention."
else
  echo "✅ You're all caught up for today!"
fi
