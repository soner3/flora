<div align="center">
  <img src="https://github.com/user-attachments/assets/fc403eb3-7b4c-45e6-b57a-10502865f98e" alt="Flora DI Framework Banner" width="65%" />

  <br />
  <br />

  <h1>🌿 Flora</h1>
  <p><b>Compile-time, reflection-free Dependency Injection for Go.</b><br/>
  <i>Spring-like developer experience, but with zero runtime overhead and absolute type safety.</i></p>

  [![Go Reference](https://pkg.go.dev/badge/github.com/soner3/flora.svg)](https://pkg.go.dev/github.com/soner3/flora)
  [![Go Report Card](https://goreportcard.com/badge/github.com/soner3/flora)](https://goreportcard.com/report/github.com/soner3/flora)
  [![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
</div>

<br />

## 🤔 The Problem

Dependency Injection (DI) in Go traditionally forces you to choose between two painful extremes:

1. **Manual Wiring (Boilerplate Hell):** You manually instantiate every struct and pass dependencies down the chain. In projects with 50+ services, your `main.go` becomes a massive, unmaintainable 1000-line initialization block. Tools like Google Wire help, but still require you to manually write and maintain massive `ProviderSets`.
2. **Reflection (Runtime Magic):** You use dynamic frameworks (like Uber Dig or Spring in Java) that resolve dependencies at runtime. This causes slower startup times, circumvents Go's strict type system, and worst of all: missing dependencies cause your application to panic *at runtime* instead of failing during compilation.

## 💡 What is Flora? (The Solution)

Flora is an organic, highly automated dependency injection framework that acts as the missing "Convention over Configuration" layer for Go. 

It parses your source code's Abstract Syntax Tree (AST), discovers components natively via struct tags or magic comments, and automatically generates a strongly-typed, readable DI container using **[Google Wire](https://github.com/google/wire)** under the hood.

**How it differs:** Flora gives you the incredible Developer Experience (DX) of Spring Boot annotations, but evaluates everything strictly at **compile-time**. No reflection. No runtime proxies. If your code compiles, your dependency graph is 100% safe.

<img src="https://github.com/user-attachments/assets/3d64cd84-891d-4238-98e5-3c0fb7a2448e" alt="Flora Icon" width="150" align="right" style="margin-left: 20px; margin-bottom: 20px;" />

## ✨ What Flora Offers

* **🚀 Zero Runtime Overhead:** The generated code is exactly as fast and memory-efficient as manually written Go code.
* **🔍 Auto-Discovery:** Embed `flora.Component` in your structs. Flora finds them, reads their constructors, and wires them automatically.
* **📦 Third-Party Integration:** Use `flora.Configuration` to effortlessly integrate external libraries (like Database drivers or Loggers) into your DI container.
* **🔌 Automatic Interface Binding:** If a component implements an interface, Flora binds it automatically. If multiple implement it, you can define a `primary` component.
* **📚 Multi-Binding (Slices):** Ask for a slice of an interface (`[]Plugin`), and Flora automatically collects all implementations across your codebase and injects them.
* **🏭 Prototypes (Factories):** Need a fresh instance instead of a Singleton? Flora generates clean closure factories (`func() (*Service, error)`) out of the box.
* **🧹 Native Go Idioms:** Full, automatic support for constructors returning initialization `error`s and `cleanup func()` routines.

---

## 📦 Installation

Flora consists of two parts: the CLI tool (to generate the container) and the core library (for the markers).

```bash
# 1. Install the CLI tool globally
go install [github.com/soner3/flora/cmd/flora@latest](https://github.com/soner3/flora/cmd/flora@latest)

# 2. Add the library to your project
go get [github.com/soner3/flora@latest](https://github.com/soner3/flora@latest)

```

---

## 🛠️ How it Works (Detailed Examples)

### 1. The Basics: Components & Auto-Wiring

Instead of maintaining huge initialization scripts, you define dependencies declaratively right where they belong. Just embed `flora.Component` and provide a standard Go constructor.

```go
package domain

import "[github.com/soner3/flora](https://github.com/soner3/flora)"

// 1. Define your interface
type UserRepository interface {
    GetUserName() string
}

// 2. Mark your implementation
type PostgresRepo struct {
    flora.Component `flora:"primary"` // "primary" resolves collisions if multiple repos exist
}

func NewPostgresRepo() *PostgresRepo {
    return &PostgresRepo{}
}
func (r *PostgresRepo) GetUserName() string { return "Alice" }

// 3. Consume the interface
type UserService struct {
    flora.Component `flora:"constructor=BuildUserService"`
    repo UserRepository
}

// Flora scans this signature, notices you need a UserRepository, 
// and automatically injects the PostgresRepo!
func BuildUserService(repo UserRepository) *UserService {
    return &UserService{repo: repo}
}

```

### 2. Third-Party Integrations (`@Configuration`)

You cannot (and should not) embed `flora.Component` into external structs like `*sql.DB` or `*redis.Client`. To prevent "wrapper struct pollution", Flora offers Configurations.

Because Go does not allow struct tags on functions, Flora uses **Magic Comments** (an idiomatic Go pattern used by `//go:embed` or `//go:generate`) to configure these providers.

```go
package config

import (
    "database/sql"
    "[github.com/soner3/flora](https://github.com/soner3/flora)"
)

type DatabaseConfig struct {
    flora.Configuration // Marks this struct as a config provider
}

// flora:primary
func (c *DatabaseConfig) ProvidePostgres(logger *domain.Logger) (*sql.DB, func(), error) {
    db, err := sql.Open("postgres", "...")
    if err != nil {
        return nil, nil, err 
    }
    
    // Flora handles the cleanup function automatically during graceful shutdown!
    cleanup := func() { db.Close() }
    
    return db, cleanup, nil 
}

```

### 3. Multi-Binding (The Plugin Pattern)

Building extensible systems usually requires tedious array wiring. With Flora, you simply define an interface, implement it multiple times, and request a slice `[]YourInterface`. Flora handles the aggregation.

```go
package plugin

import "[github.com/soner3/flora](https://github.com/soner3/flora)"

type Plugin interface { Execute() }

type LoggerPlugin struct{ flora.Component }
func NewLoggerPlugin() *LoggerPlugin { return &LoggerPlugin{} }
func (p *LoggerPlugin) Execute() {}

type MetricsPlugin struct{ flora.Component }
func NewMetricsPlugin() *MetricsPlugin { return &MetricsPlugin{} }
func (p *MetricsPlugin) Execute() {}

// --- The Consumer ---

type PluginManager struct {
    flora.Component
    plugins []Plugin 
}

// Flora automatically discovers LoggerPlugin and MetricsPlugin,
// bundles them into a slice, and injects them here!
func NewPluginManager(plugins []Plugin) *PluginManager {
    return &PluginManager{plugins: plugins}
}

```

### 4. Prototypes (Dynamic Instantiation)

By default, Flora treats every component as a **Singleton** (one instance per container). If you need a fresh instance for every HTTP request, use the `prototype` scope.

Flora will then provide a Factory closure, automatically resolving all deep dependencies inside the closure!

```go
package report

import "[github.com/soner3/flora](https://github.com/soner3/flora)"

type PdfGenerator struct {
    flora.Component `flora:"scope=prototype"` // Change scope to prototype
    db *sql.DB
}

func NewPdfGenerator(db *sql.DB) (*PdfGenerator, func(), error) {
    return &PdfGenerator{db: db}, func() { println("Cleaned up temp files") }, nil
}

// --- The Consumer ---
type ReportService struct {
    flora.Component
    // Ask for a factory function instead of the struct directly!
    pdfFactory func() (*PdfGenerator, func(), error) 
}

func NewReportService(factory func() (*PdfGenerator, func(), error)) *ReportService {
    return &ReportService{pdfFactory: factory}
}

func (s *ReportService) HandleRequest() {
    // Calling the factory creates a fresh PdfGenerator, but Flora 
    // already wired the *sql.DB into it behind the scenes!
    pdf, cleanup, err := s.pdfFactory()
    defer cleanup() 
}

```

---

## 🚀 Generating the Container

Once your components are marked, run the Flora CLI in your project root:

```bash
# Scans the current directory and places flora_container.go in ./cmd/server
flora generate --input ./ --output ./cmd/server

```

Flora acts as the brain. It resolves the AST, validates the graph, and orchestrates Google Wire to generate a flawless, human-readable `flora_container.go`.

Now, simply boot your app:

```go
package main

import "yourproject/cmd/server"

func main() {
    // 100% statically typed, no reflection, full performance.
    container, cleanup, err := server.InitializeContainer()
    if err != nil {
        panic(err)
    }
    defer cleanup() // Safely closes all DBs, connections, and files

    container.UserService.Start()
}

```

---

<div align="center">
<img src="https://github.com/user-attachments/assets/78c5ff4d-7441-4ba0-a476-50754f96de55" alt="Flora Ecosystem" width="100%" />
</div>

## 🌍 The Flora Ecosystem

Whether you are building a small CLI tool or a massive microservice architecture—Flora scales with you. By moving dependency definitions directly to the components, your architecture remains clean, decoupled, and easily testable.

## ⚙️ Configuration Reference

### Struct Tags (`flora.Component`)

| Tag | Example | Description |
| --- | --- | --- |
| `constructor` | `flora:"constructor=BuildApp"` | Overrides the default `New<StructName>` lookup. |
| `primary` | `flora:"primary"` | Resolves interface collisions. The primary struct wins. |
| `scope` | `flora:"scope=prototype"` | Sets the lifecycle. Default is `singleton`. |
| `order` | `flora:"order=1"` | Defines sorting order when injected via Slice (`[]Interface`). |
| (Empty) | `flora:""` | Explicitly marks a component with default rules. |

### Magic Comments (`flora.Configuration`)

Must be placed directly above the configuration method.

| Comment | Description |
| --- | --- |
| `// flora:primary` | Marks the returned type as the primary implementation to resolve collisions. |
| `// flora:scope=prototype` | Changes the lifecycle to a factory function (fresh instance per call). |
| `// flora:order=1` | Defines the sorting order when the type is injected via Slice (`[]Interface`). |
| `// flora:primary,scope=prototype` | You can combine multiple instructions separated by commas. |

---

## 📜 License

Flora is released under the **Apache 2.0 License**.

Copyright © 2026 Soner Astan.

## 🙏 Acknowledgments

Flora builds its foundation upon [Google Wire](https://github.com/google/wire). While Flora provides the auto-discovery, AST parsing, and the highly automated developer experience, the actual generation of the static dependency graph is powered by Wire.
Lies es dir in Ruhe durch. Trifft diese Struktur genau den "Aha!"-Moment, den Entwickler beim Lesen deines Repositories haben sollen?

```
