package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	v1 "github.com/hashicorp-demoapp/coffee-service/service/v1"
)

var runTest *bool = flag.Bool("run.test", false, "Should we run the tests")

var opt = godog.Options{
	Output: colors.Colored(os.Stdout),
	Format: "progress", // can define default values
}

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opt)
}

func TestMain(m *testing.M) {
	flag.Parse()
	if !*runTest {
		return
	}

	format := "progress"
	for _, arg := range os.Args[1:] {
		fmt.Println(arg)
		if arg == "-test.v=true" { // go test transforms -v option
			format = "pretty"
			break
		}
	}

	status := godog.RunWithOptions("godog", func(s *godog.Suite) {
		FeatureContext(s)
	}, godog.Options{
		Format: format,
		Paths:  []string{"features"},
	})

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}

type V1APIFeature struct {
	svc *v1.CoffeeService
	rw  *httptest.ResponseRecorder
	r   *http.Request
}

func FeatureContext(s *godog.Suite) {
	v1api := V1APIFeature.initService()

	api.initHandlers()

	s.Step(`^the server is running$`, api.theServerIsRunning)

	s.Step(`^I make a "([^"]*)" request to "([^"]*)"$`, api.iMakeARequestTo)
	s.Step(`^I make a "([^"]*)" request to "([^"]*)" where "([^"]*)" is "([^"]*)"$`, api.iMakeARequestToWhereIs)
	s.Step(`^I make a "([^"]*)" request to "([^"]*)" with the following request body:$`, api.iMakeARequestToWithTheFollowingRequestBody)

	s.Step(`^a list of products should be returned$`, api.aListOfProductsShouldBeReturned)
	s.Step(`^a list of the product\'s ingredients should be returned$`, api.thatProductsIngredientsShouldBeReturned)

	s.Step(`^the response status should be "([^"]*)"$`, api.theResponseStatusShouldBe)
}
