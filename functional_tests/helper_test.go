package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"

	"github.com/cucumber/messages-go/v10"
	"github.com/gorilla/mux"
	"github.com/hashicorp-demoapp/product-api-go/data"
	"github.com/hashicorp-demoapp/product-api-go/data/model"
	"github.com/hashicorp-demoapp/product-api-go/handlers"
	"github.com/hashicorp/go-hclog"
)

func (api *apiFeature) initHandlers() {
	// Coffee
	mc := &data.MockConnection{}
	mc.On("GetProducts").Return(model.Coffees{model.Coffee{ID: 1, Name: "Test"}}, nil)
	mc.On("GetIngredientsForCoffee").Return(model.Ingredients{
		model.Ingredient{ID: 1, Name: "Coffee"},
		model.Ingredient{ID: 2, Name: "Milk"},
		model.Ingredient{ID: 2, Name: "Sugar"},
	})

	l := hclog.Default()

	api.mc = mc
	api.hc = handlers.NewCoffee(mc, l)
}

func (api *apiFeature) initRouter(method, endpoint string, userID *string) error {
	if strings.Contains(endpoint, "/coffees") {
		api.hc.ServeHTTP(api.rw, api.r)
		return nil
	}

	return nil
}

func (api *apiFeature) theServerIsRunning() error {
	connected, err := api.mc.IsConnected()
	if err != nil {
		return err
	}
	if connected == false {
		return fmt.Errorf("Mock connection is not connected")
	}
	return nil
}

func (api *apiFeature) iMakeARequestTo(method, endpoint string) error {
	api.rw = httptest.NewRecorder()
	api.r = httptest.NewRequest(method, endpoint, nil)

	err := api.initRouter(method, endpoint, nil)
	if err != nil {
		return err
	}

	return nil
}

func (api *apiFeature) iMakeARequestToWhereIs(method, endpoint string, attribute, value string) error {
	api.rw = httptest.NewRecorder()
	api.r = httptest.NewRequest(method, endpoint, nil)

	vars := map[string]string{attribute: value}
	api.r = mux.SetURLVars(api.r, vars)

	err := api.initRouter(method, endpoint, nil)
	if err != nil {
		return err
	}

	return nil
}

func (api *apiFeature) iMakeARequestToWithTheFollowingRequestBody(method, endpoint string, body *messages.PickleStepArgument_PickleDocString) error {
	api.rw = httptest.NewRecorder()
	api.r = httptest.NewRequest(method, endpoint, nil)

	rb := strings.NewReader(body.Content)
	api.r.Body = ioutil.NopCloser(rb)

	err := api.initRouter(method, endpoint, nil)
	if err != nil {
		return err
	}

	return nil
}

func (api *apiFeature) aListOfProductsShouldBeReturned() error {
	bd := model.Coffees{}

	err := json.Unmarshal(api.rw.Body.Bytes(), &bd)
	if err != nil {
		return err
	}
	return nil
}

func (api *apiFeature) thatProductsIngredientsShouldBeReturned() error {
	bd := model.Ingredients{}
	err := json.Unmarshal(api.rw.Body.Bytes(), &bd)
	if err != nil {
		return err
	}
	return nil
}

func (api *apiFeature) theResponseStatusShouldBe(statusCode string) error {
	switch statusCode {
	case "OK":
		if api.rw.Code != http.StatusOK {
			return fmt.Errorf("expected status code does not match actual, %v vs. %v", http.StatusOK, api.rw.Code)
		}
	default:
		return fmt.Errorf("Status Code is not valid, %s", statusCode)
	}
	return nil
}
