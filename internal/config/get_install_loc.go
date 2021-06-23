package config

import (
	"github.com/BRO3886/gtasks/internal/utils"
	"github.com/kardianos/osext"
)

//GetInstallLocation to get the install loc of binary
func GetInstallLocation() string {
	folderPath, err := osext.ExecutableFolder()
	if err != nil {
		utils.ErrorP("Get install location: %s", err.Error())
	}
	return folderPath
}
