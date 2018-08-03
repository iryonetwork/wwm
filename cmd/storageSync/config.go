package main

import (
	"time"

	"github.com/caarlos0/env"

	"github.com/iryonetwork/wwm/config"
)

// Config represents configuration of storageSync
type Config struct {
	config.Config
	CloudStorageHost   string        `env:"CLOUD_STORAGE_HOST" envDefault:"cloudStorage"`
	CloudStoragePath   string        `env:"CLOUD_STORAGE_PATH" envDefault:"storage"`
	NatsAddr           string        `env:"NATS_ADDR" envDefault:"localNats:4242"`
	NatsClusterID      string        `env:"NATS_CLUSTER_ID" envDefault:"localNats"`
	NatsClientID       string        `env:"NATS_CLIENT_ID" envDefault:"storageSync"`
	NatsUsername       string        `env:"NATS_USERNAME" envDefault:"nats"`
	NatsSecret         string        `env:"NATS_SECRET,required"`
	NatsConnRetries    int           `env:"NATS_CONN_RETRIES" envDefault:"10"`
	NatsConnWait       time.Duration `env:"NATS_CONN_WAIT" envDefault:"500ms"`
	NatsConnWaitFactor float32       `env:"NATS_CONN_WAIT_FACTOR" envDefault:"3.0"`
	AckWait            time.Duration `env:"ACK_WAIT" envDefault:"10000ms"`
	MaxInflight        int           `env:"MAX_INFLIGHT" envDefault:"10"`
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
