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

			components, err := ParseComponents(packages)

			if tc.expErr != nil {
				if !errors.Is(err, tc.expErr) {
					t.Fatalf("expected error %v, got %v", tc.expErr, err)
				}

			} else {
				if err != nil {
					t.Fatalf("ParseComponents failed: %v", err)
				}

				if len(components) < 1 {
					t.Fatalf("ParseComponents returned no components")
				}
			}

		})
	}
}
