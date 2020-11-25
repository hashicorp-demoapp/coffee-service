package service

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/go-hclog"
)

// TODO: Move this to hckit.

// HealthService is an HTTP Handler for health checking
type HealthService struct {
	logger hclog.Logger
}

// NewHealth creates a new Health handler
func NewHealth(l hclog.Logger) *HealthService {
	return &HealthService{l}
}

// ServeHTTP implements the handler interface
func (h *HealthService) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(rw, "%s", "ok")
}
