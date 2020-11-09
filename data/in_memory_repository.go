package data

import (
	"time"

	"github.com/hashicorp/go-memdb"

	"github.com/hashicorp-demoapp/coffee-service/config"
	"github.com/hashicorp-demoapp/coffee-service/data/model"
)

// TableNameKey is a typesafe discriminator for table names
type TableNameKey string

func (t TableNameKey) String() string {
	return string(t)
}

const (
	// Ingredient is the ingredient table name
	Ingredient TableNameKey = "ingredient"
	// Coffee is the coffee table name
	Coffee TableNameKey = "coffee"
	// CoffeeIngredient is the coffee_ingredient table name
	CoffeeIngredient TableNameKey = "coffee_ingredient"
)

// InMemoryRepository implements the coffee-service.data.Repository interface
// uisng go-membdb instead of postgres.
type InMemoryRepository struct {
	db *memdb.MemDB
}

// NewInMemoryDB is the InMemoryRepository factory method. It fulfills the same
// interface as Repository, but uses go-membdb internally to provide data. NOTE,
// this interface requires build time tooling.
func NewInMemoryDB(config *config.Config) (Repository, error) {
	// Create a new data base
	db, err := memdb.NewMemDB(createSchema())
	if err != nil {
		panic(err)
	}

	repository := &InMemoryRepository{db}

	repository.loadIngredients()
	repository.loadCoffees()

	return repository, nil

}

// Find returns all coffees from the database
// Used to accept ctx opentracing.SpanContext
func (r *InMemoryRepository) Find() (model.Coffees, error) {
	txn := r.db.Txn(true)

	iter, err := txn.Get(Coffee.String(), "id")
	if err != nil {
		return nil, err
	}

	coffees := make([]model.Coffee, 0)

	for coffee := iter.Next(); coffee != nil; coffee = iter.Next() {
		coffees = append(coffees, coffee.(model.Coffee))
	}

	for _, coffee := range coffees {
		coffeeIngredients := make([]model.CoffeeIngredients, 0)

		innerIter, err := txn.Get(CoffeeIngredient.String(), "id")
		if err != nil {
			return nil, err
		}

		for ingredient := innerIter.Next(); ingredient != nil; ingredient = innerIter.Next() {
			coffeeIngredients = append(coffeeIngredients, ingredient.(model.CoffeeIngredients))
		}

		coffee.Ingredients = coffeeIngredients
	}

	return coffees, nil
}

func createSchema() *memdb.DBSchema {
	// Create the DB schema
	// TODO Update to this model with tooling.
	return &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"coffees": {
				Name: "coffee",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.IntFieldIndex{Field: "id"},
					},
				},
			},
			"indredient": {
				Name: "ingredients",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.IntFieldIndex{Field: "id"},
					},
				},
			},
			"coffee_ingredient": {
				Name: "coffee_ingredients",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.IntFieldIndex{Field: "id"},
					},
				},
			},
		},
	}
}

func (r *InMemoryRepository) loadIngredients() error {
	timestamp := time.Now().String()
	txn := r.db.Txn(true)

	// Insert some people
	ingredients := []*model.Ingredient{
		{ID: 1, Name: "Espresso'", CreatedAt: timestamp, UpdatedAt: timestamp},
		{ID: 2, Name: "Semi Skimmed Milk", CreatedAt: timestamp, UpdatedAt: timestamp},
		{ID: 3, Name: "Hot Water", CreatedAt: timestamp, UpdatedAt: timestamp},
		{ID: 4, Name: "Pumpkin Spice", CreatedAt: timestamp, UpdatedAt: timestamp},
		{ID: 5, Name: "Steamed Milk", CreatedAt: timestamp, UpdatedAt: timestamp},
	}

	for _, row := range ingredients {
		if err := txn.Insert(Ingredient.String(), row); err != nil {
			return err
		}
	}

	txn.Commit()
	return nil
}

func (r *InMemoryRepository) loadCoffees() error {
	timestamp := time.Now().String()
	txn := r.db.Txn(true)

	coffees := []*model.Coffee{
		{
			ID:          1,
			Name:        "Packer Spiced Latte",
			Teaser:      "Packed with goodness to spice up your images",
			Description: "",
			Price:       350,
			Image:       "/packer.png",
			CreatedAt:   timestamp,
			UpdatedAt:   timestamp,
		},
		{
			ID:          2,
			Name:        "Vaulatte",
			Teaser:      "Nothing gives you a safe and secure feeling like a Vaulatte",
			Description: "",
			Price:       200,
			Image:       "/vault.png",
			CreatedAt:   timestamp,
			UpdatedAt:   timestamp,
		},
		{
			ID:          3,
			Name:        "Nomadicano",
			Teaser:      "Drink one today and you will want to schedule another",
			Description: "",
			Price:       150,
			Image:       "/nomad.png",
			CreatedAt:   timestamp,
			UpdatedAt:   timestamp,
		},
		{
			ID:          4,
			Name:        "Terraspresso",
			Teaser:      "Nothing kickstarts your day like a provision of Terraspresso",
			Description: "",
			Price:       150,
			Image:       "/terraform.png",
			CreatedAt:   timestamp,
			UpdatedAt:   timestamp,
		},
		{
			ID:          5,
			Name:        "Vagrante espresso",
			Teaser:      "Stdin is not a tty",
			Description: "",
			Price:       200,
			Image:       "/vagrant.png",
			CreatedAt:   timestamp,
			UpdatedAt:   timestamp,
		},
		{
			ID:          6,
			Name:        "Connectaccino",
			Teaser:      "Discover the wonders of our meshy service",
			Description: "",
			Price:       250,
			Image:       "/consul.png",
			CreatedAt:   timestamp,
			UpdatedAt:   timestamp,
		},
	}

	for _, c := range coffees {
		if err := txn.Insert(Coffee.String(), c); err != nil {
			return err
		}
	}

	txn.Commit()
	return nil
}
func (r *InMemoryRepository) loadCoffeeIngredients() error {
	timestamp := time.Now().String()
	txn := r.db.Txn(true)

	coffeeIngredients := []*model.CoffeeIngredients{
		{
			ID:           1,
			CoffeeID:     1,
			IngredientID: 1,
			CreatedAt:    timestamp,
			UpdatedAt:    timestamp,
		},
		{
			ID:           2,
			CoffeeID:     1,
			IngredientID: 2,
			CreatedAt:    timestamp,
			UpdatedAt:    timestamp,
		},
		{
			ID:           3,
			CoffeeID:     1,
			IngredientID: 4,
			CreatedAt:    timestamp,
			UpdatedAt:    timestamp,
		},
		{
			ID:           4,
			CoffeeID:     2,
			IngredientID: 1,
			CreatedAt:    timestamp,
			UpdatedAt:    timestamp,
		},
		{
			ID:           5,
			CoffeeID:     2,
			IngredientID: 2,
			CreatedAt:    timestamp,
			UpdatedAt:    timestamp,
		},
		{
			ID:           6,
			CoffeeID:     3,
			IngredientID: 1,
			CreatedAt:    timestamp,
			UpdatedAt:    timestamp,
		},
		{
			ID:           7,
			CoffeeID:     3,
			IngredientID: 3,
			CreatedAt:    timestamp,
			UpdatedAt:    timestamp,
		},
		{
			ID:           8,
			CoffeeID:     4,
			IngredientID: 1,
			CreatedAt:    timestamp,
			UpdatedAt:    timestamp,
		},
		{
			ID:           9,
			CoffeeID:     5,
			IngredientID: 1,
			CreatedAt:    timestamp,
			UpdatedAt:    timestamp,
		},
		{
			ID:           10,
			CoffeeID:     6,
			IngredientID: 1,
			CreatedAt:    timestamp,
			UpdatedAt:    timestamp,
		},
		{
			ID:           11,
			CoffeeID:     6,
			IngredientID: 5,
			CreatedAt:    timestamp,
			UpdatedAt:    timestamp,
		},
	}

	for _, ci := range coffeeIngredients {
		if err := txn.Insert(CoffeeIngredient.String(), ci); err != nil {
			return err
		}
	}

	txn.Commit()
	return nil
}
