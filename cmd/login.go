package cmd

import (
	"github.com/BRO3886/gtasks/api"
	"github.com/BRO3886/gtasks/internal/utils"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with Google Tasks",
	Long: `Authenticate with Google Tasks using OAuth2 + PKCE flow.

This command will:
1. Open your default browser for Google authentication
2. Start a local server to handle the OAuth2 callback
3. Save your authentication token for future use

If the browser doesn't open automatically, you'll be provided with a URL to visit manually.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := api.Login()
		if err != nil {
			utils.ErrorP("Login failed: %v\n", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}