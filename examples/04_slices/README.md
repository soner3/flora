# Flora Example: Slices & Multi-Binding (Plugins)

This example highlights one of Flora's most powerful features: **Automatic Slice Aggregation**. It is the perfect pattern for building extensible plugin systems, middlewares, or event listeners in Go.

## The Concept
Normally, when you build an extensible system, you have to manually maintain an array or slice of all your implementations and pass it to your manager or router. 

With Flora, you simply ask for a slice of an interface (`[]Plugin`) in your consumer's constructor. 

Flora will automatically:
1. Scan your entire codebase.
2. Find every component that implements the `Plugin` interface.
3. Bundle them together into a slice.
4. Sort them based on the `flora:"order=X"` struct tag.
5. Inject the fully sorted slice into your consumer.

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
-------------------------------------------------
-> PluginManager automatically discovered 2 plugins.
   [1] LoggerPlugin: Setting up logging...
   [2] MetricsPlugin: Starting metrics collection...
-------------------------------------------------

```