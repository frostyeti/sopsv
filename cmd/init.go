package cmd

import (
	"github.com/fatih/color"
	"github.com/frostyeti/sopsv/internal/config"
	"github.com/frostyeti/sopsv/internal/vault"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize sopsv with a default vault",
	Long:  "Creates a 'default' vault under the user's profile and configures it to be the default.",
	Run: func(cmd *cobra.Command, args []string) {
		name := "default"

		// Attempt to create the vault
		err := vault.NewVault(name)
		if err != nil {
			// If it already exists, that's fine, we will just set it as default
			if err.Error() != "" {
				color.Yellow("Vault '%s' already exists or could not be created: %v", name, err)
			}
		} else {
			color.Green("Vault '%s' created successfully at %s", name, vault.ResolveVaultPath(name))
		}

		// Set it as default
		err = config.SetDefaultVault(name)
		if err != nil {
			color.Red("Error setting default vault: %v", err)
			return
		}

		color.Green("Default vault set to '%s'", name)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
