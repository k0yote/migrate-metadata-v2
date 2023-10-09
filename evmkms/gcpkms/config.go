package gcpkms

import (
	"encoding/json"
	"fmt"
	"os"
)

type Key struct {
	Keyring string `json:"Keyring"`
	Name    string `json:"Name"`
	Version string `json:"Version"`
}

func (k Key) isValid() bool {
	return k.Keyring != "" && k.Name != "" && k.Version != ""
}

type Config struct {
	ProjectID          string `json:"ProjectID"`
	LocationID         string `json:"LocationID"`
	CredentialLocation string `json:"CredentialLocation,omitempty"`
	Key                Key    `json:"Key"`
	ChainID            uint64 `json:"ChainID"`
}

func (cfg Config) IsValid() (bool, error) {
	if cfg.ProjectID == "" {
		return false, fmt.Errorf("empty ProjectID")
	}

	if cfg.LocationID == "" {
		return false, fmt.Errorf("empty LocationID")
	}

	if cfg.CredentialLocation == "" {
		cfg.CredentialLocation = os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	}

	if cfg.CredentialLocation == "" {
		return false, fmt.Errorf("empty CredentialLocation")
	}

	if !cfg.Key.isValid() {
		return false, fmt.Errorf("invalid Key")
	}

	return true, nil
}

func LoadConfigFromFile(filePath string) (*Config, error) {
	f, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = json.Unmarshal(f, &cfg)
	if err != nil {
		return nil, err
	}

	if cfg.CredentialLocation == "" {
		cfg.CredentialLocation = os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")

	}

	if _, err := cfg.IsValid(); err != nil {
		return nil, err
	}

	return &cfg, nil
}
