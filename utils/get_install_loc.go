package utils

import (
	"log"

	"github.com/kardianos/osext"
)

//GetInstallLocation to get the install loc of binary
func GetInstallLocation() string {
	folderPath, err := osext.ExecutableFolder()
	if err != nil {
		log.Fatal(err)
	}
	return folderPath
}
