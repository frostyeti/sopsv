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

// ResolveVaultPath returns the absolute path to a vault.
// If it has path separators or ends in .yaml, it resolves to local or absolute path.
// Otherwise it looks in the data directory (~/.local/share/sopsv/).
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

	pubKey, err := age.EnsureKeyExists()
	if err != nil {
		return fmt.Errorf("failed to ensure age key: %w", err)
	}

	cleartext := []byte("default: value\n")
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

func sopsCmd(args ...string) *exec.Cmd {
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

	cleartext, err := sops.DecryptFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt vault %s: %w", name, err)
	}

	var data map[string]interface{}
	if err := yaml.Unmarshal(cleartext, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml from vault %s: %w", name, err)
	}

	return data, nil
}
