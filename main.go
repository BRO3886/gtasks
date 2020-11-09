/*
Copyright © 2020 SIDDHARTHA VARMA <siddverma1999@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"io/ioutil"
	"os"

	"github.com/BRO3886/google-tasks-cli/cmd"
	"github.com/BRO3886/google-tasks-cli/utils"
)

func generateConfig() {
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
	folderPath := utils.GetInstallLocation()
	ioutil.WriteFile(folderPath+"/config.json", []byte(credString), os.FileMode(mode))
}

func main() {
	folderPath := utils.GetInstallLocation()
	_, err := ioutil.ReadFile(folderPath + "/config.json")
	if err != nil {
		generateConfig()
	}
	cmd.Execute()
}
