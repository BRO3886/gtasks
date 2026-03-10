# GTasks - Google Tasks CLI Tool

## Project Overview

GTasks is a command-line interface tool for managing Google Tasks, written in Go. This project allows users to interact with their Google Tasks directly from the terminal, providing functionality to view, create, update, and delete tasks and task lists.

## Architecture & Structure

### Main Components

- **Main Entry Point**: `main.go` - Initializes configuration and executes CLI commands
- **CLI Commands**: `cmd/` directory contains all Cobra-based command implementations
- **API Layer**: `api/` directory handles Google Tasks API interactions
- **Configuration**: `internal/config/` manages OAuth2 credentials and app settings
- **Utilities**: `internal/utils/` provides logging and sorting functionality

### Key Dependencies

- **Cobra**: CLI framework for command structure and parsing
- **Google Tasks API**: Official Google API client for Go
- **OAuth2**: Authentication with Google services
- **PromptUI**: Interactive prompts for user input
- **TableWriter**: Formatted table output for task listings
- **koanf**: Layered config loading (env vars > config file > build-time defaults)
- **go-keyring**: System keyring storage for OAuth2 tokens (macOS Keychain, Linux Secret Service, Windows Credential Manager)

## Core Functionality

### Authentication (`api/auth.go`)

- OAuth2 flow implementation using PKCE + localhost callback
- Token stored in system keyring (go-keyring); falls back to `token.json` on headless systems
- Existing `token.json` is auto-migrated to keyring on first use
- Login/logout functionality
- Service client creation for API calls

### Task Management (`api/tasks.go`, `cmd/tasks.go`)

- **View Tasks**: Display tasks in table format with sorting options
- **Create Tasks**: Add new tasks with title, notes, and due dates
- **Update Tasks**: Mark tasks as completed
- **Delete Tasks**: Remove tasks from task lists
- **Sort Options**: By due date, title, or position

### Task List Management (`api/tasklists.go`, `cmd/tasklists.go`)

- View available task lists
- Create new task lists
- Delete task lists
- Interactive selection when not specified via flags

## Configuration & Setup

### OAuth2 Configuration

- Credentials supplied via env vars (`GTASKS_CLIENT_ID`, `GTASKS_CLIENT_SECRET`), config file, or build-time `-ldflags`
- Tokens stored in system keyring; falls back to `token.json` in config dir

### Configuration File

Supports `config.toml`, `config.yaml`, or `config.json` in the config directory:

```toml
[credentials]
client_id     = "..."
client_secret = "..."

[tasks]
default_task_list = "My Tasks"
```

### Installation Locations (XDG)

- New installs: `~/.config/gtasks/` (XDG Base Directory)
- Existing installs: `~/.gtasks/` (legacy, kept for backwards compatibility)
- If both exist, XDG wins and a migration warning is printed

## Command Structure

```
gtasks
├── login                    # Authenticate with Google
├── logout                   # Remove stored credentials
├── tasklists               # Task list operations
│   ├── view                # List all task lists
│   ├── add -t "title"      # Create new task list
│   └── rm                  # Delete task list
└── tasks [-l "list-name"]  # Task operations
    ├── view [--sort=due]   # View tasks (with sorting)
    ├── add [-t "title"]    # Create new task
    ├── done                # Mark task complete
    └── rm                  # Delete task
```

## Build System

### Makefile Targets

- **Cross-platform builds**: Windows, Linux, macOS (Intel & ARM)
- **Release packaging**: Automated tarball creation
- **GitHub releases**: Integration with `gh` CLI tool

### Build Commands

- `make linux` - Build for Linux
- `make windows` - Build for Windows
- `make mac` - Build for macOS
- `make all` - Build for all platforms
- `make release` - Build release archives

### Release Process
Steps in order — do not skip or reorder:
1. `git push` — push all commits to main **first**. Never tag unpushed commits.
2. `git tag vX.Y.Z` — tag after push so the tag points to a commit already on remote main
3. `git push origin vX.Y.Z` — push the tag explicitly
4. `make release` — builds all platform archives
5. `gh release create vX.Y.Z <archives>`

- **CRITICAL: Push before tag.** Tagging an unpushed commit then running `gh release create` pushes the tag + that commit but leaves `main` behind on remote — release binary is built from code not reachable from main.

## Development Guidelines

### Code Style

- Standard Go formatting and conventions
- Error handling with user-friendly messages
- Consistent use of utilities for logging and output

### Dev Pipeline

When implementing a new feature, follow this workflow:

1. **Implement feature** - Add the feature code following existing patterns
2. **Build and test** - Run `make dev EMBED_CREDS=1` to build with credentials from `.env`, then test with `./gtasks <command>`
3. **Update docs** - Update `docs/` Hugo site if the feature adds/changes commands
4. **Update README** - Update `README.md` if needed for user-facing changes
5. **Commit** - Create a commit with a clear message describing the change
6. **Do NOT push or release** - Wait for explicit approval before pushing or creating releases

### Testing Commands

- Build for dev: `make dev EMBED_CREDS=1` (creates `./gtasks` binary with embedded credentials)
- Test authentication: `./gtasks login`
- Test task operations: `./gtasks tasks view`

### Key Files to Understand

- `cmd/root.go` - Main CLI structure and help text
- `api/auth.go` - Login flow, keyring token storage, migration
- `cmd/tasks.go` - Task viewing with table formatting, default list resolution
- `internal/config/credentials.go` - OAuth2 config generation
- `internal/config/file.go` - koanf-based layered config loader
- `internal/config/get_install_loc.go` - XDG + legacy config dir resolution

## API Integration

### Google Tasks API Usage

- Scopes: `tasks.TasksScope` (read/write access to tasks)
- Endpoints used: Tasks.List, Tasks.Insert, Tasks.Patch, Tasks.Delete
- Task lists: TaskLists.List, TaskLists.Insert, TaskLists.Delete

### Error Handling

- Graceful handling of authentication failures
- User-friendly error messages for API failures
- Proper cleanup of resources

## Notable Features

- **Interactive Mode**: Prompts for task list selection when not specified; auto-selects when only one list exists
- **Default Task List**: Configurable via `-l` flag, `GTASKS_DEFAULT_TASKLIST` env var, or config file
- **Date Parsing**: Flexible date input using `dateparse` library
- **Cross-platform**: Builds for multiple operating systems
- **Table Output**: Formatted display of tasks with status indicators
- **Sorting**: Multiple sort options for task views
- **Documentation**: Hugo-based documentation website
- **Layered Config**: koanf-based config supports toml/yaml/json, env vars, and build-time defaults

## Security Considerations

- OAuth2 tokens stored in system keyring (not plaintext); falls back to file with 0600 perms
- Uses official Google OAuth2 flow with PKCE
- No sensitive data logged or exposed
- Config file containing credentials should be chmod 600

## Personal Project Context

This is one of Siddhartha Varma's notable personal projects, created in December 2021. It demonstrates:

- Go programming proficiency
- OAuth2 implementation experience
- CLI tool development skills
- Google API integration
- Cross-platform build systems
- Documentation and user experience focus

The project has garnered 80+ GitHub stars, indicating community adoption and usefulness.
