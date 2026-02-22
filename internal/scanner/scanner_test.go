/*
Copyright © 2026 Soner Astan astansoner@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
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
