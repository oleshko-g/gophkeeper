package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

func initConfigDir() {
	userCfgDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	app.cfgDir = path.Join(userCfgDir, "gophkeeper")

	err = os.MkdirAll(app.cfgDir, dirPerm)
	if err != nil {
		panic(err)
	}

	app.state = cfgDirInitialized
	fmt.Println(app.cfgDir)
}

func initConfig() {
	cfgData, err := os.ReadFile(path.Join(app.cfgDir, ".cfg"))
	if err != nil {
		panic(err)
	}

	var cfg config
	err = json.Unmarshal(cfgData, &cfg)
	if err != nil {
		panic(err)
	}
}

type config struct {
	GophKeeperURL string `json:"keeper_url"`
}

func openCfgFile(path string) error {
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
	dirPerm  os.FileMode = 0o700
)
