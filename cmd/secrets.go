package cmd

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
	"github.com/frostyeti/sopsv/internal/config"
	"github.com/frostyeti/sopsv/internal/vault"
	"github.com/spf13/cobra"
)

var secretsCmd = &cobra.Command{
	Use:   "secrets",
	Short: "Manage secrets within a vault",
}

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set a secret in a vault",
	Run: func(cmd *cobra.Command, args []string) {
		vaultName := getVaultName(cmd)

		key, _ := cmd.Flags().GetString("key")
		if key == "" {
			color.Red("Key is required")
			os.Exit(1)
		}

		var value string

		generate, _ := cmd.Flags().GetBool("generate")
		useStdin, _ := cmd.Flags().GetBool("stdin")

		if generate {
			size, _ := cmd.Flags().GetInt("size")
			b := make([]byte, size)
			_, err := rand.Read(b)
			if err != nil {
				color.Red("Failed to generate secret: %v", err)
				os.Exit(1)
			}
			value = base64.RawURLEncoding.EncodeToString(b)[:size]
		} else if useStdin {
			bytes, err := io.ReadAll(os.Stdin)
			if err != nil {
				color.Red("Failed to read from stdin: %v", err)
				os.Exit(1)
			}
			value = string(bytes)
		} else {
			valFlag, _ := cmd.Flags().GetString("value")
			if valFlag != "" {
				value = valFlag
			} else {
				color.Red("Value, --generate, or --stdin is required")
				os.Exit(1)
			}
		}

		err := vault.SetSecret(vaultName, key, value)
		if err != nil {
			color.Red("Error setting secret: %v", err)
			os.Exit(1)
		}

		color.Green("Secret '%s' set successfully in vault '%s'", key, vaultName)
	},
}

var rmCmd = &cobra.Command{
	Use:     "rm",
	Aliases: []string{"remove", "delete"},
	Short:   "Remove a secret from a vault",
	Run: func(cmd *cobra.Command, args []string) {
		vaultName := getVaultName(cmd)

		key, _ := cmd.Flags().GetString("key")
		if key == "" {
			color.Red("Key is required")
			os.Exit(1)
		}

		err := vault.RemoveSecret(vaultName, key)
		if err != nil {
			color.Red("Error removing secret: %v", err)
			os.Exit(1)
		}

		color.Green("Secret '%s' removed from vault '%s'", key, vaultName)
	},
}

var ensureCmd = &cobra.Command{
	Use:   "ensure",
	Short: "Get a secret, or generate and save it if it does not exist",
	Run: func(cmd *cobra.Command, args []string) {
		vaultName := getVaultName(cmd)

		key, _ := cmd.Flags().GetString("key")
		if key == "" {
			color.Red("Key is required")
			os.Exit(1)
		}

		data, err := vault.ReadVault(vaultName)
		if err == nil {
			if val, ok := data[key]; ok {
				fmt.Println(val)
				return
			}
		}

		size, _ := cmd.Flags().GetInt("size")
		b := make([]byte, size)
		_, err = rand.Read(b)
		if err != nil {
			color.Red("Failed to generate secret: %v", err)
			os.Exit(1)
		}
		value := base64.RawURLEncoding.EncodeToString(b)[:size]

		err = vault.SetSecret(vaultName, key, value)
		if err != nil {
			color.Red("Error ensuring secret: %v", err)
			os.Exit(1)
		}

		fmt.Println(value)
	},
}

var lsCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list"},
	Short:   "List all secret keys in a vault",
	Run: func(cmd *cobra.Command, args []string) {
		vaultName := getVaultName(cmd)

		keys, err := vault.ListKeys(vaultName)
		if err != nil {
			color.Red("Error reading vault: %v", err)
			os.Exit(1)
		}

		for _, key := range keys {
			fmt.Println(key)
		}
	},
}

func getVaultName(cmd *cobra.Command) string {
	vaultName, _ := cmd.Flags().GetString("vault")
	if vaultName == "" {
		vaultName = config.GetDefaultVault()
	}
	if vaultName == "" {
		color.Red("No vault specified and no default vault set. Use --vault or set a default vault.")
		os.Exit(1)
	}
	return vaultName
}

func init() {
	rootCmd.AddCommand(secretsCmd)

	secretsCmd.AddCommand(setCmd)
	secretsCmd.AddCommand(rmCmd)
	secretsCmd.AddCommand(ensureCmd)
	secretsCmd.AddCommand(lsCmd)

	// Map core commands to root
	rootCmd.AddCommand(setCmd)
	rootCmd.AddCommand(rmCmd)
	rootCmd.AddCommand(lsCmd)

	setCmd.Flags().StringP("key", "k", "", "Key name")
	setCmd.Flags().StringP("value", "v", "", "Secret value")
	setCmd.Flags().Bool("generate", false, "Generate a random value")
	setCmd.Flags().Int("size", 16, "Size of generated secret")
	setCmd.Flags().Bool("stdin", false, "Read value from stdin")

	rmCmd.Flags().StringP("key", "k", "", "Key name")

	ensureCmd.Flags().StringP("key", "k", "", "Key name")
	ensureCmd.Flags().Int("size", 16, "Size of generated secret if it doesn't exist")
}
