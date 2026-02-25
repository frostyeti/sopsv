---
title: Executing Commands
description: How to inject secrets into commands
---

The most powerful feature of `sopsv` is the `exec` command. It allows you to run any command or script with your secrets automatically injected as environment variables.

## Basic Execution

```bash
sopsv exec --vaults myvault -- npm start
```

This will decrypt `myvault`, load all key-value pairs into environment variables, and then run `npm start`. The secrets will be available to the application, but they are never written to disk unencrypted, and they won't show up in your terminal history.

## Multiple Vaults

You can merge secrets from multiple vaults by providing a comma-separated list. If there are duplicate keys, the right-most vault takes precedence.

```bash
sopsv exec --vaults default,production -- ./deploy.sh
```

## Explicit File Paths

You don't have to use managed vaults. You can pass explicit paths to SOPS-encrypted YAML files.

```bash
sopsv exec --vaults ./config/secrets.yaml -- python main.py
```
