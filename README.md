# Google Tasks CLI

`gtasks`: A CLI Tool for Google Tasks

![gtasks](https://raw.githubusercontent.com/BRO3886/google-tasks-cli/master/.github/gtasks.png?token=AJQCPITXNHRYONWR4WB3RZC7WMHIY)

## Currently available commands

- [x] Login
- [x] View Task-List
- [x] Create Task-List
- [ ] Update Task-List
- [ ] Delete Task-List
- [x] View Tasks
- [x] Create Tasks
- [ ] Edit Task
- [ ] Mark as completed
- [ ] Delete Task


## Instructions to Run:
  - Pre-requisites
    - Go
  - Directions to install
  ```bash
  git clone https://github.com/BRO3886/google-tasks-cli
  ```
  - Directions to execute
  ```bash
  go run .
  ```
  Or, you can check out the pre-compiled binaries under **Releases**
  - Usage
```
Usage:
  gtasks [command]

Available Commands:
  help        Help about any command
  login       Logging into Google Tasks
  tasklists   View and create tasklists for currently signed-in account
  tasks       View, create, list and delete tasks in a tasklist

Flags:
  -h, --help     help for gtasks
  -t, --toggle   Help message for toggle

Use "gtasks [command] --help" for more information about a command.
```

<div align="center">
Made with :coffee: & <a href="https://cobra.dev">Cobra</a>
</div>
