package main

import (
	"fmt"
	"github.com/hashicorp-demoapp/coffee-service/config"
	"github.com/hashicorp-demoapp/coffee-service/service"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/nicholasjackson/env"

	// opentracing "github.com/opentracing/opentracing-go"
)

func main() {
	// Lifecycle event
	hclog.Default().Info("Starting coffee-service")

	// Lifecycle event
	hclog.Default().Info("Parsing environment variables")
	err := env.Parse()
	if err != nil {
		// Unrecoverable error
		hclog.Default().Error("Error parsing environment variables", "error", err)
		os.Exit(1)
	}
	// Lifecycle event
	hclog.Default().Info("Finished parsing environment variables")

	var cfg *config.Config
	// Lifecycle event
	hclog.Default().Info("Loading configuration from environment")
	if cfg, err = config.NewFromEnv(); err != nil {
		// Unrecoverable error
		hclog.Default().Error("Error reading configuration", "error", err)
		os.Exit(1)
	}
	// Lifecycle event
	cfg.Logger.Info("Finished loading configuration from environment")

	// Lifecycle event
	cfg.Logger.Info("Initializing router")
	router := mux.NewRouter()

	/*
	   Configure middleware here
	*/

	// Lifecycle event
	cfg.Logger.Info("Router initialized")

	// Lifecycle event
	cfg.Logger.Info("Registering not found handler")
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// Component initialization
	cfg.Logger.Info("Initializing HealthService")
	healthService := service.NewHealth(cfg.Logger)
	// Component initialized
	cfg.Logger.Info("HealthService initialized")

	// Lifecycle event
	cfg.Logger.Info("Registering health handler")
	router.Handle("/health", healthService).Methods("GET")
	// Lifecycle event
	cfg.Logger.Info("Health handler registered")

	// Component initialization
	cfg.Logger.Info(fmt.Sprintf("Initializing CoffeeService version %s", cfg.Version))
	coffeeService, err := service.NewCoffee(cfg)
	if err != nil {
		// Unrecoverable error
		cfg.Logger.Error("Unable to initialize CoffeeService", "error", err)
		os.Exit(1)
	}
	// Component initialized
	cfg.Logger.Info("CoffeeService initialized")

	// Lifecycle event
	cfg.Logger.Info("Registering coffee handler")
	router.Handle("/coffees", coffeeService).Methods("GET")
	// Lifecycle event
	cfg.Logger.Info("Coffee handler registered")

	// Lifecycle event
	cfg.Logger.Info("Starting service listener", "bind", cfg.BindAddress)
	err = http.ListenAndServe(cfg.BindAddress, router)
	if err != nil {
		// Unrecoverable error
		cfg.Logger.Error("Unable to start server.", "error", err)
		os.Exit(1)
	}
}
