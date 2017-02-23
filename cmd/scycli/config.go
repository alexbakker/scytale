package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"path"
)

var (
	configPath     string
	configFilename string
	configPrepared = false
	configDefaults = &config{
		Endpoint: "http://localhost:8080",
	}
)

type config struct {
	Endpoint string `json:"endpoint"`
}

func init() {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	configPath = path.Join(usr.HomeDir, ".config/scytale")
	configFilename = path.Join(configPath, "config.json")
}

func prepareConfig() error {
	if _, err := os.Stat(configPath); err != nil {
		if os.IsNotExist(err) {
			if err = os.MkdirAll(configPath, 0777); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	if _, err := os.Stat(configFilename); err != nil {
		if os.IsNotExist(err) {
			return saveConfig(configDefaults)
		}
		return err
	}
	return nil
}

func loadConfig() (*config, error) {
	if !configPrepared {
		if err := prepareConfig(); err != nil {
			return nil, err
		}
		configPrepared = true
	}

	bytes, err := ioutil.ReadFile(configFilename)
	if err != nil {
		return nil, err
	}

	cfg := config{}
	if err := json.Unmarshal(bytes, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func saveConfig(cfg *config) error {
	bytes, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(configFilename, bytes, 0666)
}
