package main

import (
	"github.com/caarlos0/env"

	"github.com/iryonetwork/wwm/config"
)

// Config represents configuration of the waitlist service
type Config struct {
	config.Config

	StorageEncryptionKey string `env:"STORAGE_ENCRYPTION_KEY,required"`

	BoltDBFilepath string `env:"BOLT_DB_FILEPATH" envDefault:"/data/waitlist.db"`

	DefaultListID   string `env:"DEFAULT_LIST_ID" envDefault:"22afd921-0630-49f4-89a8-d1ad7639ee83"`
	DefaultListName string `env:"DEFAULT_LIST_NAME" envDefault:"default"`
}

func getConfig() (*Config, error) {
	common, err := config.New()
	if err != nil {
		return nil, err
	}

	cfg := &Config{Config: *common}

	return cfg, env.Parse(cfg)
}
