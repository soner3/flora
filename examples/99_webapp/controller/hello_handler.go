package controller

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/soner3/flora"
)

type HelloHandler struct {
	flora.Component
}

func NewHelloHandler() *HelloHandler {
	return &HelloHandler{}
}

// InitRoutes registers the /hello endpoint.
func (h *HelloHandler) InitRoutes(r chi.Router) {
	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World from Flora!"))
	})
}
