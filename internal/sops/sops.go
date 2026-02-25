package sops

import (
	"fmt"

	"github.com/getsops/sops/v3/decrypt"
)

// DecryptFile decrypts a sops encrypted file and returns the decrypted content
func DecryptFile(filePath string) ([]byte, error) {
	cleartext, err := decrypt.File(filePath, "yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt file %s: %w", filePath, err)
	}
	return cleartext, nil
}
