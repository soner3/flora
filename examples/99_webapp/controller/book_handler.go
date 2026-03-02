// Copyright © 2026 Soner Astan astansoner@gmail.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
