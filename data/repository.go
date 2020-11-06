package data

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"
	"contrib.go.opencensus.io/integrations/ocsql"

	"github.com/hashicorp-demoapp/coffee-service/data/model"
)

// CoffeesRepository is the command/query interface this respository supports.
type CoffeesRepository interface {
	FindCoffees(ctx opentracing.SpanContext) (model.Coffees, error)
}

// DBConnection is a connection to the DB.
type DBConnection struct {
	db *sqlx.DB
}

// New is the CoffeeRepository factory method. It encapsulates the Postgres DB.
// It will attempt to create a connection, and keep retrying the database connection
// until successful or it timeuts. When running the application on a scheduler it
// is possible (likely) that the app will come up before the database, this can
// cause the app to go into a CrashLoopBackoff cycle.
// TODO: Read git history to see if this retry. I'm suspecting this is in place
// to allow behavioral tests to not fail while the environment spins up.
func New(connection string, withTracing bool) (data.CoffeesRepository, error) {
	st := time.Now()
	dt := 1 * time.Second  // this should be an exponential backoff
	mt := 60 * time.Second // max time to wait of the DB connection

	for {
		logger.Info("Using connection: " + conf.DBConnection)
		var repository *DBConnection
		var err error

		if withTracing {
			repository, err = data.NewWithTracing(connection)
		} else {
			repository, err = data.New(connection)
		}
		if err == nil {
			return repository, nil
		}

		logger.Error("Unable to connect to database", "error", err)

		// check if max time has elapsed
		if time.Now().Sub(st) > mt {
			return nil, err
		}

		// retry
		time.Sleep(dt)
	}
}

// new creates a new connection to the database
func new(connection string) (*DBConnection, error) {
	db, err := sqlx.Connect("postgres", connection)
	if err != nil {
		return nil, err
	}

	return &DBConnection{db}, nil
}

// newWithTracing wraps the connection with OpenCensus instrumentation
// to allow db query inspection from OpenCensus compliant backends.
func newWithTracing(connection string) (*DBConnection, error) {
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
	dbx := sqlx.NewDB(db, "postgres")

	return &DBConnection{dbx}, nil
}

// FindCoffees returns all products from the database
func (c *DBConnection) FindCoffees(ctx opentracing.SpanContext) (model.Coffees, error) {
	// tracer := opentracing.GlobalTracer()
	// span := tracer.StartSpan("repository-find-coffees", ext.RPCServerOption(ctx))
	// defer span.Finish()

	coffees := model.Coffees{}

	// selectCoffeesSpan := tracer.StartSpan("db-select-coffees", ext.RPCServerOption(ctx))
	// defer selectCoffeesSpan.Finish()
	err := c.db.Select(&coffees, "SELECT * FROM coffees")
	if err != nil {
		return nil, err
	}
	// selectCoffeesSpan.LogFields(
	// 	otlog.String("event", "coffees-result-count"),
	// 	otlog.String("value", fmt.Sprintf("%d", len(coffees))),
	// )

	// fetch the ingredients for each coffee
	// selectIngredientsSpan := tracer.StartSpan("db-select-ingredients", ext.RPCServerOption(ctx))
	// defer selectIngredientsSpan.Finish()

	for n, coffee := range coffees {
		coffeeIngredients := []model.CoffeeIngredients{}

		err := c.db.Select(&coffeeIngredients, "SELECT ingredient_id FROM coffee_ingredients WHERE coffee_id=$1", coffee.ID)
		if err != nil {
			return nil, err
		}

		coffees[n].Ingredients = coffeeIngredients

		json, err := json.Marshal(coffees[n])
		if err != nil {
			ext.LogError(selectCoffeesSpan, err)
		}

		// selectIngredientsSpan.LogFields(
		// 	otlog.String("event", "ingredients-populated"),
		// 	otlog.String("value", string(json)))
	}

	return coffees, nil
}
