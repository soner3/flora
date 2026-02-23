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

func TestParseComponents(t *testing.T) {
	testcases := []struct {
		name         string
		testdataPath string
		expErr       error
	}{
		{
			name:         "TestParseComponentsSuccessful",
			testdataPath: "testdata/happy",
			expErr:       nil,
		},
		{
			name:         "TestParseComponentsMissingConstructor",
			testdataPath: "testdata/err_no_constructor",
			expErr:       ErrConstructorNotFound,
		},
		{
			name:         "TestParseComponentsNotAFunc",
			testdataPath: "testdata/err_not_func",
			expErr:       ErrConstructorNotFunc,
		},
		{
			name:         "TestParseComponentsNoReturn",
			testdataPath: "testdata/err_no_return",
			expErr:       ErrInvalidConstructor,
		},
		{
			name:         "TestParseComponentsWrongType",
			testdataPath: "testdata/err_wrong_type",
			expErr:       ErrInvalidConstructor,
		},
		{
			name:         "TestParseComponentsNoImplementation",
			testdataPath: "testdata/err_no_impl",
			expErr:       ErrNoImplementation,
		},
		{
			name:         "TestParseComponentsInterfaceCollisionNoPrimary",
			testdataPath: "testdata/err_collision_no_primary",
			expErr:       ErrInterfaceCollision,
		},
		{
			name:         "TestParseComponentsInterfaceCollisionMultiPrimary",
			testdataPath: "testdata/err_collision_multi_primary",
			expErr:       ErrInterfaceCollision,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			packages, err := ScanPackages(tc.testdataPath)

			if err != nil {
				t.Fatalf("ScanPackages failed: %v", err)
			}

			genCtx, err := ParseComponents(packages)

			if tc.expErr != nil {
				if !errors.Is(err, tc.expErr) {
					t.Fatalf("expected error %v, got %v", tc.expErr, err)
				}

			} else {
				if err != nil {
					t.Fatalf("ParseComponents failed: %v", err)
				}

				if len(genCtx.Components) < 1 {
					t.Fatalf("ParseComponents returned no components")
				}

				if len(genCtx.SliceBindings) < 1 {
					t.Fatalf("ParseComponents returned no slice bindings")
				}
			}

		})
	}
}

func TestValidateConstructor(t *testing.T) {

	testcases := []struct {
		name         string
		testdataPath string
		expErr       error
	}{
		{
			name:         "TestValidateConstructorSuccessful",
			testdataPath: "testdata/happy",
			expErr:       nil,
		},
		{
			name:         "TestParseComponentsTooManyReturns",
			testdataPath: "testdata/err_too_many_returns",
			expErr:       ErrInvalidConstructor,
		},
		{
			name:         "TestParseComponentsTwoReturnsWrongSecond",
			testdataPath: "testdata/err_two_returns_wrong_second",
			expErr:       ErrInvalidConstructor,
		},
		{
			name:         "TestParseComponentsThreeReturnsWrongSecond",
			testdataPath: "testdata/err_three_returns_wrong_second",
			expErr:       ErrInvalidConstructor,
		},
		{
			name:         "TestParseComponentsThreeReturnsWrongThird",
			testdataPath: "testdata/err_three_returns_wrong_third",
			expErr:       ErrInvalidConstructor,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			packages, err := ScanPackages(tc.testdataPath)

			if err != nil {
				t.Fatalf("ScanPackages failed: %v", err)
			}

			genCtx, err := ParseComponents(packages)

			if tc.expErr != nil {
				if !errors.Is(err, tc.expErr) {
					t.Errorf("expected error %v, got %v", tc.expErr, err)
				}

			} else {
				if err != nil {
					t.Errorf("ParseComponents failed: %v", err)
				}

				if len(genCtx.Components) < 1 {
					t.Errorf("ParseComponents returned no components")
				}

				if len(genCtx.SliceBindings) < 1 {
					t.Errorf("ParseComponents returned no slice bindings")
				}
			}

		})
	}
}
