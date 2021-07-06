---
title: 'Task commands'
draft: false
weight: 4
summary: View, create, and delete tasks in a tasklist
---

## Help command
* to view inline help for all the commands
```
❯ gtasks tasks help

        View, create, list and delete tasks in a tasklist
        for the currently signed in account.
        Usage:
        [WITH LIST FLAG] 
        gtasks tasks -l "<task-list name>" view|add|rm|done

        [WITHOUT LIST FLAG]
        gtasks tasks view|add|rm|done
        * You would be prompted to select a tasklist

Usage:
  gtasks tasks [command]

Available Commands:
  add         Add task in a tasklist
  done        Mark tasks as done
  rm          Delete a task in a tasklist
  view        View tasks in a tasklist

Flags:
  -h, --help              help for tasks
  -l, --tasklist string   use this flag to specify a tasklist

Use "gtasks tasks [command] --help" for more information about a command.
```

## Add Task
* First select the tasklist
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
* Then add task
```
❯ gtasks tasks add
✔ DSC VIT
Creating task in DSC VIT
Title: testing
Note: testing
Due Date: 12 July 2021
```
* For a shorthand syntax use:
```
gtasks tasks add -l "DSC VIT"
```

## View all tasks in a tasklist
* First select tasklist
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
* Then you'll be able to see tasks in a tabular format
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
* To include completed tasks:
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

## Mark task as done
* With prompt:
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
* For a shorter syntax:
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

## Delete a task
* With prompt:
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
* For a shorter syntax:
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
