package config

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
)

const DefaultConfigFileName = ".gatorconfig.json"

func getConfigFilePath() (string, error) {
	path, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configFilePath := path + "/" + DefaultConfigFileName
	return configFilePath, err
}

func getCurrentUser() (*user.User, error) {
	currentUser, err := user.Current()
	if err != nil {
		return nil, err
	}
	return currentUser, nil
}

type Config struct {
	DBUrl           string `json:"db_url"`
	CurrentUsername string `json:"current_username"`
}

func (c *Config) Write() error {
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		return fmt.Errorf("config file %s doesn't exist", path)
	} else if err != nil {
		return err
	}

	// Open file in write mode (overwrite)
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(c)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) SetUser(Username string) error {
	oldUsername := c.CurrentUsername
	if oldUsername != Username {
		c.CurrentUsername = Username
		err := c.Write()
		if err != nil {
			return err
		}
	}
	return nil

}

func (c *Config) SetCurrentUser() error {
	currentUser, err := getCurrentUser()
	if err != nil {
		return err
	}
	err = c.SetUser(currentUser.Username)
	return err
}

func NewConfig() (Config, error) {
	path, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	jsonData, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = json.Unmarshal(jsonData, &config)
	if err != nil {
		return Config{}, err
	}

	return config, err
}
