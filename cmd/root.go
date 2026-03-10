package cmd

import (
	"github.com/BRO3886/gtasks/internal/config"
	"github.com/BRO3886/gtasks/internal/utils"
	"github.com/spf13/cobra"
)

// version is set during build
var Version = "DEV"

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
