# Implementation Plan: refactor-tests

**Branch**: `001-refactor-tests` | **Date**: 2026-06-23 | **Spec**: [spec.md](file:///home/soner/Dev/flora/specs/001-refactor-tests/spec.md)

**Input**: Feature specification from `/specs/001-refactor-tests/spec.md`

## Summary

Refactor all existing test files across the Flora project to eliminate messy `testdata` usages, strictly adhere to idiomatic Golang testing guidelines (standard `testing` package), and achieve 100% test coverage across all packages (excluding `examples`). Invalid AST inputs will be replaced with inline string literals, and valid mock code will be consolidated.

## Technical Context

**Language/Version**: Go 1.26.1

**Primary Dependencies**: Standard Library `testing` package

**Storage**: N/A

**Testing**: Standard Go testing (table-driven tests, inline string fixtures for errors)

**Target Platform**: Cross-platform (Linux/macOS/Windows)

**Project Type**: CLI Tool / Go Library

**Performance Goals**: Fast test execution time

**Constraints**: 
- Strictly 100% test coverage (excluding `examples`)
- No external assertion libraries (e.g., `testify`)
- Invalid AST inputs MUST use inline string literals, not `testdata`
- Valid AST fixtures MUST be consolidated

**Scale/Scope**: All test files in `./internal` and `./cmd`

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- [x] **Idiomatic Go**: Uses standard `testing` package and table-driven tests.
- [x] **TDD & Best Practices**: Improves test maintainability and removes bad practices (fragmented `testdata`).
- [x] **Clean Architecture**: Consolidates fixtures.

## Execution Roadmap & Implementation Details Mapping

This section defines the sequence of tasks required for the core implementation and refinement phases. It references the appropriate implementation artifacts to ensure context is clear at each step.

### Phase 1: Audit & Setup
1. **Audit Existing `testdata`**: Identify files containing invalid/intentionally incorrect AST inputs vs. valid complex Go files across all packages.
2. **Setup Consolidated `testdata` Root**: Create a single `testdata/` directory at the repository root.
   - *Reference*: See **Structure Decision** below.

### Phase 2: Core Implementation (Package-by-Package Refactoring)
*For all steps below, strictly adhere to the **Testing Framework** constraint (Standard Library `testing` ONLY) and the **Testdata Handling** constraint (Inline strings for invalid ASTs, root `testdata` for valid ASTs) defined in [`research.md`](file:///home/soner/Dev/flora/specs/001-refactor-tests/research.md).*

1. **Refactor `internal/scanner` Tests** (`scanner_test.go`, `parser_test.go`): Convert invalid AST parsing inputs to inline string literals. Move valid AST files to the root `testdata/`.
2. **Refactor `internal/engine` Tests** (`engine_test.go`, `wiregen/generator_test.go`): Remove external mock dependencies, use inline strings, and ensure table-driven patterns.
3. **Refactor `internal/app` Tests** (`generate_test.go`): Standardize to idiomatic Go.
4. **Refactor `cmd` Tests** (`root_test.go`, `generate_test.go`): Remove `cmd/testdata` and adapt tests.
5. **Refactor `internal/errs` Tests** (`error_test.go`): Verify coverage and idiomatic structure.

### Phase 3: Refinement & Validation
1. **Clean Up Fragmented Directories**: Delete the old, scattered `testdata` directories.
2. **Coverage Verification**: Verify 100% test coverage (excluding `examples`).
   - *Reference*: Execute the validation script defined in [`quickstart.md`](file:///home/soner/Dev/flora/specs/001-refactor-tests/quickstart.md).
3. **Idiomatic Check**: Ensure zero usage of `testify` and that all code passes `go fmt` and `go vet`.

## Project Structure

### Documentation (this feature)

```text
specs/001-refactor-tests/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
└── tasks.md             # Phase 2 output (to be generated)
```

### Source Code (repository root)

```text
# Consolidated Test Structure
testdata/                # New consolidated root for valid Go fixtures
├── valid_app.go
├── valid_engine.go
└── ...

internal/
├── app/
│   └── generate_test.go # Updated to use inline strings for errors
├── engine/
│   ├── engine_test.go
│   └── wiregen/
│       └── generator_test.go
├── errs/
│   └── error_test.go
└── scanner/
    ├── scanner_test.go
    └── parser_test.go

cmd/
├── root_test.go
└── generate_test.go
```

**Structure Decision**: We will remove all scattered `testdata` directories (`internal/app/testdata`, `internal/engine/wiregen/testdata`, `internal/scanner/testdata`, `cmd/testdata`). Invalid test cases will use inline string literals within their respective `*_test.go` files. Any valid, complex Go files needed for AST parsing tests will be consolidated into a single `/testdata` directory at the repository root.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

*(No violations)*
