---
title: Introduction
description: Introduction to sopsv
---

`sopsv` is a powerful and simple CLI tool that manages local encrypted secrets using `sops` and `age`. It is designed to act as a local secret vault for managing, encrypting, and decrypting simple key-value pairs, providing a developer experience mirroring existing tools like `kpv`, but with the robustness of SOPS under the hood.

## Why sopsv?

Developers often need to manage local secrets (API keys, database passwords, etc.) securely without committing them to version control. While SOPS is an excellent tool for encrypting files, its standard usage requires understanding complex command-line arguments and configuration files. `sopsv` provides a wrapper around SOPS and `age` to give you a simple, intuitive developer experience for managing simple key-value vaults.

## Features

- **Automatic `age` key management:** No need to manually generate or configure `age` keys. `sopsv` creates and manages a local X25519 identity for you automatically.
- **Simple Vault CRUD:** Easily create, list, delete, and manage multiple secret vaults.
- **Key-Value Secret Management:** Set, get, remove, and list secrets within your vaults using intuitive commands.
- **Command Execution:** Run your applications with secrets automatically injected as environment variables.
- **Automatic Tooling:** Built-in `sopsv tools install` command downloads and configures the required `sops` and `age` binaries for your platform.
