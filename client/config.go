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
	cfgFile, err := openFileForRW(path.Join(app.cfgDir, ".cfg"))
	if err != nil {
		panic(err)
	}

	cfgData, err := os.ReadFile(cfgFile.Name())
	if err != nil {
		panic(err)
	}

	if len(cfgData) == 0 {
		cfgData, err = json.Marshal(config{
			GophKeeperURL: ":8080",
		})
		if err != nil {
			panic(err)
		}

		_, err = cfgFile.Write(cfgData)
		if err != nil {
			panic(err)
		}
		app.state = cfgFileInitialized
	}

	err = json.Unmarshal(cfgData, app.config)
	if err != nil {
		panic(err)
	}
	fmt.Println(app.config)
	app.state = cfgFileInitialized
}

type config struct {
	GophKeeperURL string `json:"keeper_url"`
}

func openFileForRW(path string) (*os.File, error) {
	cfgFile, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, filePerm)
	if err != nil {
		return nil, err
	}

	return cfgFile, nil
}

// UNIX permissions
const (
	// user:  Write Read _
	// group: _     Read _
	// other: _     Read _
	filePerm os.FileMode = 0o600
	dirPerm  os.FileMode = 0o700
)
