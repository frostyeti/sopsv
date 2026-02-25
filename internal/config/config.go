package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

var (
	AppDir       string
	ConfigPath   string
	KeysPath     string
	DataDir      string
	DefaultViper *viper.Viper
)

func UserDataDir() (string, error) {
	switch runtime.GOOS {
	case "windows":
		dir := os.Getenv("LocalAppData")
		if dir == "" {
			return "", fmt.Errorf("%%LocalAppData%% is not defined")
		}
		return dir, nil
	case "darwin":
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, "Library", "Application Support"), nil
	default: // unix
		dir := os.Getenv("XDG_DATA_HOME")
		if dir == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			return filepath.Join(home, ".local", "share"), nil
		}
		return dir, nil
	}
}

func InitConfig() error {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get user config dir: %w", err)
	}
	AppDir = filepath.Join(userConfigDir, "sopsv")
	ConfigPath = filepath.Join(AppDir, "config.yaml")
	KeysPath = filepath.Join(AppDir, "keys.txt")

	userDataDir, err := UserDataDir()
	if err != nil {
		return fmt.Errorf("failed to get user data dir: %w", err)
	}
	DataDir = filepath.Join(userDataDir, "sopsv")

	// Ensure directories exist
	if err := os.MkdirAll(AppDir, 0700); err != nil {
		return fmt.Errorf("failed to create app config dir: %w", err)
	}
	if err := os.MkdirAll(DataDir, 0700); err != nil {
		return fmt.Errorf("failed to create app data dir: %w", err)
	}

	DefaultViper = viper.New()
	DefaultViper.SetConfigFile(ConfigPath)
	DefaultViper.SetConfigType("yaml")

	if err := DefaultViper.ReadInConfig(); err != nil {
		if os.IsNotExist(err) {
			// Write an empty default config if not exists
			DefaultViper.Set("default_vault", "")
			if writeErr := DefaultViper.WriteConfigAs(ConfigPath); writeErr != nil {
				return fmt.Errorf("failed to write default config: %w", writeErr)
			}
		} else {
			return fmt.Errorf("failed to read config: %w", err)
		}
	}

	// Make keys available via env var for SOPS
	os.Setenv("SOPS_AGE_KEY_FILE", KeysPath)

	return nil
}

// GetDefaultVault gets the default vault from config
func GetDefaultVault() string {
	return DefaultViper.GetString("default_vault")
}

// SetDefaultVault sets the default vault in config
func SetDefaultVault(vault string) error {
	DefaultViper.Set("default_vault", vault)
	return DefaultViper.WriteConfig()
}
