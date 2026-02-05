package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username
	err := write(c)
	if err != nil {
		return err
	}

	return nil
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configFilePath := filepath.Join(homeDir, configFileName)
	return configFilePath, nil
}

func write(cfg *Config) error {
	filePath, err := getConfigFilePath()

	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer file.Close()
	encoder := json.NewEncoder(file)
	err = encoder.Encode(cfg)

	if err != nil {
		return err
	}

	return nil
}

func Read() (Config, error) {
	filePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	data, err := os.ReadFile(filePath)

	if err != nil {
		return Config{}, nil
	}

	var cfg Config
	err = json.Unmarshal(data, &cfg)

	if err != nil {
		return Config{}, nil
	}

	return cfg, nil
}
