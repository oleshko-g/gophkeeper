package client

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/viper"
)

var cfgDir string

func initConfig() {
	cfgDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	viper.SetConfigFile(path.Join(cfgDir, "gophkeeper/config.yaml"))
	fmt.Println(viper.ConfigFileUsed())
}

func initConfigFile(path string) error {
	_, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, filePerm)
	if err != nil {
		return err
	}
	return nil
}

// UNIX permissions
const (
	// user:  Write Read _
	// group: _     Read _
	// other: _     Read _
	filePerm os.FileMode = 0o644
)
