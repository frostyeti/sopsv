package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/frostyeti/sopsv/internal/config"
	"github.com/frostyeti/sopsv/internal/vault"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get [key]",
	Short: "Get a specific key from the default or specified vault",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		vaultName, _ := cmd.Flags().GetString("vault")

		if vaultName == "" {
			vaultName = config.GetDefaultVault()
		}

		if vaultName == "" {
			color.Red("No vault specified and no default vault set.")
			os.Exit(1)
		}

		data, err := vault.ReadVault(vaultName)
		if err != nil {
			color.Red("Error reading vault: %v", err)
			os.Exit(1)
		}

		val, ok := data[key]
		if !ok {
			color.Red("Key '%s' not found in vault '%s'", key, vaultName)
			os.Exit(1)
		}

		fmt.Println(val)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().StringP("vault", "v", "", "Vault to read from")
}
