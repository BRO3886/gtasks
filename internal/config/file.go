package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"

	"github.com/BRO3886/gtasks/internal/utils"
)

// k is the package-level koanf instance loaded by LoadAppConfig.
var k = koanf.New(".")

// LoadAppConfig loads configuration using the following priority order (highest first):
//
//  1. Environment variables (GTASKS_* prefix, stripped before loading)
//  2. .env file in the config directory
//  3. Config file: config.toml / config.yaml / config.json in the config directory
//
// A missing config file or .env file is silently ignored.
// Malformed files log a warning and fall through to lower-priority sources.
func LoadAppConfig() {
	cfgDir := GetInstallLocation()

	// 3. Config file (lowest priority — loaded first, overridden by layers above)
	for _, candidate := range []struct {
		name   string
		parser koanf.Parser
	}{
		{"config.toml", toml.Parser()},
		{"config.yaml", yaml.Parser()},
		{"config.yml", yaml.Parser()},
		{"config.json", json.Parser()},
	} {
		cfgPath := filepath.Join(cfgDir, candidate.name)
		if _, err := os.Stat(cfgPath); err == nil {
			if err := k.Load(file.Provider(cfgPath), candidate.parser); err != nil {
				utils.Warn("Could not parse config file %s: %v\n", cfgPath, err)
			}
			break // use the first one found
		}
	}

	// 2. .env file
	dotenvPath := filepath.Join(cfgDir, ".env")
	if err := godotenv.Load(dotenvPath); err != nil && !os.IsNotExist(err) {
		utils.Warn("Could not read .env file %s: %v\n", dotenvPath, err)
	}

	// 1. Environment variables — GTASKS_ prefix, mapped to dotted keys
	// e.g. GTASKS_CLIENT_ID -> credentials.client_id
	//      GTASKS_DEFAULT_TASKLIST -> tasks.default_task_list
	k.Load(env.Provider("GTASKS_", ".", func(s string) string {
		s = strings.TrimPrefix(s, "GTASKS_")
		s = strings.ToLower(s)
		switch s {
		case "client_id":
			return "credentials.client_id"
		case "client_secret":
			return "credentials.client_secret"
		case "default_tasklist":
			return "tasks.default_task_list"
		}
		return s
	}), nil)
}

// GetDefaultTaskList returns the default task list from config/env, or empty string.
func GetDefaultTaskList() string {
	return k.String("tasks.default_task_list")
}

// GetCredentials returns client ID and secret from config/env.
func GetCredentials() (clientID, clientSecret string) {
	return k.String("credentials.client_id"), k.String("credentials.client_secret")
}
