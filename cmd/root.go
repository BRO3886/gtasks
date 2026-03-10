package cmd

import (
	"fmt"

	"github.com/BRO3886/gtasks/internal/config"
	"github.com/BRO3886/gtasks/internal/utils"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

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
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	viper.SetDefault("license", "apache")
}

// initConfig loads the gtasks config file and reads environment variables.
func initConfig() {
	// Load config.toml from the XDG/gtasks config directory.
	config.LoadAppConfig()

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			utils.ErrorP("%v", err)
		}
		viper.AddConfigPath(home)
		viper.SetConfigName(".google-tasks-cli")
	}

	viper.AutomaticEnv()

	// Suppress "config file not found" — the legacy viper config (.google-tasks-cli) is
	// not used; gtasks reads its own config.toml via config.LoadAppConfig() above.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
