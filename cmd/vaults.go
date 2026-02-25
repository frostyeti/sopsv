package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/frostyeti/sopsv/internal/config"
	"github.com/frostyeti/sopsv/internal/vault"
	"github.com/spf13/cobra"
)

var vaultsCmd = &cobra.Command{
	Use:   "vaults",
	Short: "Manage sopsv vaults",
}

var newVaultCmd = &cobra.Command{
	Use:   "new [name]",
	Short: "Create a new vault",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		err := vault.NewVault(name)
		if err != nil {
			color.Red("Error creating vault: %v", err)
			return
		}
		color.Green("Vault '%s' created successfully at %s", name, vault.ResolveVaultPath(name))
	},
}

var listVaultsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all managed vaults",
	Run: func(cmd *cobra.Command, args []string) {
		vaults, err := vault.ListVaults()
		if err != nil {
			color.Red("Error listing vaults: %v", err)
			return
		}

		defaultVault := config.GetDefaultVault()
		for _, v := range vaults {
			if v == defaultVault {
				color.Green("* %s", v)
			} else {
				fmt.Println("  " + v)
			}
		}
	},
}

var useVaultCmd = &cobra.Command{
	Use:   "use [name]",
	Short: "Set the default vault",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		// Verify it exists (optional but good idea)

		err := config.SetDefaultVault(name)
		if err != nil {
			color.Red("Error setting default vault: %v", err)
			return
		}
		color.Green("Default vault set to '%s'", name)
	},
}

var showVaultCmd = &cobra.Command{
	Use:   "show [name]",
	Short: "Show the decrypted contents of a vault",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := config.GetDefaultVault()
		if len(args) > 0 {
			name = args[0]
		}

		if name == "" {
			color.Red("No vault specified and no default vault set.")
			return
		}

		data, err := vault.ReadVault(name)
		if err != nil {
			color.Red("Error reading vault: %v", err)
			return
		}

		for k, v := range data {
			fmt.Printf("%s: %v\n", k, v)
		}
	},
}

var deleteVaultCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete a vault",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		err := vault.DeleteVault(name)
		if err != nil {
			color.Red("Error deleting vault: %v", err)
			return
		}
		color.Green("Vault '%s' deleted successfully", name)

		if config.GetDefaultVault() == name {
			_ = config.SetDefaultVault("")
		}
	},
}

var editVaultCmd = &cobra.Command{
	Use:   "edit [name]",
	Short: "Edit a vault interactively in your $EDITOR",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := config.GetDefaultVault()
		if len(args) > 0 {
			name = args[0]
		}

		if name == "" {
			color.Red("No vault specified and no default vault set.")
			return
		}

		err := vault.EditVault(name)
		if err != nil {
			color.Red("Error editing vault: %v", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(vaultsCmd)
	vaultsCmd.AddCommand(newVaultCmd)
	vaultsCmd.AddCommand(listVaultsCmd)
	vaultsCmd.AddCommand(useVaultCmd)
	vaultsCmd.AddCommand(showVaultCmd)
	vaultsCmd.AddCommand(editVaultCmd)
	vaultsCmd.AddCommand(deleteVaultCmd)
}
