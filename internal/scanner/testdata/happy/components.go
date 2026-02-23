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
package happy

import "github.com/soner3/flora"

type Greeter interface {
	Greet() string
}

type SimpleLogger struct {
	flora.Component
}

func NewSimpleLogger() SimpleLogger {
	return SimpleLogger{}
}

type GermanGreeter struct {
	flora.Component `flora:"constructor=BuildGermanGreeter,"`
}

func BuildGermanGreeter() *GermanGreeter {
	return &GermanGreeter{}
}

func (g *GermanGreeter) Greet() string {
	return "Hallo"
}

type App struct {
	flora.Component
}

func NewApp(g Greeter, l SimpleLogger) *App {
	return &App{}
}

type JustANormalStruct struct {
	SomeConfig string
	Value      int
}

type UntaggedComponent struct {
	flora.Component
}

func NewUntaggedComponent() *UntaggedComponent {
	return nil
}

type Plugin interface {
	Execute()
}

type AuthPlugin struct {
	flora.Component
}

func NewAuthPlugin() *AuthPlugin { return nil }
func (p *AuthPlugin) Execute()   {}

type MetricsPlugin struct {
	flora.Component
}

func NewMetricsPlugin() *MetricsPlugin { return nil }
func (p *MetricsPlugin) Execute()      {}

type PluginManager struct {
	flora.Component
}

func NewPluginManager(plugins []Plugin) *PluginManager { return nil }

type ComponentWithErr struct {
	flora.Component
}

func NewComponentWithErr() (*ComponentWithErr, error) { return nil, nil }

type ComponentWithCleanup struct {
	flora.Component
}

func NewComponentWithCleanup() (*ComponentWithCleanup, func()) { return nil, nil }

type ComponentWithBoth struct {
	flora.Component
}

func NewComponentWithBoth() (*ComponentWithBoth, func(), error) { return nil, nil, nil }
