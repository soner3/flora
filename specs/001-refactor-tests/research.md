# Research: Refactoring Tests

## Testing Framework
- **Decision**: Standard Library `testing` package exclusively. No `testify` or external assertion libraries.
- **Rationale**: Aligns with the Flora Constitution and the specification which mandates idiomatic Go code ("Effective Go"). This minimizes dependencies and enforces clean, simple testing patterns (like table-driven tests).
- **Alternatives considered**: `stretchr/testify` (rejected to maintain zero-dependency test suite).

## Testdata Handling
- **Decision**: Inline string literals for invalid/error AST parsing cases. Consolidate minimal `testdata` for valid, complex Go source fixtures into a single `testdata` directory at the project root.
- **Rationale**: Specified in the clarification phase. Reduces fragmentation and messy "intentionally incorrect" files on disk.
- **Alternatives considered**: Ban `testdata` completely (rejected to allow complex, valid Go files to be easily maintained as real files).
