package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// createTempModule creates a temporary Go module with the provided source code,
// correctly linking it to the flora module so that packages.Load succeeds.
func createTempModule(t *testing.T, src string) string {
	t.Helper()
	dir := t.TempDir()

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	floraPath := filepath.Dir(filepath.Dir(cwd))

	modContent := fmt.Sprintf(`module errtest

go 1.26.1

require github.com/soner3/flora v0.0.0
replace github.com/soner3/flora => %s
`, floraPath)

	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(modContent), 0644); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte(src), 0644); err != nil {
		t.Fatal(err)
	}

	return dir
}
