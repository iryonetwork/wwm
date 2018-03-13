package main

import (
	"github.com/caarlos0/env"

	"github.com/iryonetwork/wwm/config"
)

type CommonConfig config.Config

type Config struct {
	config.Config
	CloudStorageHost             string `env:"CLOUD_STORAGE_HOST" envDefault:"cloudStorage"`
	CloudStoragePath             string `env:"CLOUD_STORAGE_PATH" envDefault:"storage"`
	BoltDBFilepath               string `env:"BOLT_DB_FILEPATH" envDefault:"/data/batchStorageSync.db"`
	PrometheusPushGatewayAddress string `env:"PROMETHEUS_PUSH_GATEWAY_ADDRESS" envDefault:"http://localPrometheusPushGateway:9091"`
}

func GetConfig() (*Config, error) {
	common, err := config.New()
	if err != nil {
		return nil, err
	}

	cfg := &Config{Config: *common}

	return cfg, env.Parse(cfg)
}
