package main

import (
	"io/ioutil"
	"reflect"
	"time"

	"github.com/caarlos0/env"
	"gopkg.in/yaml.v2"

	"github.com/iryonetwork/wwm/config"
	"github.com/iryonetwork/wwm/service/statusReporter/polling"
)

// Config represents configuration of localStatusReporter
type Config struct {
	config.Config
	DefaultTimeout time.Duration        `env:"DEFAULT_TIMEOUT" envDefault:"1s"`
	CountThreshold int                  `env:"COUNT_THRESHOLD" envDefault:"3"`
	Interval       time.Duration        `env:"INTERVAL" envDefault:"3s"`
	StatusValidity time.Duration        `env:"STSTUS_VALIDITY" envDefault:"20s"`
	Components     ComponentsCollection `env:"COMPONENTS_FILEPATH" envDefault:"/components.yml"`
}

// ComponentsCollection represents collection of components that are checked by statusReporter
type ComponentsCollection struct {
	Local    map[string]Component `yaml:"Local"`
	Cloud    map[string]Component `yaml:"Cloud"`
	External map[string]Component `yaml:"External"`
}

// Component represent component that is checked by statusReporter
type Component struct {
	URLType        polling.URLType `yaml:"urlType"`
	URL            string          `yaml:"url"`
	Timeout        time.Duration   `yaml:"timeout"`
	CountThreshold int             `yaml:"countThreshold"`
	Interval       time.Duration   `yaml:"interval"`
	StatusValidity time.Duration   `yaml:"statusValidity"`
}

// GetConfig parses environment variables and returns pointer to config and error
func GetConfig() (*Config, error) {
	common, err := config.New()
	if err != nil {
		return nil, err
	}

	cfg := &Config{Config: *common}

	parsers := map[reflect.Type]env.ParserFunc{
		reflect.TypeOf(cfg.Components): parseComponents,
	}
	return cfg, env.ParseWithFuncs(cfg, parsers)
}

func parseComponents(filepath string) (interface{}, error) {
	components := ComponentsCollection{}

	yamlFile, err := ioutil.ReadFile(filepath)
	if err != nil {
		return components, nil
	}

	err = yaml.Unmarshal(yamlFile, &components)
	if err != nil {
		return nil, err
	}

	return components, nil
}

func (c Component) getConfig(defaultTimeout time.Duration, defaultCfg *polling.Cfg) (polling.Cfg, time.Duration) {
	timeout := defaultTimeout
	pollingCfg := polling.Cfg{}
	if defaultCfg != nil {
		pollingCfg = *defaultCfg
	}

	if c.Interval != time.Duration(0) {
		pollingCfg.Interval = &c.Interval
	}
	if c.CountThreshold != 0 {
		pollingCfg.CountThreshold = &c.CountThreshold
	}
	if c.StatusValidity != time.Duration(0) {
		pollingCfg.Interval = &c.StatusValidity
	}
	if c.Timeout != time.Duration(0) {
		timeout = c.Timeout
	}
	return pollingCfg, timeout
}
