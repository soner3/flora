# Flora Example: Slices & Multi-Binding

This example highlights one of Flora's most powerful features: **Automatic Slice Aggregation**. It is the perfect pattern for building extensible plugin systems, middlewares, event listeners, or route registries in Go.

## The Concept

Normally, when you build an extensible system, you have to manually maintain an array or slice of all your implementations and pass it to your manager or router.

With Flora, you simply ask for a slice of _any_ type in your consumer's constructor. This works for:

- **Interfaces** (e.g., `[]Plugin`)
- **Concrete Pointers** (e.g., `[]*Route`)
- **Structs & Type Aliases** (e.g., `[]MyAlias`)

Flora will automatically:

1. Scan your entire codebase.
2. Find every component that matches the requested type (or implements the interface).
3. Bundle them together into a slice.
4. Sort them based on the `flora:"order=X"` struct tag or magic comment.
5. Inject the fully sorted slice into your consumer.

## How to run

1. Generate the DI container (`flora_container.go`) in the current directory:
   ```bash
   flora gen -o .
   ```
