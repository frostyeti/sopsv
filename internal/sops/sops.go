package sops

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
)

// TODO: write only the portions of sops that we need directly in go,
// so we don't have to shell out to the large sops binary.

// DecryptFile decrypts a sops encrypted file and returns the decrypted content
func DecryptFile(filePath string) ([]byte, error) {
	cmd := exec.Command("sops", "-d", filePath)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to decrypt file %s: %s: %w", filePath, stderr.String(), err)
	}

	return out.Bytes(), nil
}

// ExtractValue decrypts and extracts a specific key from a sops encrypted file
func ExtractValue(filePath string, key string) (string, error) {
	extractArg := fmt.Sprintf(`["%s"]`, key)
	cmd := exec.Command("sops", "-d", "--extract", extractArg, filePath)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to extract key %s from %s: %s: %w", key, filePath, stderr.String(), err)
	}

	return out.String(), nil
}

// SetValue sets a specific key to a value in a sops encrypted file
func SetValue(filePath string, key string, value interface{}) error {
	keyArg := fmt.Sprintf(`["%s"]`, key)

	valBytes, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to encode value as JSON: %w", err)
	}
	valueArg := string(valBytes)

	cmd := exec.Command("sops", "set", filePath, keyArg, valueArg)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to set key %s in %s: %s: %w", key, filePath, stderr.String(), err)
	}
	return nil
}
