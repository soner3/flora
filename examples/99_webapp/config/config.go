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
