---
title: "Configuration"
description: "Configure gtasks with a config file, environment variables, or build-time flags. Supports TOML, YAML, and JSON."
draft: false
weight: 5
sitemap:
  priority: 0.8
---

GTasks supports configuration through a config file, environment variables, and CLI flags.
Each layer can override the one below it:

```
CLI flag  >  environment variable  >  config file  >  build-time default
```

## Config file location

GTasks follows the [XDG Base Directory Specification](https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html).
The config file is looked for in these directories in order:

| Priority | Path | Notes |
|----------|------|-------|
| 1 | `$XDG_CONFIG_HOME/gtasks/` | XDG standard path; `XDG_CONFIG_HOME` defaults to `~/.config` |
| 2 | `~/.gtasks/` | Legacy path; used automatically when that directory already exists |

New installations use `~/.config/gtasks/` by default. Existing installations that already have `~/.gtasks/` continue using it. If both exist, XDG wins and gtasks prints a migration warning with the exact command to move your files.

> **Note:** Authentication tokens are stored in the **system keyring**, not on disk. On headless systems where no keyring is available, the token falls back to a `token.json` file in the config directory.

## Creating the config file

```bash
mkdir -p ~/.config/gtasks
touch ~/.config/gtasks/config.toml
chmod 600 ~/.config/gtasks/config.toml  # recommended if storing credentials
```

## Config file format

gtasks supports **TOML**, **YAML**, and **JSON** — the first file found wins:

- `config.toml`
- `config.yaml` / `config.yml`
- `config.json`

All keys are optional — omit any section or key you do not need.

```toml
# GTasks configuration file (~/.config/gtasks/config.toml)

[credentials]
# Google OAuth2 client ID — required to use gtasks.
# Overridden by: GTASKS_CLIENT_ID environment variable.
client_id = "your-client-id.apps.googleusercontent.com"

# Google OAuth2 client secret — required to use gtasks.
# Overridden by: GTASKS_CLIENT_SECRET environment variable.
client_secret = "your-client-secret"

[tasks]
# Default task list to use when the -l / --tasklist flag is not provided.
# When set, gtasks skips the interactive task list prompt.
# Overridden by: GTASKS_DEFAULT_TASKLIST environment variable, then the -l flag.
# default_task_list = "My Tasks"
```

## Settings reference

### `[credentials]`

Required for all users. gtasks does not ship with embedded credentials — you must supply your own Google OAuth2 client ID and secret.

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

### Supply credentials

```toml
[credentials]
client_id     = "123456789-abc.apps.googleusercontent.com"
client_secret = "GOCSPX-xxxxxxxxxxxxxxxxxxxxxxx"
```

This is equivalent to setting `GTASKS_CLIENT_ID` and `GTASKS_CLIENT_SECRET` as environment
variables, but stored persistently in the config file. Recommended: `chmod 600` the file.

## Environment variables

All settings in the config file can be overridden with environment variables:

| Variable | Overrides |
|----------|-----------|
| `GTASKS_CLIENT_ID` | `credentials.client_id` |
| `GTASKS_CLIENT_SECRET` | `credentials.client_secret` |
| `GTASKS_DEFAULT_TASKLIST` | `tasks.default_task_list` |
| `XDG_CONFIG_HOME` | Base directory for the config folder (XDG spec) |
