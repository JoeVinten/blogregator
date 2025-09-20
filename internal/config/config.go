package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

func getConfigFilePath() (string, error) {

	homeDir, err := os.UserHomeDir()

	if err != nil {
		return "", err
	}

	configFile := filepath.Join(homeDir, configFileName)

	return configFile, nil
}

func write(cfg Config) error {

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
