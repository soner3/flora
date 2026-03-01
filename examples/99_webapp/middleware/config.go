package middleware

import (
	"fmt"
	"net/http"

	"github.com/soner3/flora"
	"github.com/soner3/flora/examples/99_full_tutorial/repository"
)

// ---------------------------------------------------------
// 1. The Interface
// Flora uses this to group all implementations into a slice.
// ---------------------------------------------------------
type Middleware interface {
	Handle(next http.Handler) http.Handler
}

// ---------------------------------------------------------
// 2. Functional Adapter (Written only once!)
// ---------------------------------------------------------
type MiddlewareFunc func(http.Handler) http.Handler

func (f MiddlewareFunc) Handle(next http.Handler) http.Handler {
	return f(next)
}

// ---------------------------------------------------------
// 3. Unique Types for Google Wire (1 line of boilerplate)
// By embedding MiddlewareFunc, these structs automatically
// implement the Middleware interface!
// ---------------------------------------------------------
type LoggerMw struct{ MiddlewareFunc }
type AuthMw struct{ MiddlewareFunc }

// ---------------------------------------------------------
// 4. The Configuration
// ---------------------------------------------------------
type MiddlewareConfig struct {
	flora.Configuration
}

// flora:order=1
func (c *MiddlewareConfig) ProvideLogger() *LoggerMw {
	// We return our unique type, containing the closure
	return &LoggerMw{func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("-> [Logger] Intercepted %s %s\n", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	}}
}

// flora:order=2
// Flora injects the database automatically!
func (c *MiddlewareConfig) ProvideAuth(repo *repository.BookRepository) *AuthMw {
	// We return our unique type, containing the closure
	return &AuthMw{func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			books, err := repo.GetAll()
			if err == nil {
				fmt.Printf("-> [Auth] Checked DB: %d books available. Auth passed!\n", len(books))
			}

			next.ServeHTTP(w, r)
		})
	}}
}
