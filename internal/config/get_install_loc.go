package config

import (
	"os"
	"path/filepath"

	"github.com/BRO3886/gtasks/internal/utils"
)

// GetInstallLocation returns the directory used for gtasks configuration and token storage.
//
// Discovery order (first existing directory wins):
//  1. $XDG_CONFIG_HOME/gtasks/  (or ~/.config/gtasks/ when XDG_CONFIG_HOME is unset)
//  2. ~/.gtasks/                 (legacy path, kept for backward compatibility)
//
// If neither directory exists, $XDG_CONFIG_HOME/gtasks/ is created for new installations.
// The legacy path ~/.gtasks/ is never created for new installs; only existing installs continue
// using it.
func GetInstallLocation() string {
	xdgDir := xdgConfigDir()
	legacyDir := legacyConfigDir()

	// Prefer XDG dir if it already exists
	if xdgDir != "" {
		if _, err := os.Stat(xdgDir); err == nil {
			return xdgDir
		}
	}

	// Fall back to legacy ~/.gtasks if it exists
	if legacyDir != "" {
		if _, err := os.Stat(legacyDir); err == nil {
			return legacyDir
		}
	}

	// Neither exists — create XDG dir for new installations
	if xdgDir != "" {
		if err := os.MkdirAll(xdgDir, 0755); err == nil {
			return xdgDir
		}
		utils.ErrorP("Create XDG config directory %s: %s", xdgDir, "failed")
	}

	// Final fallback: create legacy dir
	if legacyDir != "" {
		if err := os.MkdirAll(legacyDir, 0755); err == nil {
			return legacyDir
		}
		utils.ErrorP("Create config directory %s: %s", legacyDir, "failed")
	}

	return ".gtasks" // unreachable in practice
}

// xdgConfigDir returns $XDG_CONFIG_HOME/gtasks, or ~/.config/gtasks if XDG_CONFIG_HOME is
// unset or empty (per XDG Base Directory Specification). Returns empty string if the home
// directory cannot be determined.
func xdgConfigDir() string {
	base := os.Getenv("XDG_CONFIG_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		base = filepath.Join(home, ".config")
	}
	return filepath.Join(base, "gtasks")
}

// legacyConfigDir returns ~/.gtasks, or empty string if the home directory cannot be
// determined.
func legacyConfigDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".gtasks")
}
