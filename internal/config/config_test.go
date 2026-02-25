package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserDataDir(t *testing.T) {
	// Simple test to ensure UserDataDir doesn't panic and returns a string
	dir, err := UserDataDir()
	require.NoError(t, err)
	assert.NotEmpty(t, dir)
}

func TestConfigInitAndVault(t *testing.T) {
	// Save original env vars
	origConfigDir := os.Getenv("XDG_CONFIG_HOME")
	origDataDir := os.Getenv("XDG_DATA_HOME")
	defer func() {
		os.Setenv("XDG_CONFIG_HOME", origConfigDir)
		os.Setenv("XDG_DATA_HOME", origDataDir)
	}()

	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, "config")
	dataDir := filepath.Join(tempDir, "data")

	// Override paths for testing
	os.Setenv("XDG_CONFIG_HOME", configDir)
	os.Setenv("XDG_DATA_HOME", dataDir)

	err := InitConfig()
	require.NoError(t, err)

	// Verify paths
	assert.Contains(t, AppDir, "sopsv")
	assert.Contains(t, DataDir, "sopsv")
	assert.Equal(t, filepath.Join(AppDir, "config.yaml"), ConfigPath)
	assert.Equal(t, filepath.Join(AppDir, "keys.txt"), KeysPath)

	// Verify env var
	assert.Equal(t, KeysPath, os.Getenv("SOPS_AGE_KEY_FILE"))

	// Test Set/Get Default Vault
	err = SetDefaultVault("test_vault")
	require.NoError(t, err)

	vault := GetDefaultVault()
	assert.Equal(t, "test_vault", vault)

	// Re-init should read the saved config
	err = InitConfig()
	require.NoError(t, err)
	assert.Equal(t, "test_vault", GetDefaultVault())
}
