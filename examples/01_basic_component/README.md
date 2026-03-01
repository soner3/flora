# Flora Example: Basic Component Auto-Wiring

This example demonstrates the absolute core of the Flora DI framework: **Zero-configuration auto-wiring.**

## The Concept
Normally, you would have to manually instantiate your dependencies and pass them into your constructors. With Flora, you simply embed the `flora.Component` tag into your structs.

Flora will automatically:
1. Discover the struct.
2. Find its constructor (`New...`).
3. Notice that `App` requires a `Greeter` interface.
4. Automatically bind `EnglishGreeter` to the `Greeter` interface because it is the only implementation.
5. Generate the boilerplate code to wire them together.

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
Starting the application...
Hello! Flora has successfully injected this component.

```