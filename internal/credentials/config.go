package credentials

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/cockroachdb/errors"
)

type Config struct {
	BridgeIP            string  `json:"bridge_ip"`
	OnePasswordItemName string  `json:"onepassword_item_name"`
	Latitude            float64 `json:"latitude,omitempty"`
	Longitude           float64 `json:"longitude,omitempty"`
}

func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrap(err, "getting user home directory")
	}

	configDir := filepath.Join(homeDir, ".config", "limelight")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", errors.Wrap(err, "creating config directory")
	}

	return filepath.Join(configDir, "config.json"), nil
}

func LoadConfig() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, errors.Wrap(err, "getting config path")
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "reading config file")
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, errors.Wrap(err, "unmarshaling config")
	}

	return &config, nil
}

func SaveConfig(config *Config) error {
	configPath, err := GetConfigPath()
	if err != nil {
		return errors.Wrap(err, "getting config path")
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return errors.Wrap(err, "marshaling config")
	}

	tmpPath := configPath + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0600); err != nil {
		return errors.Wrap(err, "writing temp config file")
	}

	if err := os.Rename(tmpPath, configPath); err != nil {
		os.Remove(tmpPath)
		return errors.Wrap(err, "moving config file")
	}

	return nil
}
