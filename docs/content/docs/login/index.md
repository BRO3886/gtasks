---
title: "Logging In"
draft: false
weight: 3
summary: Login to Google Tasks CLI with your Google Account
---

## Prerequisites

Before logging in, gtasks needs Google OAuth2 credentials. Set them via environment variables or a config file:

```bash
# Option A — environment variables
export GTASKS_CLIENT_ID="your-client-id.apps.googleusercontent.com"
export GTASKS_CLIENT_SECRET="your-client-secret"
```

```toml
# Option B — ~/.config/gtasks/config.toml  (or ~/.gtasks/config.toml for legacy installs)
[credentials]
client_id     = "your-client-id.apps.googleusercontent.com"
client_secret = "your-client-secret"
```

See the [Configuration](../configuration/) page for details on obtaining credentials.

## Login

```
gtasks login
```

- This opens your browser for Google OAuth2 authentication and starts a local callback server.
- If the browser does not open automatically, the CLI prints a URL you can visit manually.
- After you grant access, the browser shows a success page — close it and return to the terminal.
- Your token is saved to the **system keyring** (macOS Keychain, Linux Secret Service, Windows Credential Manager). On headless systems without a keyring, it falls back to a file in the config directory.

## Logout

```
gtasks logout
```

Removes the stored token from the keyring (and token file if present).
