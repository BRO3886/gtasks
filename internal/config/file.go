package config

import (
	"os"
	"path/filepath"

	"github.com/BRO3886/gtasks/internal/utils"
	"github.com/pelletier/go-toml"
)

// AppConfig holds settings loaded from the gtasks configuration file (config.toml).
//
// The file is looked up in GetInstallLocation() — typically
// $XDG_CONFIG_HOME/gtasks/config.toml (or ~/.gtasks/config.toml for legacy installs).
//
// Environment variables and CLI flags override values set here; see each field for details.
type AppConfig struct {
	// Credentials holds OAuth2 client credentials for the Google Tasks API.
	// These are only needed when building or running gtasks with your own
	// Google Cloud project; official releases have credentials embedded at build time.
	Credentials struct {
		// ClientID is the Google OAuth2 client ID.
		// Overridden by the GTASKS_CLIENT_ID environment variable.
		ClientID string `toml:"client_id"`

		// ClientSecret is the Google OAuth2 client secret.
		// Overridden by the GTASKS_CLIENT_SECRET environment variable.
		ClientSecret string `toml:"client_secret"`
	} `toml:"credentials"`

	// Tasks holds preferences for task operations.
	Tasks struct {
		// DefaultTaskList is the task list name used when the -l flag is not provided.
		// Overridden by the GTASKS_DEFAULT_TASKLIST environment variable, then by -l flag.
		DefaultTaskList string `toml:"default_task_list"`
	} `toml:"tasks"`
}

// appCfg is the package-level config loaded by LoadAppConfig.
var appCfg AppConfig

// LoadAppConfig reads config.toml from the gtasks configuration directory and stores
// the result in the package-level appCfg variable. A missing config file is silently
// ignored; a malformed file logs a warning and uses zero-value defaults.
func LoadAppConfig() {
	cfgPath := filepath.Join(GetInstallLocation(), "config.toml")

	data, err := os.ReadFile(cfgPath)
	if err != nil {
		if !os.IsNotExist(err) {
			utils.Warn("Could not read config file %s: %v\n", cfgPath, err)
		}
		return
	}

	if err := toml.Unmarshal(data, &appCfg); err != nil {
		utils.Warn("Could not parse config file %s: %v\n", cfgPath, err)
	}
}

// GetAppConfig returns the loaded application configuration.
func GetAppConfig() AppConfig {
	return appCfg
}

// GetDefaultTaskList returns the default task list name from the config file,
// or an empty string if none is set.
func GetDefaultTaskList() string {
	return appCfg.Tasks.DefaultTaskList
}
