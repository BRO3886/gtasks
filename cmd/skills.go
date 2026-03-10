package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/BRO3886/gtasks/internal/skills"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var skillsAgentFlag string

var skillsCmd = &cobra.Command{
	Use:   "skills",
	Short: "Manage AI agent skills for gtasks",
	Long: `Install, uninstall, and check the status of the gtasks agent skill.

The gtasks skill teaches AI coding agents how to use gtasks effectively.
It includes command references, workflows, and usage examples.`,
}

var skillsInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install gtasks skill for AI agents",
	Long: `Installs the gtasks agent skill to the selected AI agent's skill directory.

Supported agents:
  claude    -> ~/.claude/skills/gtasks-cli/    (Claude Code, Copilot, Cursor, OpenCode, Augment)
  codex     -> ~/.agents/skills/gtasks-cli/    (Codex CLI, Copilot, Windsurf, OpenCode, Augment)
  openclaw  -> ~/.openclaw/skills/gtasks-cli/  (OpenClaw)

Without --agent, automatically detects installed agents.`,
	RunE: runSkillsInstall,
}

var skillsUninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall gtasks skill from AI agents",
	Long:  `Removes the gtasks agent skill from the selected AI agent's skill directory.`,
	RunE:  runSkillsUninstall,
}

var skillsStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show skill installation status",
	RunE:  runSkillsStatus,
}

func init() {
	skillsInstallCmd.Flags().StringVar(&skillsAgentFlag, "agent", "", "Agent target: claude, codex, openclaw, or all")
	skillsUninstallCmd.Flags().StringVar(&skillsAgentFlag, "agent", "", "Agent target: claude, codex, openclaw, or all")

	skillsCmd.AddCommand(skillsInstallCmd, skillsUninstallCmd, skillsStatusCmd)
	rootCmd.AddCommand(skillsCmd)
}

func runSkillsInstall(cmd *cobra.Command, args []string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	allTargets := skills.DefaultTargets(homeDir)
	targets, err := resolveTargets(allTargets, skillsAgentFlag, "install")
	if err != nil {
		return err
	}
	if len(targets) == 0 {
		return nil
	}

	return installToTargets(embeddedSkills, targets, homeDir)
}

func installToTargets(embeddedFS fs.FS, targets []skills.AgentTarget, homeDir string) error {
	green := color.New(color.FgGreen, color.Bold)

	for _, t := range targets {
		written, err := skills.Install(embeddedFS, t, Version)
		if err != nil {
			return fmt.Errorf("failed to install for %s: %w", t.Name, err)
		}

		green.Print("OK ")
		fmt.Printf("Installed gtasks-cli skill to %s\n", skills.DisplayPath(skills.SkillDir(t), homeDir))
		fmt.Printf("  Files: %s\n", strings.Join(written, ", "))
	}

	fmt.Println("\nThe skill will be available in your next session.")
	return nil
}

func runSkillsUninstall(cmd *cobra.Command, args []string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	allTargets := skills.DefaultTargets(homeDir)
	targets, err := resolveTargets(allTargets, skillsAgentFlag, "uninstall")
	if err != nil {
		return err
	}
	if len(targets) == 0 {
		return nil
	}

	green := color.New(color.FgGreen, color.Bold)
	yellow := color.New(color.FgYellow)

	for _, t := range targets {
		removed, err := skills.Uninstall(t)
		if err != nil {
			return fmt.Errorf("failed to uninstall from %s: %w", t.Name, err)
		}
		if removed {
			green.Print("OK ")
			fmt.Printf("Removed gtasks-cli skill from %s\n", skills.DisplayPath(skills.SkillDir(t), homeDir))
		} else {
			yellow.Print("- ")
			fmt.Printf("Not installed at %s\n", skills.DisplayPath(skills.SkillDir(t), homeDir))
		}
	}

	return nil
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
			red.Printf("  x ")
			fmt.Printf("%-12s %s (not installed)\n", t.Name, displayDir)
			continue
		}

		installed := skills.InstalledVersion(t)
		if installed == "" {
			yellow.Printf("  ? ")
			fmt.Printf("%-12s %s (installed, unknown version)\n", t.Name, displayDir)
		} else if installed != Version {
			yellow.Printf("  ! ")
			fmt.Printf("%-12s %s (installed %s, outdated)\n", t.Name, displayDir, installed)
		} else {
			green.Printf("  OK ")
			fmt.Printf("%-12s %s (installed %s)\n", t.Name, displayDir, installed)
		}
	}

	return nil
}

func resolveTargets(allTargets []skills.AgentTarget, agentFlag, action string) ([]skills.AgentTarget, error) {
	if agentFlag != "" {
		return resolveAgentFlag(allTargets, agentFlag)
	}

	if !isTerminal() {
		detected := skills.DetectAgents(allTargets)
		if len(detected) == 0 {
			return allTargets[:1], nil
		}
		return detected, nil
	}

	return runAgentPicker(allTargets, action)
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

func runAgentPicker(allTargets []skills.AgentTarget, action string) ([]skills.AgentTarget, error) {
	detected := skills.DetectAgents(allTargets)
	defaults := make([]string, 0, len(detected))
	for _, t := range detected {
		defaults = append(defaults, t.Key)
	}
	if len(defaults) == 0 && action == "install" {
		for _, t := range allTargets {
			defaults = append(defaults, t.Key)
		}
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "> {{ .Name | cyan }}",
		Inactive: "  {{ .Name }}",
		Selected: "{{ .Name }}",
		Details:  "",
	}

	prompt := promptui.Select{
		Label:     fmt.Sprintf("%s gtasks skill for which agent?", strings.Title(action)),
		Items:     allTargets,
		Templates: templates,
		Size:      len(allTargets),
		Searcher: func(input string, index int) bool {
			t := allTargets[index]
			name := strings.ToLower(t.Name + " " + t.Key)
			return strings.Contains(name, strings.ToLower(input))
		},
	}

	index, _, err := prompt.Run()
	if err != nil {
		if err == promptui.ErrInterrupt || err == promptui.ErrEOF {
			fmt.Println("Cancelled")
			return nil, nil
		}
		return nil, fmt.Errorf("selection error: %w", err)
	}

	selected := []skills.AgentTarget{allTargets[index]}
	if len(defaults) > 0 && action == "install" && contains(defaults, selected[0].Key) {
		return selected, nil
	}

	return selected, nil
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func isTerminal() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (info.Mode() & os.ModeCharDevice) != 0
}
