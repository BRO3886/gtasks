package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/BRO3886/gtasks/internal/skills"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var skillsAgentFlag string

var skillsCmd = &cobra.Command{
	Use:   "skills",
	Short: "Manage AI agent skills for gtasks",
	Long: `Install, uninstall, and check the status of the gtasks agent skill.

The gtasks skill teaches AI coding agents (Claude Code, Codex CLI, etc.)
how to use gtasks effectively. It includes command references and usage examples.`,
}

func init() {
	rootCmd.AddCommand(skillsCmd)
}

// --- skills install ---

var skillsInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install gtasks skill for AI agents",
	Long: `Installs the gtasks agent skill to the selected AI agent's skill directory.

Supported agents:
  claude    → ~/.claude/skills/gtasks-cli/    (Claude Code, Copilot, Cursor, OpenCode, Augment)
  codex     → ~/.agents/skills/gtasks-cli/    (Codex CLI, Copilot, Windsurf, OpenCode, Augment)
  openclaw  → ~/.openclaw/skills/gtasks-cli/  (OpenClaw)

Without --agent, shows an interactive picker to select which agents to install for.`,
	RunE: runSkillsInstall,
}

func init() {
	skillsInstallCmd.Flags().StringVar(&skillsAgentFlag, "agent", "", "Agent target: claude, codex, openclaw, or all")
	skillsCmd.AddCommand(skillsInstallCmd)
}

func runSkillsInstall(cmd *cobra.Command, args []string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	allTargets := skills.DefaultTargets(homeDir)
	targets, err := resolveTargets(allTargets, skillsAgentFlag, homeDir, "install")
	if err != nil {
		return err
	}
	if len(targets) == 0 {
		return nil // user cancelled
	}

	return installToTargets(EmbeddedSkills, targets, homeDir)
}

func installToTargets(embeddedFS fs.FS, targets []skills.AgentTarget, homeDir string) error {
	green := color.New(color.FgGreen, color.Bold)

	for _, t := range targets {
		written, err := skills.Install(embeddedFS, t, Version)
		if err != nil {
			return fmt.Errorf("failed to install for %s: %w", t.Name, err)
		}

		green.Print("✓ ")
		fmt.Printf("Installed gtasks-cli skill to %s\n", skills.DisplayPath(skills.SkillDir(t), homeDir))
		fmt.Printf("  Files: %s\n", strings.Join(written, ", "))
	}

	fmt.Println("\nThe skill will be available in your next session.")
	return nil
}

// --- skills uninstall ---

var skillsUninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall gtasks skill from AI agents",
	Long:  `Removes the gtasks agent skill from the selected AI agent's skill directory.`,
	RunE:  runSkillsUninstall,
}

func init() {
	skillsUninstallCmd.Flags().StringVar(&skillsAgentFlag, "agent", "", "Agent target: claude, codex, openclaw, or all")
	skillsCmd.AddCommand(skillsUninstallCmd)
}

func runSkillsUninstall(cmd *cobra.Command, args []string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	allTargets := skills.DefaultTargets(homeDir)
	targets, err := resolveTargets(allTargets, skillsAgentFlag, homeDir, "uninstall")
	if err != nil {
		return err
	}
	if len(targets) == 0 {
		return nil // user cancelled
	}

	green := color.New(color.FgGreen, color.Bold)
	yellow := color.New(color.FgYellow)

	for _, t := range targets {
		removed, err := skills.Uninstall(t)
		if err != nil {
			return fmt.Errorf("failed to uninstall from %s: %w", t.Name, err)
		}
		if removed {
			green.Print("✓ ")
			fmt.Printf("Removed gtasks-cli skill from %s\n", skills.DisplayPath(skills.SkillDir(t), homeDir))
		} else {
			yellow.Print("- ")
			fmt.Printf("Not installed at %s\n", skills.DisplayPath(skills.SkillDir(t), homeDir))
		}
	}

	return nil
}

// --- skills status ---

var skillsStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show skill installation status",
	RunE:  runSkillsStatus,
}

func init() {
	skillsCmd.AddCommand(skillsStatusCmd)
}

func runSkillsStatus(cmd *cobra.Command, args []string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	allTargets := skills.DefaultTargets(homeDir)
	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)
	yellow := color.New(color.FgYellow)

	fmt.Printf("gtasks-cli skill (binary %s):\n", Version)
	for _, t := range allTargets {
		displayDir := skills.DisplayPath(skills.SkillDir(t), homeDir)
		if !skills.IsInstalled(t) {
			red.Printf("  ✗ ")
			fmt.Printf("%-12s %s (not installed)\n", t.Name, displayDir)
			continue
		}
		installed := skills.InstalledVersion(t)
		if installed == "" {
			yellow.Printf("  ? ")
			fmt.Printf("%-12s %s (installed, unknown version)\n", t.Name, displayDir)
		} else if installed != Version {
			yellow.Printf("  ⚠ ")
			fmt.Printf("%-12s %s (installed %s, outdated)\n", t.Name, displayDir, installed)
		} else {
			green.Printf("  ✓ ")
			fmt.Printf("%-12s %s (installed %s)\n", t.Name, displayDir, installed)
		}
	}

	return nil
}

// --- shared helpers ---

// resolveTargets determines which agent targets to operate on.
// If --agent is specified, uses that directly.
// Otherwise, shows an interactive picker.
func resolveTargets(allTargets []skills.AgentTarget, agentFlag, homeDir, action string) ([]skills.AgentTarget, error) {
	// If --agent flag provided, resolve directly
	if agentFlag != "" {
		return resolveAgentFlag(allTargets, agentFlag)
	}

	// Non-interactive mode: detect agents or default to claude
	detected := skills.DetectAgents(allTargets)
	if len(detected) == 0 {
		// Default to claude
		return allTargets[:1], nil
	}
	return detected, nil
}

func resolveAgentFlag(allTargets []skills.AgentTarget, flag string) ([]skills.AgentTarget, error) {
	flag = strings.ToLower(strings.TrimSpace(flag))
	if flag == "all" {
		return allTargets, nil
	}
	for _, t := range allTargets {
		if t.Key == flag {
			return []skills.AgentTarget{t}, nil
		}
	}
	return nil, fmt.Errorf("unknown agent %q (valid: claude, codex, openclaw, all)", flag)
}
