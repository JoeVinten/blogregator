package config

import (
	"encoding/json"
	"fmt"
	"os"
)

func ReadConfig() (Config, error) {

	filePath, err := getConfigFilePath()

	if err != nil {
		return Config{}, fmt.Errorf("getting the file path %w", err)
	}

	file, err := os.Open(filePath)

	if err != nil {
		if os.IsNotExist(err) {
			return Config{}, fmt.Errorf("config path not found at %s: %w", filePath, err)
		}
		return Config{}, fmt.Errorf("reading the config file %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	cfg := Config{}

	err = decoder.Decode(&cfg)

	if err != nil {
		return Config{}, fmt.Errorf("parsing the config JSON, %w", err)
	}

	return cfg, nil
}
