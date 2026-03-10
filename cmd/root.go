package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/BRO3886/gtasks/internal/config"
	"github.com/BRO3886/gtasks/internal/update"
	"github.com/BRO3886/gtasks/internal/utils"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// Version is set during build
var Version = "DEV"

// updateResultCh receives the background update check result (if any).
var updateResultCh = make(chan *update.Result, 1)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "gtasks",
	Short:   "A CLI Tool for Google Tasks",
	Version: Version,
	Long: `
	A CLI Tool for managing your Google Tasks:

	* Run gtasks help for checking out inline help
	* Run gtasks login to log-in with your Google account

	Made with ❤ by https://github.com/BRO3886
`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if !shouldCheckForUpdate(cmd) {
			updateResultCh <- nil
			return
		}

		homeDir, err := os.UserHomeDir()
		if err != nil {
			updateResultCh <- nil
			return
		}

		// If cache is fresh, resolve synchronously (just a file read — instant).
		// If cache is stale, fetch from GitHub in a goroutine and wait up to 2s in PostRun.
		cached := update.ReadCache(homeDir)
		if cached != nil && update.IsCacheFresh(cached, time.Now()) {
			if update.CompareVersions(Version, cached.Latest) {
				updateResultCh <- &update.Result{Latest: cached.Latest, HasUpdate: true}
			} else {
				updateResultCh <- nil
			}
		} else {
			go func() {
				updateResultCh <- update.Check(homeDir, Version)
			}()
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		printUpdateNotice()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		utils.ErrorP("%s\n", err.Error())
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	config.LoadAppConfig()
}

// shouldCheckForUpdate returns false for commands/contexts where the check should be skipped.
func shouldCheckForUpdate(cmd *cobra.Command) bool {
	if os.Getenv("GTASKS_NO_UPDATE_CHECK") != "" {
		return false
	}

	if Version == "" || Version == "DEV" {
		return false
	}

	name := cmd.Name()
	if name == "completion" || name == "skills" {
		return false
	}

	// Skip if stdout is not a TTY (piped output or scripting context)
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	if fi.Mode()&os.ModeCharDevice == 0 {
		return false
	}

	return true
}

// printUpdateNotice prints an update notice to stderr if a newer version is available.
// On a cache-hit the channel is already filled so this returns instantly.
// On a cache-miss (first run per 24h) waits up to 2s for the GitHub fetch goroutine.
func printUpdateNotice() {
	var result *update.Result
	select {
	case result = <-updateResultCh:
	case <-time.After(2 * time.Second):
		result = nil
	}

	if result != nil && result.HasUpdate {
		yellow := color.New(color.FgYellow)
		fmt.Fprintln(os.Stderr)
		yellow.Fprintf(os.Stderr, "A new version of gtasks is available: %s → %s\n", Version, result.Latest)
		fmt.Fprintf(os.Stderr, "Update: curl -fsSL https://gtasks.sidv.dev/install | bash\n")
	}
}
