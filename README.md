# 🌿 Flora

**Compile-time Dependency Injection for Go using struct tags.** *Spring-like convenience with zero runtime overhead.*

Flora is an organic, highly automated dependency injection framework for Go. It scans your code, discovers components via struct tags, and automatically generates a strongly-typed DI container using Google Wire under the hood. 

No more manual wiring. No more massive `Initialize()` functions. Just write your code, and let Flora grow the ecosystem.

## ✨ Key Features

* **Zero Runtime Overhead:** Uses Google Wire to generate pure, static Go code.
* **Auto-Discovery:** Just embed `flora.Component` in your struct, and Flora handles the rest.
* **Automatic Interface Binding:** If your struct implements an interface, Flora binds it automatically.
* **Multi-Binding (Group Injection):** Request a slice of an interface (e.g., `[]Plugin`), and Flora will inject *all* matching implementations automatically!
* **Conflict Resolution:** Mark default implementations with the `flora:"primary"` tag.

---

## 📖 Quick Start

Let's build a simple app with a Database service and a Plugin system.

### 1. Define your Components

Just embed `flora.Component` in your structs. Flora will look for a `New<StructName>` constructor by default.

```go
package database

import "[github.com/soner3/flora](https://github.com/soner3/flora)"

type Config struct {
    flora.Component `flora:"constructor=LoadConfig"`
    URL string
}

func LoadConfig() Config {
    return Config{URL: "postgres://localhost:5432/db"}
}

type PostgresRepo struct {
    // We mark this as primary in case there are multiple databases
    flora.Component `flora:"primary"` 
    cfg Config
}

func NewPostgresRepo(cfg Config) *PostgresRepo {
    return &PostgresRepo{cfg: cfg}
}

func (r *PostgresRepo) Fetch() string { return "Data from DB" }

```

### 2. The Magic of Multi-Binding (Plugins)

Want to build an extensible system? Just define an interface and ask for a slice of it. Flora will gather all implementations automatically.

```go
package plugin

import "[github.com/soner3/flora](https://github.com/soner3/flora)"

type Plugin interface {
    Execute()
}

// Plugin A
type LoggerPlugin struct{ flora.Component }
func NewLoggerPlugin() *LoggerPlugin { return &LoggerPlugin{} }
func (p *LoggerPlugin) Execute() { println("Logger running") }

// Plugin B
type MetricsPlugin struct{ flora.Component }
func NewMetricsPlugin() *MetricsPlugin { return &MetricsPlugin{} }
func (p *MetricsPlugin) Execute() { println("Metrics running") }

// The Consumer
type PluginManager struct {
    flora.Component
    plugins []Plugin // Flora will inject LoggerPlugin AND MetricsPlugin here!
}

func NewPluginManager(plugins []Plugin) *PluginManager {
    return &PluginManager{plugins: plugins}
}

```

### 3. Generate the Container

Run the Flora CLI in your project root:

```bash
flora gen -d ./ -o ./

```

*Flora will scan your code and generate a `flora_container.go` file containing the `FloraContainer`.*

### 4. Run your App

Now, just initialize the generated container in your `main.go` and access your fully wired ecosystem!

```go
package main

import "fmt"

func main() {
    container, err := InitializeContainer()
    if err != nil {
        panic(err)
    }

    // Single Binding
    data := container.PostgresRepo.Fetch()
    fmt.Println(data)

    // Multi-Binding
    for _, p := range container.SliceOfPlugin {
        p.Execute()
    }
}

```

---

## 🛠 Struct Tags Reference

You can customize how Flora wires your components using the `flora` tag:

| Tag | Example | Description |
| --- | --- | --- |
| `constructor` | ``flora:"constructor=BuildApp"`` | Tells Flora to use a specific function instead of `New<StructName>`. |
| `primary` | ``flora:"primary"`` | Resolves collisions. If multiple structs implement the same interface, the primary one is injected. |
| (Empty) | ``flora:""`` | Explicitly marks a component without applying special rules. |

---

## 📜 License

Flora is released under the [Apache 2.0 License](https://www.google.com/search?q=LICENSE).

