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

package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/soner3/flora"
	"github.com/soner3/flora/examples/99_webapp/config"
	"github.com/soner3/flora/examples/99_webapp/controller"
	"github.com/soner3/flora/examples/99_webapp/middleware"
)

// WebConfig utilizes flora.Configuration to provide third-party types
// like the Chi router and the standard library http.Server to our DI container.
type WebConfig struct {
	flora.Configuration
}

// ProvideRouter sets up our HTTP router.
// flora:primary
func (c *WebConfig) ProvideRouter(controllers []controller.Controller, middleware []middleware.Middleware) http.Handler {
	r := chi.NewRouter()

	// Add some standard middleware
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)

	// Apply middleware in the correct order
	for _, mw := range middleware {
		r.Use(mw.Handle)
	}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the Flora Bookstore API!"))
	})

	for _, ctrl := range controllers {
		ctrl.InitRoutes(r)
	}

	return r
}

// ProvideServer wires the AppConfig and the Router together to create the HTTP Server.
// flora:primary
func (c *WebConfig) ProvideServer(cfg *config.AppConfig, handler http.Handler) (*http.Server, func(), error) {
	addr := fmt.Sprintf(":%d", cfg.Port)

	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	// Flora handles graceful shutdowns natively!
	// When defer cleanup() is called in main, this function will execute.
	cleanup := func() {
		fmt.Println("-> Gracefully shutting down the web server...")

		// Give active requests 5 seconds to finish before forcing shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_ = srv.Shutdown(ctx)
	}

	return srv, cleanup, nil
}
