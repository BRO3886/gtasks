/*
Copyright © 2022 SIDDHARTHA VARMA <siddverma1999@gmail.com>

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

	"github.com/BRO3886/gtasks/cmd"
	"github.com/BRO3886/gtasks/internal/config"
)

func main() {
	folderPath := config.GetInstallLocation()
	_, err := ioutil.ReadFile(folderPath + "/config.json")
	if err != nil {
		config.GenerateConfig()
	}
	cmd.Execute()
}
