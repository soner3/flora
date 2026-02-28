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

	"github.com/soner3/flora/internal/engine"
)

func TestParsePackages(t *testing.T) {
	testcases := []struct {
		name         string
		testdataPath string
		expErr       error
	}{
		{
			name:         "TestParsePackagesSuccessful",
			testdataPath: "testdata/happy",
			expErr:       nil,
		},
		{
			name:         "TestParsePackagesMissingProvider",
			testdataPath: "testdata/err_no_constructor",
			expErr:       ErrProviderFuncNotFound,
		},
		{
			name:         "TestParsePackagesNotAFunc",
			testdataPath: "testdata/err_not_func",
			expErr:       ErrInvalidProviderFunc,
		},
		{
			name:         "TestParsePackagesNoReturn",
			testdataPath: "testdata/err_no_return",
			expErr:       ErrInvalidProviderFunc,
		},
		{
			name:         "TestParsePackagesWrongType",
			testdataPath: "testdata/err_wrong_type",
			expErr:       ErrInvalidProviderFunc,
		},
		{
			name:         "TestParsePackagesTooManyReturns",
			testdataPath: "testdata/err_too_many_returns",
			expErr:       ErrInvalidProviderFunc,
		},
		{
			name:         "TestParsePackagesTwoReturnsWrongSecond",
			testdataPath: "testdata/err_two_returns_wrong_second",
			expErr:       ErrInvalidProviderFunc,
		},
		{
			name:         "TestParsePackagesThreeReturnsWrongSecond",
			testdataPath: "testdata/err_three_returns_wrong_second",
			expErr:       ErrInvalidProviderFunc,
		},
		{
			name:         "TestParsePackagesThreeReturnsWrongThird",
			testdataPath: "testdata/err_three_returns_wrong_third",
			expErr:       ErrInvalidProviderFunc,
		},
		{
			name:         "TestParsePackagesNoImplementation",
			testdataPath: "testdata/err_no_impl",
			expErr:       ErrNoImplementation,
		},
		{
			name:         "TestParsePackagesInterfaceCollisionNoPrimary",
			testdataPath: "testdata/err_collision_no_primary",
			expErr:       ErrInterfaceCollision,
		},
		{
			name:         "TestParsePackagesInterfaceCollisionMultiPrimary",
			testdataPath: "testdata/err_collision_multi_primary",
			expErr:       ErrInterfaceCollision,
		},
		{
			name:         "TestParsePackagesAnonSlice",
			testdataPath: "testdata/err_anon_slice",
			expErr:       ErrInvalidSlice,
		},
		{
			name:         "TestParsePackagesFirstReturnErr",
			testdataPath: "testdata/err_first_return_err",
			expErr:       ErrInvalidProviderFunc,
		},
		{
			name:         "TestParsePackagesSelfReferential",
			testdataPath: "testdata/err_self_ref",
			expErr:       ErrInvalidProviderFunc,
		},
		{
			name:         "TestParsePackagesAnonIfaceSingle",
			testdataPath: "testdata/err_anon_iface_single",
			expErr:       ErrInvalidInterface,
		},
		{
			name:         "TestParsePackagesAnonIfacePrimary",
			testdataPath: "testdata/err_anon_iface_primary",
			expErr:       ErrInvalidInterface,
		},
		{
			name:         "TestParsePackagesInvalidScope",
			testdataPath: "testdata/err_invalid_scope",
			expErr:       ErrInvalidMetadata,
		},
		{
			name:         "TestParsePackagesHappyQualifier",
			testdataPath: "testdata/happy_qualifier",
			expErr:       nil,
		},
		{
			name:         "TestParsePackagesPrototypeWithParams",
			testdataPath: "testdata/err_prototype_param",
			expErr:       ErrInvalidProviderFunc,
		},
		{
			name:         "TestParsePackagesPrototypeInvalidReturn",
			testdataPath: "testdata/err_prototype_return",
			expErr:       ErrInvalidProviderFunc,
		},
		{
			name:         "TestParsePackagesInvalidOrder",
			testdataPath: "testdata/err_invalid_order",
			expErr:       ErrInvalidMetadata,
		},
		{
			name:         "TestParsePackagesConfigInvalidScope",
			testdataPath: "testdata/err_config_scope",
			expErr:       ErrInvalidMetadata,
		},
		{
			name:         "TestParsePackagesUnexportedPrefix",
			testdataPath: "testdata/err_unexported_prefix",
			expErr:       ErrInvalidMetadata,
		},
		{
			name:         "TestParsePackagesUnexportedPositional",
			testdataPath: "testdata/err_unexported_pos",
			expErr:       ErrInvalidMetadata,
		},
		{
			name:         "TestParsePackagesErrConfigProvider",
			testdataPath: "testdata/err_config_provider",
			expErr:       ErrInvalidProviderFunc,
		},
		{
			name:         "TestParsePackagesHappyConfig",
			testdataPath: "testdata/happy",
			expErr:       nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			packages, err := ScanPackages(tc.testdataPath)
			if err != nil {
				t.Fatalf("ScanPackages failed: %v", err)
			}

			genCtx, err := ParsePackages(packages)

			if tc.expErr != nil {
				if !errors.Is(err, tc.expErr) {
					t.Errorf("expected error %v, got %v", tc.expErr, err)
				}
			} else {
				if err != nil {
					t.Errorf("ParsePackages failed: %v", err)
				}

				if len(genCtx.Components) < 1 {
					t.Errorf("ParsePackages returned no components")
				}

			}
		})
	}
}

func TestIsExported(t *testing.T) {
	testcases := []struct {
		name      string
		component *engine.ComponentMetadata
		expErr    error
	}{
		{name: "TestSuccessfull", component: &engine.ComponentMetadata{ConstructorName: "Success"}, expErr: nil},
		{name: "TestNotExportet", component: &engine.ComponentMetadata{ConstructorName: "fail"}, expErr: ErrInvalidMetadata},
		{name: "TestEmpty", component: &engine.ComponentMetadata{ConstructorName: ""}, expErr: ErrInvalidMetadata},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := isExported(tc.component)
			if tc.expErr == nil {
				if err != nil {
					t.Errorf("expected %v, got %v instead", tc.expErr, err)

				}
			} else {
				if !errors.Is(err, ErrInvalidMetadata) {
					t.Errorf("expected %v, got %v instead", tc.expErr, err)

				}
			}
		})

	}

}
