package config

import (
	"encoding/json"
	"os"
)

func write(cfg Config) error {

	filePath, err := getConfigFilePath()

	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")

	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0600)

}

func (c *Config) SetUser(username string) error {

	c.CurrentUsername = username

	err := write(*c)

	if err != nil {
		return err
	}

	return nil

}
