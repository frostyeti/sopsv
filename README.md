# sopsv

`sopsv` is a simple CLI tool that uses [sops](https://github.com/getsops/sops) and [age](https://github.com/FiloSottile/age) as a local secret vault. It manages the encryption and decryption of simple key-value pairs stored in YAML files.

## Setup

1. Make sure you have Go installed.
2. Build the project:
   ```bash
   go build -o sopsv
   ```
3. Move `sopsv` to a location in your `$PATH`.

*Note: `sopsv` requires `sops` and `age` to be installed and available in your `$PATH`. You can install them automatically using the built-in installer tool:*

```bash
# Install to ~/.local/bin (Mac/Linux) or %LOCALAPPDATA%\Programs\bin (Windows)
sopsv tools install

# Install globally to /usr/local/bin (Mac/Linux) or C:\Program Files\bin (Windows)
sudo sopsv tools install --global

# Install to a custom directory
sopsv tools install --dest /path/to/my/bin
```

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

### Manage Secrets

The `secrets` subcommand group allows you to manage entries inside a vault.

- **Set a secret:**
  ```bash
  sopsv secrets set --key api-token --value "secret123"
  sopsv secrets set --key api-token --generate --size 32
  echo "secret" | sopsv secrets set --key api-token --stdin
  ```

- **Get a secret:**
  ```bash
  sopsv secrets get --key api-token
  ```
  You can also get multiple keys and format the output:
  ```bash
  sopsv secrets get --key DB_USER --key DB_PASS --format dotenv
  sopsv secrets get --key DB_USER --key DB_PASS --format sh
  sopsv secrets get --key DB_USER --key DB_PASS --format json
  ```

- **Ensure a secret exists:**
  *(Returns the secret if it exists, otherwise generates and saves a new random one)*
  ```bash
  sopsv secrets ensure --key my-random-key --size 24
  ```

- **Remove a secret:**
  ```bash
  sopsv secrets rm --key api-token
  ```

- **List all keys in a vault:**
  ```bash
  sopsv secrets ls
  ```

- **Export and Import:**
  ```bash
  sopsv secrets export --json --pretty --file secrets.json
  sopsv secrets import --file secrets.json
  ```

*Note: All `secrets` commands operate on the default vault unless overridden by the `--vault` (or `-V`) flag.*

### Executing Commands with Secrets

- **Run a command with secrets injected as Environment Variables:**
  ```bash
  sopsv exec --vaults myvault,myothervault -- env
  ```
  *(If no vaults are specified, it defaults to the configured default vault)*
