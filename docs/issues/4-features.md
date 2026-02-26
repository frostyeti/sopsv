---
type: docs
tags: ["docs", "enhancement"]
---

# Feature Proposals

## Context
Proposed enhancements for the `sopsv` CLI to improve security, usability, and integrations after the initial setup.

## Implementation Plan

1. **Environment Variable Sync (`sopsv sync`):**
   - *Proposal*: A command to watch a vault and automatically sync its decrypted secrets to a local `.env` file for local development convenience without needing to run `exec` every time.
   - *Implementation*: Use `fsnotify` to listen for file changes on the vault and write to a `.env` destination.

2. **Secret Rotation Assistance:**
   - *Proposal*: Commands to help bulk-rotate `age` keys, re-encrypting all stored vaults seamlessly without data loss.
   - *Implementation*: Introduce `sopsv rotate` that reads the current configuration, prompts for a new `age` recipient public key, and iterates through all managed vaults, decrypting with the old key and encrypting with the new one.

3. **Cloud KMS Integration (AWS/GCP/Azure):**
   - *Proposal*: While currently focused on local `age` keys, allowing SOPS to fall back or integrate with cloud KMS for team-based sharing.
   - *Implementation*: Integrate `getsops/sops/v3` cloud key providers and expose CLI flags/config options in `~/.config/sopsv/config.yaml` to configure ARN/Project/Vault URLs.

4. **Interactive Vault Editor (`sopsv edit`):**
   - *Proposal*: An interactive command that securely opens a vault in the user's `$EDITOR`, decrypting it on the fly and re-encrypting upon save, similar to `sops edit` but with `sopsv`'s context/profiles.
   - *Implementation*: Read the `$EDITOR` environment variable, decrypt the vault to a secure temporary file (e.g., in memory or `os.TempDir`), wait for the editor process to exit, and re-encrypt the file if modified.

5. **Git Pre-commit Hook Integration:**
   - *Proposal*: Provide a command `sopsv hook install` that sets up a Git pre-commit hook to prevent accidental commits of unencrypted secret files or plaintext `.env` files.
   - *Implementation*: Write a shell script template to `.git/hooks/pre-commit` that checks staged files against `sops` metadata or known plaintext `.env` patterns.

6. **Interactive `init` command (`sopsv init`):**
   - *Proposal*: A setup wizard that asks the user where they want to store secrets, what keys to use, and helps them migrate existing plaintext `.env` files into a managed vault.
   - *Implementation*: Use a prompt library like `charmbracelet/huh` or `AlecAivazis/survey` to walk the user through setup and migration.
