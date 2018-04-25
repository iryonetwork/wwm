package main

import (
	"github.com/caarlos0/env"

	"github.com/iryonetwork/wwm/config"
)

// Config represents configuration of the waitlist service
type Config struct {
	config.Config

	BoltDBFilepath string `env:"BOLT_DB_FILEPATH" envDefault:"/data/waitlist.db"`
}

func getConfig() (*Config, error) {
	common, err := config.New()
	if err != nil {
		return nil, err
	}

	cfg := &Config{Config: *common}

	return cfg, env.Parse(cfg)
}
