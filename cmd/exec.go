package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/fatih/color"
	"github.com/frostyeti/sopsv/internal/config"
	"github.com/frostyeti/sopsv/internal/vault"
	"github.com/spf13/cobra"
)

var execCmd = &cobra.Command{
	Use:   "exec [command] [args...]",
	Short: "Execute a command with secrets from vaults as environment variables",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		vaults, _ := cmd.Flags().GetStringSlice("vaults")

		if len(vaults) == 0 {
			defVault := config.GetDefaultVault()
			if defVault != "" {
				vaults = []string{defVault}
			}
		}

		env := os.Environ()

		// Read each vault and append to env
		for _, v := range vaults {
			data, err := vault.ReadVault(v)
			if err != nil {
				color.Red("Warning: Error reading vault %s: %v", v, err)
				continue
			}

			for key, val := range data {
				// simple string representation
				env = append(env, fmt.Sprintf("%s=%v", key, val))
			}
		}

		command := exec.Command(args[0], args[1:]...)
		command.Env = env
		command.Stdin = os.Stdin
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr

		if err := command.Run(); err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				os.Exit(exitError.ExitCode())
			}
			color.Red("Error executing command: %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
	execCmd.Flags().StringSliceP("vaults", "v", []string{}, "Comma separated list of vaults to use")
}
