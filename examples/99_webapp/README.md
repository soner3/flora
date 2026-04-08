# Flora Example: Full Web Application (Tutorial)

This example demonstrates how to build a fully-structured, production-ready Go web application using Flora. It combines multiple Flora concepts like components, configurations, interfaces, and multi-binding (slice injection) to create a Clean Architecture API.

## The Concept
This tutorial showcases how Flora orchestrates a complex dependency graph without any manual wiring:
* **Master & Replica Database (`@Configuration`)**: Connects to two separate SQLite instances (Master for writes, Replica for reads) and returns two identical types (`*sql.DB`) safely resolved via **Qualifiers**.
* **Repository & Service (`@Component`)**: Implements the data layer. The repository consumes both Master and Replica DBs using the `flora:"inject(...)"` tag. The service depends on the repository via an **Interface**, decoupling the layers.
* **Middlewares (Slice Injection)**: Uses Flora's `// flora:order=X` feature to automatically inject and sort a slice of pure Go middlewares (like Logger and Auth). 
* **Controllers (Slice Injection)**: Automatically discovers all HTTP handlers implementing the `Controller` interface and registers them to the router.

With Flora, adding a new endpoint or middleware is as simple as creating a new struct with `flora.Component` - no need to ever touch the routing configuration again!

## How to run

1.  Generate the DI container (`flora_container.go`) in the current directory. We use `-i .` to tell Flora to scan all directories in the root dir:
    ```bash
    flora gen -i . -o .
    ```

2.  Run the application:
    ```bash
    go run .
    ```

3.  Test the API in your browser or via `curl`:
    * `curl http://localhost:8080/hello`
    * `curl http://localhost:8080/books`

## Expected Output

**Server Startup:**
```text
-> [DB] Connecting to database (file::memory:?cache=shared)...
-> [Repo] BookRepository initialized
-> [Service] BookService initialized with Repository Interface
-> [Web] Registering 2 custom middlewares...
-> [Web] Registering 2 controllers...

```

**Upon requesting `/books`:**

```text
-> [Logger] Intercepted GET /books
-> [Auth] Checked DB: 1 books available. Auth passed!

```

**Upon stopping the server (Ctrl+C):**

```text
-> Gracefully shutting down the web server...
-> [DB] Closing database connection...

```
