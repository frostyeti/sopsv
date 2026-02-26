package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

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

		result := make(map[string]interface{})
		for _, key := range keys {
			val, err := vault.GetSecret(vaultName, key)
			if err != nil {
				// if extract fails, it's usually because the key isn't there or decryption failed
				color.Red("Key '%s' not found or error extracting in vault '%s': %v", key, vaultName, err)
				os.Exit(1)
			}
			// Sops extract returns the value as a JSON string or literal string
			// We will try to unmarshal it as JSON, or keep as string
			var unmarshaled interface{}
			if err := json.Unmarshal([]byte(val), &unmarshaled); err == nil {
				result[key] = unmarshaled
			} else {
				// sops extract might just print the unquoted string depending on version,
				// or maybe we just want to trim spaces if it's not valid json
				result[key] = strings.TrimSpace(val)
			}
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
	rootCmd.AddCommand(getSecretsCmd)

	getSecretsCmd.Flags().StringSliceP("key", "k", []string{}, "Key name(s)")
	getSecretsCmd.Flags().StringP("format", "f", "text", "Output format (text, dotenv, sh, json)")
}
