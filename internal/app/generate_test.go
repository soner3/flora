package app

import (
	"os"
	"testing"
)

func TestRunGenerate(t *testing.T) {
	testcases := []struct {
		name    string
		dir     string
		outDir  string
		wantErr bool
	}{
		{
			name:    "TestScanError",
			dir:     "./testdata/scan_err",
			outDir:  t.TempDir(),
			wantErr: true,
		},
		{
			name:    "TestParseError",
			dir:     "./testdata/parse_err",
			outDir:  t.TempDir(),
			wantErr: true,
		},
		{
			name:    "TestZeroComponents",
			dir:     "./testdata/empty",
			outDir:  t.TempDir(),
			wantErr: false,
		},
		{
			name:    "TestGenerateError",
			dir:     "./testdata/happy",
			outDir:  "invalid\x00path",
			wantErr: true,
		},
		{
			name:    "TestSuccess",
			dir:     "./testdata/happy",
			outDir:  "",
			wantErr: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			outDir := tc.outDir

			if outDir == "" {
				tmpDir, err := os.MkdirTemp(".", "weld_app_test_*")
				if err != nil {
					t.Fatal(err)
				}
				defer os.RemoveAll(tmpDir)
				outDir = tmpDir
			}

			err := RunGenerate(tc.dir, outDir)

			if tc.wantErr {
				if err == nil {
					t.Errorf("expected an error, but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("did not expect an error, but got: %v", err)
				}
			}
		})
	}
}
