package main

import (
	"github.com/caarlos0/env"

	"github.com/iryonetwork/wwm/config"
)

// Config represents configuration of localAuth
type Config struct {
	config.Config
	BoltDBFilepath   string `env:"BOLT_DB_FILEPATH" envDefault:"/data/localAuth.db"`
	CloudAuthHost    string `env:"CLOUD_AUTH_HOST" envDefault:"cloudAuth"`
	CloudAuthPath    string `env:"CLOUD_AUTH_PATH" envDefault:"auth"`
	AuthSyncKeyPath  string `env:"AUTH_SYNC_KEY_PATH" envDefault:"/certs/localAuthSync-key.pem"`
	AuthSyncCertPath string `env:"AUTH_SYNC_CERT_PATH" envDefault:"/certs/localAuthSync.pem"`
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
