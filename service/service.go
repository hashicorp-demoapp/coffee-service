package service

import (
	"fmt"
	"net/http"

	"github.com/hashicorp-demoapp/coffee-service/config"
	"github.com/hashicorp-demoapp/coffee-service/data"
	v1 "github.com/hashicorp-demoapp/coffee-service/service/v1"
	v2 "github.com/hashicorp-demoapp/coffee-service/service/v2"
	v3 "github.com/hashicorp-demoapp/coffee-service/service/v3"
	"github.com/hashicorp/go-hclog"
)

// TODO: Refactor to IoC

// CoffeeService is the service implementation for this microservice.
type CoffeeService struct {
	repository data.Repository
	logger     hclog.Logger
}

// NewFromConfig is a factory method that returns a configured handler for the
// configured ServiceVersion
func NewFromConfig(cfg *config.Config) (http.Handler, error) {
	var repository data.Repository
	var err error

	cfg.Logger.Debug(fmt.Sprintf("Resolving repository for version %v", cfg.Version))
	if cfg.Version == config.V1 || cfg.Version == config.V2 {
		cfg.Logger.Debug("Loading Postgres")
		if repository, err = data.NewFromConfig(cfg); err != nil {
			cfg.Logger.Debug(fmt.Sprintf("Error loading postgres %+v", err))
			return nil, err
		}
	} else if cfg.Version == config.V3 {
		cfg.Logger.Debug("Loading in memory db")
		if repository, err = data.NewInMemoryDB(cfg); err != nil {
			cfg.Logger.Debug(fmt.Sprintf("Error loading in memory db %+v", err))
			return nil, err
		}
	}

	cfg.Logger.Debug(fmt.Sprintf("Resolving service for version %v", cfg.Version))
	var handler http.Handler
	switch cfg.Version {
	case config.V1:
		handler = v1.NewCoffeeService(repository, cfg.Logger)
	case config.V2:
		handler = v2.NewCoffeeService(repository, cfg.Logger)
	case config.V3:
		handler = v3.NewCoffeeService(repository, cfg.Logger)
	}

	return handler, nil
}
