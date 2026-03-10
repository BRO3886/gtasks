---
title: "Installation"
description: "Install gtasks via curl, go install, or manual binary download. Available for macOS, Linux, and Windows."
draft: false
weight: 2
sitemap:
  priority: 0.8
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

1. Download the tarball for your platform from [GitHub Releases](https://github.com/BRO3886/gtasks/releases)
   (e.g. `gtasks_mac_arm64_vX.Y.Z.tar.gz`)
2. Extract and move to a directory in your `PATH`:

   ```bash
   tar -xzf gtasks_*.tar.gz
   mv gtasks ~/.local/bin/   # or /usr/local/bin/
   ```

   **Windows**: extract the `.tar.gz`, move `gtasks.exe` to a folder in your PATH
3. Verify: `gtasks --version`

## Go install

If you have [Go](https://golang.org/) installed:

```bash
go install github.com/BRO3886/gtasks@latest
```

