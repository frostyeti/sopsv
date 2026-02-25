---
title: Installation
description: How to install sopsv
---

You can download `sopsv` from the [GitHub Releases](https://github.com/frostyeti/sopsv/releases) page for your operating system.

## Installing Dependencies

`sopsv` requires `sops` and `age` to function. We've built an installer to automatically download and configure them for you!

```bash
sopsv tools install
```

This command will:
1. Detect your operating system and architecture.
2. Download the latest `sops` and `age` binaries from GitHub.
3. Install them into your local user directory (e.g., `~/.local/bin` on Linux/macOS, or `%LocalAppData%\Programs` on Windows).
4. Automatically update your `PATH` environment variable so you can use them anywhere.

*Note: If your PATH was updated, you may need to restart your terminal or run `source ~/.bashrc` (or equivalent) to use the new commands.*

### Global Installation

To install the tools globally (e.g., `/usr/local/bin`), run:

```bash
sudo sopsv tools install --global
```
