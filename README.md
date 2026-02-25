# 🌿 Flora

**Compile-time Dependency Injection for Go. Spring-like convenience, zero runtime overhead, and absolutely no magic.**

Flora is an organic, highly automated dependency injection framework for Go. It bridges the gap between the phenomenal developer experience (DX) of enterprise frameworks like Spring Boot and the strict, explicit, performance-oriented philosophy of Go.

It scans your code, discovers components via struct tags or magic comments, and automatically generates a strongly-typed, readable DI container using **Google Wire** under the hood.

No more manual wiring. No more massive `Initialize()` functions. Just write your code, and let Flora grow the ecosystem.

---

## 🤔 Why Flora? (The Philosophy)

The Go community traditionally dislikes the "magic" of frameworks like Spring (Reflection, runtime proxies, implicit behavior) because it contradicts Go's philosophy of explicitness and causes runtime panics. However, manually maintaining DI graphs (even with tools like Google Wire) becomes a massive boilerplate burden in projects with 50+ microservices.

**Flora is the golden middle ground:**

1. **Developer Experience:** You define dependencies declaratively right where they belong (at the component level).
2. **Code Generation, Not Reflection:** Flora acts as a transpiler. It parses your AST (Abstract Syntax Tree) and generates a static, human-readable `flora_container.go` file.
3. **Fail Fast:** Missing dependencies or circular loops are caught at **compile time**. If your code compiles, your DI graph is 100% safe.

---

## ✨ Key Features

* 🚀 **Zero Runtime Overhead:** No `reflect` package used at runtime. Full execution speed.
* 🔍 **Auto-Discovery:** Embed `flora.Component` in your struct, and Flora handles the rest.
* 📦 **Third-Party Integration:** Use `flora.Configuration` to wrap external dependencies (like `*sql.DB` or Loggers) cleanly into the DI graph.
* 🔌 **Automatic Interface Binding:** If a component implements an interface, Flora binds it automatically when requested.
* 📚 **Multi-Binding (Slices):** Request a slice of an interface (e.g., `[]Plugin`), and Flora automatically discovers and injects *all* implementations!
* 🏭 **Prototypes (Factories):** Need a fresh instance per HTTP request? Request a factory function (`func() (*Service, func(), error)`) and Flora handles the closure cleanly.
* 🧹 **Native Cleanup & Error Handling:** Fully supports Go-idiomatic constructors returning `error` and `cleanup func()`.

---

## 📖 Comprehensive Guide & Examples

### 1. The Standard Component (`@Component`)

For your own domain services, simply embed `flora.Component` and use struct tags. Flora will automatically look for a `New<StructName>` constructor.

```go
package domain

import "github.com/soner3/flora"

// 1. The Interface
type UserRepository interface {
    GetUserName() string
}

// 2. The Implementation
type PostgresRepo struct {
    flora.Component `flora:"primary"` // "primary" resolves collisions if multiple repos exist
}

func NewPostgresRepo() *PostgresRepo {
    return &PostgresRepo{}
}
func (r *PostgresRepo) GetUserName() string { return "Alice" }

// 3. The Consumer
type UserService struct {
    flora.Component `flora:"constructor=BuildUserService"`
    repo UserRepository
}

// Flora will automatically inject PostgresRepo as the UserRepository here!
func BuildUserService(repo UserRepository) *UserService {
    return &UserService{repo: repo}
}

```

### 2. Third-Party Integrations (`@Configuration` & `@Bean`)

You cannot embed `flora.Component` into an external struct like `*sql.DB`. For these cases, use `flora.Configuration`. Because Go does not support tags on methods, Flora uses idiomatic **Magic Comments** here.

```go
package config

import (
    "database/sql"
    "github.com/soner3/flora"
)

type DatabaseConfig struct {
    flora.Configuration // Marks this struct as a configuration provider
}

// flora:primary
func (c *DatabaseConfig) ProvidePostgres() (*sql.DB, func(), error) {
    db, err := sql.Open("postgres", "...")
    if err != nil {
        return nil, nil, err // Handled cleanly at startup
    }
    
    cleanup := func() { db.Close() }
    
    // Flora registers *sql.DB in the container and handles the cleanup!
    return db, cleanup, nil 
}

```

### 3. The Magic of Multi-Binding (Plugins)

Want to build an extensible system? Just define an interface and ask for a slice of it. Flora will gather all implementations automatically without any manual array wiring.

```go
package plugin

import "github.com/soner3/flora"

type Plugin interface { Execute() }

type LoggerPlugin struct{ flora.Component }
func NewLoggerPlugin() *LoggerPlugin { return &LoggerPlugin{} }
func (p *LoggerPlugin) Execute() { /* ... */ }

type MetricsPlugin struct{ flora.Component }
func NewMetricsPlugin() *MetricsPlugin { return &MetricsPlugin{} }
func (p *MetricsPlugin) Execute() { /* ... */ }

type PluginManager struct {
    flora.Component
    plugins []Plugin // Flora automatically injects BOTH LoggerPlugin and MetricsPlugin!
}

func NewPluginManager(plugins []Plugin) *PluginManager {
    return &PluginManager{plugins: plugins}
}

```

### 4. Prototypes (Runtime Instantiation)

Sometimes you don't want a Singleton. You want a fresh instance for every request, but you still want the DI container to resolve its deep dependencies (like DB connections).

```go
package report

import "github.com/soner3/flora"

type PdfGenerator struct {
    flora.Component `flora:"scope=prototype"` // Note the scope!
    db *sql.DB
}

// The DB is injected by Flora, but the instance is created on demand
func NewPdfGenerator(db *sql.DB) (*PdfGenerator, func(), error) {
    return &PdfGenerator{db: db}, func() { println("Cleaned up temp files") }, nil
}

// --- In your HTTP Handler ---

type ReportService struct {
    flora.Component
    // We request a FACTORY function, not the struct!
    pdfFactory func() (*PdfGenerator, func(), error) 
}

func NewReportService(factory func() (*PdfGenerator, func(), error)) *ReportService {
    return &ReportService{pdfFactory: factory}
}

func (s *ReportService) HandleRequest() {
    // Generate a fresh instance with its dependencies already wired!
    pdf, cleanup, err := s.pdfFactory()
    if err != nil { /* handle */ }
    defer cleanup() 
    
    // ... use pdf ...
}

```

---

## 🚀 Generating and Running

Run the Flora CLI in your project root to scan your codebase and generate the Wire container:

```bash
flora gen -d ./ -o ./

```

*Flora generates a `flora_container.go` file. Wire is automatically executed under the hood.*

Initialize your app in your `main.go`:

```go
package main

func main() {
    // 100% statically typed, no reflection, full performance.
    container, cleanup, err := InitializeContainer()
    if err != nil {
        panic(err)
    }
    defer cleanup() // Cleans up all DB connections and resources gracefully

    container.UserService.DoSomething()
}

```

---

## 🛠 Configuration Reference

### Struct Tags (For `flora.Component`)

| Tag | Example | Description |
| --- | --- | --- |
| `constructor` | `flora:"constructor=BuildApp"` | Overrides the default `New<StructName>` lookup. |
| `primary` | `flora:"primary"` | Resolves interface collisions. The primary struct wins. |
| `scope` | `flora:"scope=prototype"` | Sets the lifecycle. Default is `singleton`. |
| (Empty) | `flora:""` | Explicitly marks a component with default rules. |

### Magic Comments (For `flora.Configuration` methods)

Comments must be placed directly above the method definition.
| Comment | Description |
| --- | --- |
| `// flora:primary` | Marks the provided type as the primary implementation. |
| `// flora:scope=prototype` | Changes the lifecycle of the provided object to a factory. |

---

## 📜 License

Flora is released under the [Apache 2.0 License](https://www.google.com/search?q=LICENSE).

## 🙏 Acknowledgments

Flora is built on top of the incredible [Google Wire](https://github.com/google/wire) project.

While Flora provides the auto-discovery, parsing, and the Spring-like developer experience, the actual heavy lifting of safely generating the static, compile-time dependency graph is powered by Wire.