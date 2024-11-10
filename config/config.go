package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Routes map[string]string
}

var config Config

func LoadConfig(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	defer file.Close()

	decoder := json.NewDecoder(file)

	if err := decoder.Decode(&config); err != nil {
		return err
	}

	return nil
}

func GetConfig() Config {
	return config
}
