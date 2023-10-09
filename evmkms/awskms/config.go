package awskms

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	KeyID   string `json:"KeyID"`
	ChainID uint64 `json:"ChainID"`
}

func (cfg Config) IsValid() (bool, error) {
	if cfg.KeyID == "" {
		return false, fmt.Errorf("empty KeyID")
	}

	return true, nil
}

type StaticCredentialsConfig struct {
	Config
	Region          string `json:"Region"`
	AccessKeyID     string `json:"AccessKeyID"`
	SecretAccessKey string `json:"SecretAccessKey"`
	SessionToken    string `json:"SessionToken,omitempty"`
}

func (cfg StaticCredentialsConfig) IsValid() (bool, error) {
	if cfg.Region == "" {
		return false, fmt.Errorf("empty Region")
	}

	if cfg.AccessKeyID == "" {
		return false, fmt.Errorf("empty AccessKeyID")
	}

	if cfg.SecretAccessKey == "" {
		return false, fmt.Errorf("empty SecretAccessKey")
	}

	return cfg.Config.IsValid()
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

	if _, err = cfg.IsValid(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func LoadStaticCredentialsConfigConfigFromFile(filePath string) (*StaticCredentialsConfig, error) {
	f, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var cfg StaticCredentialsConfig
	err = json.Unmarshal(f, &cfg)
	if err != nil {
		return nil, err
	}

	if _, err = cfg.IsValid(); err != nil {
		return nil, err
	}

	return &cfg, nil
}
