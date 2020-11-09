package v1

import (
	"fmt"
	"net/http"

	hclog "github.com/hashicorp/go-hclog"

	"github.com/hashicorp-demoapp/coffee-service/data"
)

// CoffeeService is the service implementation for this microservice.
type CoffeeService struct {
	repository data.Repository
	logger     hclog.Logger
}

// NewCoffeeService is a factory method that returns a new instance of the CoffeeService.
func NewCoffeeService(repository data.Repository, l hclog.Logger) *CoffeeService {
	return &CoffeeService{repository, l}
}

// ServeHTTP handles incoming requests for the api coffees route
func (c *CoffeeService) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	c.logger.Info("Handle Coffees")

	coffees, err := c.repository.FindCoffees()
	if err != nil {
		c.logger.Error("Unable to get coffees from database", "error", err)
		http.Error(rw, "Unable to get coffees from database", http.StatusInternalServerError)
	}
	c.logger.Info(fmt.Sprintf("Found %d coffees", len(coffees)))

	coffeesJSON, err := coffees.ToJSON()
	if err != nil {
		c.logger.Error("Unable to convert coffees to JSON", "error", err)
		http.Error(rw, "Unable to convert coffees to JSON", http.StatusInternalServerError)
	}

	rw.Write(coffeesJSON)
}
