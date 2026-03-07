<div align="center">
  <img src="https://github.com/user-attachments/assets/70725b89-beef-4ad2-a8c3-211411a4b0c6" alt="Flora DI Framework Banner" width="65%" />

  <br />
  <br />

  <h1>Flora</h1>
  <p><b>Compile-time Dependency Injection for Go.</b><br/>
  <i>Automated wiring, zero runtime overhead, and absolutely no magic.</i></p>

  <p>
    <a href="https://pkg.go.dev/github.com/soner3/flora"><img src="https://pkg.go.dev/badge/github.com/soner3/flora.svg" alt="Go Reference"></a>
    <a href="https://github.com/soner3/flora/releases"><img src="https://img.shields.io/github/v/release/soner3/flora" alt="GitHub Release"></a>
    <a href="https://github.com/soner3/flora"><img src="https://img.shields.io/github/go-mod/go-version/soner3/flora" alt="Go Version"></a>
    <a href="https://github.com/soner3/flora/actions"><img src="https://img.shields.io/github/actions/workflow/status/soner3/flora/release.yml?branch=main" alt="Build Status"></a>
    <a href="https://goreportcard.com/report/github.com/soner3/flora"><img src="https://goreportcard.com/badge/github.com/soner3/flora" alt="Go Report Card"></a>
    <img src="https://img.shields.io/badge/Coverage-100%25-brightgreen.svg" alt="Test Coverage">
    <a href="https://opensource.org/licenses/Apache-2.0"><img src="https://img.shields.io/badge/License-Apache_2.0-blue.svg" alt="License"></a>
  </p>
</div>

<br />

## Why Flora?

Dependency Injection in Go traditionally forces developers to choose between two painful extremes:

1. **Manual Wiring (Boilerplate):** You manually construct the dependency graph or maintain massive `ProviderSets` (e.g., in Google Wire). As your application scales to dozens of services, wiring becomes a tedious, unmaintainable chore.
2. **Reflection (Runtime Magic):** You use dynamic frameworks that resolve dependencies at runtime. This causes slower startup times, circumvents Go's strict compiler, and worst of all: missing dependencies cause your application to panic *at runtime* instead of failing during compilation.

Flora solves this by acting as a "Convention over Configuration" layer. It parses your source code's Abstract Syntax Tree (AST) to natively discover your components, and then automatically generates a strongly-typed DI container using [Google Wire](https://github.com/google/wire) under the hood.

You get the developer experience of a modern, automated framework with the safety and performance of purely static code.

* **Zero Runtime Overhead:** The generated container is exactly as fast and memory-efficient as manually written Go code.
* **No Reflection:** Everything is evaluated at compile-time. If your code compiles, your dependency graph is 100% safe.
* **Auto-Discovery:** Just embed a marker into your structs. Flora reads your constructors and wires the entire graph automatically.
* **Native Go Idioms:** Full, automatic support for constructors returning initialization `error`s and `cleanup func()` routines for graceful shutdowns.

---

## Installation

Flora consists of a CLI tool (to generate the container) and a core library (for the markers).

```bash
# 1. Install the CLI tool globally
go install github.com/soner3/flora/cmd/flora@latest

# 2. Add the library to your project
go get github.com/soner3/flora@latest

```

---

## Quick Start

Instead of maintaining huge initialization scripts, you define dependencies declaratively right where they belong.

```go
package main

import "github.com/soner3/flora"

// 1. Mark your struct as a component
type Greeter struct {
    flora.Component
}

// 2. Provide a standard constructor
func NewGreeter() *Greeter {
    return &Greeter{}
}

func (g *Greeter) Greet() string { return "Hello from Flora!" }

type App struct {
    flora.Component
    greeter *Greeter
}

// Flora automatically discovers the *Greeter dependency and injects it.
func NewApp(g *Greeter) *App {
    return &App{greeter: g}
}

```

Generate the container by running the CLI in your project root:

```bash
flora gen -i . -o .

```

---

## Core Concepts & Reference

Flora uses two primary markers to understand your architecture. Knowing when to use which is the key to building clean applications.

### 1. `flora.Component` (For your Domain Logic)

Use this for structs that you own and write yourself (e.g., Services, Handlers, Repositories). By embedding `flora.Component`, you tell Flora to manage this struct.

**Default Behaviors:**

* **Constructor:** Flora automatically looks for an exported function named `New<StructName>` (e.g., `NewUserService`).
* **Scope:** `singleton` (one instance per application container).
* **Interfaces:** If the struct implements any interfaces, Flora binds them automatically.

**Tag Reference:**
You can override the default behavior using the `flora` struct tag. Multiple tags can be combined using commas.

| Tag | Example | Description |
| :--- | :--- | :--- |
| *(Shorthand)* | ``flora:"BuildApp"`` | Any string without a key-value pair is automatically treated as the `constructor` name. |
| `constructor` | ``flora:"constructor=BuildApp"`` | Explicitly overrides the default `New<StructName>` constructor lookup. |
| `primary` | ``flora:"primary"`` | Resolves interface collisions. If multiple structs implement the same interface, the primary one is injected. |
| `scope` | ``flora:"scope=prototype"`` | Changes the lifecycle to a Factory function (a fresh instance is created per injection). |
| `order` | ``flora:"order=1"`` | Defines the sorting order when the component is injected via Slice (`[]Interface`). |
| *(Empty)* | ``flora:""`` | Explicitly marks a component with default rules (optional, embedding the struct is usually enough). |

### 2. `flora.Configuration` (For Third-Party & Adapters)

Use this when you cannot (or should not) modify the target struct. This is strictly meant for **Third-Party Integrations** (like `*sql.DB`, `*redis.Client`) or functional paradigms (like pure HTTP Middlewares).

Instead of embedding a marker into the target struct, you create a config struct and use **Magic Comments** (`// flora:...`) above its provider methods.

```go
package config

import (
    "database/sql"
    "github.com/soner3/flora"
)

type DatabaseConfig struct {
    flora.Configuration
}

// flora:primary
func (c *DatabaseConfig) ProvidePostgres() (*sql.DB, func(), error) {
    db, err := sql.Open("postgres", "...")
    if err != nil {
      return nil, nil, err
    }
    
    // Flora handles the cleanup function automatically during graceful shutdown!
    cleanup := func() { db.Close() }
    
    return db, cleanup, nil 
}

```

**Magic Comment Reference:**
Because Go does not allow struct tags on functions, Flora reads the comments immediately preceding the provider method.

| Comment | Description |
| --- | --- |
| `// flora:primary` | Marks the returned type as the primary implementation. |
| `// flora:scope=prototype` | Changes the lifecycle to a Factory function. |
| `// flora:order=1` | Defines the sorting order when injected via Slice (`[]Interface`). |

---

## Advanced Usage

### Clean Architecture & Type Aliasing

In Hexagonal or Clean Architecture, your core domain should not have direct dependencies on external frameworks ("Framework Pollution"). Flora fully supports Go Type Aliases (`=`) to keep your domain pure.

```go
package domain

import "github.com/soner3/flora"

// 1. Define a neutral alias in your domain layer
type DIComponent = flora.Component

// 2. Use the alias in your domain structs. No direct framework imports!
type UserService struct {
    DIComponent `flora:"primary"`
    repo UserRepository
}

```

### Multi-Binding (Slice Injection)

Building extensible systems usually requires tedious array wiring. With Flora, you simply define an interface, implement it multiple times, and request a slice `[]YourInterface`. Flora handles the aggregation and sorting based on the `order` tag.

```go
// flora:order=1
func (c *MiddlewareConfig) ProvideLogger() Middleware {
    return &LoggerMiddleware{}
}

// flora:order=2
func (c *MiddlewareConfig) ProvideAuth(repo *Repository) Middleware {
    return &AuthMiddleware{repo: repo}
}

// Flora automatically bundles Logger and Auth into a sorted slice!
func NewServer(middlewares []Middleware) *Server {
    return &Server{chain: middlewares}
}

```

---

## Learn by Example

The best way to master Flora is by looking at the provided examples. Check out the **[`/examples`](https://github.com/soner3/flora/tree/master/examples)** directory in this repository. It contains focused, runnable scenarios ranging from basic to advanced:

* **`01_basic_component`**: The absolute basics. Struct embedding, interfaces, and auto-wiring.
* **`02_configuration`**: How to use `@Configuration` to integrate external packages and manage graceful shutdowns (`cleanup func()`).
* **`03_prototypes`**: How to generate and inject Factory closures (`func() (*Struct, error)`) instead of Singletons.
* **`04_slices`**: How to utilize Multi-Binding to automatically collect and sort plugins/middlewares into a slice.
* **`99_webapp`**: **(Recommended)** A simple REST API implementing routing, middlewares, and a mock database.

---

## License & Acknowledgments

Flora is released under the **Apache 2.0 License**.

Flora builds its foundation upon [Google Wire](https://github.com/google/wire). While Flora provides the auto-discovery, AST parsing, and the highly automated developer experience, the actual generation of the static dependency graph is powered by Wire.
