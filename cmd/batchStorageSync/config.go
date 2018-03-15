package main

import (
	"github.com/caarlos0/env"

	"github.com/iryonetwork/wwm/config"
)

// Config represents configuration of batchStorageSync
type Config struct {
	config.Config
	CloudStorageHost             string `env:"CLOUD_STORAGE_HOST" envDefault:"cloudStorage"`
	CloudStoragePath             string `env:"CLOUD_STORAGE_PATH" envDefault:"storage"`
	BoltDBFilepath               string `env:"BOLT_DB_FILEPATH" envDefault:"/data/batchStorageSync.db"`
	PrometheusPushGatewayAddress string `env:"PROMETHEUS_PUSH_GATEWAY_ADDRESS" envDefault:"http://localPrometheusPushGateway:9091"`
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
