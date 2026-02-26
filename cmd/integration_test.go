//go:build integration
// +build integration

package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/frostyeti/sopsv/internal/config"
	"github.com/stretchr/testify/require"
)

// We won't test cobra commands directly executing as it can be tricky to capture output,
// but we will test the actual end-to-end flow with the real vault/sops/age packages.
// Since it's an integration test, we will touch real files and invoke real sops/age binaries
// (assuming they are installed in the environment).

func TestIntegration_VaultAndSecretsLifecycle(t *testing.T) {
	// 1. Setup mock directories
	tempDir := t.TempDir()

	// Set XDG dirs to point to tempDir so real config uses them
	os.Setenv("XDG_CONFIG_HOME", filepath.Join(tempDir, "config"))
	os.Setenv("XDG_DATA_HOME", filepath.Join(tempDir, "data"))

	err := config.InitConfig()
	require.NoError(t, err)

	// In real environment, age keys will be generated automatically in tempDir/config/sopsv/keys.txt

	vaultName := "test_integration_vault"

	// Ensure we're starting fresh
	vaultPath := filepath.Join(config.DataDir, vaultName+".yaml")
	_, err = os.Stat(vaultPath)
	require.True(t, os.IsNotExist(err), "vault should not exist yet")

	// Create a new vault (this will also generate age keys)
	err = ExecuteCommand("vault", "new", vaultName)
	require.NoError(t, err, "vault creation should succeed")

	// Ensure file was created
	_, err = os.Stat(vaultPath)
	require.NoError(t, err, "vault file should be created")

	// Test Set Secret
	err = ExecuteCommand("set", "--key", "DB_USER", "--value", "admin", "--vault", vaultName)
	require.NoError(t, err, "setting secret should succeed")

	err = ExecuteCommand("secrets", "set", "--key", "DB_PASS", "--value", "supersecret", "--vault", vaultName)
	require.NoError(t, err, "setting another secret should succeed")

	// We can't easily capture ExecuteCommand stdout in a simple way without redirecting os.Stdout,
	// so we will just test the internal ReadVault directly to verify it was set properly
	// Or we can mock os.Stdout

	// Let's test the exec command
	// We'll run a shell script that prints the env vars
	scriptPath := filepath.Join(tempDir, "test.sh")
	os.WriteFile(scriptPath, []byte("#!/bin/sh\necho \"USER=$DB_USER PASS=$DB_PASS\""), 0755)

	// Since ExecuteCommand does not return output, we can't assert the output of `exec` command easily.
	// But we can check that it doesn't return an error.
	err = ExecuteCommand("exec", "--vaults", vaultName, "--", scriptPath)
	require.NoError(t, err, "exec command should succeed")

	// Test Remove Secret
	err = ExecuteCommand("rm", "--key", "DB_USER", "--vault", vaultName)
	require.NoError(t, err, "removing secret should succeed")

	// Test Delete Vault
	err = ExecuteCommand("vault", "rm", vaultName)
	require.NoError(t, err, "deleting vault should succeed")

	_, err = os.Stat(vaultPath)
	require.True(t, os.IsNotExist(err), "vault file should be deleted")
}

// ExecuteCommand is a helper function to run the root command with given args
func ExecuteCommand(args ...string) error {
	rootCmd.SetArgs(args)
	return rootCmd.Execute()
}
