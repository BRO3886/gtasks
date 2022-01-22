---
title: "Tasklist commands"
draft: false
weight: 3
summary: View and create tasklists for currently signed-in account.
---

## Help command

- to view inline help for all the commands

```
❯ gtasks tasklists --help

        View and create tasklists for currently signed-in account

        View tasklists:
        gtasks tasklists view

        Create tasklist:
        gtasks tasklists create -t <TITLE>
        gtasks tasklists create --title <TITLE>

        Remove tasklist
        gtasks tasklists rm

Usage:
  gtasks tasklists [flags]
  gtasks tasklists [command]

Available Commands:
  create      create tasklist
  rm          remove tasklist
  update      update tasklist title
  view        view tasklists

Flags:
  -h, --help   help for tasklists

Use "gtasks tasklists [command] --help" for more information about a command.
```

## Create Tasklist

Examples:

```
❯ gtasks tasklists add --title "some title"

❯ gtasks tasklists add -t "some title"
```

## View all Tasklists

Example:

```
❯ gtasks tasklists view
```

## Update a tasklist title

Examples:

```
❯ gtasks tasklists update --title "some title"

❯ gtasks tasklists update  -t "some title"
```

## Delete a tasklist

Examples:

```
❯ gtasks tasklists rm
Use the arrow keys to navigate: ↓ ↑ → ←
? Select Tasklist:
  ▸ VIT
    Daily todo
    personal projects
    To watch
↓   DSC VIT

```
