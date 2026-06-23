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
	"reflect"
	"testing"

	"github.com/soner3/flora/internal/engine"
)

func TestParsePackages(t *testing.T) {
	testcases := []struct {
		name   string
		path   string
		src    string
		expErr error
	}{
		{
			name:   "TestParsePackagesSuccessful",
			path:   "../../testdata/scanner/happy",
			expErr: nil,
		},
		{
			name: "TestParsePackagesMissingProvider",
			src: `package errnoconstructor

import "github.com/soner3/flora"

type Bad struct{ flora.Component }
`,
			expErr: ErrProviderFuncNotFound,
		},
		{
			name: "TestParsePackagesNotAFunc",
			src: `package errnotfunc

import "github.com/soner3/flora"

type BadComponent struct {
	flora.Component
}

var NewBadComponent = "not a function"
`,
			expErr: ErrInvalidProviderFunc,
		},
		{
			name: "TestParsePackagesNoReturn",
			src: `package errnoreturn

import "github.com/soner3/flora"

type Bad struct{ flora.Component }

func NewBad() {}
`,
			expErr: ErrInvalidProviderFunc,
		},
		{
			name: "TestParsePackagesWrongType",
			src: `package errwrongtype

import "github.com/soner3/flora"

type Bad struct{ flora.Component }

func NewBad() string { return "" }
`,
			expErr: ErrInvalidProviderFunc,
		},
		{
			name: "TestParsePackagesTooManyReturns",
			src: `package errtoomanyreturns

import "github.com/soner3/flora"

type Bad struct{ flora.Component }

func NewBad() (*Bad, func(), error, string) { return nil, nil, nil, "" }
`,
			expErr: ErrInvalidProviderFunc,
		},
		{
			name: "TestParsePackagesTwoReturnsWrongSecond",
			src: `package errtworeturnswrongsecond

import "github.com/soner3/flora"

type Bad struct{ flora.Component }

func NewBad() (*Bad, string) { return nil, "" }
`,
			expErr: ErrInvalidProviderFunc,
		},
		{
			name: "TestParsePackagesThreeReturnsWrongSecond",
			src: `package errthreereturnswrongsecond

import "github.com/soner3/flora"

type Bad struct{ flora.Component }

func NewBad() (*Bad, string, error) { return nil, "", nil }
`,
			expErr: ErrInvalidProviderFunc,
		},
		{
			name: "TestParsePackagesThreeReturnsWrongThird",
			src: `package errthreereturnswrongthird

import "github.com/soner3/flora"

type Bad struct{ flora.Component }

func NewBad() (*Bad, func(), string) { return nil, nil, "" }
`,
			expErr: ErrInvalidProviderFunc,
		},
		{
			name: "TestParsePackagesNoImplementation",
			src: `package errnoimpl

import "github.com/soner3/flora"

type Iface interface{ Do() }

type Consumer struct{ flora.Component }

func NewConsumer(i Iface) *Consumer { return nil }
`,
			expErr: ErrNoImplementation,
		},
		{
			name: "TestParsePackagesInterfaceCollisionNoPrimary",
			src: `package errcollisionnoprimary

import "github.com/soner3/flora"

type Greeter interface {
	Greet()
}
type GreeterA struct {
	flora.Component
}

func NewGreeterA() *GreeterA { return nil }
func (g *GreeterA) Greet()   {}

type GreeterB struct {
	flora.Component
}

func NewGreeterB() *GreeterB { return nil }
func (g *GreeterB) Greet()   {}

type Consumer struct {
	flora.Component
}

func NewConsumer(g Greeter) *Consumer { return nil }
`,
			expErr: ErrInterfaceCollision,
		},
		{
			name: "TestParsePackagesInterfaceCollisionMultiPrimary",
			src: `package errcollisionmultiprimary

import "github.com/soner3/flora"

type Greeter interface {
	Greet()
}
type GreeterA struct {
	flora.Component ` + "`" + `flora:"primary"` + "`" + `
}

func NewGreeterA() *GreeterA { return nil }
func (g *GreeterA) Greet()   {}

type GreeterB struct {
	flora.Component ` + "`" + `flora:"primary"` + "`" + `
}

func NewGreeterB() *GreeterB { return nil }
func (g *GreeterB) Greet()   {}

type Consumer struct {
	flora.Component
}

func NewConsumer(g Greeter) *Consumer { return nil }
`,
			expErr: ErrInterfaceCollision,
		},
		{
			name: "TestParsePackagesAnonSlice",
			src: `package erranonslice

import "github.com/soner3/flora"

type Bad struct{ flora.Component }

func NewBad(s []interface{ Do() }) *Bad { return nil }
`,
			expErr: ErrInvalidSlice,
		},
		{
			name: "TestParsePackagesFirstReturnErr",
			src: `package errfirstreturnerr

import "github.com/soner3/flora"

type Bad struct{ flora.Component }

func NewBad() error { return nil }
`,
			expErr: ErrInvalidProviderFunc,
		},
		{
			name: "TestParsePackagesSelfReferential",
			src: `package errselfref

import "github.com/soner3/flora"

type Bad struct{ flora.Component }

func NewBad(b *Bad) *Bad { return nil }
`,
			expErr: ErrInvalidProviderFunc,
		},
		{
			name: "TestParsePackagesAnonIfaceSingle",
			src: `package erranonifacesingle

import "github.com/soner3/flora"

type Impl struct{ flora.Component }

func NewImpl() *Impl { return nil }
func (i *Impl) Do()  {}

type Bad struct{ flora.Component }

func NewBad(req interface{ Do() }) *Bad { return nil }
`,
			expErr: ErrInvalidInterface,
		},
		{
			name: "TestParsePackagesAnonIfacePrimary",
			src: `package erranonifaceprimary

import "github.com/soner3/flora"

type Impl1 struct {
	flora.Component ` + "`" + `flora:"primary"` + "`" + `
}

func NewImpl1() *Impl1 { return nil }
func (i *Impl1) Do()   {}

type Impl2 struct{ flora.Component }

func NewImpl2() *Impl2 { return nil }
func (i *Impl2) Do()   {}

type Bad struct{ flora.Component }

func NewBad(req interface{ Do() }) *Bad { return nil }
`,
			expErr: ErrInvalidInterface,
		},
		{
			name: "TestParsePackagesInvalidScope",
			src: `package errinvalidscope

import "github.com/soner3/flora"

type FakeComponent struct {
	flora.Component ` + "`" + `flora:"scope=invalid"` + "`" + `
}

func NewFakeComponent() *FakeComponent {
	return &FakeComponent{}
}
`,
			expErr: ErrInvalidMetadata,
		},
		{
			name: "TestParsePackagesPrototypeWithParams",
			src: `package errprototypeparam

import "github.com/soner3/flora"

type Bad struct{ flora.Component }

func NewBad(f func(id int) *Bad) *Bad { return nil }
`,
			expErr: ErrInvalidProviderFunc,
		},
		{
			name: "TestParsePackagesPrototypeInvalidReturn",
			src: `package errprototypereturn

import "github.com/soner3/flora"

type Bad struct{ flora.Component }

func NewBad(f func() (*Bad, string)) *Bad { return nil }
`,
			expErr: ErrInvalidProviderFunc,
		},
		{
			name: "TestParsePackagesInvalidOrder",
			src: `package errinvalidorder

import "github.com/soner3/flora"

type BadOrderComponent struct {
	flora.Component ` + "`" + `flora:"order=one"` + "`" + `
}

func NewBadOrderComponent() *BadOrderComponent {
	return &BadOrderComponent{}
}
`,
			expErr: ErrInvalidMetadata,
		},
		{
			name: "TestParsePackagesConfigInvalidScope",
			src: `package errconfigscope

import "github.com/soner3/flora"

type BadConfig struct {
	flora.Configuration
}

// flora:scope=super_singleton
func (c *BadConfig) ProvideFloat() float32 {
	return 1.0
}
`,
			expErr: ErrInvalidMetadata,
		},
		{
			name: "TestParsePackagesUnexportedPrefix",
			src: `package errunexportedprefix

import "github.com/soner3/flora"

type MyComponent struct {
	flora.Component ` + "`" + `flora:"constructor=badConstructor"` + "`" + `
}

func badConstructor() *MyComponent { return nil }
`,
			expErr: ErrInvalidMetadata,
		},
		{
			name: "TestParsePackagesUnexportedPositional",
			src: `package errunexportedpos

import "github.com/soner3/flora"

type MyComponent struct {
	flora.Component ` + "`" + `flora:"badConstructor"` + "`" + `
}

func badConstructor() *MyComponent { return nil }
`,
			expErr: ErrInvalidMetadata,
		},
		{
			name: "TestParsePackagesErrConfigProvider",
			src: `package errconfigprovider

import "github.com/soner3/flora"

type BadConfig struct {
	flora.Configuration
}

func (c *BadConfig) ProvideInvalid() error {
	return nil
}
`,
			expErr: ErrInvalidProviderFunc,
		},
		{
			name:   "TestParsePackagesHappyConfig",
			path:   "../../testdata/scanner/happy",
			expErr: nil,
		},
		{
			name: "TestParsePackagesErrQualifierNotFound",
			src: `package errqualifiernotfound

import "github.com/soner3/flora"

type MyService struct {
	flora.Component ` + "`" + `flora:"inject(db=doesNotExist)"` + "`" + `
}

func NewMyService(db string) *MyService { return nil }
`,
			expErr: ErrInvalidMetadata,
		},
		{
			name: "TestParsePackagesErrInjectUnknownParam",
			src: `package errinjectunknownparam

import "github.com/soner3/flora"

type MyService struct {
	flora.Component ` + "`" + `flora:"inject(wrongParamName=masterDB)"` + "`" + `
}

func NewMyService(db string) *MyService { return nil }
`,
			expErr: ErrInvalidMetadata,
		},
		{
			name: "TestParsePackagesErrSliceAny",
			src: `package errsliceany

import "github.com/soner3/flora"

type Bad struct{ flora.Component }

func NewBad(s []any) *Bad { return nil }
`,
			expErr: ErrInvalidSlice,
		},
		{
			name:   "TestParsePackagesHappySlicePtr",
			path:   "../../testdata/scanner/happy_slice_ptr",
			expErr: nil,
		},
		{
			name: "TestParsePackagesErrQualifierCollisionType",
			src: `package errqualifiercollisiontype

import "github.com/soner3/flora"

type A struct {
	flora.Component ` + "`" + `flora:"name=myQual"` + "`" + `
}

func NewA() *A { return nil }

type B struct {
	flora.Component ` + "`" + `flora:"name=myQual"` + "`" + `
}

func NewB() *B { return nil }
`,
			expErr: ErrInvalidMetadata,
		},
		{
			name: "TestParsePackagesErrInjectTypeMismatch",
			src: `package errinjecttypemismatch

import "github.com/soner3/flora"

type DB struct {
	flora.Component ` + "`" + `flora:"name=myDB"` + "`" + `
}

func NewDB() *DB { return nil }

type Consumer struct {
	flora.Component ` + "`" + `flora:"inject(dep=myDB)"` + "`" + `
}

func NewConsumer(dep string) *Consumer { return nil }
`,
			expErr: ErrInvalidMetadata,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			path := tc.path
			if tc.src != "" {
				path = createTempModule(t, tc.src)
			}
			packages, err := ScanPackages(path)
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

func TestParseFloraTag(t *testing.T) {
	testcases := []struct {
		name       string
		rawTag     string
		expName    string
		expInject  map[string]string
		expPrimary bool
		expScope   string
		expErr     error
	}{
		{
			name:       "TestEmptyTag",
			rawTag:     ``,
			expName:    "",
			expInject:  map[string]string{},
			expPrimary: false,
			expScope:   ScopeSingleton,
			expErr:     nil,
		},
		{
			name:       "TestBasicAttributes",
			rawTag:     `flora:"primary, scope=prototype"`,
			expName:    "",
			expInject:  map[string]string{},
			expPrimary: true,
			expScope:   ScopePrototype,
			expErr:     nil,
		},
		{
			name:       "TestQualifierName",
			rawTag:     `flora:"name=myCustomService"`,
			expName:    "myCustomService",
			expInject:  map[string]string{},
			expPrimary: false,
			expScope:   ScopeSingleton,
			expErr:     nil,
		},
		{
			name:       "TestSingleInject",
			rawTag:     `flora:"inject(db=masterDB)"`,
			expName:    "",
			expInject:  map[string]string{"db": "masterDB"},
			expPrimary: false,
			expScope:   ScopeSingleton,
			expErr:     nil,
		},
		{
			name:       "TestMultipleInjectWithSpaces",
			rawTag:     `flora:"inject( db = masterDB , logger= fileLogger )"`,
			expName:    "",
			expInject:  map[string]string{"db": "masterDB", "logger": "fileLogger"},
			expPrimary: false,
			expScope:   ScopeSingleton,
			expErr:     nil,
		},
		{
			name:       "TestComplexCombination",
			rawTag:     `flora:"primary, name=myService, inject(db=masterDB, cache=redisCache), scope=prototype"`,
			expName:    "myService",
			expInject:  map[string]string{"db": "masterDB", "cache": "redisCache"},
			expPrimary: true,
			expScope:   ScopePrototype,
			expErr:     nil,
		},
		{
			name:       "TestInvalidInjectFormat",
			rawTag:     `flora:"inject(db:masterDB)"`,
			expName:    "",
			expInject:  map[string]string{},
			expPrimary: false,
			expScope:   ScopeSingleton,
			expErr:     ErrInvalidMetadata,
		},
		{
			name:       "Empty Parts in Main Tag (Triggers continue)",
			rawTag:     `flora:", primary, , name=myService, "`,
			expName:    "myService",
			expInject:  map[string]string{},
			expPrimary: true,
			expScope:   ScopeSingleton,
			expErr:     nil,
		},
		{
			name:       "Empty Parts in Inject Tag (Triggers continue)",
			rawTag:     `flora:"inject(db=master, , cache=redis, )"`,
			expName:    "",
			expInject:  map[string]string{"db": "master", "cache": "redis"},
			expPrimary: false,
			expScope:   ScopeSingleton,
			expErr:     nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			metadata := &engine.ComponentMetadata{
				StructName: "TestStruct",
			}

			err := parseFloraTag(tc.rawTag, metadata)

			if tc.expErr != nil {
				if !errors.Is(err, tc.expErr) {
					t.Errorf("expected error %v, got %v", tc.expErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if metadata.QualifierName != tc.expName {
				t.Errorf("expected QualifierName '%s', got '%s'", tc.expName, metadata.QualifierName)
			}
			if metadata.IsPrimary != tc.expPrimary {
				t.Errorf("expected IsPrimary %v, got %v", tc.expPrimary, metadata.IsPrimary)
			}
			if metadata.Scope != tc.expScope {
				t.Errorf("expected Scope '%s', got '%s'", tc.expScope, metadata.Scope)
			}
			if !reflect.DeepEqual(metadata.InjectParams, tc.expInject) {
				t.Errorf("expected InjectParams %v, got %v", tc.expInject, metadata.InjectParams)
			}
		})
	}
}
