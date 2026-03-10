---
title: 'Installation'
draft: false
weight: 2
summary: "Download binary from [releases](https://github.com/BRO3886/google-tasks-cli/releases) if you know what you're doing. Otherwise, check detailed instructions."
---


## macOS / Linux (recommended)

```bash
curl -fsSL https://gtasks.sidv.dev/install | bash
```

Installs the latest release to `~/.local/bin`. To install elsewhere:

```bash
INSTALL_DIR=/usr/local/bin curl -fsSL https://gtasks.sidv.dev/install | bash
```

## Manual install

1. Download the binary for your system from [GitHub Releases](https://github.com/BRO3886/gtasks/releases)
2. Move to a directory in your `PATH`:
   - **macOS/Linux**: `mv gtasks ~/.local/bin/` (or `/usr/local/bin/`) and `chmod +x gtasks`
   - **Windows**: Move to a folder already in your PATH, or add its folder to PATH
3. Verify: `gtasks --version`

## Go install

If you have [Go](https://golang.org/) installed:

```bash
go install github.com/BRO3886/gtasks@latest
```

