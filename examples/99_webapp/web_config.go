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
	"database/sql"
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
func (c *WebConfig) ProvideRouter(controllers []controller.Controller, middlewares []middleware.Middleware) http.Handler {
	r := chi.NewRouter()

	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)

	for _, mw := range middlewares {
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

// flora:primary, inject(startupDB=masterDB)
func (c *WebConfig) ProvideServer(cfg *config.AppConfig, handler http.Handler, startupDB *sql.DB) (*http.Server, func(), error) {
	fmt.Println("-> [Web] Pinging Master DB before starting server...")
	if err := startupDB.Ping(); err != nil {
		return nil, nil, err
	}

	addr := fmt.Sprintf(":%d", cfg.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	cleanup := func() {
		fmt.Println("-> Gracefully shutting down the web server...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = srv.Shutdown(ctx)
	}

	return srv, cleanup, nil
}
