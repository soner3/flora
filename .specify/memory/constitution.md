<!-- Sync Impact Report
- Version: 0.0.0 -> 1.0.0
- Modified Principles: Created new principles
- Added Sections: Core Principles, Additional Constraints, Development Workflow, Governance
- Removed Sections: None
- Templates requiring updates:
  - .specify/templates/plan-template.md (⚠ pending)
  - .specify/templates/spec-template.md (⚠ pending)
  - .specify/templates/tasks-template.md (⚠ pending)
- Follow-up TODOs: None
-->
# Flora Constitution

## Core Principles

### I. Idiomatic Go (Effective Go)
All code MUST adhere to the conventions and idioms described in the official "Effective Go" documentation. This includes proper naming conventions (MixedCaps), clear error handling, appropriate use of pointers vs. values, and leveraging Go's built-in concurrency primitives correctly. Code MUST be formatted using `gofmt` or `goimports` prior to any commit to ensure consistent style across the project.

### II. Test-Driven Development (TDD)
TDD is mandatory and NON-NEGOTIABLE. The Red-Green-Refactor cycle MUST be strictly enforced:
1. Write a failing test for the desired behavior.
2. Ensure the test fails for the right reason.
3. Write the minimum amount of code required to make the test pass.
4. Refactor the code for clarity and maintainability while keeping the tests green.
Tests should describe the behavior and contracts of the code rather than its internal implementation details.

### III. Go Development Best Practices
- **Interfaces**: Keep interfaces small (1-3 methods) and define them where they are consumed, not where they are implemented.
- **Composition over Inheritance**: Use embedding and composition to build complex types and behaviors.
- **Error Handling**: Handle errors explicitly. Do not ignore errors or use panic for normal error handling. Wrap errors with contextual information where appropriate.
- **Concurrency**: "Do not communicate by sharing memory; instead, share memory by communicating." Prefer channels over shared memory with mutexes when passing data between goroutines.

### IV. Simplicity and Readability
Keep the design simple. Avoid premature optimization and over-engineering (YAGNI). Write clear, concise, and readable code. Comments should explain the "why" (intent and context) rather than the "what" (which the code should express clearly). Exported identifiers MUST have documentation comments.

### V. Clean Architecture and Modularity
Separate concerns into distinct, focused packages. A package should have a clear, single responsibility. Minimize dependencies between packages to prevent tight coupling and cyclic imports.

## Additional Constraints

- **Tooling**: Standard Go tooling (`go test`, `go build`, `go vet`) MUST be used.
- **Linting**: Code MUST pass a standard linter (e.g., `golangci-lint`) without errors or warnings before being merged.
- **Dependencies**: Keep external dependencies to a minimum. Any new dependency must be justified and vetted for security, maintenance, and license compliance.

## Development Workflow

- All new features and bug fixes MUST be developed using TDD.
- Every Pull Request (PR) MUST include comprehensive tests demonstrating the behavior of the change.
- Code reviews MUST explicitly verify adherence to "Effective Go" guidelines, the Go best practices defined in this constitution, and the presence of adequate tests.

## Governance

- This Constitution supersedes all other development practices or informal agreements.
- All code, PRs, and reviews MUST verify compliance with these principles.
- Amendments to this constitution require documentation of the rationale, team approval, and, if necessary, a migration plan for existing code.

**Version**: 1.0.0 | **Ratified**: 2026-06-23 | **Last Amended**: 2026-06-23
