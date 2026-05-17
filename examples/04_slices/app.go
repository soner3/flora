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
	"fmt"

	"github.com/soner3/flora"
)

// =========================================================
// EXAMPLE 1: MULTI-BINDING WITH INTERFACES
// =========================================================

type Plugin interface {
	Execute()
}

type LoggerPlugin struct {
	flora.Component `flora:"order=1"` // This should run first
}

func NewLoggerPlugin() *LoggerPlugin { return &LoggerPlugin{} }
func (p *LoggerPlugin) Execute() {
	fmt.Println("   [1] LoggerPlugin: Setting up logging...")
}

type MetricsPlugin struct {
	flora.Component `flora:"order=2"` // This should run second
}

func NewMetricsPlugin() *MetricsPlugin { return &MetricsPlugin{} }
func (p *MetricsPlugin) Execute() {
	fmt.Println("   [2] MetricsPlugin: Starting metrics collection...")
}

type PluginManager struct {
	flora.Component
	plugins []Plugin // Flora bundles all structs that implement 'Plugin'
}

func NewPluginManager(plugins []Plugin) *PluginManager {
	return &PluginManager{plugins: plugins}
}

func (m *PluginManager) RunAll() {
	fmt.Printf("-> PluginManager automatically discovered %d interface plugins.\n", len(m.plugins))
	for _, p := range m.plugins {
		p.Execute()
	}
}

// =========================================================
// EXAMPLE 2: MULTI-BINDING WITH CONCRETE POINTERS/STRUCTS
// =========================================================

// A standard struct (not an interface!)
type Route struct {
	Path   string
	Method string
}

// A configuration that provides multiple instances of the same type
type RouteConfig struct {
	flora.Configuration
}

// flora:name=homeRoute, order=1
func (c *RouteConfig) ProvideHomeRoute() *Route {
	return &Route{Path: "/", Method: "GET"}
}

// flora:name=apiRoute, order=2
func (c *RouteConfig) ProvideApiRoute() *Route {
	return &Route{Path: "/api/v1/users", Method: "POST"}
}

type Router struct {
	flora.Component
	routes []*Route // Flora bundles all concrete *Route pointers here!
}

func NewRouter(routes []*Route) *Router {
	return &Router{routes: routes}
}

func (r *Router) PrintRoutes() {
	fmt.Printf("-> Router automatically discovered %d concrete routes.\n", len(r.routes))
	for _, rt := range r.routes {
		fmt.Printf("   Route: [%s] %s\n", rt.Method, rt.Path)
	}
}
