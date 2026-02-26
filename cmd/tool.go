package cmd

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/fatih/color"
	"github.com/frostyeti/sopsv/internal/envpath"
	"github.com/frostyeti/sopsv/internal/installer"
	"github.com/spf13/cobra"
)

var toolCmd = &cobra.Command{
	Use:   "tool",
	Short: "Manage dependencies for sopsv (sops and age)",
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the latest versions of sops and age",
	Run: func(cmd *cobra.Command, args []string) {
		global, _ := cmd.Flags().GetBool("global")
		dest, _ := cmd.Flags().GetString("dest")

		var targetDir string

		if dest != "" {
			targetDir = dest
		} else {
			if runtime.GOOS == "windows" {
				if global {
					targetDir = `C:\Program Files\bin`
				} else {
					home, _ := os.UserHomeDir()
					targetDir = filepath.Join(home, "AppData", "Local", "Programs", "bin")
				}
			} else {
				if global {
					targetDir = "/usr/local/bin"
				} else {
					home, _ := os.UserHomeDir()
					targetDir = filepath.Join(home, ".local", "bin")
				}
			}
		}

		if err := os.MkdirAll(targetDir, 0755); err != nil {
			color.Red("Failed to create destination directory %s: %v", targetDir, err)
			os.Exit(1)
		}

		color.Cyan("Installing tools to %s", targetDir)

		if err := installer.InstallSops(targetDir); err != nil {
			color.Red("Error installing sops: %v", err)
			os.Exit(1)
		}

		if err := installer.InstallAge(targetDir); err != nil {
			color.Red("Error installing age: %v", err)
			os.Exit(1)
		}

		if err := envpath.EnsurePath(targetDir, global); err != nil {
			color.Yellow("Warning: Could not automatically update PATH: %v", err)
			color.Yellow("Please manually ensure %s is in your PATH.", targetDir)
		}
	},
}

func init() {
	rootCmd.AddCommand(toolCmd)
	toolCmd.AddCommand(installCmd)

	installCmd.Flags().BoolP("global", "g", false, "Install globally")
	installCmd.Flags().StringP("dest", "d", "", "Custom destination directory")
}
