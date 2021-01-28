package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/cucumber/messages-go/v10"
	"github.com/gorilla/mux"
	"github.com/hashicorp-demoapp/coffee-service/data"
	"github.com/hashicorp-demoapp/coffee-service/data/entities"
	v1 "github.com/hashicorp-demoapp/coffee-service/service/v1"
	"github.com/hashicorp/go-hclog"
)

func (api *V1APIFeature) initService() {
	repo := data.MockRepository{}
	repo.On("Find").Return(entities.Coffees{entities.Coffee{ID: 1, Name: "Test"}}, nil)
	api.svc = v1.NewCoffeeService(&repo, hclog.Default())
}

func (api *V1APIFeature) initHandlers() error {
	mockRepo := &data.MockRepository{}
	mockRepo.On("Find").Return(entities.Coffees{entities.Coffee{ID: 1, Name: "Test"}}, nil)

	logger := hclog.Default()

	api.svc = v1.NewCoffeeService(mockRepo, logger)

	return nil
}

func (api *V1APIFeature) iMakeARequestTo(method, endpoint string) error {
	api.rw = httptest.NewRecorder()
	api.r = httptest.NewRequest(method, endpoint, nil)

	err := api.initHandlers()
	if err != nil {
		return err
	}

	return nil
}

func (api *V1APIFeature) iMakeARequestToWhereIs(method, endpoint string, attribute, value string) error {
	api.rw = httptest.NewRecorder()
	api.r = httptest.NewRequest(method, endpoint, nil)

	vars := map[string]string{attribute: value}
	api.r = mux.SetURLVars(api.r, vars)

	err := api.initHandlers()
	if err != nil {
		return err
	}

	return nil
}

func (api *V1APIFeature) iMakeARequestToWithTheFollowingRequestBody(method, endpoint string, body *messages.PickleStepArgument_PickleDocString) error {
	api.rw = httptest.NewRecorder()
	api.r = httptest.NewRequest(method, endpoint, nil)

	rb := strings.NewReader(body.Content)
	api.r.Body = ioutil.NopCloser(rb)

	err := api.initHandlers()
	if err != nil {
		return err
	}

	return nil
}

func (api *V1APIFeature) aListOfProductsShouldBeReturned() error {
	bd := entities.Coffees{}

	err := json.Unmarshal(api.rw.Body.Bytes(), &bd)
	if err != nil {
		return err
	}
	return nil
}

func (api *V1APIFeature) thatProductsIngredientsShouldBeReturned() error {
	bd := entities.Ingredients{}
	err := json.Unmarshal(api.rw.Body.Bytes(), &bd)
	if err != nil {
		return err
	}
	return nil
}

func (api *V1APIFeature) theResponseStatusShouldBe(statusCode string) error {
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
