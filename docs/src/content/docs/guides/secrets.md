---
title: Managing Secrets
description: How to add, remove, and read secrets
---

Within a vault, you can store key-value pairs representing your secrets.

## Set a secret

To set a secret, use the `secrets set` command:

```bash
sopsv secrets set --key API_KEY --value "my-super-secret-key" --vault myvault
```

If you don't specify `--vault`, it will use your default vault. If the vault doesn't exist, it will be automatically created.

### Generate a random secret

You can automatically generate a secure random value:

```bash
sopsv secrets set --key JWT_SECRET --generate --size 32
```

### Read from stdin

For multiline secrets or piping from another command:

```bash
cat my-cert.pem | sopsv secrets set --key SSL_CERT --stdin
```

## Get a secret

To retrieve a single secret:

```bash
sopsv secrets get --key API_KEY
```

## List all secrets

To view all keys (and optionally values) in a vault:

```bash
sopsv secrets ls
```

To see the values as well, use `--show-values`.

## Remove a secret

To delete a secret from the vault:

```bash
sopsv secrets rm --key API_KEY
```
