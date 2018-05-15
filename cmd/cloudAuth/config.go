package main

import (
	"io/ioutil"
	"reflect"
	"strings"

	"github.com/caarlos0/env"
	"gopkg.in/yaml.v2"

	"github.com/iryonetwork/wwm/config"
	"github.com/iryonetwork/wwm/gen/auth/models"
	"github.com/iryonetwork/wwm/storage/auth"
)

// Config represents configuration of cloudAuth
type Config struct {
	config.Config

	StorageEncryptionKey string `env:"STORAGE_ENCRYPTION_KEY,required"`

	BoltDBFilepath string `env:"BOLT_DB_FILEPATH" envDefault:"/data/cloudAuth.db"`

	// filepath to yaml
	ServiceCertsAndPaths Services `env:"SERVICES_FILEPATH" envDefault:"/serviceCertsAndPaths.yml"`

	// filepath to yaml
	StorageInitData auth.InitData `env:"STORAGE_INIT_DATA_FILEPATHS" envDefault:"/rolesAndRules.yml"`
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
		reflect.TypeOf(cfg.StorageInitData):      parseStorageInitData,
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

func parseStorageInitData(filepaths string) (interface{}, error) {
	filepathsSlice := strings.Split(filepaths, ",")

	storageInitData := auth.InitData{
		Locations:     []*models.Location{},
		Organizations: []*models.Organization{},
		Clinics:       []*models.Clinic{},
		Roles:         []*models.Role{},
		Rules:         []*models.Rule{},
		Users:         []*models.User{},
		UserRoles:     []*models.UserRole{},
	}

	for _, filepath := range filepathsSlice {
		initData := auth.InitData{}
		yamlFile, err := ioutil.ReadFile(filepath)

		if err != nil {
			return storageInitData, nil
		}

		err = yaml.Unmarshal(yamlFile, &initData)
		if err != nil {
			return nil, err
		}

		storageInitData.Locations = append(storageInitData.Locations, initData.Locations...)
		storageInitData.Organizations = append(storageInitData.Organizations, initData.Organizations...)
		storageInitData.Clinics = append(storageInitData.Clinics, initData.Clinics...)
		storageInitData.Roles = append(storageInitData.Roles, initData.Roles...)
		storageInitData.Rules = append(storageInitData.Rules, initData.Rules...)
		storageInitData.Users = append(storageInitData.Users, initData.Users...)
		storageInitData.UserRoles = append(storageInitData.UserRoles, initData.UserRoles...)

	}

	return storageInitData, nil
}
