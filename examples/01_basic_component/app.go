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
