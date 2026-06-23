# Tasks: refactor-tests

**Input**: Design documents from `/specs/001-refactor-tests/`

**Prerequisites**: plan.md, spec.md, research.md, quickstart.md

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2)
- Exact file paths are included in descriptions.

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [x] T001 Create consolidated `testdata/` directory at the repository root
- [x] T002 Identify all valid complex Go files currently in fragmented `testdata` directories (e.g., `internal/scanner/testdata`)

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T003 Move identified valid complex Go fixtures to the root `testdata/` directory

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Clean and Idiomatic Test Architecture (Priority: P1) 🎯 MVP

**Goal**: Refactor tests to follow idiomatic Golang guidelines, removing messy fragmented `testdata` and replacing them with inline string literals or root `testdata` fixtures.

**Independent Test**: All tests pass (`go test ./...`) without relying on the old `testdata` directories.

### Implementation for User Story 1

- [x] T004 [P] [US1] Refactor `internal/scanner/scanner_test.go` and `internal/scanner/parser_test.go` to use inline string literals for invalid ASTs and the new root `testdata` for valid ASTs
- [x] T005 [P] [US1] Refactor `internal/engine/engine_test.go` and `internal/engine/wiregen/generator_test.go` to use inline strings and table-driven patterns, avoiding external mocks
- [x] T006 [P] [US1] Refactor `internal/app/generate_test.go` to standardize to idiomatic Go and inline strings
- [x] T007 [P] [US1] Refactor `cmd/root_test.go` and `cmd/generate_test.go` to adapt to new test patterns
- [x] T008 [P] [US1] Refactor `internal/errs/error_test.go` for idiomatic consistency
- [x] T009 [US1] Delete fragmented `testdata` directories (`internal/app/testdata`, `internal/engine/wiregen/testdata`, `internal/scanner/testdata`, `cmd/testdata`)

**Checkpoint**: At this point, User Story 1 should be fully functional and all test architectures clean.

---

## Phase 4: User Story 2 - Achieve 100% Test Coverage (Priority: P2)

**Goal**: Ensure strictly 100% test coverage for all packages (excluding `examples`) to prevent regressions.

**Independent Test**: `go test` with coverage reports exactly 100% for all targeted packages.

### Implementation for User Story 2

- [x] T010 [US2] Run coverage report and identify gaps across all packages (excluding `examples`)
- [x] T011 [P] [US2] Write missing tests for `internal/scanner/...` to reach 100%
- [x] T012 [P] [US2] Write missing tests for `internal/engine/...` to reach 100%
- [x] T013 [P] [US2] Write missing tests for `internal/app/...` to reach 100%
- [x] T014 [P] [US2] Write missing tests for `internal/errs/...` to reach 100%
- [x] T015 [P] [US2] Write missing tests for `cmd/...` to reach 100%
- [x] T016 [US2] Verify 100% coverage globally (excluding `examples`) using `go tool cover`

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently and test coverage is perfect.

---

## Phase 5: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [x] T017 [P] Run `go fmt ./...` and `go vet ./...` to ensure idiomatic style across the codebase
- [x] T018 Verify no external assertions (e.g., `stretchr/testify`) are used in any test file

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
  - User Story 2 (P2) naturally follows User Story 1 (P1) as coverage requires the tests to be refactored first.
- **Polish (Final Phase)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - Core architecture clean up.
- **User Story 2 (P2)**: Best started after User Story 1 to avoid writing tests against the old architecture.

### Parallel Opportunities

- All refactoring tasks within Phase 3 (T004-T008) can run in parallel by different team members or agents, as they target isolated packages.
- All coverage enforcement tasks within Phase 4 (T011-T015) can be run in parallel.

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL)
3. Complete Phase 3: User Story 1 (Refactor all tests)
4. **STOP and VALIDATE**: Tests pass and fragmented directories are gone.

### Incremental Delivery

1. Complete MVP (US1).
2. Add User Story 2 (Coverage) → Fix gaps → Run `quickstart.md` validation script.
3. Polish and conclude.
