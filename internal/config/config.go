// Read and supply configuration values
package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Database struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
		Schema   string `json:"schema"`
		Database string `json:"database"`
	} `json:"database"`
	User struct {
		Secret     string `json:"secret"`
		PrivateKey string `json:"private_key"`
		PublicKey  string `json:"public_key"`
	} `json:"user"`
}

// Read the file listed as "f" and return a config object
func New(f string) (*Config, error) {
	var c Config

	if _, err := os.Stat(f); err != nil {
		return nil, fmt.Errorf("file does not exist: %v", f)
	}

	if file, err := os.ReadFile(f); err != nil {
		return nil, fmt.Errorf("unable to read config file: %v", err)
	} else {
		if err := json.Unmarshal(file, &c); err != nil {
			return nil, fmt.Errorf("error parsing json: %v", err)
		} else {
			return &c, nil
		}
	}
}
