# sopsv

`sopsv` is a simple CLI tool that uses [sops](https://github.com/getsops/sops) and [age](https://github.com/FiloSottile/age) as a local secret vault. It manages the encryption and decryption of simple key-value pairs stored in YAML files.

## Setup

1. Make sure you have Go installed.
2. Build the project:
   ```bash
   go build -o sopsv
   ```
3. Move `sopsv` to a location in your `$PATH`.

*Note: `sopsv` requires `sops` to be installed and available in your `$PATH`.*

## Initialization

You do not need to manually initialize the vault. `sopsv` will automatically:
- Generate an `age` key if one does not exist at `~/.config/sopsv/keys.txt`
- Set `SOPS_AGE_KEY_FILE` internally during execution
- Store vaults in your user data directory (`~/.local/share/sopsv/` on Linux, `%LocalAppData%/sopsv/` on Windows)

## Commands

### Manage Vaults

- **Create a new vault:**
  ```bash
  sopsv vaults new myvault
  ```
  *(Creates a new encrypted file in the data directory)*

- **List all managed vaults:**
  ```bash
  sopsv vaults list
  ```

- **Set the default vault:**
  ```bash
  sopsv vaults use myvault
  ```

- **Show the decrypted contents of a vault:**
  ```bash
  sopsv vaults show myvault
  ```
  *(If no vault is specified, it uses the default vault)*

- **Delete a vault:**
  ```bash
  sopsv vaults delete myvault
  ```

### Accessing Secrets

- **Get a specific key:**
  ```bash
  sopsv get mykey
  ```
  *(Reads from the default vault. You can override it with `-v myvault`)*

### Executing Commands with Secrets

- **Run a command with secrets injected as Environment Variables:**
  ```bash
  sopsv exec --vaults myvault,myothervault -- env
  ```
  *(If no vaults are specified, it defaults to the configured default vault)*
