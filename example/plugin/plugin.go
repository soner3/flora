/*
Copyright © 2026 Soner Astan astansoner@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package plugin

import (
	"fmt"

	"github.com/soner3/flora"
)

type Plugin interface {
	Name() string
	Execute()
}

type LoggerPlugin struct {
	flora.Component `flora:"order=1"`
}

func NewLoggerPlugin() *LoggerPlugin {
	return &LoggerPlugin{}
}

func (p *LoggerPlugin) Name() string {
	return "Logger"
}

func (p *LoggerPlugin) Execute() {
	fmt.Println("Logging system activated.")
}

type MetricsPlugin struct {
	flora.Component `flora:"order=2"`
}

func NewMetricsPlugin() *MetricsPlugin {
	return &MetricsPlugin{}
}

func (p *MetricsPlugin) Name() string {
	return "Metrics"
}

func (p *MetricsPlugin) Execute() {
	fmt.Println("Metrics system activated.")
}

type PluginManager struct {
	flora.Component
	plugins []Plugin
}

func NewPluginManager(plugins []Plugin) *PluginManager {
	return &PluginManager{
		plugins: plugins,
	}
}

func (m *PluginManager) RunAll() {
	fmt.Printf("\n--- Plugin Manager ---\n")
	fmt.Printf("Found %d plugins!\n", len(m.plugins))
	for _, p := range m.plugins {
		fmt.Printf(" -> Starting '%s' plugin...\n", p.Name())
		p.Execute()
	}
	fmt.Printf("----------------------\n")
}
