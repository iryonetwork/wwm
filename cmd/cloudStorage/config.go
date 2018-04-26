package main

import (
	"github.com/caarlos0/env"

	"github.com/iryonetwork/wwm/config"
)

// Config represents configuration of cloudStorage
type Config struct {
	config.Config

	S3Endpoint  string `env:"S3_ENDPOINT" envDefault:"cloudMinio:9000"`
	S3AccessKey string `env:"S3_ACCESS_KEY" envDefault:"cloud"`
	S3Region    string `env:"S3_REGION" envDefault:"us-east-1"`
	S3Secret    string `env:"S3_SECRET,required"`

	StorageEncryptionKey string `env:"STORAGE_ENCRYPTION_KEY,required"`
}

// GetConfig parses environment variables and returns pointer to config and error
func GetConfig() (*Config, error) {
	common, err := config.New()
	if err != nil {
		return nil, err
	}

	cfg := &Config{Config: *common}

	return cfg, env.Parse(cfg)
}
