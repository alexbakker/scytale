package config

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
	configDefaults = &Config{
		URL: "http://localhost:8080",
	}
)

type Config struct {
	URL string `json:"url"`
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
			return Save(configDefaults)
		}
		return err
	}
	return nil
}

func Load() (*Config, error) {
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

	cfg := Config{}
	if err := json.Unmarshal(bytes, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func Save(cfg *Config) error {
	bytes, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(configFilename, bytes, 0666)
}
