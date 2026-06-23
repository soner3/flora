# Feature Specification: refactor-tests

**Feature Branch**: `001-refactor-tests`

**Created**: 2026-06-23

**Status**: Draft

**Input**: User description: "Review and refactor all test files across the project to establish a clean, scalable testing architecture that strictly adheres to idiomatic Golang guidelines. Currently, adding new tests has become overly complicated due to past workarounds. Numerous directories containing intentionally incorrect or mock code were created inside various testdata folders. These were primarily used to test AST parsing, code generation, and the core application logic. This approach has led to a violation of clean testing principles and created a fragmented, hard-to-maintain test suite. Comprehensive Refactoring: Analyze and refactor all existing test files to resolve the structural issues caused by the current testdata usage. Best Practices: Implement sustainable, idiomatic Go testing solutions for AST parsing, code generation, and application testing that replace the current messy directory structures. 100% Test Coverage: Ensure that every package in the codebase achieves strictly 100% test coverage after the refactoring is complete. The examples directory is explicitly excluded from this task. Do not refactor its contents or include it in the coverage requirements."

## Clarifications

### Session 2026-06-23
- Q: To adhere to idiomatic Golang guidelines, should we strictly use the standard library `testing` package for assertions, or are specific external libraries allowed? → A: Option A - Strictly use the standard `testing` package.
- Q: To resolve the structural issues with `testdata`, what is the exact preferred approach for handling AST parsing test inputs? → A: Option B - Use inline string literals for invalid/error AST cases, but allow minimal consolidated `testdata` strictly for valid, complex Go source files.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Clean and Idiomatic Test Architecture (Priority: P1)

As a developer, I want all tests to follow idiomatic Golang guidelines without relying on messy and overly complicated `testdata` folders containing intentionally incorrect code, so that adding new tests is straightforward and the test suite is easy to maintain.

**Why this priority**: The current fragmented architecture acts as a significant bottleneck for developers, making it difficult to maintain existing functionality and add new features.

**Independent Test**: Can be fully tested by verifying that all tests run successfully and confirming that the problematic `testdata` directories have been removed and replaced with standard inline mock code or structured test utilities.

**Acceptance Scenarios**:

1. **Given** the current fragmented test suite, **When** the developer refactors the AST parsing and code generation tests, **Then** the tests pass using idiomatic Go testing patterns instead of external `testdata` mock files.
2. **Given** a new feature to be tested, **When** a developer adds a new test case, **Then** they can do so easily without needing to navigate or manipulate convoluted external test directories.

---

### User Story 2 - Achieve 100% Test Coverage (Priority: P2)

As a maintainer, I want strictly 100% test coverage for all packages (explicitly excluding the `examples` directory), so that I have complete confidence in the application's correctness after the extensive test refactoring.

**Why this priority**: Comprehensive refactoring introduces the risk of regressions. Ensuring 100% coverage provides a safety net and guarantees that no edge cases are missed during the restructuring.

**Independent Test**: Can be fully tested by running `go test -cover` and verifying that the output reports 100% for all packages outside of the `examples` folder.

**Acceptance Scenarios**:

1. **Given** the fully refactored test suite, **When** the coverage report is generated, **Then** every package (excluding `examples`) must report exactly 100% coverage.
2. **Given** the `examples` directory, **When** the coverage report is generated, **Then** this directory is not included in the coverage checks or enforcement.

### Edge Cases

- What happens when a core component genuinely requires an external file to test AST parsing? Valid, complex Go source files should be placed in a minimal, consolidated `testdata` directory. Invalid or error-inducing AST code MUST NOT use `testdata` and must instead be defined as inline string literals within the test files.
- How does the system handle coverage for intentionally unreachable error paths or panic conditions? (These must be tested or refactored out if they cannot be reached).

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST eliminate the reliance on intentionally incorrect or scattered mock code inside fragmented `testdata` folders. Invalid AST and error cases MUST be tested using inline string literals. Minimal, consolidated `testdata` is only permitted for valid, complex Go source files.
- **FR-002**: System MUST strictly use the standard library `testing` package (without external assertion dependencies like `testify`) to implement idiomatic Go testing solutions (e.g., table-driven tests) for AST parsing, code generation, and core application logic.
- **FR-003**: System MUST achieve strictly 100% test coverage for all packages.
- **FR-004**: System MUST explicitly exclude the `examples` directory from refactoring and test coverage requirements.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: `go test ./...` passes successfully with no errors or regressions.
- **SC-002**: Go test coverage report shows strictly 100% for all packages excluding `examples`.
- **SC-003**: The overall number of messy, fragmented files inside `testdata` directories is reduced to zero for invalid/error cases, and consolidated into a single structured root for valid fixtures.
- **SC-004**: Adding a new test requires significantly less context switching and file manipulation compared to the previous state.

## Assumptions

- We assume that the `examples` directory does not contain critical core functionality that must be covered.
- We assume that the existing functionality of the application is correct, and the refactoring aims only to change *how* it is tested, not *what* it does.
- We assume that standard Go tools (`go test`, `go tool cover`) will be used to enforce coverage and idiomatic practices.
