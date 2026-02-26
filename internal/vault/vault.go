package vault

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/frostyeti/sopsv/internal/age"
	"github.com/frostyeti/sopsv/internal/config"
	"github.com/frostyeti/sopsv/internal/sops"
	"gopkg.in/yaml.v3"
)

var (
	ageEnsureKeyExists = age.EnsureKeyExists
	ageGetPublicKey    = age.GetPublicKey
	sopsDecryptFile    = sops.DecryptFile
	sopsExtractValue   = sops.ExtractValue
	sopsSetValue       = sops.SetValue
)

// ResolveVaultPath returns the absolute path to a vault.
func ResolveVaultPath(name string) string {
	if strings.Contains(name, string(filepath.Separator)) || strings.HasSuffix(name, ".yaml") {
		absPath, err := filepath.Abs(name)
		if err == nil {
			return absPath
		}
		return name
	}
	return filepath.Join(config.DataDir, name+".yaml")
}

// NewVault creates a new empty vault.
func NewVault(name string) error {
	path := ResolveVaultPath(name)

	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("vault %s already exists at %s", name, path)
	}

	pubKey, err := ageEnsureKeyExists()
	if err != nil {
		return fmt.Errorf("failed to ensure age key: %w", err)
	}

	cleartext := []byte("sopsv_init: true\n")
	return encryptFileWithSops(path, cleartext, pubKey)
}

func encryptFileWithSops(path string, cleartext []byte, pubKey string) error {
	tempFile, err := os.CreateTemp("", "sopsv-*.yaml")
	if err != nil {
		return err
	}
	defer os.Remove(tempFile.Name())

	if _, err := tempFile.Write(cleartext); err != nil {
		return err
	}
	tempFile.Close()

	cmd := sopsCmd("--encrypt", "--age", pubKey, tempFile.Name())

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("sops encrypt failed: %s: %w", string(out), err)
	}

	return os.WriteFile(path, out, 0600)
}

var sopsCmd = func(args ...string) *exec.Cmd {
	return exec.Command("sops", args...)
}

// ListVaults lists all vaults in the data directory.
func ListVaults() ([]string, error) {
	entries, err := os.ReadDir(config.DataDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read data dir: %w", err)
	}

	var vaults []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".yaml") {
			vaults = append(vaults, strings.TrimSuffix(entry.Name(), ".yaml"))
		}
	}
	return vaults, nil
}

// DeleteVault deletes a vault.
func DeleteVault(name string) error {
	path := ResolveVaultPath(name)
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("failed to delete vault %s: %w", name, err)
	}
	return nil
}

// ReadVault reads and decrypts a vault returning key/value pairs.
func ReadVault(name string) (map[string]interface{}, error) {
	path := ResolveVaultPath(name)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("vault %s not found at %s", name, path)
	}

	cleartext, err := sopsDecryptFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt vault %s: %w", name, err)
	}

	var data map[string]interface{}
	if err := yaml.Unmarshal(cleartext, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml from vault %s: %w", name, err)
	}

	if data == nil {
		data = make(map[string]interface{})
	}

	return data, nil
}

// WriteVault updates the vault with the provided data (encrypting the whole file).
func WriteVault(name string, data map[string]interface{}) error {
	path := ResolveVaultPath(name)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("vault %s not found at %s", name, path)
	}

	pubKey, err := ageGetPublicKey()
	if err != nil {
		return fmt.Errorf("failed to get public key for encryption: %w", err)
	}

	cleartext, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal yaml: %w", err)
	}

	return encryptFileWithSops(path, cleartext, pubKey)
}

// EditVault opens the vault interactively in the user's editor using sops.
func EditVault(name string) error {
	path := ResolveVaultPath(name)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("vault %s not found at %s", name, path)
	}

	cmd := sopsCmd(path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// SetSecret sets a single secret in the vault using sops set
func SetSecret(vaultName string, key string, value interface{}) error {
	path := ResolveVaultPath(vaultName)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Auto initialize if missing
		err = NewVault(vaultName)
		if err != nil {
			return err
		}
	}

	return sopsSetValue(path, key, value)
}

// RemoveSecret removes a single secret from the vault by setting it to empty
func RemoveSecret(vaultName string, key string) error {
	path := ResolveVaultPath(vaultName)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("vault %s not found", vaultName)
	}

	return sopsSetValue(path, key, "")
}

// GetSecret extracts a single secret from the vault
func GetSecret(vaultName string, key string) (string, error) {
	path := ResolveVaultPath(vaultName)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", fmt.Errorf("vault %s not found", vaultName)
	}

	return sopsExtractValue(path, key)
}

// ListKeys returns all keys in a vault without decrypting it
func ListKeys(vaultName string) ([]string, error) {
	path := ResolveVaultPath(vaultName)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("vault %s not found", vaultName)
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read vault file: %w", err)
	}

	var data map[string]interface{}
	if err := yaml.Unmarshal(bytes, &data); err != nil {
		return nil, fmt.Errorf("failed to parse vault yaml: %w", err)
	}

	var keys []string
	for k := range data {
		if k != "sops" { // ignore sops metadata key
			keys = append(keys, k)
		}
	}
	return keys, nil
}
