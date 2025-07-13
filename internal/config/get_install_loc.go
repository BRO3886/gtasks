package config

import (
	"os"
	"path/filepath"

	"github.com/BRO3886/gtasks/internal/utils"
)

// GetInstallLocation returns the ~/.gtasks directory for storing config and tokens
func GetInstallLocation() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		utils.ErrorP("Get home directory: %s", err.Error())
		return ".gtasks" // fallback to current directory
	}
	
	configDir := filepath.Join(homeDir, ".gtasks")
	
	// Create directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		utils.ErrorP("Create config directory: %s", err.Error())
		return homeDir // fallback to home directory
	}
	
	return configDir
}
