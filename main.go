package main

import (
	"net/http"
	"os"
	"time"

	tracing "github.com/DerekStrickland/learn-consul-jaeger/go-hckit/tracing"
	"github.com/gorilla/mux"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/nicholasjackson/env"
	opentracing "github.com/opentracing/opentracing-go"

	"github.com/hashicorp-demoapp/coffee-service/config"
	"github.com/hashicorp-demoapp/coffee-service/data"
	"github.com/hashicorp-demoapp/coffee-service/service/v1"
	"github.com/hashicorp-demoapp/coffee-service/service/v2"
	"github.com/hashicorp-demoapp/coffee-service/service/v3"
)

// Config format for application
type Config struct {
	DBConnection   string `json:"db_connection"`
	BindAddress    string `json:"bind_address"`
	MetricsAddress string `json:"metrics_address"`
}

var logger hclog.Logger
var logFormat = env.String("LOG_FORMAT", false, "text", "Log file format. [text|json]")
var logLevel = env.String("LOG_LEVEL", false, "DEBUG", "Log level for output. [info|debug|trace|warn|error]")
var logOutput = env.String("LOG_OUTPUT", false, "stdout", "Location to write log output, default is stdout, e.g. /var/log/web.log")
var dbTraceEnabled = env.Bool("DB_TRACE_ENABLED", false, false, "Add instrumentation to DB facade to generate spans for all db calls")

func main() {
	logger = hclog.Default()
	logger.Info("Starting coffee-service")

	err := env.Parse()
	if err != nil {
		hclog.Default().Error("Error parsing flags", "error", err)
		os.Exit(1)
	}

	conf := config.NewFromEnv()

	tracer, closer := tracing.Init("coffee-service")
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	// TODO: Only do this is API_VERSION is > 2
	// load the db connection
	repository, err := retryDBUntilReady()
	if err != nil {
		logger.Error("Timeout waiting for database connection")
		os.Exit(1)
	}

	router := mux.NewRouter()

	// coffeeService := service.NewCoffeeService(repository, logger)
	// router.Handle("/coffees", coffeeService).Methods("GET")

	var coffeesRouter = router.PathPrefix("/coffees").Subrouter()
	coffeesRouter.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusNotFound)
    })

	var v1Router = api.PathPrefix("/v1").Subrouter()
    v1Router.Handle("", v1.NewCoffeeService(repository, logger)).Methods("GET")

	var v2Router = api.PathPrefix("/v2").Subrouter()
    v2Router.Handle("", v2.NewCoffeeService(repository, logger)).Methods("GET")

	var v3Router = api.PathPrefix("/v3").Subrouter()
    v3Router.Handle("", v3.NewCoffeeService(logger)).Methods("GET")

	err = http.ListenAndServe(config.BindAddress, r)
	if err != nil {
		logger.Error("Unable to start server.", "error", err)
		os.Exit(1)
	}
}
