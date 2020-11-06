package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/nicholasjackson/env"
)

// File defines a config file
type Config struct {
	ConnectionString string
	BindAddress string
	MetricsAddress string
}

func NewFromEnv() *Config {
	username := env.String("USERNAME", false, "postgres", "Postgress username")
	password := env.String("PASSWORD", false, "password", "Postgress password")
	formatString := "host=localhost port=5432 user=%s password=%s dbname=products sslmode=disable"

	return &Config{
		ConnectionString: fmt.Sprintf(formatString, username, password)
		BindAddress: env.String("BIND_ADDRESS", false, ":9090", "Address to bind the service instance to")
		MetricsAddress: env.String("METRICS_ADDRESS", false, ":9103", "Postgress password")
	}
}

