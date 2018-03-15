package main

import (
	"time"

	"github.com/caarlos0/env"

	"github.com/iryonetwork/wwm/config"
)

// Config represents configuration of localStorage
type Config struct {
	config.Config
	S3Endpoint         string        `env:"S3_ENDPOINT" envDefault:"localMinio:9000"`
	S3AccessKey        string        `env:"S3_ACCESS_KEY" envDefault:"local"`
	S3Region           string        `env:"S3_REGION" envDefault:"us-east-1"`
	NatsAddr           string        `env:"NATS_ADDR" envDefault:"localNats:4242"`
	NatsClusterID      string        `env:"NATS_CLUSTER_ID" envDefault:"localNats"`
	NatsClientID       string        `env:"NATS_CLIENT_ID" envDefault:"localStorage"`
	NatsUsername       string        `env:"NATS_USERNAME" envDefault:"nats"`
	NatsConnRetries    int           `env:"NATS_CONN_RETRIES" envDefault:"5"`
	NatsConnWait       time.Duration `env:"NATS_CONN_WAIT" envDefault:"500ms"`
	NatsConnWaitFactor float32       `env:"NATS_CONN_WAIT_FACTOR" envDefault:"3.0"`
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
