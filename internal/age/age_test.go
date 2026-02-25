package age

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/frostyeti/sopsv/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnsureKeyExists(t *testing.T) {
	tempDir := t.TempDir()
	config.KeysPath = filepath.Join(tempDir, "keys.txt")

	// 1. Should generate a new key when file doesn't exist
	pubKey1, err := EnsureKeyExists()
	require.NoError(t, err)
	assert.True(t, strings.HasPrefix(pubKey1, "age1"))

	// Verify file was written
	content, err := os.ReadFile(config.KeysPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), pubKey1)
	assert.Contains(t, string(content), "AGE-SECRET-KEY-")

	// 2. Should return existing key when file exists
	pubKey2, err := EnsureKeyExists()
	require.NoError(t, err)
	assert.Equal(t, pubKey1, pubKey2)
}

func TestGetPublicKey(t *testing.T) {
	tempDir := t.TempDir()
	config.KeysPath = filepath.Join(tempDir, "keys.txt")

	// Should fail if no file
	_, err := GetPublicKey()
	require.Error(t, err)

	// Write mock key
	mockKey := "AGE-SECRET-KEY-1H3H82R29GZ46XQU7H9QJ8Q2P680YV34W439WQZQ9Y9Z4QY848GXS0T2"
	err = os.WriteFile(config.KeysPath, []byte(mockKey), 0600)
	require.NoError(t, err)

	// Will fail because the fake key is not actually a valid age key
	// (ParseX25519Identity will fail)
	_, err = GetPublicKey()
	require.Error(t, err)
}
