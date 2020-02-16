package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	DbArgs string `json:"dbargs"`
	Token  string `json:"token"`
}

func getConfig(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}

	decoder := json.NewDecoder(file)
	config := new(Config)
	if decodeErr := decoder.Decode(config); decodeErr != nil {
		return Config{}, decodeErr
	}

	return *config, nil
}
