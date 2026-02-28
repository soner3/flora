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
package wiregen

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/soner3/flora/internal/engine"
	"github.com/soner3/flora/internal/scanner"
)

func TestGenerate(t *testing.T) {

	loadHappyComponents := func(t *testing.T) *engine.GeneratorContext {
		packages, err := scanner.ScanPackages("testdata/happy")
		if err != nil {
			t.Fatalf("ScanPackages failed: %v", err)
		}
		genCtx, err := scanner.ParsePackages(packages)
		if err != nil {
			t.Fatalf("ParsePkgs failed: %v", err)
		}
		return genCtx
	}

	testcases := []struct {
		name     string
		setupDir func(t *testing.T) string
		genCtx   *engine.GeneratorContext
		expErr   error
	}{
		{
			name: "TestGenerateSuccessfully",
			setupDir: func(t *testing.T) string {
				tmpDir, err := os.MkdirTemp(".", "flora_test_out_*")
				if err != nil {
					t.Fatal(err)
				}
				return tmpDir
			},
			genCtx: nil,
			expErr: nil,
		},
		{
			name: "TestNoComponentsProvided",
			setupDir: func(t *testing.T) string {
				return t.TempDir()
			},
			genCtx: &engine.GeneratorContext{},
			expErr: nil,
		},
		{
			name: "TestComponentInMainLeak",
			setupDir: func(t *testing.T) string {
				return t.TempDir()
			},
			genCtx: &engine.GeneratorContext{
				Components: []*engine.ComponentMetadata{
					{
						PackageName:     "main",
						PackagePath:     "github.com/test/main",
						StructName:      "App",
						ConstructorName: "NewApp",
					},
				},
			},
			expErr: ErrMainComponentLeak,
		},
		{
			name: "TestInterfaceInMainLeak",
			setupDir: func(t *testing.T) string {
				return t.TempDir()
			},
			genCtx: &engine.GeneratorContext{
				Components: []*engine.ComponentMetadata{
					{
						PackageName:     "otherpkg",
						PackagePath:     "github.com/test/otherpkg",
						StructName:      "Service",
						ConstructorName: "NewService",
						Implements: []engine.InterfaceMetadata{
							{
								PackageName:   "main",
								PackagePath:   "github.com/test/main",
								InterfaceName: "MyInterface",
							},
						},
					},
				},
			},
			expErr: ErrMainInterfaceLeak,
		},
		{
			name: "TestPrototypeInterfaceInMainLeak",
			setupDir: func(t *testing.T) string {
				return t.TempDir()
			},
			genCtx: &engine.GeneratorContext{
				Components: []*engine.ComponentMetadata{
					{
						PackageName:     "otherpkg",
						PackagePath:     "github.com/test/otherpkg",
						StructName:      "ProtoService",
						ConstructorName: "NewProtoService",
						Scope:           "prototype",
						Implements: []engine.InterfaceMetadata{
							{
								PackageName:   "main",
								PackagePath:   "github.com/test/main",
								InterfaceName: "MyInterface",
							},
						},
					},
				},
			},
			expErr: ErrMainInterfaceLeak,
		},
		{
			name: "TestPrototypeLocalParamsTrimming",
			setupDir: func(t *testing.T) string {
				tmpDir, err := os.MkdirTemp(".", "flora_test_out_*")
				if err != nil {
					t.Fatal(err)
				}
				dummyFile := filepath.Join(tmpDir, "dummy.go")
				os.WriteFile(dummyFile, []byte("package main\n"), 0644)

				return tmpDir
			},
			genCtx: &engine.GeneratorContext{
				Components: []*engine.ComponentMetadata{
					{
						PackageName:     "main",
						PackagePath:     "github.com/test/main",
						StructName:      "Service",
						ConstructorName: "NewService",
						Scope:           "prototype",
						Params: []engine.ParamMetadata{
							{Name: "p0", Type: "*main.MyDep", Imports: []string{"context"}},
							{Name: "p1", Type: "[]main.MyDep"},
							{Name: "p2", Type: "main.MyDep"},
						},
					},
				},
			},

			expErr: ErrWireExecution,
		},
		{
			name: "TestInvalidOutputDir",
			setupDir: func(t *testing.T) string {
				return "invalid\x00path"
			},
			genCtx: &engine.GeneratorContext{
				Components: []*engine.ComponentMetadata{
					{PackageName: "pkg", StructName: "A", ConstructorName: "NewA"},
				},
			},
			expErr: ErrCreateOutputDir,
		},
		{
			name: "TestWireExecutionFailed",
			setupDir: func(t *testing.T) string {
				tmpDir, err := os.MkdirTemp(".", "flora_test_out_*")
				if err != nil {
					t.Fatal(err)
				}
				return tmpDir
			},
			genCtx: &engine.GeneratorContext{
				Components: []*engine.ComponentMetadata{
					{
						PackageName:     "happy",
						PackagePath:     "github.com/soner3/flora/internal/engine/wiregen/testdata/happy",
						StructName:      "GhostComponent",
						ConstructorName: "NewGhostComponent",
						IsPointer:       true,
					},
				},
			},
			expErr: ErrWireExecution,
		},
		{
			name: "TestResolveOutputDirFailed",
			setupDir: func(t *testing.T) string {
				originalWD, err := os.Getwd()
				if err != nil {
					t.Fatal(err)
				}
				tempDir, err := os.MkdirTemp("", "flora_del_*")
				if err != nil {
					t.Fatal(err)
				}

				if err := os.Chdir(tempDir); err != nil {
					t.Fatal(err)
				}
				t.Cleanup(func() {
					os.Chdir(originalWD)
				})

				os.RemoveAll(tempDir)

				return "."
			},
			genCtx: &engine.GeneratorContext{
				Components: []*engine.ComponentMetadata{
					{PackageName: "pkg", StructName: "A", ConstructorName: "NewA"},
				},
			},
			expErr: ErrResolveOutputDir,
		},
		{
			name: "TestReadExistingPackageName",
			setupDir: func(t *testing.T) string {
				tmpDir, err := os.MkdirTemp(".", "flora_test_out_*")
				if err != nil {
					t.Fatal(err)
				}
				dummyFile := filepath.Join(tmpDir, "dummy.go")
				os.WriteFile(dummyFile, []byte("package custompkg\n"), 0644)
				return tmpDir
			},
			genCtx: &engine.GeneratorContext{
				Components: []*engine.ComponentMetadata{
					{PackageName: "custompkg", StructName: "A", ConstructorName: "NewA"},
				},
			},
			expErr: ErrWireExecution,
		},
		{
			name: "TestRootDirectoryFallbackToMain",
			setupDir: func(t *testing.T) string {
				return "/"
			},
			genCtx: &engine.GeneratorContext{
				Components: []*engine.ComponentMetadata{
					{PackageName: "main", StructName: "A", ConstructorName: "NewA"},
				},
			},
			expErr: ErrWriteTempFile,
		},
		{
			name: "TestParseTemplateFailed",
			setupDir: func(t *testing.T) string {
				tmpDir, _ := os.MkdirTemp(".", "flora_test_out_*")

				originalTmpl := wireTemplate
				wireTemplate = "{{ unclosed"

				t.Cleanup(func() {
					wireTemplate = originalTmpl
				})
				return tmpDir
			},
			genCtx: &engine.GeneratorContext{
				Components: []*engine.ComponentMetadata{
					{PackageName: "pkg", StructName: "A", ConstructorName: "NewA"},
				},
			},
			expErr: ErrParseTemplate,
		},
		{
			name: "TestExecuteTemplateFailed",
			setupDir: func(t *testing.T) string {
				tmpDir, _ := os.MkdirTemp(".", "flora_test_out_*")

				originalTmpl := wireTemplate
				wireTemplate = `{{template "ghost_template"}}`

				t.Cleanup(func() {
					wireTemplate = originalTmpl
				})
				return tmpDir
			},
			genCtx: &engine.GeneratorContext{
				Components: []*engine.ComponentMetadata{
					{PackageName: "pkg", StructName: "A", ConstructorName: "NewA"},
				},
			},
			expErr: ErrExecuteTemplate,
		},
		{
			name: "TestWriteTempFileFailed",
			setupDir: func(t *testing.T) string {
				tmpDir, _ := os.MkdirTemp(".", "flora_test_out_*")

				os.Mkdir(filepath.Join(tmpDir, "flora_injector.go"), 0755)

				return tmpDir
			},
			genCtx: &engine.GeneratorContext{
				Components: []*engine.ComponentMetadata{
					{PackageName: "pkg", StructName: "A", ConstructorName: "NewA"},
				},
			},
			expErr: ErrWriteTempFile,
		},
		{
			name: "TestEnsureWireDependencyFailed",
			setupDir: func(t *testing.T) string {
				return t.TempDir()
			},
			genCtx: &engine.GeneratorContext{
				Components: []*engine.ComponentMetadata{
					{PackageName: "pkg", StructName: "A", ConstructorName: "NewA"},
				},
			},
			expErr: ErrEnsureWireDependency,
		},
		{
			name: "TestRenameGeneratedFileFailed",
			setupDir: func(t *testing.T) string {
				tmpDir, err := os.MkdirTemp(".", "flora_test_out_*")
				if err != nil {
					t.Fatal(err)
				}

				blockerPath := filepath.Join(tmpDir, "flora_container.go")
				if err := os.Mkdir(blockerPath, 0755); err != nil {
					t.Fatal(err)
				}

				return tmpDir
			},
			genCtx: nil,
			expErr: ErrRenameGeneratedFile,
		},
		{
			name: "TestSelfImportContinue",
			setupDir: func(t *testing.T) string {
				tmpDir, err := os.MkdirTemp(".", "flora_test_out_*")
				if err != nil {
					t.Fatal(err)
				}
				dummyFile := filepath.Join(tmpDir, "dummy.go")
				os.WriteFile(dummyFile, []byte("package main\n"), 0644)
				return tmpDir
			},
			genCtx: &engine.GeneratorContext{
				Components: []*engine.ComponentMetadata{
					{
						PackageName:     "main",
						PackagePath:     "github.com/test/main",
						StructName:      "App",
						ConstructorName: "NewApp",
					},
					{
						PackageName:     "other",
						PackagePath:     "github.com/test/other",
						StructName:      "ExternalService",
						ConstructorName: "NewExternalService",
						Params: []engine.ParamMetadata{
							{
								Name:    "p0",
								Type:    "*main.App",
								Imports: []string{"github.com/test/main"},
							},
						},
					},
				},
			},
			expErr: ErrWireExecution,
		},
		{
			name: "TestConfigInMainLeak",
			setupDir: func(t *testing.T) string {
				return t.TempDir()
			},
			genCtx: &engine.GeneratorContext{
				Components: []*engine.ComponentMetadata{
					{
						StructName:        "DB",
						PackageName:       "sql",
						ConfigStructName:  "AppConfig",
						ConfigPackageName: "main",
					},
				},
			},
			expErr: ErrMainComponentLeak,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {

			outDir := tc.setupDir(t)

			if strings.HasPrefix(filepath.Base(outDir), "flora_test_out_") {
				defer os.RemoveAll(outDir)
			}

			genCtx := tc.genCtx
			if genCtx == nil {
				genCtx = loadHappyComponents(t)
			}

			err := NewWireGenerator().Generate(outDir, genCtx)

			if tc.expErr != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tc.expErr)
				}
				if !errors.Is(err, tc.expErr) {
					t.Errorf("expected error %v, got %v", tc.expErr, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}
