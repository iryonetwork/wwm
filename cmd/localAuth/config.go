package main

import (
	"io/ioutil"
	"reflect"

	"github.com/caarlos0/env"
	"gopkg.in/yaml.v2"

	"github.com/iryonetwork/wwm/config"
)

// Config represents configuration of localAuth
type Config struct {
	config.Config

	StorageEncryptionKey string `env:"STORAGE_ENCRYPTION_KEY,required"`

	BoltDBFilepath string `env:"BOLT_DB_FILEPATH" envDefault:"/data/localAuth.db"`

	CloudAuthHost string `env:"CLOUD_AUTH_HOST" envDefault:"cloudAuth"`
	CloudAuthPath string `env:"CLOUD_AUTH_PATH" envDefault:"auth"`

	AuthSyncKeyPath  string `env:"AUTH_SYNC_KEY_PATH,required"`
	AuthSyncCertPath string `env:"AUTH_SYNC_CERT_PATH,required"`

	// filepath to yaml
	ServiceCertsAndPaths Services `env:"SERVICES_FILEPATH" envDefault:"/serviceCertsAndPaths.yml"`
}

// Services is a wrapper struct for map of allowed services certs and paths
// to make env parser to execute custom parser without "type not suppoerted" error
type Services struct {
	Map map[string][]string
}

// GetConfig parses environment variables and returns pointer to config and error
func GetConfig() (*Config, error) {
	common, err := config.New()
	if err != nil {
		return nil, err
	}

	cfg := &Config{Config: *common}

	parsers := map[reflect.Type]env.ParserFunc{
		reflect.TypeOf(cfg.ServiceCertsAndPaths): parseServiceCertsAndPaths,
	}

	return cfg, env.ParseWithFuncs(cfg, parsers)
}

func parseServiceCertsAndPaths(filepath string) (interface{}, error) {
	serviceCertsAndPaths := Services{
		Map: make(map[string][]string),
	}

	yamlFile, err := ioutil.ReadFile(filepath)
	if err != nil {
		return serviceCertsAndPaths, nil
	}

	err = yaml.Unmarshal(yamlFile, &serviceCertsAndPaths.Map)
	if err != nil {
		return nil, err
	}

	return serviceCertsAndPaths, nil
}
