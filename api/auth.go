package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/BRO3886/gtasks/internal/config"
	"github.com/BRO3886/gtasks/internal/utils"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/tasks/v1"
)

func Login(c *oauth2.Config) error {
	folderPath := config.GetInstallLocation()
	// fmt.Println(folderPath)
	tokFile := folderPath + "/token.json"
	_, err := tokenFromFile(tokFile)
	if err != nil {
		tok := getTokenFromWeb(c)
		saveToken(tokFile, tok)
		return nil
	}
	return fmt.Errorf("already logged in")
}

func Logout() error {
	folderPath := config.GetInstallLocation()
	// fmt.Println(folderPath)
	tokFile := folderPath + "/token.json"
	return os.Remove(tokFile)
}

// gets the tasks service
func GetService() *tasks.Service {
	c := config.ReadCredentials()
	client := getClient(c)
	srv, err := tasks.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		utils.ErrorP("Unable to retrieve tasks Client %v", err)
	}

	return srv
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
	utils.Warn("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n\nEnter the code: ", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		utils.ErrorP("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		utils.ErrorP("Unable to retrieve token from web: %v", err)
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
	utils.Warn("Saving credential file\n")
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		utils.ErrorP("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
