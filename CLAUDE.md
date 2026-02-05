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

## Core Functionality

### Authentication (`api/auth.go`)

- OAuth2 flow implementation for Google Tasks API
- Token storage and retrieval from local filesystem
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

- Uses hardcoded Google OAuth2 credentials for the CLI app
- Stores user tokens in `~/.gtasks/token.json`
- Configuration file: `~/.gtasks/config.json`

### Installation Locations

- Detects installation path dynamically
- Creates necessary directories for config and token storage

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
- `make release` - Create GitHub release

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

- `cmd/root.go:19-31` - Main CLI structure and help text
- `api/auth.go:17-28` - Login flow implementation
- `cmd/tasks.go:47-93` - Task viewing with table formatting
- `internal/config/credentials.go:27-32` - OAuth2 config generation

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

- **Interactive Mode**: Prompts for task list selection when not specified
- **Date Parsing**: Flexible date input using `dateparse` library
- **Cross-platform**: Builds for multiple operating systems
- **Table Output**: Formatted display of tasks with status indicators
- **Sorting**: Multiple sort options for task views
- **Documentation**: Hugo-based documentation website

## Security Considerations

- OAuth2 tokens stored locally with appropriate file permissions (0600)
- Uses official Google OAuth2 flow
- No sensitive data logged or exposed
- Client credentials are for a registered Google Cloud project

## Personal Project Context

This is one of Siddhartha Varma's notable personal projects, created in December 2021. It demonstrates:

- Go programming proficiency
- OAuth2 implementation experience
- CLI tool development skills
- Google API integration
- Cross-platform build systems
- Documentation and user experience focus

The project has garnered 80+ GitHub stars, indicating community adoption and usefulness.
