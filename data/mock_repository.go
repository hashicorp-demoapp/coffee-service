package data

import (
	"github.com/stretchr/testify/mock"

	"github.com/hashicorp-demoapp/coffee-service/data/model"
)

// MockRepository is a mock connection object for unit tests.
type MockRepository struct {
	mock.Mock
}

// Find mock stub
func (r *MockRepository) Find() (model.Coffees, error) {
	args := r.Called()

	if m, ok := args.Get(0).(model.Coffees); ok {
		return m, args.Error(1)
	}

	return nil, args.Error(1)
}
