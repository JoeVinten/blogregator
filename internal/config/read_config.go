package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func getConfigFilePath() (string, error) {

	homeDir, err := os.UserHomeDir()

	if err != nil {
		return "", err
	}

	configFile := filepath.Join(homeDir, configFileName)

	return configFile, nil
}

func ReadConfig() (Config, error) {

	filePath, err := getConfigFilePath()

	if err != nil {
		return Config{}, fmt.Errorf("getting the file path %w", err)
	}

	data, err := os.ReadFile(filePath)

	if err != nil {
		if os.IsNotExist(err) {
			return Config{}, fmt.Errorf("config path not found at %s: %w", filePath, err)
		}
		return Config{}, fmt.Errorf("reading the config file %w", err)
	}

	config := Config{}

	err = json.Unmarshal(data, &config)

	if err != nil {
		return Config{}, fmt.Errorf("parsing the config JSON, %w", err)
	}

	return config, nil
}
