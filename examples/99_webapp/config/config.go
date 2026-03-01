package config

import "github.com/soner3/flora"

// AppConfig holds the global configuration for our application.
type AppConfig struct {
	flora.Component // Managed by Flora as a Singleton

	AppName string
	Port    int
	DSN     string // Database Connection String
}

// NewAppConfig creates and loads the configuration (e.g., from ENV variables).
// For this tutorial, we will use hardcoded default values.
func NewAppConfig() *AppConfig {
	return &AppConfig{
		AppName: "Flora Bookstore API",
		Port:    8080,
		DSN:     "file::memory:?cache=shared",
	}
}
