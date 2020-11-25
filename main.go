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
	hclog.Default().Info("Starting coffee-service")

	err := env.Parse()
	if err != nil {
		hclog.Default().Error("Error parsing flags", "error", err)
		os.Exit(1)
	}

	var cfg *config.Config
	if cfg, err = config.NewFromEnv(); err != nil {
		hclog.Default().Error("Error reading configuration", "error", err)
		os.Exit(1)
	}

	var closer io.Closer
	if closer, err = hckit.InitGlobalTracer("coffee-service"); err != nil {
		cfg.Logger.Error("Unable to initialize Tracer", "error", err)
		os.Exit(1)
	}
	defer closer.Close()
	cfg.Logger.Debug("Tracing initialized")

	router := mux.NewRouter()
	router.Use(hckit.TracingMiddleware)

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	cfg.Logger.Debug("Registering health service")
	healthService := service.NewHealth(cfg.Logger)
	router.Handle("/health", healthService).Methods("GET")

	cfg.Logger.Debug("Initializing coffee service")
	coffeeService, err := service.NewFromConfig(cfg)
	if err != nil {
		cfg.Logger.Error("Unable to initialize coffeeService", "error", err)
		os.Exit(1)
	}
	cfg.Logger.Info("CoffeeService initialized")

	router.Handle("/coffees", coffeeService).Methods("GET")

	err = http.ListenAndServe(cfg.BindAddress, router)
	if err != nil {
		cfg.Logger.Error("Unable to start server.", "error", err)
		os.Exit(1)
	}

	cfg.Logger.Info("Started service", "bind", cfg.BindAddress, "metrics", cfg.MetricsAddress)
}
