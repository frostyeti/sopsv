package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/frostyeti/sopsv/internal/vault"
	"github.com/spf13/cobra"
)

var getSecretsCmd = &cobra.Command{
	Use:   "get",
	Short: "Get secrets from a vault",
	Run: func(cmd *cobra.Command, args []string) {
		vaultName := getVaultName(cmd)

		keys, _ := cmd.Flags().GetStringSlice("key")
		if len(keys) == 0 {
			color.Red("At least one key is required")
			os.Exit(1)
		}

		format, _ := cmd.Flags().GetString("format")

		data, err := vault.ReadVault(vaultName)
		if err != nil {
			color.Red("Error reading vault: %v", err)
			os.Exit(1)
		}

		result := make(map[string]interface{})
		for _, key := range keys {
			val, ok := data[key]
			if !ok {
				color.Red("Key '%s' not found in vault '%s'", key, vaultName)
				os.Exit(1)
			}
			result[key] = val
		}

		switch format {
		case "json":
			bytes, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(bytes))
		case "dotenv":
			for k, v := range result {
				fmt.Printf("%s=%v\n", k, v)
			}
		case "sh", "bash", "zsh":
			for k, v := range result {
				fmt.Printf("export %s=\"%v\"\n", k, v)
			}
		case "text":
			fallthrough
		default:
			if len(keys) == 1 {
				fmt.Println(result[keys[0]])
			} else {
				for k, v := range result {
					fmt.Printf("%s: %v\n", k, v)
				}
			}
		}
	},
}

func init() {
	secretsCmd.AddCommand(getSecretsCmd)

	getSecretsCmd.Flags().StringSliceP("key", "k", []string{}, "Key name(s)")
	getSecretsCmd.Flags().StringP("format", "f", "text", "Output format (text, dotenv, sh, json)")
}
