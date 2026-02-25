package age

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"filippo.io/age"
	"github.com/frostyeti/sopsv/internal/config"
)

// EnsureKeyExists checks if the age keys file exists, and creates a new one if it doesn't.
// It returns the public key.
func EnsureKeyExists() (string, error) {
	keysPath := config.KeysPath

	if _, err := os.Stat(keysPath); os.IsNotExist(err) {
		identity, err := age.GenerateX25519Identity()
		if err != nil {
			return "", fmt.Errorf("failed to generate age identity: %w", err)
		}

		pubKey := identity.Recipient().String()
		keyData := fmt.Sprintf("# created by sopsv\n# public key: %s\n%s\n", pubKey, identity.String())

		if err := os.WriteFile(keysPath, []byte(keyData), 0600); err != nil {
			return "", fmt.Errorf("failed to write age key file: %w", err)
		}

		return pubKey, nil
	}

	// File exists, let's parse the first valid identity
	return GetPublicKey()
}

// GetPublicKey reads the age key file and extracts the public key
func GetPublicKey() (string, error) {
	keysPath := config.KeysPath
	file, err := os.Open(keysPath)
	if err != nil {
		return "", fmt.Errorf("failed to open age key file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "AGE-SECRET-KEY-") {
			identity, err := age.ParseX25519Identity(line)
			if err != nil {
				continue // Try next if invalid
			}
			return identity.Recipient().String(), nil
		}
	}

	return "", fmt.Errorf("no valid age identity found in %s", keysPath)
}
