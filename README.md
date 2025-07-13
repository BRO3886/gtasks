# Google Tasks CLI

`gtasks`: A CLI Tool for Google Tasks

![gtasks image](docs/static/images/screenshot.png)

---

## Docs

Refer to the [docs website](https://gtasks.sidv.dev) to read about available commands.

## Instructions to install

1. Download the binary for your system (check [releases](https://github.com/BRO3886/gtasks/releases))
2. Move to a particular folder, for eg Documents
3. Append the absolute path (use `pwd`) of the folder to `PATH`
4. Execute `gtasks` from anywhere

## Instructions to install using go install

> Not working yet. Will be fixed soon.

```bash
go install github.com/BRO3886/gtasks@latest
```



## Instructions to Run and Build from Source:

### Prerequisites

- Go 1.24+
- Google Cloud Console OAuth2 credentials (see Configuration section)

### Setup

1. Clone the repository:

```bash
git clone https://github.com/BRO3886/gtasks
cd gtasks
```

2. Create a `.env` file with your OAuth2 credentials:

```bash
# .env
GTASKS_CLIENT_ID=your-client-id.apps.googleusercontent.com
GTASKS_CLIENT_SECRET=your-client-secret
```

### Build Commands

```bash
# Development build (uses .env file)
make dev

# Build for specific platforms
make linux    # Linux (amd64 + arm64)
make windows  # Windows (amd64)
make mac      # macOS (amd64 + arm64)

# Build for all platforms
make all

# Create release packages
make release
```

### Configuration

To use GTasks, you need to set up Google OAuth2 credentials:

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select existing one
3. Enable the Google Tasks API
4. Create OAuth2 credentials:

   - Application type: "Desktop application"
   - Add authorized redirect URIs:
     - `http://localhost:8080/callback`
     - `http://localhost:8081/callback`
     - `http://localhost:8082/callback`
     - `http://localhost:9090/callback`
     - `http://localhost:9091/callback`

5. Set environment variables:

```bash
export GTASKS_CLIENT_ID="your-client-id.apps.googleusercontent.com"
export GTASKS_CLIENT_SECRET="your-client-secret"
```

Or create a `.env` file (for building from source).

### Token Storage

GTasks stores authentication tokens in `~/.gtasks/token.json`. This directory is created automatically on first login.

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

## Commands

### Help

- To see details about a command

```bash
gtasks <COMMAND> help
```

### Auth

- Login

```bash
gtasks login
```

- Logout

```bash
gtasks logout
```

### Tasklists

- Viewing Tasklists

```bash
gtasks tasklists view
```

- Creating a Tasklist

```bash
gtasks tasklists add -t 'title'
gtasks tasklists add --title 'title'
```

- Deleting a Tasklist

```bash
gtasks tasklists rm
```

### Tasks

- To pre-select tasklist, provide it's title as follows:

```bash
gtasks tasks -l <title> subcommand [--subcommand-flags]
```

Examples:

```bash
gtasks tasks [--tasklist|-l] "DSC VIT" view [--include-completed | -i]
```

**Note:** If the `-l` flag is not provided you will be able to choose a tasklist from the prompt

- Viewing tasks

```bash
gtasks tasks view
```

- Include completed tasks

```bash
gtasks tasks view -i
gtasks tasks view --include-completed
```

- Sort options

```bash
gtasks tasks view ... --sort [due,title,position, default=position]
```

- Adding a task

```bash
gtasks tasks add
```

- Mark task as completed

```bash
gtasks tasks done
```

- Deleting a task

```bash
gtasks tasks rm
```

<div align="center">
Made with :coffee: & <a href="https://cobra.dev">Cobra</a>
</div>
