package main

import (
	"encoding/json"
	"io/ioutil"
	"reflect"

	"github.com/caarlos0/env"

	"github.com/iryonetwork/wwm/config"
	"github.com/iryonetwork/wwm/reports/filesDataExporter"
)

// Config represents configuration of batchDataExporter
type Config struct {
	config.Config
	BucketsRateLimit int `env:"BUCKETS_RATE_LIMIT" envDefault:"2"`

	// filepath to yaml
	FieldsToSanitize  SanitizerCfg `env:"SANITIZER_CONFIG_FILEPATH" envDefault:"/sanitizerConfig.json"`
	DataEncryptionKey string       `env:"DATA_ENCRYPTION_KEY,required"`

	DbUsername    string `env:"DB_USERNAME,required"`
	DbPassword    string `env:"DB_PASSWORD,required"`
	PGHost        string `env:"POSTGRES_HOST" envDefault:"postgres"`
	PGDatabase    string `env:"POSTGRES_DATABASE" envDefault:"reports"`
	PGRole        string `env:"POSTGRES_ROLE" envDefault:"dataexportservice"`
	DbDetailedLog bool   `env:"DB_DETAILED_LOG" envDefault:"false"`

	BoltDBFilepath string `env:"BOLT_DB_FILEPATH" envDefault:"/data/batchDataExporter.db"`

	PrometheusPushGatewayAddress string `env:"PROMETHEUS_PUSH_GATEWAY_ADDRESS" envDefault:"http://localPrometheusPushGateway:9091"`
}

// SanitizerCfg is a wrapper struct for slice with list of fields to sanitize
// to make env parser to execute custom parser without "type not supported" error
type SanitizerCfg struct {
	Slice []filesDataExporter.FieldToSanitize
}

// getConfig parses environment variables and returns pointer to config and error
func getConfig() (*Config, error) {
	common, err := config.New()
	if err != nil {
		return nil, err
	}

	cfg := &Config{Config: *common}

	parsers := map[reflect.Type]env.ParserFunc{
		reflect.TypeOf(cfg.FieldsToSanitize): parseFieldsToSanitize,
	}

	return cfg, env.ParseWithFuncs(cfg, parsers)
}

func parseFieldsToSanitize(filepath string) (interface{}, error) {
	sanitizerCfg := SanitizerCfg{
		Slice: []filesDataExporter.FieldToSanitize{},
	}

	jsonFile, err := ioutil.ReadFile(filepath)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonFile, &sanitizerCfg.Slice)
	if err != nil {
		return nil, err
	}

	return sanitizerCfg, nil
}
