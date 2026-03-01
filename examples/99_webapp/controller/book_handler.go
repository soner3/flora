package controller

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/soner3/flora"
	"github.com/soner3/flora/examples/99_full_tutorial/service"
)

type BookHandler struct {
	flora.Component
	service *service.BookService
}

func NewBookHandler(s *service.BookService) *BookHandler {
	return &BookHandler{service: s}
}

// InitRoutes registers the /books endpoint.
func (h *BookHandler) InitRoutes(r chi.Router) {
	r.Get("/books", h.ListBooks)
}

func (h *BookHandler) ListBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.service.GetAllBooks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}
