//go:build unix
// +build unix

package envpath

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

// EnsurePath makes sure the given directory is in the user's PATH environment variable.
// On Unix, this means checking ~/.bashrc and ~/.zshrc and appending it if missing.
func EnsurePath(dir string, global bool) error {
	// If it's already in the PATH environment, we might not need to do anything.
	// But appending to profiles ensures it's persistent across sessions.

	if global {
		// global usually means /usr/local/bin which is almost always in PATH.
		color.Yellow("Note: Assuming %s is already in system PATH.", dir)
		return nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	profiles := []string{
		filepath.Join(home, ".bashrc"),
		filepath.Join(home, ".zshrc"),
		filepath.Join(home, ".profile"),
	}

	exportLine := fmt.Sprintf("export PATH=\"$PATH:%s\"", dir)
	added := false

	for _, profile := range profiles {
		if _, err := os.Stat(profile); err == nil {
			if !containsLine(profile, dir) {
				if err := appendToFile(profile, exportLine); err != nil {
					color.Red("Failed to update %s: %v", profile, err)
				} else {
					color.Green("Updated %s with %s", profile, dir)
					added = true
				}
			} else {
				added = true // already there
			}
		}
	}

	if added {
		color.Cyan("Please restart your terminal or run `source ~/.bashrc` (or equivalent) to apply PATH changes.")
	}

	return nil
}

func containsLine(filePath, searchStr string) bool {
	f, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), searchStr) {
			return true
		}
	}
	return false
}

func appendToFile(filePath, text string) error {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString("\n# Added by sopsv installer\n" + text + "\n")
	return err
}
