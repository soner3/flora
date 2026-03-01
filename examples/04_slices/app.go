package main

import (
	"fmt"

	"github.com/soner3/flora"
)

// ---------------------------------------------------------
// 1. The Interface
// ---------------------------------------------------------
type Plugin interface {
	Execute()
}

// ---------------------------------------------------------
// 2. The Implementations (Plugins)
// ---------------------------------------------------------

// We use the "order" tag to tell Flora exactly in which position
// this plugin should appear in the slice.
type LoggerPlugin struct {
	flora.Component `flora:"order=1"` // This should run first
}

func NewLoggerPlugin() *LoggerPlugin {
	return &LoggerPlugin{}
}

func (p *LoggerPlugin) Execute() {
	fmt.Println("   [1] LoggerPlugin: Setting up logging...")
}

type MetricsPlugin struct {
	flora.Component `flora:"order=2"` // This should run second
}

func NewMetricsPlugin() *MetricsPlugin {
	return &MetricsPlugin{}
}

func (p *MetricsPlugin) Execute() {
	fmt.Println("   [2] MetricsPlugin: Starting metrics collection...")
}

// ---------------------------------------------------------
// 3. The Consumer
// ---------------------------------------------------------
type PluginManager struct {
	flora.Component

	// We don't ask for a single Plugin, we ask for a SLICE of Plugins!
	plugins []Plugin
}

// Flora scans the entire codebase, finds both LoggerPlugin and MetricsPlugin,
// bundles them into a slice, sorts them by their "order" tag, and injects them here.
func NewPluginManager(plugins []Plugin) *PluginManager {
	return &PluginManager{plugins: plugins}
}

func (m *PluginManager) RunAll() {
	fmt.Printf("-> PluginManager automatically discovered %d plugins.\n", len(m.plugins))

	for _, p := range m.plugins {
		p.Execute() // They will be executed in the exact order we defined!
	}
}
