---
title: "Task commands"
draft: false
weight: 4
summary: View, create, and delete tasks in a tasklist
---

## Help command

- to view inline help for all the commands

```
❯ gtasks tasks help

        View, create, list and delete tasks in a tasklist
        for the currently signed in account.
        Usage:
        [WITH LIST FLAG]
        gtasks tasks -l "<task-list name>" view|add|rm|done|undo|clear|info|update

        [WITHOUT LIST FLAG]
        gtasks tasks view|add|rm|done|undo|clear|info|update
        * You would be prompted to select a tasklist

Usage:
  gtasks tasks [command]

Available Commands:
  add         Add task in a tasklist
  clear       Hide all completed tasks from the list
  done        Mark tasks as done
  info        View detailed information about a task
  rm          Delete a task in a tasklist
  undo        Mark a completed task as incomplete
  update      Update an existing task
  view        View tasks in a tasklist

Flags:
  -h, --help              help for tasks
  -l, --tasklist string   use this flag to specify a tasklist

Use "gtasks tasks [command] --help" for more information about a command.
```

## Add Task

- First select the tasklist

```
❯ gtasks tasks add
Use the arrow keys to navigate: ↓ ↑ → ←
? Select Tasklist:
  ▸ DSC VIT
    Daily todo
    Life
    Placement todo
↓   To watch
```

- Then add task

```
❯ gtasks tasks add
✔ DSC VIT
Creating task in DSC VIT
Title: testing
Note: testing
Due Date: 12 July 2021
```

- For a shorthand syntax use:

```
gtasks tasks add -l "DSC VIT" --title <some title> [--note <some note> | --due <some due date>]
```

### Recurring Tasks

Create multiple tasks with a repeating schedule using the `--repeat` flag:

```
❯ gtasks tasks add -l "DSC VIT" -t "Daily standup" -d "2025-02-10" --repeat daily --repeat-count 5
Creating task in DSC VIT
Creating 5 recurring tasks...
Created 5 tasks
```

This creates 5 tasks for Feb 10, 11, 12, 13, 14.

Available repeat patterns:
- `daily` or `day`
- `weekly` or `week`
- `monthly` or `month`
- `yearly` or `year`

You can use `--repeat-count` to specify the number of occurrences:

```
gtasks tasks add -t "Weekly sync" -d "2025-02-10" --repeat weekly --repeat-count 4
```

Or use `--repeat-until` to specify an end date:

```
gtasks tasks add -t "Weekly sync" -d "2025-02-10" --repeat weekly --repeat-until "2025-03-10"
```

Both can be combined - the command stops at whichever limit is reached first.

## View all tasks in a tasklist

- First select tasklist

```
❯ gtasks tasks view
Use the arrow keys to navigate: ↓ ↑ → ←
? Select Tasklist:
  ▸ DSC VIT
    Daily todo
    Life
    Placement todo
↓   To watch
```

- Then you'll be able to see tasks in a tabular format

```
❯ gtasks tasks view
✔ DSC VIT
Tasks in DSC VIT:
| NO |        TITLE         |          DESCRIPTION           | STATUS |     DUE      |
|----|----------------------|--------------------------------|--------|--------------|
|  1 | testing              | testing                        | ✖      | 12 July 2021 |
|  2 | HopeHouse            | Checkout the app. Maybe        | ✖      | 06 July 2021 |
|    |                      | migrate to Flutter 2.0         |        |              |
|  3 | Vitty App Publishing | Get Appbundle for publishing   | ✖      | 07 July 2021 |
|    |                      | Vitty                          |        |              |
|  4 | Cadence              | App status - Yajat             | ✖      | 07 July 2021 |
|  5 | Keats android        | Take update on webview from    | ✖      | 11 July 2021 |
|    |                      | hishaam                        |        |              |
|  6 | Keats ios            | Check up on the apple dev      | ✖      | 08 July 2021 |
|    |                      | account status - Swamita       |        |              |
```

- Output formats (table, json, csv)

Use `--format` to change the output format. The default is `table`.

```
❯ gtasks tasks view --format table

❯ gtasks tasks view --format json

❯ gtasks tasks view --format csv
```

JSON example (pipe to `jq`):

```
❯ gtasks tasks view -l "DSC VIT" --format json | jq '.[] | {title, status, due}'
```

CSV example (redirect to a file):

```
❯ gtasks tasks view -l "DSC VIT" --format csv > tasks.csv
```

- To include completed tasks:

```
❯ gtasks tasks view --include-completed

❯ gtasks tasks -l "DSC VIT" view -i
```

Example:

```
❯ gtasks tasks -l "DSC VIT" view -i
Tasks in DSC VIT:
| NO |          TITLE           |          DESCRIPTION           | STATUS |       DUE        |
|----|--------------------------|--------------------------------|--------|------------------|
|  1 | testing                  | testing                        | ✖      | 12 July 2021     |
|  2 | Gidget fixes             | Push updated appbundle to play | ✔      | 04 July 2021     |
|    |                          | store                          |        |                  |
|  3 | Gidget fixes             | take new aab from Rishav       | ✔      | 06 July 2021     |
|  4 | HopeHouse                | Checkout the app. Maybe        | ✖      | 06 July 2021     |
|    |                          | migrate to Flutter 2.0         |        |                  |
|  5 | Vitty App Publishing     | Get Appbundle for publishing   | ✖      | 07 July 2021     |
|    |                          | Vitty                          |        |                  |
|  6 | Cadence                  | App status - Yajat             | ✖      | 07 July 2021     |
|  7 | Keats android            | Take update on webview from    | ✖      | 11 July 2021     |
|    |                          | hishaam                        |        |                  |
|  8 | Keats ios                | Check up on the apple dev      | ✖      | 08 July 2021     |
|    |                          | account status - Swamita       |        |                  |
|  9 | Testing                  | Something testing ono          | ✔      | 12 July 2021     |
| 10 | asjla                    | sjasj                          | ✔      | 12 July 2021     |
| 11 | testing                  | testing 1 2 3                  | ✔      | No Due Date      |
| 12 | abdcd                    | ahfje                          | ✔      | 10 July 2021     |
```

- To show completed tasks:

```
❯ gtasks tasks view --completed

❯ gtasks tasks -l "DSC VIT" view --completed
```

- To change sort order (due date, title, position, defeault=position)

```
❯ gtasks tasks view --sort due

❯ gtasks tasks -l "DSC VIT" view --sort title
```

- To limit the number of results:

```
❯ gtasks tasks view --max 5

❯ gtasks tasks -l "DSC VIT" view --max 10
```

## Mark task as done

- With prompt:

```
❯ gtasks tasks done
✔ DSC VIT
Tasks in DSC VIT:
Use the arrow keys to navigate: ↓ ↑ → ←
? Select Task:
  ▸ testing
    HopeHouse
    Vitty App Publishing
    Cadence
↓   Keats android
```

- For a shorter syntax:

```
❯ gtasks tasks view -l "DSC VIT"
Tasks in DSC VIT:
| NO |        TITLE         |          DESCRIPTION           | STATUS |     DUE      |
|----|----------------------|--------------------------------|--------|--------------|
|  1 | testing              | testing                        | ✖      | 12 July 2021 |
|  2 | HopeHouse            | Checkout the app. Maybe        | ✖      | 06 July 2021 |
|    |                      | migrate to Flutter 2.0         |        |              |
|  3 | Vitty App Publishing | Get Appbundle for publishing   | ✖      | 07 July 2021 |
|    |                      | Vitty                          |        |              |
|  4 | Cadence              | App status - Yajat             | ✖      | 07 July 2021 |
|  5 | Keats android        | Take update on webview from    | ✖      | 11 July 2021 |
|    |                      | hishaam                        |        |              |
|  6 | Keats ios            | Check up on the apple dev      | ✖      | 08 July 2021 |
|    |                      | account status - Swamita       |        |              |

❯ gtasks tasks done -l "DSC VIT" 1
Marked as complete: testing
```

## Undo a completed task

Mark a completed task as incomplete again.

- With prompt:

```
❯ gtasks tasks undo
✔ DSC VIT
Tasks in DSC VIT:
Use the arrow keys to navigate: ↓ ↑ → ←
? Select Task:
  ▸ testing (completed)
    HopeHouse (completed)
```

- For a shorter syntax (first view completed tasks to get the number):

```
❯ gtasks tasks view -l "DSC VIT" --include-completed
Tasks in DSC VIT:
| NO |        TITLE         | STATUS    |
|----|----------------------|-----------|
|  1 | testing              | completed |
|  2 | HopeHouse            | completed |

❯ gtasks tasks undo -l "DSC VIT" 1
Marked as incomplete: testing
```

## Clear completed tasks

Hide all completed tasks from the list. This marks completed tasks as hidden so they won't be returned by the API (primarily affects tasks completed via the CLI).

```
❯ gtasks tasks clear -l "DSC VIT"
✔ Clear all completed tasks from 'DSC VIT'? [y/N]: y
Cleared completed tasks from DSC VIT
```

- Use `--force` or `-f` to skip the confirmation prompt:

```
❯ gtasks tasks clear -l "DSC VIT" --force
Cleared completed tasks from DSC VIT
```

## View detailed task information

The `info` command displays detailed information about a task, including links/URLs that may have been shared to Google Tasks (e.g., from Android's "Share With..." feature).

By default, `info` only considers pending tasks (matching `view` behavior). Use `-i` to include completed tasks.

- With prompt:

```
❯ gtasks tasks info
✔ DSC VIT
Tasks in DSC VIT:
Use the arrow keys to navigate: ↓ ↑ → ←
? Select Task:
  ▸ testing
    HopeHouse
    Vitty App Publishing
```

- For a shorter syntax using task number:

```
❯ gtasks tasks view -l "DSC VIT"
Tasks in DSC VIT:
| NO |        TITLE         |          DESCRIPTION           | STATUS |     DUE      |
|----|----------------------|--------------------------------|--------|--------------|
|  1 | testing              | testing                        | ✖      | 12 July 2021 |
|  2 | HopeHouse            | Checkout the app. Maybe        | ✖      | 06 July 2021 |

❯ gtasks tasks info -l "DSC VIT" 1

Task: testing
Status: Needs action
Due: 12 July 2021
Notes: testing

Links:
  - https://example.com/some-link

View in Google Tasks: https://tasks.google.com/...
```

- To get info on a completed task, use `-i` (must match how you viewed the list):

```
❯ gtasks tasks view -l "DSC VIT" -i
❯ gtasks tasks info -l "DSC VIT" 3 -i
```

The info command is particularly useful for viewing:

- Full task notes (not truncated)
- Links/URLs attached to the task
- WebViewLink to open the task in Google Tasks web interface
- Complete due date information
- Task completion status

## Update a task

Update an existing task's title, note, or due date.

### Interactive mode

When no flags are provided, you'll be prompted for each field with the current value displayed. Press Enter to keep the current value, or type a new value.

```
❯ gtasks tasks update 1
Updating task: testing

Title [testing]: new title
Note [testing notes]: 
Due [12 July 2021]: 

Updated: new title
```

### Flag mode

Use flags to update specific fields without prompts:

```
❯ gtasks tasks update 1 --title "New title"
Updating task: testing

Updated: New title

❯ gtasks tasks update 1 --note "Updated note" --due "tomorrow"
Updating task: New title

Updated: New title
```

Available flags:
- `-t, --title` - New title for the task
- `-n, --note` - New note for the task  
- `-d, --due` - New due date for the task

## Delete a task

- With prompt:

```
❯ gtasks tasks rm
✔ DSC VIT
Tasks in DSC VIT:
Use the arrow keys to navigate: ↓ ↑ → ←
? Select Task:
  ▸ testing
    HopeHouse
    Vitty App Publishing
    Cadence
↓   Keats android
```

- For a shorter syntax:

```
❯ gtasks tasks view -l "DSC VIT"
Tasks in DSC VIT:
| NO |        TITLE         |          DESCRIPTION           | STATUS |     DUE      |
|----|----------------------|--------------------------------|--------|--------------|
|  1 | testing              | testing                        | ✖      | 12 July 2021 |
|  2 | HopeHouse            | Checkout the app. Maybe        | ✖      | 06 July 2021 |
|    |                      | migrate to Flutter 2.0         |        |              |
|  3 | Vitty App Publishing | Get Appbundle for publishing   | ✖      | 07 July 2021 |
|    |                      | Vitty                          |        |              |
|  4 | Cadence              | App status - Yajat             | ✖      | 07 July 2021 |
|  5 | Keats android        | Take update on webview from    | ✖      | 11 July 2021 |
|    |                      | hishaam                        |        |              |
|  6 | Keats ios            | Check up on the apple dev      | ✖      | 08 July 2021 |
|    |                      | account status - Swamita       |        |              |

❯ gtasks tasks rm -l "DSC VIT" 1
Deleted: testing
```
