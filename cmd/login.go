package cmd

import (
	"github.com/BRO3886/gtasks/api"
	"github.com/BRO3886/gtasks/internal/config"
	"github.com/BRO3886/gtasks/internal/utils"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Logging into Google Tasks",
	Long:  `This command uses the credentials.json file and makes a request to get your tokens`,
	Run: func(cmd *cobra.Command, args []string) {
		c := config.ReadCredentials()
		err := api.Login(c)
		if err != nil {
			utils.ErrorP("%v\n", err)
			return
		}
		utils.Info("Logged in successfully\n")
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
