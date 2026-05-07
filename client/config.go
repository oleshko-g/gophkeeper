package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

func newConfig() (config, error) {
	cfgData, err := os.ReadFile(path.Join(cfgDir, ".cfg"))
	if err != nil {
		return config{}, err
	}

	var cfg config
	err = json.Unmarshal(cfgData, &cfg)
	if err != nil {
		return config{}, err
	}

	return cfg, nil
}

type config struct {
	GophKeeperURL string `json:"keeper_url"`
}

var cfgDir string

func initConfig() {
	userCfgDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	cfgDir := path.Join(userCfgDir, "gophkeeper")

	err = os.MkdirAll(cfgDir, dirPerm)
	if err != nil {
		panic(err)
	}

	cfg, err := newConfig()
	a.config = cfg
	if err != nil {
		panic(err)
	}

	fmt.Println(cfgDir)
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
	dirPerm  os.FileMode = 0o700
)
