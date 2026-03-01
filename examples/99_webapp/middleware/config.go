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
