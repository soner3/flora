package cmd

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"
	"time"
)

func TestGenerateCmd(t *testing.T) {
	testcases := []struct {
		name        string
		args        []string
		expectedOut string
		cleanup     func()
		expErr      bool
	}{
		{
			name:        "TestInvalidInputDirectory",
			args:        []string{"gen", "-i", "./invalid"},
			expectedOut: "invalid directory provided for flag 'input':",
			expErr:      true,
		},
		{
			name:        "TestInvalidInputDirectoryNotADirectory",
			args:        []string{"gen", "-i", "./root.go"},
			expectedOut: "invalid path provided for flag 'input':",
			expErr:      true,
		},
		{
			name:        "TestInvalidWatchDirectory",
			args:        []string{"gen", "-i", "./testdata/generate/happy", "-o", "./testdata/generate/happy", "-w", "-d", "./invalid"},
			expectedOut: "invalid directory provided for flag 'watch-dir':",
			expErr:      true,
		},
		{
			name:        "TestInvalidWatchDirectoryNotADirectory",
			args:        []string{"gen", "-i", "./testdata/generate/happy", "-o", "./testdata/generate/happy", "-w", "-d", "./root.go"},
			expectedOut: "invalid path provided for flag 'watch-dir':",
			expErr:      true,
		},
		{
			name:        "TestRunGenerateSuccess",
			args:        []string{"gen", "-i", "./testdata/generate/happy", "-o", "./testdata/generate/happy"},
			expectedOut: "Successfully generated flora container!",
			cleanup: func() {
				err := os.Remove("./testdata/generate/happy/flora_container.go")
				if err != nil && !os.IsNotExist(err) {
					t.Errorf("expected no error deleting generated file, got %v", err)
				}
			},
			expErr: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			inputDir = "."
			outputDir = "flora"
			watch = false
			watchDir = "."

			generateCmd.Flags().Set("input", ".")
			generateCmd.Flags().Set("output", "flora")
			generateCmd.Flags().Set("watch", "false")
			generateCmd.Flags().Set("watch-dir", ".")

			b := new(bytes.Buffer)
			rootCmd.SetOut(b)
			rootCmd.SetErr(b)
			rootCmd.SetArgs(tc.args)

			err := rootCmd.Execute()

			if tc.expErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if !strings.Contains(err.Error(), tc.expectedOut) {
					t.Fatalf("expected error to contain %q, got %q", tc.expectedOut, err.Error())
				}
				return
			}

			defer tc.cleanup()

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if tc.expectedOut != "" && !strings.Contains(b.String(), tc.expectedOut) {
				t.Fatalf("expected output to contain %q, got %q", tc.expectedOut, b.String())
			}
		})
	}
}

func TestGenerateCmd_WatchModeShutdown(t *testing.T) {
	inputDir = "."
	outputDir = "flora"
	watch = false
	watchDir = "."

	generateCmd.Flags().Set("input", ".")
	generateCmd.Flags().Set("output", "flora")
	generateCmd.Flags().Set("watch", "false")
	generateCmd.Flags().Set("watch-dir", ".")

	inDirTmp := t.TempDir()
	outDirTmp := t.TempDir()
	watchDirTmp := t.TempDir()

	generateCmd.Flags().Set("input", inDirTmp)
	generateCmd.Flags().Set("output", outDirTmp)
	generateCmd.Flags().Set("watch", "true")
	generateCmd.Flags().Set("watch-dir", watchDirTmp)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	errCh := make(chan error, 1)

	go func() {
		errCh <- generateCmd.ExecuteContext(ctx)
	}()

	select {
	case err := <-errCh:
		if err != nil {
			t.Errorf("Expected graceful shutdown without error, but got: %v", err)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("Command did not shut down in time! Context cancellation was ignored.")
	}
}
