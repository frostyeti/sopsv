package vault

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/frostyeti/sopsv/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveVaultPath(t *testing.T) {
	config.DataDir = "/mock/data"

	// Just a name
	assert.Equal(t, "/mock/data/myvault.yaml", ResolveVaultPath("myvault"))

	// Explicit yaml
	assert.True(t, strings.HasSuffix(ResolveVaultPath("myvault.yaml"), "myvault.yaml"))

	// Absolute path
	absPath := "/tmp/testvault.yaml"
	if os.PathSeparator == '\\' {
		absPath = "C:\\tmp\\testvault.yaml"
	}
	assert.Equal(t, absPath, ResolveVaultPath(absPath))
}

func setupMockEnvironment(t *testing.T) string {
	tempDir := t.TempDir()
	config.DataDir = tempDir

	// Mock age
	ageEnsureKeyExists = func() (string, error) { return "age1mock", nil }
	ageGetPublicKey = func() (string, error) { return "age1mock", nil }

	// Mock sops cmd to just copy the input cleartext to output path
	// We will create a fake cmd that just echoes success
	sopsCmd = func(args ...string) *exec.Cmd {
		// Just run something harmless that succeeds, e.g. "echo" or "true"
		// Better yet, for encryption we just simulate a success write in encryptFileWithSops
		// The easiest way to mock exec is to use the standard library TestHelper approach,
		// but since we just need it not to fail and to output something:
		if len(args) > 0 && args[0] == "--encrypt" {
			// Find the output path (not an argument to sops, it outputs to stdout)
			// Actually our sopsCmd is expected to output encrypted text to stdout
			// Let's make it echo something
			cmd := exec.Command("echo", "encrypted_mock_data")
			return cmd
		}
		return exec.Command("echo")
	}

	// Mock DecryptFile
	sopsDecryptFile = func(filePath string) ([]byte, error) {
		return []byte("key1: value1\nkey2: value2\n"), nil
	}

	// Mock SetValue and ExtractValue
	sopsSetValue = func(filePath string, key string, value interface{}) error {
		return nil
	}
	sopsExtractValue = func(filePath string, key string) (string, error) {
		if key == "nonexistent" {
			return "", fmt.Errorf("not found")
		}
		return "mockvalue", nil
	}

	return tempDir
}

func TestNewVault(t *testing.T) {
	tempDir := setupMockEnvironment(t)

	err := NewVault("newvault")
	require.NoError(t, err)

	// Verify file was "encrypted"
	vaultPath := filepath.Join(tempDir, "newvault.yaml")
	content, err := os.ReadFile(vaultPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "encrypted_mock_data")

	// Verify it fails if already exists
	err = NewVault("newvault")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

func TestListVaults(t *testing.T) {
	tempDir := setupMockEnvironment(t)

	// Create some fake vaults
	err := os.WriteFile(filepath.Join(tempDir, "v1.yaml"), []byte("mock"), 0600)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(tempDir, "v2.yaml"), []byte("mock"), 0600)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(tempDir, "not_a_vault.txt"), []byte("mock"), 0600)
	require.NoError(t, err)

	vaults, err := ListVaults()
	require.NoError(t, err)
	assert.Len(t, vaults, 2)
	assert.Contains(t, vaults, "v1")
	assert.Contains(t, vaults, "v2")
}

func TestDeleteVault(t *testing.T) {
	tempDir := setupMockEnvironment(t)
	vaultPath := filepath.Join(tempDir, "todelete.yaml")
	err := os.WriteFile(vaultPath, []byte("mock"), 0600)
	require.NoError(t, err)

	err = DeleteVault("todelete")
	require.NoError(t, err)

	_, err = os.Stat(vaultPath)
	assert.True(t, os.IsNotExist(err))
}

func TestReadVault(t *testing.T) {
	tempDir := setupMockEnvironment(t)
	vaultPath := filepath.Join(tempDir, "readvault.yaml")
	err := os.WriteFile(vaultPath, []byte("encrypted_data"), 0600)
	require.NoError(t, err)

	data, err := ReadVault("readvault")
	require.NoError(t, err)

	assert.Equal(t, "value1", data["key1"])
	assert.Equal(t, "value2", data["key2"])
}

func TestWriteVault(t *testing.T) {
	tempDir := setupMockEnvironment(t)
	vaultPath := filepath.Join(tempDir, "writevault.yaml")
	err := os.WriteFile(vaultPath, []byte("encrypted_data"), 0600)
	require.NoError(t, err)

	data := map[string]interface{}{
		"newkey": "newvalue",
	}

	err = WriteVault("writevault", data)
	require.NoError(t, err)

	content, err := os.ReadFile(vaultPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "encrypted_mock_data")
}

func TestSetSecret(t *testing.T) {
	tempDir := setupMockEnvironment(t)

	// Test auto-init when vault doesn't exist
	err := SetSecret("setvault", "mykey", "myval")
	require.NoError(t, err)

	vaultPath := filepath.Join(tempDir, "setvault.yaml")
	_, err = os.Stat(vaultPath)
	assert.False(t, os.IsNotExist(err), "vault should have been auto-created")
}

func TestRemoveSecret(t *testing.T) {
	tempDir := setupMockEnvironment(t)
	vaultPath := filepath.Join(tempDir, "rmvault.yaml")
	err := os.WriteFile(vaultPath, []byte("mock"), 0600)
	require.NoError(t, err)

	err = RemoveSecret("rmvault", "key1")
	require.NoError(t, err)

	// Since we use mock sopsSetValue that returns nil, even nonexistent key removal "succeeds"
	err = RemoveSecret("rmvault", "nonexistent")
	require.NoError(t, err)
}
