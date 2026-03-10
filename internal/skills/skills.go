package skills

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// AgentTarget represents an AI coding agent that supports the Agent Skills standard.
type AgentTarget struct {
	Name    string
	Key     string
	BaseDir string
}

// DefaultTargets returns the supported agent targets with expanded home directory paths.
func DefaultTargets(homeDir string) []AgentTarget {
	return []AgentTarget{
		{
			Name:    "Claude Code",
			Key:     "claude",
			BaseDir: filepath.Join(homeDir, ".claude", "skills"),
		},
		{
			Name:    "Codex CLI",
			Key:     "codex",
			BaseDir: filepath.Join(homeDir, ".agents", "skills"),
		},
		{
			Name:    "OpenClaw",
			Key:     "openclaw",
			BaseDir: filepath.Join(homeDir, ".openclaw", "skills"),
		},
	}
}

const skillDirName = "gtasks-cli"
const versionFileName = ".gtasks-version"

// SkillDir returns the full path to the skill directory for a given target.
func SkillDir(target AgentTarget) string {
	return filepath.Join(target.BaseDir, skillDirName)
}

// Install writes embedded skill files to the target's skill directory.
func Install(embeddedFS fs.FS, target AgentTarget, version string) ([]string, error) {
	destDir := SkillDir(target)

	if err := os.RemoveAll(destDir); err != nil {
		return nil, fmt.Errorf("failed to remove existing skill directory: %w", err)
	}

	var written []string

	err := fs.WalkDir(embeddedFS, "assets/gtasks-cli", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel("assets/gtasks-cli", path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(destDir, relPath)

		if d.IsDir() {
			return os.MkdirAll(destPath, 0o755)
		}

		data, err := fs.ReadFile(embeddedFS, path)
		if err != nil {
			return fmt.Errorf("failed to read embedded file %s: %w", path, err)
		}

		if err := os.MkdirAll(filepath.Dir(destPath), 0o755); err != nil {
			return fmt.Errorf("failed to create directory for %s: %w", destPath, err)
		}

		if err := os.WriteFile(destPath, data, 0o644); err != nil {
			return fmt.Errorf("failed to write %s: %w", destPath, err)
		}

		written = append(written, relPath)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to install skills: %w", err)
	}

	versionPath := filepath.Join(destDir, versionFileName)
	if err := os.WriteFile(versionPath, []byte(version+"\n"), 0o644); err != nil {
		return nil, fmt.Errorf("failed to write version file: %w", err)
	}

	return written, nil
}

// Uninstall removes the skill directory from the target's location.
func Uninstall(target AgentTarget) (bool, error) {
	destDir := SkillDir(target)

	_, err := os.Stat(destDir)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check skill directory: %w", err)
	}

	if err := os.RemoveAll(destDir); err != nil {
		return false, fmt.Errorf("failed to remove skill directory: %w", err)
	}

	return true, nil
}

// InstalledVersion reads the version from an installed skill directory.
func InstalledVersion(target AgentTarget) string {
	versionPath := filepath.Join(SkillDir(target), versionFileName)
	data, err := os.ReadFile(versionPath)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

// IsInstalled checks if the skill is installed for a given target.
func IsInstalled(target AgentTarget) bool {
	skillMD := filepath.Join(SkillDir(target), "SKILL.md")
	_, err := os.Stat(skillMD)
	return err == nil
}

// DetectAgents returns targets whose parent agent directory exists.
func DetectAgents(targets []AgentTarget) []AgentTarget {
	var detected []AgentTarget
	for _, t := range targets {
		agentDir := filepath.Dir(t.BaseDir)
		if _, err := os.Stat(agentDir); err == nil {
			detected = append(detected, t)
		}
	}
	return detected
}

// DisplayPath returns a user-friendly path with ~ for home directory.
func DisplayPath(path, homeDir string) string {
	if strings.HasPrefix(path, homeDir) {
		return "~" + path[len(homeDir):]
	}
	return path
}
