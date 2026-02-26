---
title: Managing Vaults
description: How to manage secret vaults
---

Vaults are encrypted YAML files that store your secrets. By default, they are stored in `~/.local/share/sopsv` (Unix) or `%LocalAppData%/sopsv` (Windows).

## Create a new vault

To create a new vault:

```bash
sopsv vaults new myvault
```

This will generate an `age` key if you don't already have one, and encrypt an empty vault using `sops`.

## List vaults

To see all your vaults:

```bash
sopsv vaults list
```

## Set a default vault

Tired of typing `--vault myvault` on every command? Set a default:

```bash
sopsv vaults use myvault
```

## View vault contents

To decrypt and view the raw YAML contents of a vault:

```bash
sopsv vaults show myvault
```

## Edit a vault

To securely open a vault in your `$EDITOR`, decrypting it on the fly and re-encrypting it upon saving:

```bash
sopsv vaults edit myvault
```

## Delete a vault

```bash
sopsv vaults delete myvault
```
