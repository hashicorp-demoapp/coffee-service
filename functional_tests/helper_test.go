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
	"github.com/hashicorp-demoapp/coffee-service/data/model"
	v1 "github.com/hashicorp-demoapp/coffee-service/service/v1"
	"github.com/hashicorp/go-hclog"
)

func (api *V1APIFeature) newService() {
	repo := data.MockRepository{}
	repo.On("Find").Return(model.Coffees{model.Coffee{ID: 1, Name: "Test"}}, nil)
	api.svc = v1.NewCoffeeService(repo, hclog.Default())
}

func (api *V1APIFeature) initRouter(method, endpoint string, userID *string) error {
	if strings.Contains(endpoint, "/coffees") {
		api.svc.ServeHTTP(api.rw, api.r)
		return nil
	}

	return nil
}

func (api *V1APIFeature) iMakeARequestTo(method, endpoint string) error {
	api.rw = httptest.NewRecorder()
	api.r = httptest.NewRequest(method, endpoint, nil)

	err := api.initRouter(method, endpoint, nil)
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

	err := api.initRouter(method, endpoint, nil)
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

	err := api.initRouter(method, endpoint, nil)
	if err != nil {
		return err
	}

	return nil
}

func (api *V1APIFeature) aListOfProductsShouldBeReturned() error {
	bd := model.Coffees{}

	err := json.Unmarshal(api.rw.Body.Bytes(), &bd)
	if err != nil {
		return err
	}
	return nil
}

func (api *V1APIFeature) thatProductsIngredientsShouldBeReturned() error {
	bd := model.Ingredients{}
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
