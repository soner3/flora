# Flora Example: Prototypes (Factories)

This example demonstrates how to change the lifecycle of a component from **Singleton** (default) to **Prototype** (fresh instance per call).

## The Concept
By default, Flora creates exactly *one* instance of each component per container. But sometimes you need a fresh instance, like a unique `Session` for every HTTP request or a temporary file handler.

If you add the struct tag ``flora:"scope=prototype"`` to your component, Flora changes its behavior. Instead of injecting the finished struct (e.g., `*Session`) into your consumers, Flora generates and injects a **Factory Closure** (`func() *Session`). 

Calling this factory function creates a brand new instance on the fly and automatically resolves all of its deep dependencies!

## How to run

1. Generate the DI container (`flora_container.go`) in the current directory:
   ```bash
   flora gen -o .

```

2. Run the application:
```bash
go run .

```



## Expected Output

```text
-> [Init] Server initialized (Singleton)
-------------------------------------------------
-> [Factory] Created a brand new Session with ID: 4182
Server handling request for Alice using Session 4182

-> [Factory] Created a brand new Session with ID: 8917
Server handling request for Bob using Session 8917

```

*(Note: The exact Session IDs will be randomly generated on your machine)*