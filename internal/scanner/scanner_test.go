package scanner

import (
	"errors"
	"testing"
)

func TestScanPackages(t *testing.T) {
	testcases := []struct {
		name     string
		path     string
		expected int
		expErr   error
	}{
		{
			name:     "TestScanPackagesSuccessfulScan",
			path:     "testdata/happy",
			expected: 1,
			expErr:   nil,
		},
		{
			name:     "TestScanPackagesFailedScan",
			path:     "testdata/foo",
			expected: 0,
			expErr:   ErrLoadPackages,
		},
		{
			name:     "TestScanPackagesCompileError",
			path:     "testdata/sad",
			expected: 0,
			expErr:   ErrCompile,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			packages, err := ScanPackages(tc.path)

			if tc.expErr != nil {
				if err == nil {
					t.Errorf("expected an error but got nil")
				}

				if packages != nil {
					t.Errorf("expected nil packages but got %v", packages)
				}

				if !errors.Is(err, tc.expErr) {
					t.Errorf("expected error %v but got %v", tc.expErr, err)
				}

			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}

				if len(packages) != tc.expected {
					t.Errorf("expected %d packages but got %d", tc.expected, len(packages))
				}

			}

		})
	}
}
