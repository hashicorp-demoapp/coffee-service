package data

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	// otlog "github.com/opentracing/opentracing-go/log"
	"contrib.go.opencensus.io/integrations/ocsql"

	"github.com/hashicorp-demoapp/coffee-service/config"
	"github.com/hashicorp-demoapp/coffee-service/data/entities"
)

// Repository is the command/query interface this respository supports.
type Repository interface {
	Find() (entities.Coffees, error)
}

// PostgresRepository is a postgres implementation of the Repository interface.
type PostgresRepository struct {
	db *sqlx.DB
}

// NewFromConfig is the CoffeeRepository factory method. It encapsulates the Postgres DB.
// It will attempt to create a connection, and keep retrying the database connection
// until successful or times out. When running the application on a scheduler it
// is possible (likely) that the app will come up before the database, this can
// cause the app to go into a CrashLoopBackoff cycle. By defining a retry loop,
// we are implementing circuit breaker, rather than just crashing on startup
// if the db is unavailable.
// TODO: this whole thing needs to be addressed.  Probably we want to move
// the circuit breaking back to the lifecycle in main, and have this just
// test IsConnected() or just try to make the call.
func NewFromConfig(cfg *config.Config) (Repository, error) {
	st := time.Now()
	dt := 1 * time.Second  // this should be an exponential backoff
	mt := 60 * time.Second // max time to wait of the DB connection

	for {
		var repository *PostgresRepository
		var err error

		if cfg.DBTraceEnabled {
			repository, err = newPostgresWithTracing(cfg.ConnectionString)
		} else {
			repository, err = newPostgres(cfg.ConnectionString)
		}
		if err == nil {
			return repository, nil
		}

		cfg.Logger.Error("Unable to connect to database", "error", err)

		// check if max time has elapsed
		if time.Now().Sub(st) > mt {
			return nil, err
		}

		// retry
		time.Sleep(dt)
	}
}

// new creates a new connection to the database
func newPostgres(connection string) (*PostgresRepository, error) {
	db, err := sqlx.Connect("postgres", connection)
	if err != nil {
		return nil, err
	}

	return &PostgresRepository{db: db}, nil
}

// newWithTracing wraps the connection with OpenCensus instrumentation
// to allow db query inspection from OpenCensus compliant backends.
func newPostgresWithTracing(connection string) (*PostgresRepository, error) {
	// Lifted from here:  https://github.com/opencensus-integrations/ocsql#jmoironsqlx
	// Register our ocsql wrapper for the provided Postgres driver.
	driverName, err := ocsql.Register("postgres", ocsql.WithAllTraceOptions())
	if err != nil {
		return nil, err
	}

	// Connect to a Postgres database using the ocsql driver wrapper.
	// TODO: Test this - might need to use url format conn string like so
	// "postgres://localhost:5432/my_database"
	db, err := sql.Open(driverName, connection)
	if err != nil {
		return nil, err
	}

	// Wrap our *sql.DB with sqlx. use the original db driver name!!!
	dbx := sqlx.NewDb(db, "postgres")

	return &PostgresRepository{dbx}, nil
}

// Find returns all products from the database
// Used to accept ctx opentracing.SpanContext
func (r *PostgresRepository) Find() (entities.Coffees, error) {
	coffees := entities.Coffees{}

	err := r.db.Select(&coffees, "SELECT * FROM coffee")
	if err != nil {
		return nil, err
	}

	for n, coffee := range coffees {
		coffeeIngredients := []entities.CoffeeIngredients{}

		err := r.db.Select(&coffeeIngredients, "SELECT ingredient_id FROM coffee_ingredient WHERE coffee_id=$1", coffee.ID)
		if err != nil {
			return nil, err
		}

		coffees[n].Ingredients = coffeeIngredients
	}

	return coffees, nil
}
