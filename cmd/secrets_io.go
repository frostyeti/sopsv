package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/frostyeti/sopsv/internal/vault"
	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export a vault's contents as JSON",
	Run: func(cmd *cobra.Command, args []string) {
		vaultName := getVaultName(cmd)

		data, err := vault.ReadVault(vaultName)
		if err != nil {
			color.Red("Error reading vault: %v", err)
			os.Exit(1)
		}

		pretty, _ := cmd.Flags().GetBool("pretty")
		file, _ := cmd.Flags().GetString("file")

		var bytes []byte
		if pretty {
			bytes, err = json.MarshalIndent(data, "", "  ")
		} else {
			bytes, err = json.Marshal(data)
		}

		if err != nil {
			color.Red("Error exporting vault data: %v", err)
			os.Exit(1)
		}

		if file != "" {
			err = os.WriteFile(file, bytes, 0600)
			if err != nil {
				color.Red("Error writing to file: %v", err)
				os.Exit(1)
			}
			color.Green("Exported successfully to %s", file)
		} else {
			fmt.Println(string(bytes))
		}
	},
}

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import secrets from a JSON file into a vault",
	Run: func(cmd *cobra.Command, args []string) {
		vaultName := getVaultName(cmd)

		file, _ := cmd.Flags().GetString("file")
		if file == "" {
			color.Red("--file flag is required")
			os.Exit(1)
		}

		bytes, err := os.ReadFile(file)
		if err != nil {
			color.Red("Error reading file: %v", err)
			os.Exit(1)
		}

		var importedData map[string]interface{}
		err = json.Unmarshal(bytes, &importedData)
		if err != nil {
			color.Red("Error parsing JSON: %v", err)
			os.Exit(1)
		}

		// Ensure vault exists
		data, err := vault.ReadVault(vaultName)
		if err != nil {
			err = vault.NewVault(vaultName)
			if err != nil {
				color.Red("Error initializing new vault: %v", err)
				os.Exit(1)
			}
			data = make(map[string]interface{})
		}

		// Merge
		for k, v := range importedData {
			data[k] = v
		}

		err = vault.WriteVault(vaultName, data)
		if err != nil {
			color.Red("Error writing vault: %v", err)
			os.Exit(1)
		}

		color.Green("Imported %d secrets into vault '%s'", len(importedData), vaultName)
	},
}

func init() {
	secretsCmd.AddCommand(exportCmd)
	secretsCmd.AddCommand(importCmd)

	exportCmd.Flags().Bool("json", true, "Export as JSON (default)")
	exportCmd.Flags().Bool("pretty", false, "Pretty print JSON output")
	exportCmd.Flags().StringP("file", "f", "", "Output file path (default stdout)")

	importCmd.Flags().Bool("json", true, "Import from JSON format (default)")
	importCmd.Flags().StringP("file", "f", "", "Input JSON file path")
}
