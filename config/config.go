package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/hashicorp/go-hclog"
)

// VersionKey supports a type safe string discriminator for service version.
type VersionKey string

const (
	// V1 of the api
	V1 VersionKey = "v1"
	// V2 of the api
	V2 VersionKey = "v2"
	// V3 of the api
	V3 VersionKey = "v3"
	// VUnknown indicates version cannot be resolved
	VUnknown VersionKey = "Unknown"
)

// String casts a VersionKey to string
func (v VersionKey) String() string {
	return string(v)
}

// VersionKeyFromString casts a string to an VersionKey
func VersionKeyFromString(key string) VersionKey {
	switch key {
	case V1.String():
		return V1
	case V2.String():
		return V2
	case V3.String():
		return V3
	}

	return VUnknown
}

// EnvVarKey supports a type safe string discriminator for environment variables.
type EnvVarKey string

// String casts an EnvVarKey to string
func (e EnvVarKey) String() string {
	return string(e)
}

// EnvVarKeyFromString casts a string to an EnvVarKey
func EnvVarKeyFromString(key string) EnvVarKey {
	switch key {
	case Username.String():
		return Username
	case Password.String():
		return Password
	case LogFormat.String():
		return LogFormat
	case LogLevel.String():
		return LogLevel
	case BindAddress.String():
		return BindAddress
	case MetricsAddress.String():
		return MetricsAddress
	case DBTraceEnabled.String():
		return DBTraceEnabled
	case Version.String():
		return Version
	}

	return Unknown
}

const (
	// Username EnvVarKey
	Username EnvVarKey = "USERNAME"
	// Password EnvVarKey
	Password EnvVarKey = "PASSWORD"
	// LogFormat EnvVarKey
	LogFormat EnvVarKey = "LOG_FORMAT"
	// LogLevel EnvVarKey
	LogLevel EnvVarKey = "LOG_LEVEL"
	//BindAddress EnvVarKey
	BindAddress EnvVarKey = "BIND_ADDRESS"
	//MetricsAddress EnvVarKey
	MetricsAddress EnvVarKey = "METRICS_ADDRESS"
	// DBTraceEnabled EnvVarKey
	DBTraceEnabled EnvVarKey = "DB_TRACE_ENABLED"
	// Version EnvVarKey
	Version EnvVarKey = "VERSION"
	// Unknown EnvVarKey
	Unknown EnvVarKey = "UNKNOWN"
)

// Config defines the service runtime configuration
type Config struct {
	ConnectionString string
	BindAddress      string
	MetricsAddress   string
	DBTraceEnabled   bool
	Logger           hclog.Logger
	Version          VersionKey
}

// NewFromEnv aggregates the environment variables to a datastructure.
func NewFromEnv() (*Config, error) {
	// TODO: error handling
	username := os.Getenv(Username.String())
	password := os.Getenv(Password.String())
	formatString := "host=localhost port=5432 user=%s password=%s dbname=products sslmode=disable"
	bindAddress := os.Getenv(BindAddress.String())
	metricsAddress := os.Getenv(MetricsAddress.String())
	// TODO: Think about moving towards opentelemetry interfaces.
	// Output: *env.String("LOG_OUTPUT", false, "stdout", "Location to write log output, default is stdout, e.g. /var/log/web.log"),
	logLevel := os.Getenv(LogLevel.String())
	isJSONFormat := strings.ToLower(os.Getenv(LogFormat.String())) == "json"
	logger := hclog.New(&hclog.LoggerOptions{
		Name:       "coffee-service",
		JSONFormat: isJSONFormat,
		Level:      hclog.LevelFromString(logLevel),
	})

	dbTraceEnabled := false
	var err error

	dteRaw := os.Getenv(DBTraceEnabled.String())
	if len(dteRaw) < 4 {
		dbTraceEnabled = false
	} else {
		if dbTraceEnabled, err = strconv.ParseBool(os.Getenv(DBTraceEnabled.String())); err != nil {
			logger.Error(fmt.Sprintf("Unable to parse %s", DBTraceEnabled.String()), "error", err)
		}
	}
	versionKey := VersionKeyFromString(os.Getenv(Version.String()))

	return &Config{
		ConnectionString: fmt.Sprintf(formatString, username, password),
		BindAddress:      bindAddress,
		MetricsAddress:   metricsAddress,
		DBTraceEnabled:   dbTraceEnabled,
		Logger:           logger,
		Version:          versionKey,
	}, nil
}
