package cmd

import (
	"github.com/BRO3886/gtasks/api"
	"github.com/BRO3886/gtasks/internal/utils"
	"github.com/spf13/cobra"
)

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout currently signed in user",
	Long:  `Logout currently signed in user.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := api.Logout()
		if err != nil {
			utils.ErrorP("%v\n", err)
			return
		}
		utils.Info("Logged out successfully\n")
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
