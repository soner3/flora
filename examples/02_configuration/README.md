# Flora Example: Third-Party Configuration

This example shows how to integrate external types (like `*sql.DB`, standard library types, or third-party loggers) into your Flora DI container.

## The Concept
You cannot embed the `flora.Component` tag into structs that you don't own (since they are imported from other packages). Instead, Flora provides the `flora.Configuration` marker. 

By creating a Configuration struct, you can write provider methods and use **Magic Comments** (like `// flora:primary`) to tell Flora exactly how to instantiate these external types.

**Key Takeaways:**
* Use `flora.Configuration` for external dependencies.
* Return a `cleanup func()` from your provider, and Flora will automatically execute it during graceful shutdown (`defer cleanup()`).

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
-> [Config] Initializing ExternalLogger...
-------------------------------------------------
[SYSTEM] INFO: Worker has successfully started its job!
-------------------------------------------------
Application is shutting down. Watch the cleanup happen:
-> [Config] Cleaning up ExternalLogger resources...

```