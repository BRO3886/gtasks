package internal

import (
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/tasks/v1"
)

//ReadCredentials reads the config.json file
func ReadCredentials() *oauth2.Config {
	folderPath := GetInstallLocation()
	b, err := ioutil.ReadFile(folderPath + "/config.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	config, err := google.ConfigFromJSON(b, tasks.TasksScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	return config
}

func GenerateConfig() {
	credString := `
	{
		"installed": {
			"client_id": "415973160530-1onpd10rt7vl0dc79sh0hf8qb7ilc0bo.apps.googleusercontent.com",
			"project_id": "tasks-cli-tool-1604690538991",
			"auth_uri": "https://accounts.google.com/o/oauth2/auth",
			"token_uri": "https://oauth2.googleapis.com/token",
			"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
			"client_secret": "s-QOvdsBvXh7vgIGJSgojFa7",
			"redirect_uris": [
				"urn:ietf:wg:oauth:2.0:oob",
				"http://localhost"
			]
		}
	}`
	mode := int(0666)
	folderPath := GetInstallLocation()
	ioutil.WriteFile(folderPath+"/config.json", []byte(credString), os.FileMode(mode))
}
