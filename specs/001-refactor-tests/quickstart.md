# Quickstart Validation

## Prerequisites
- Go 1.26.1

## Run Tests
Run all tests to verify 100% test coverage (excluding `examples`):

```bash
go test -coverpkg=$(go list ./... | grep -v /examples | tr '\n' ',') -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

**Expected Outcome**: 
- `go test` exits with `0` (success).
- All packages (except `examples`) report `coverage: 100.0% of statements` in the coverage tool output.
- No fragmented `testdata` directories exist under `internal/` or `cmd/`. All valid fixtures are inside the root `testdata/` directory.
