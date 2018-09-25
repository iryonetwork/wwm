package config

import (
	"github.com/caarlos0/env"
)

// Config struct holds commonly used configuration options
type Config struct {
	DomainType       string `env:"DOMAIN_TYPE" envDefault:"global"`
	DomainID         string `env:"DOMAIN_ID" envDefault:"*"`
	ServerHost       string `env:"SERVER_HOST" envDefault:"0.0.0.0"`
	ServerPort       int    `env:"SERVER_PORT" envDefault:"443"`
	ServerPortHTTPS  int    `env:"SERVER_PORT_HTTPS" envDefault:"443"`
	ServerPortHTTP   int    `env:"SERVER_PORT_HTTP" envDefault:"80"`
	KeyPath          string `env:"KEY_PATH,required"`
	CertPath         string `env:"CERT_PATH,required"`
	MetricsPort      int    `env:"METRICS_PORT" envDefault:"9090"`
	MetricsNamespace string `env:"METRICS_NAMESPACE"`
	StatusPort       int    `env:"STATUS_PORT" envDefault:"4433"`
	StatusNamespace  string `env:"STATUS_NAMESPACE"`
	StorageHost      string `env:"STORAGE_HOST" envDefault:"localStorage"`
	StoragePath      string `env:"STORAGE_PATH" envDefault:"storage"`
	AuthHost         string `env:"AUTH_HOST" envDefault:"localAuth"`
	AuthPath         string `env:"AUTH_PATH" envDefault:"auth"`
	TracerAddr       string `env:"TRACER_ADDR" envDefault:"jaeger:5775"`
}

// New returns new instance of Config
func New() (*Config, error) {
	cfg := &Config{}

	return cfg, env.Parse(cfg)
}
