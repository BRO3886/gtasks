package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/BRO3886/gtasks/internal/config"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/tasks/v1"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Logging into Google Tasks",
	Long:  `This command uses the credentials.json file and makes a request to get your tokens`,
	Run: func(cmd *cobra.Command, args []string) {
		c := config.ReadCredentials()
		getClient(c)
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(c *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	folderPath := config.GetInstallLocation()
	// fmt.Println(folderPath)
	tokFile := folderPath + "/token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(c)
		saveToken(tokFile, tok)
	}
	return c.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\nEnter the code: ", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

//gets the tasks service
func getService() *tasks.Service {
	c := config.ReadCredentials()
	client := getClient(c)

	srv, err := tasks.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve tasks Client %v", err)
	}

	return srv
}
