package config

import (
	"os"
	"encoding/json"
	"io"
)

const configFileName = "/.gatorconfig.json"

type Config struct {
	DBURL string  `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func getConfigFilePath() (string, error) {
	path, err := os.UserHomeDir()

	if err != nil {
		return "", err
	}

	path += configFileName

	return path, nil
}

func Read() (Config, error) {
	path, err := getConfigFilePath()

	if err != nil {
		return Config{}, err
	}

	jsonFile, err := os.Open(path)

	if err != nil {
		return Config{}, err
	}

	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)

	if err != nil {
		return Config{}, err
	}

	var cfg Config

	err = json.Unmarshal(byteValue, &cfg)

	if err != nil {
		return Config{}, err
	}



	return cfg, nil
}

func write(cfg Config) error {
	path, err := getConfigFilePath()

	if err != nil {
		return nil
	}

	jsonFile, err := os.Create(path)

	if err != nil {
		return err
	}

	defer jsonFile.Close()

	jsonString, err := json.Marshal(cfg)

	if err != nil {
		return err
	}

	_, err = jsonFile.Write(jsonString)
	
	if err != nil {
		return err
	}

	return nil

}

func (cfg *Config) SetUser (userName string) error{
	cfg.CurrentUserName = userName
	return write(*cfg)
}
