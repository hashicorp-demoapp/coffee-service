package v3

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp-demoapp/coffee-service/data"
	"github.com/hashicorp-demoapp/coffee-service/data/entities"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
)

func setupCoffeeHandler(t *testing.T) (*CoffeeService, *httptest.ResponseRecorder, *http.Request) {
	c := &data.MockRepository{}
	c.On("Find").Return(entities.Coffees{entities.Coffee{ID: 1, Name: "Test"}}, nil)

	l := hclog.Default()

	return &CoffeeService{c, l}, httptest.NewRecorder(), httptest.NewRequest("GET", "/coffees", nil)
}

func TestCoffeesReturnsCoffees(t *testing.T) {
	c, rw, r := setupCoffeeHandler(t)

	c.ServeHTTP(rw, r)

	assert.Equal(t, http.StatusOK, rw.Code)

	bd := entities.Coffees{}
	err := json.Unmarshal(rw.Body.Bytes(), &bd)
	assert.NoError(t, err)
}
