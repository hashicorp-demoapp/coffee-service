package data

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	// otlog "github.com/opentracing/opentracing-go/log"
	"contrib.go.opencensus.io/integrations/ocsql"

	"github.com/hashicorp-demoapp/coffee-service/config"
	"github.com/hashicorp-demoapp/coffee-service/data/model"
)

// Repository is the command/query interface this respository supports.
type Repository interface {
	Find() (model.Coffees, error)
}

// PostgresRepository is a postgres implementation of the Repository interface.
type PostgresRepository struct {
	db *sqlx.DB
}

// NewFromConfig is the CoffeeRepository factory method. It encapsulates the Postgres DB.
// It will attempt to create a connection, and keep retrying the database connection
// until successful or it timeuts. When running the application on a scheduler it
// is possible (likely) that the app will come up before the database, this can
// cause the app to go into a CrashLoopBackoff cycle.
// TODO: Read git history to see if this retry. I'm suspecting this is in place
// to allow behavioral tests to not fail while the environment spins up.
func NewFromConfig(config *config.Config) (Repository, error) {
	st := time.Now()
	dt := 1 * time.Second  // this should be an exponential backoff
	mt := 60 * time.Second // max time to wait of the DB connection

	for {
		var repository *PostgresRepository
		var err error

		if config.DBTraceEnabled {
			repository, err = newPostgresWithTracing(config.ConnectionString)
		} else {
			repository, err = newPostgres(config.ConnectionString)
		}
		if err == nil {
			return repository, nil
		}

		config.Logger.Error("Unable to connect to database", "error", err)

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
func (r *PostgresRepository) Find() (model.Coffees, error) {
	coffees := model.Coffees{}

	err := r.db.Select(&coffees, "SELECT * FROM coffees")
	if err != nil {
		return nil, err
	}

	for n, coffee := range coffees {
		coffeeIngredients := []model.CoffeeIngredients{}

		err := r.db.Select(&coffeeIngredients, "SELECT ingredient_id FROM coffee_ingredients WHERE coffee_id=$1", coffee.ID)
		if err != nil {
			return nil, err
		}

		coffees[n].Ingredients = coffeeIngredients
	}

	return coffees, nil
}
