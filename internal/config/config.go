package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	URL      string `json:"db_url"`
	UserName string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func Read() Config {
	cfgPath, err := getConfigFilePath()
	if err != nil {
		return Config{}
	}

	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return Config{}
	}

	var cfg Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return Config{}
	}
	return cfg
}

func (cfg Config) SetUser(userName string) {
	cfg.UserName = userName
	if err := write(cfg); err != nil {
		fmt.Errorf("failed to write config: %v", err)
	}
}

func getConfigFilePath() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	cfgPath := homedir + "/" + configFileName
	return cfgPath, nil
}

func write(cfg Config) error {
	cfgPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	err = os.WriteFile(cfgPath, data, 0600)
	if err != nil {
		return err
	}

	return nil
}
