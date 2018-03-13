package config

import (
	"github.com/caarlos0/env"
)

type Config struct {
	ServerHost       string `env:"SERVER_HOST" envDefault:"0.0.0.0"`
	ServerPort       int    `env:"SERVER_PORT" envDefault:"443"`
	KeyPath          string `env:"KEY_PATH" envDefault:"/certs/private.key"`
	CertPath         string `env:"CERT_PATH" envDefault:"/certs/public.crt"`
	MetricsPort      int    `env:"METRICS_PORT" envDefault:"9090"`
	MetricsNamespace string `env:"METRICS_NAMESPACE"`
	StatusPort       int    `env:"STATUS_PORT" envDefault:"4433"`
	StatusNamespace  string `env:"STATUS_NAMESPACE"`
	StorageHost      string `env:"STORAGE_HOST" envDefault:"localStorage`
	StoragePath      string `env:"STORAGE_PATH" envDefault:"storage"`
	AuthHost         string `env:"AUTH_HOST" envDefault:"localAuth"`
	AuthPath         string `env:"AUTH_PATH" envDefault:"auth"`
}

func New() (*Config, error) {
	cfg := &Config{}

	return cfg, env.Parse(cfg)
}
