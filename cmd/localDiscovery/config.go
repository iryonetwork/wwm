package main

import (
	"github.com/caarlos0/env"

	"github.com/iryonetwork/wwm/config"
)

// Config represents configuration of localDiscovery
type Config struct {
	config.Config

	VaultToken  string `env:"VAULT_TOKEN,required"`
	VaultDBRole string `env:"VAULT_DB_ROLE,required"`

	PGHost     string `env:"POSTGRES_HOST" envDefault:"postgres"`
	PGDatabase string `env:"POSTGRES_DATABASE" envDefault:"localdiscovery"`
	PGRole     string `env:"POSTGRES_ROLE" envDefault:"localdiscoveryservice"`

	CloudDiscoveryHost string `env:"CLOUD_DISCOVERY_HOST" envDefault:"cloudDiscovery"`
}

func getConfig() (*Config, error) {
	common, err := config.New()
	if err != nil {
		return nil, err
	}

	cfg := &Config{Config: *common}

	return cfg, env.Parse(cfg)
}
