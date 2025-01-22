package scan

import (
	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
)

type cargoLockPackage struct {
	Name    string `toml:"name"`
	Version string `toml:"version"`
}

type cargoLockFile struct {
	Version  int                `toml:"version"`
	Packages []cargoLockPackage `toml:"package"`
}

func GetCrates(cargoLockPath string) (map[string]string, error) {
	log.Debugf("GetCrates %s", cargoLockPath)
	gathered := make(map[string]string)

	var cargoLock cargoLockFile
	_, err := toml.DecodeFile(cargoLockPath, &cargoLock)
	if err != nil {
		log.Errorf("Failed to parse %s", cargoLockPath)
		return nil, err
	}

	for _, lockPackage := range cargoLock.Packages {
		gathered[lockPackage.Name] = lockPackage.Version
	}

	return gathered, nil
}
