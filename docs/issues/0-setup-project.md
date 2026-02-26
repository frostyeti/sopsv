---
type: task
tags: ["feature", "cli"]
---

# setup project

look at github.com/frostyeti/kpv and github.com/frostyeti/osv

create a cli that uses sops and age as a secret vault locally where
the yaml, json, and .env files are simple key value pairs that can be encrypted and decrypted with sops and age.

the command should be able to pull one value out of the store and decrypt it with sops and age. 

the cli should be able to determine if sops and ages are configured and if not, create the age file and sops.yaml file. 

the default store should be in the user profile. named stores/vaults should also be in the user's profile under ~/.local/share/sopsv (and the similar for windows).

if the user wants to use a local path, then they would need to use the full
or relative path to the store e.g. ./myvault.yaml or /home/user/myvault.yaml.  

Use viper to store configuration settings allow the user to specify the default store.

There should a vaults sub command that has a use, show, list, delete, and new command that stores meta in the config file and if a new vault is created, store it in the user's profile under ~/.local/share/sopsv with the name of the vault and a .yaml extension.

For an exec command, the cli should be able to supply multiple stores
to pull all the secrets from and then execute the command with the secrets as environment variables.

create a readme with instructions on how to use the cli and how to set up the project.

## Implementation Plan
1. **Initialize Project:** Run `go mod init github.com/frostyeti/sopsv`.
2. **Setup CLI Framework:** Add Cobra for command routing and Viper for config (`~/.config/sopsv/config.yaml`).
3. **Core Dependencies:** Integrate `getsops/sops/v3` and `FiloSottile/age` to handle in-memory encryption/decryption where possible, or shell out if required. Include `fatih/color` for terminal styling.
4. **Auto-Initialization:** Implement logic to detect missing `.sops.yaml` or `age` keys on startup and generate them interactively or automatically in `~/.config/sopsv/keys.txt`.
5. **Vaults Command Group (`sopsv vaults`):**
   - `new`: Create a new encrypted file in OS-specific data dir (`~/.local/share/sopsv/`).
   - `list`: Show all managed vaults.
   - `use`: Set the default vault in Viper config.
   - `show`: Display decrypted contents of a vault.
   - `delete`: Safely remove a vault.
6. **Key/Value Retrieval:** Implement logic to parse decrypted yaml/json/env and fetch single values (`sopsv get <key>`).
7. **Exec Command (`sopsv exec`):** Implement sub-process spawning that aggregates secrets from provided vaults and injects them into the child process's environment variables.
8. **Documentation:** Write a comprehensive `README.md` covering setup, initialization, and usage examples.
