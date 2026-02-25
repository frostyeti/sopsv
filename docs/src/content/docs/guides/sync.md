---
title: Syncing Secrets
description: How to sync secrets to a local .env file
---

Often, you might have legacy tools or local servers that expect a `.env` file instead of supporting `sopsv exec` directly.

The `sync` command solves this by continuously decrypting a vault and writing its secrets to a local file (like `.env`).

## One-time Sync

To generate a `.env` file once:

```bash
sopsv sync myvault --dest .env
```

## Continuous Sync (Watch mode)

You can run `sync` in watch mode to automatically keep the `.env` file updated whenever the vault changes:

```bash
sopsv sync myvault --dest .env --watch
```

This is perfect for running in a background terminal tab while you develop.

*Note: Always remember to add the generated `.env` file to your `.gitignore` to prevent accidentally committing plaintext secrets!*
