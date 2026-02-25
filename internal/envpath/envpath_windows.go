//go:build windows
// +build windows

package envpath

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"golang.org/x/sys/windows/registry"
)

// EnsurePath makes sure the given directory is in the user's or system PATH environment variable.
func EnsurePath(dir string, global bool) error {
	var k registry.Key
	var err error

	if global {
		// System PATH
		k, err = registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\Session Manager\Environment`, registry.QUERY_VALUE|registry.SET_VALUE)
	} else {
		// User PATH
		k, err = registry.OpenKey(registry.CURRENT_USER, `Environment`, registry.QUERY_VALUE|registry.SET_VALUE)
	}

	if err != nil {
		return fmt.Errorf("failed to open registry key: %v", err)
	}
	defer k.Close()

	pathVal, _, err := k.GetStringValue("Path")
	if err != nil {
		return fmt.Errorf("failed to get Path value: %v", err)
	}

	paths := strings.Split(pathVal, ";")
	for _, p := range paths {
		if strings.EqualFold(strings.TrimSpace(p), dir) {
			color.Yellow("Directory %s is already in PATH.", dir)
			return nil
		}
	}

	newPath := pathVal
	if !strings.HasSuffix(newPath, ";") {
		newPath += ";"
	}
	newPath += dir

	err = k.SetStringValue("Path", newPath)
	if err != nil {
		return fmt.Errorf("failed to set Path value: %v", err)
	}

	color.Green("Successfully added %s to PATH.", dir)
	color.Cyan("You may need to restart your terminal or computer for changes to take effect.")

	return nil
}
