package cmd

import (
	"os"

	"github.com/frostyeti/sopsv/internal/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sopsv",
	Short: "A local secret vault CLI using sops and age",
	Long: `sopsv is a CLI that uses sops and age as a secret vault locally.
It encrypts and decrypts simple key-value pairs stored in YAML, JSON, or .env files.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return config.InitConfig()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Flags will be added here
}
