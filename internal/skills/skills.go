package skills

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// AgentTarget represents an AI agent platform that can have skills installed.
type AgentTarget struct {
	Name     string // Display name (e.g., "Claude Code")
	Key      string // CLI flag key (e.g., "claude")
	BaseDir  string // Agent's base directory (~/.claude)
	SkillDir string // Skill installation directory (relative to BaseDir)
}

// DefaultTargets returns the standard set of agent targets with paths resolved to the user's home directory.
func DefaultTargets(homeDir string) []AgentTarget {
	return []AgentTarget{
		{
			Name:     "Claude Code",
			Key:      "claude",
			BaseDir:  filepath.Join(homeDir, ".claude"),
			SkillDir: "skills/gtasks-cli",
		},
		{
			Name:     "Codex/Windsurf",
			Key:      "codex",
			BaseDir:  filepath.Join(homeDir, ".agents"),
			SkillDir: "skills/gtasks-cli",
		},
		{
			Name:     "OpenClaw",
			Key:      "openclaw",
			BaseDir:  filepath.Join(homeDir, ".openclaw"),
			SkillDir: "skills/gtasks-cli",
		},
		{
			Name:     "VS Code Copilot",
			Key:      "copilot",
			BaseDir:  filepath.Join(homeDir, ".vscode"),
			SkillDir: "extensions/copilot/skills/gtasks-cli",
		},
	}
}

// SkillDir returns the full path to the skill directory for a target.
func SkillDir(t AgentTarget) string {
	return filepath.Join(t.BaseDir, t.SkillDir)
}

// DisplayPath returns a user-friendly display path with home directory abbreviated.
func DisplayPath(path, homeDir string) string {
	if strings.HasPrefix(path, homeDir) {
		return "~" + strings.TrimPrefix(path, homeDir)
	}
	return path
}

// Install installs the embedded skill files to the target's skill directory.
func Install(embeddedFS fs.FS, target AgentTarget, version string) ([]string, error) {
	skillDir := SkillDir(target)

	// Create the skill directory
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create skill directory: %w", err)
	}

	var written []string

	// Walk the embedded FS and copy all files
	err := fs.WalkDir(embeddedFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory
		if path == "." {
			return nil
		}

		relPath := strings.TrimPrefix(path, "./")
		targetPath := filepath.Join(skillDir, relPath)

		if d.IsDir() {
			if err := os.MkdirAll(targetPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", targetPath, err)
			}
			return nil
		}

		// Read and write the file
		content, err := fs.ReadFile(embeddedFS, path)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", path, err)
		}

		if err := os.WriteFile(targetPath, content, 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", targetPath, err)
		}

		written = append(written, relPath)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to install skill: %w", err)
	}

	// Write version file
	versionFile := filepath.Join(skillDir, ".version")
	if err := os.WriteFile(versionFile, []byte(version), 0644); err != nil {
		return nil, fmt.Errorf("failed to write version file: %w", err)
	}

	return written, nil
}

// Uninstall removes the skill from the target's skill directory.
func Uninstall(target AgentTarget) (bool, error) {
	skillDir := SkillDir(target)

	if _, err := os.Stat(skillDir); os.IsNotExist(err) {
		return false, nil
	}

	if err := os.RemoveAll(skillDir); err != nil {
		return false, fmt.Errorf("failed to remove skill directory: %w", err)
	}

	return true, nil
}

// IsInstalled checks if the skill is installed for a target.
func IsInstalled(target AgentTarget) bool {
	skillDir := SkillDir(target)
	info, err := os.Stat(skillDir)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// InstalledVersion returns the installed version string, or empty if not installed or no version file.
func InstalledVersion(target AgentTarget) string {
	versionFile := filepath.Join(SkillDir(target), ".version")
	content, err := os.ReadFile(versionFile)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(content))
}

// DetectAgents returns the targets that have their base agent directory present.
func DetectAgents(targets []AgentTarget) []AgentTarget {
	var detected []AgentTarget
	for _, t := range targets {
		if _, err := os.Stat(t.BaseDir); err == nil {
			detected = append(detected, t)
		}
	}
	return detected
}
