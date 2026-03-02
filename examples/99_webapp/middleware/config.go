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

package middleware

import (
	"fmt"
	"net/http"

	"github.com/soner3/flora"
	"github.com/soner3/flora/examples/99_full_tutorial/repository"
)

// ---------------------------------------------------------
// 1. The Interface
// ---------------------------------------------------------
type Middleware interface {
	Handle(next http.Handler) http.Handler
}

// ---------------------------------------------------------
// 2. Functional Adapter
// Prevents boilerplate! We only write this once.
// ---------------------------------------------------------
type MiddlewareFunc func(http.Handler) http.Handler

func (f MiddlewareFunc) Handle(next http.Handler) http.Handler {
	return f(next)
}

// ---------------------------------------------------------
// 3. The Configuration
// ---------------------------------------------------------
type MiddlewareConfig struct {
	flora.Configuration
}

// flora:order=1
func (c *MiddlewareConfig) ProvideLogger() Middleware {
	return MiddlewareFunc(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("-> [Logger] Intercepted %s %s\n", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	})
}

// flora:order=2
func (c *MiddlewareConfig) ProvideAuth(repo *repository.BookRepository) Middleware {
	return MiddlewareFunc(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			books, err := repo.GetAll()
			if err == nil {
				fmt.Printf("-> [Auth] Checked DB: %d books available. Auth passed!\n", len(books))
			}

			next.ServeHTTP(w, r)
		})
	})
}
