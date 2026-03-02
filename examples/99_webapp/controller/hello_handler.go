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
