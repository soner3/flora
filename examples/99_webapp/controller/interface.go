package controller

import (
	"github.com/go-chi/chi/v5"
)

// Controller is the common interface for all our HTTP endpoints.
type Controller interface {
	// InitRoutes is called by the Server to register the handler's endpoints.
	InitRoutes(r chi.Router)
}
