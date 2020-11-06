package service

import (
	"net/http"

	"github.com/hashicorp-demoapp/coffee-service/config"
	"github.com/hashicorp-demoapp/coffee-service/data"
	v1 "github.com/hashicorp-demoapp/coffee-service/service/v1"
	v2 "github.com/hashicorp-demoapp/coffee-service/service/v2"
	v3 "github.com/hashicorp-demoapp/coffee-service/service/v3"
)

// TODO: Refactor to IoC

// NewFromConfig is a factory method that returns a configured handler for the
// configured ServiceVersion
func NewFromConfig(cfg *config.Config) (http.Handler, error) {
	var repository data.CoffeesRepository

	if cfg.Version == config.V1 || cfg.Version == config.V2 {
		repository, err := data.NewFromConfig(cfg)
		if err != nil {
			cfg.Logger.Error("Timeout waiting for database connection")
			return nil, err
		}
	}

	var handler http.Handler
	switch cfg.Version {
	case config.V1:
		handler = v1.NewCoffeeService(repository, cfg.Logger)
	case config.V2:
		handler = v2.NewCoffeeService(repository, cfg.Logger)
	case config.V3:
		handler = v3.NewCoffeeService(cfg.Logger)
	}

	return handler, nil
}
