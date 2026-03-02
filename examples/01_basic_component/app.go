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

// ---------------------------------------------------------
// 1. The Interface
// ---------------------------------------------------------
type Greeter interface {
	Greet()
}

// ---------------------------------------------------------
// 2. The Implementation
// ---------------------------------------------------------
type EnglishGreeter struct {
	flora.Component // Tells Flora: "Manage this struct in the DI container"
}

// Flora automatically finds this constructor based on the struct name
func NewEnglishGreeter() *EnglishGreeter {
	return &EnglishGreeter{}
}

func (g *EnglishGreeter) Greet() {
	fmt.Println("Hello! Flora has successfully injected this component.")
}

// ---------------------------------------------------------
// 3. The Consumer
// ---------------------------------------------------------
type App struct {
	flora.Component // This is also managed by Flora
	greeter         Greeter
}

// Flora scans this constructor and sees that 'App' needs a 'Greeter'.
// Since 'EnglishGreeter' is the only struct implementing it,
// Flora injects it automatically without any configuration!
func NewApp(g Greeter) *App {
	return &App{greeter: g}
}

func (a *App) Run() {
	a.greeter.Greet()
}
