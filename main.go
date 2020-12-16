package main

import (
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	hckit "github.com/hashicorp-demoapp/go-hckit"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/nicholasjackson/env"

	// opentracing "github.com/opentracing/opentracing-go"

	"github.com/hashicorp-demoapp/coffee-service/config"
	"github.com/hashicorp-demoapp/coffee-service/service"
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
	hclog.Default().Info("Finished loading configuration from environment")

	var closer io.Closer
	// Lifecycle event
	hclog.Default().Info("Initializing tracing")
	if closer, err = hckit.InitGlobalTracer("coffee-service"); err != nil {
		// Unrecoverable error
		cfg.Logger.Error("Unable to initialize Tracer", "error", err)
		os.Exit(1)
	}
	defer closer.Close()
	// Lifecycle event
	cfg.Logger.Info("Tracing initialized")

	// Lifecycle event
	cfg.Logger.Info("Initializing router")
	router := mux.NewRouter()

	// Add middleware config here
	router.Use(hckit.TracingMiddleware)

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	// Lifecycle event
	cfg.Logger.Info("Router initialized")

	// Lifecycle event
	cfg.Logger.Info("Registering health service")
	healthService := service.NewHealth(cfg.Logger)
	router.Handle("/health", healthService).Methods("GET")
	// Lifecycle event
	cfg.Logger.Info("Health service registered")

	// Component initialization
	cfg.Logger.Info("Initializing coffee service")
	coffeeService, err := service.NewFromConfig(cfg)
	if err != nil {
		// Unrecoverable error
		cfg.Logger.Error("Unable to initialize coffeeService", "error", err)
		os.Exit(1)
	}
	// Component initialized
	cfg.Logger.Info("CoffeeService initialized")

	// Lifecycle event
	cfg.Logger.Info("Registering coffee service")
	router.Handle("/coffees", coffeeService).Methods("GET")
	// Lifecycle event
	cfg.Logger.Info("Coffee service registered")

	err = http.ListenAndServe(cfg.BindAddress, router)
	if err != nil {
		// Unrecoverable error
		cfg.Logger.Error("Unable to start server.", "error", err)
		os.Exit(1)
	}

	// Lifecycle event
	cfg.Logger.Info("Started service", "bind", cfg.BindAddress, "metrics", cfg.MetricsAddress)
}
