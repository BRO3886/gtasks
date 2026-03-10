---
title: "Configuration"
draft: false
weight: 5
summary: Configure gtasks with a config file, environment variables, or build-time flags
---

GTasks supports configuration through a TOML file, environment variables, and CLI flags.
Each layer can override the one below it:

```
CLI flag  >  environment variable  >  config file  >  build-time default
```

## Config file location

GTasks follows the [XDG Base Directory Specification](https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html).
The config file is named **`config.toml`** and is found by checking these directories in order:

| Priority | Path | Notes |
|----------|------|-------|
| 1 | `$XDG_CONFIG_HOME/gtasks/config.toml` | XDG standard path; `XDG_CONFIG_HOME` defaults to `~/.config` |
| 2 | `~/.gtasks/config.toml` | Legacy path; used automatically when the `~/.gtasks/` directory already exists |

New installations always use the XDG path (`~/.config/gtasks/` by default). Existing
installations that already have a `~/.gtasks/` directory continue using it without any
changes. To migrate, move the directory and set `XDG_CONFIG_HOME` if needed.

Authentication tokens (`token.json`) are stored alongside `config.toml` in the same directory.

## Creating the config file

Create the file manually:

```bash
mkdir -p ~/.config/gtasks
touch ~/.config/gtasks/config.toml
```

## Config file format

The file uses [TOML](https://toml.io) syntax. All keys are optional — omit any section
or key you do not need.

```toml
# GTasks configuration file
# Location: $XDG_CONFIG_HOME/gtasks/config.toml  (usually ~/.config/gtasks/config.toml)
#           ~/.gtasks/config.toml                 (legacy path, used when ~/.gtasks/ exists)

[credentials]
# OAuth2 client ID for the Google Tasks API.
# Only required when you are building gtasks from source with your own Google Cloud project.
# Overridden by: GTASKS_CLIENT_ID environment variable.
# client_id = "your-client-id.apps.googleusercontent.com"

# OAuth2 client secret for the Google Tasks API.
# Only required when you are building gtasks from source with your own Google Cloud project.
# Overridden by: GTASKS_CLIENT_SECRET environment variable.
# client_secret = "your-client-secret"

[tasks]
# Default task list to use when the -l / --tasklist flag is not provided.
# When set, gtasks skips the interactive task list prompt.
# Overridden by: GTASKS_DEFAULT_TASKLIST environment variable, then the -l flag.
# default_task_list = "My Tasks"
```

## Settings reference

### `[credentials]`

These settings are needed only when building gtasks from source with your own Google Cloud
project. Released binaries have credentials embedded at build time.

| Key | Type | Env var override | Description |
|-----|------|-----------------|-------------|
| `client_id` | string | `GTASKS_CLIENT_ID` | Google OAuth2 client ID |
| `client_secret` | string | `GTASKS_CLIENT_SECRET` | Google OAuth2 client secret |

### `[tasks]`

| Key | Type | Env var override | CLI flag override | Description |
|-----|------|-----------------|-------------------|-------------|
| `default_task_list` | string | `GTASKS_DEFAULT_TASKLIST` | `-l` / `--tasklist` | Task list selected automatically when no flag is given |

## Examples

### Set a default task list

To always operate on "Work" without typing `-l "Work"` every time:

```toml
[tasks]
default_task_list = "Work"
```

Then:

```bash
# No prompt, uses "Work" automatically
gtasks tasks view
gtasks tasks add -t "Finish report"

# Override for a single command with the flag
gtasks tasks view -l "Personal"
```

### Supply credentials for a custom build

```toml
[credentials]
client_id     = "123456789-abc.apps.googleusercontent.com"
client_secret = "GOCSPX-xxxxxxxxxxxxxxxxxxxxxxx"
```

This is equivalent to setting the `GTASKS_CLIENT_ID` and `GTASKS_CLIENT_SECRET` environment
variables, but stored persistently in the config file.

## Environment variables

All settings in the config file can be overridden with environment variables:

| Variable | Overrides |
|----------|-----------|
| `GTASKS_CLIENT_ID` | `credentials.client_id` |
| `GTASKS_CLIENT_SECRET` | `credentials.client_secret` |
| `GTASKS_DEFAULT_TASKLIST` | `tasks.default_task_list` |
| `XDG_CONFIG_HOME` | Base directory for the config folder (XDG spec) |
