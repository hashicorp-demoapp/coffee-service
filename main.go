package main

import (
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

// Config format for application
type Config struct {
	DBConnection   string `json:"db_connection"`
	BindAddress    string `json:"bind_address"`
	MetricsAddress string `json:"metrics_address"`
}

func main() {
	hclog.Default().Info("Starting coffee-service")

	err := env.Parse()
	if err != nil {
		hclog.Default().Error("Error parsing flags", "error", err)
		os.Exit(1)
	}

	cfg := config.NewFromEnv()

	closer, err := hckit.InitGlobalTracer("coffee-service")
	if err != nil {
		cfg.Logger.Error("Unable to initialize Tracer", "error", err)
		os.Exit(1)
	}
	defer closer.Close()

	router := mux.NewRouter()
	router.Use(hckit.TracingMiddleware)

	// coffeeService := service.NewCoffeeService(repository, logger)
	//

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	coffeeService, err := service.NewFromConfig(cfg)
	if err != nil {
		cfg.Logger.Error("Unable to initialize coffeeService", "error", err)
		os.Exit(1)
	}

	router.Handle("/coffees", coffeeService).Methods("GET")

	err = http.ListenAndServe(cfg.BindAddress, router)
	if err != nil {
		cfg.Logger.Error("Unable to start server.", "error", err)
		os.Exit(1)
	}
}
