package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"
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
			name:        "TestInvalidOutputDirectory",
			args:        []string{"gen", "-i", "./root.go"},
			expectedOut: "invalid path provided for flag 'input':",
			expErr:      true,
		},
		{
			name:        "TestSuccess",
			args:        []string{"gen", "-i", "./testdata/generate/happy", "-o", "./testdata/generate/happy"},
			expectedOut: "Successfully generated flora container!",
			cleanup: func() {
				err := os.Remove("./testdata/generate/happy/flora_container.go")
				if err != nil {
					t.Fatalf("expected no error deleting generated file, got %v", err)
				}
			},
			expErr: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
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
