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
	"bytes"
	"cmp"
	"errors"
	"fmt"
	"go/build"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"text/template"

	"github.com/soner3/flora/internal/engine"
	"github.com/soner3/flora/internal/errs"
)

var (
	ErrResolveOutputDir     = errors.New("failed to resolve absolute output directory")
	ErrCreateOutputDir      = errors.New("failed to create output directory")
	ErrMainComponentLeak    = errors.New("component belongs to package 'main' (Go forbids importing main)")
	ErrMainInterfaceLeak    = errors.New("interface belongs to package 'main' (Go forbids importing main)")
	ErrParseTemplate        = errors.New("failed to parse wire template")
	ErrExecuteTemplate      = errors.New("failed to execute wire template")
	ErrWriteTempFile        = errors.New("failed to write temporary wire file")
	ErrEnsureWireDependency = errors.New("failed to ensure google/wire dependency")
	ErrWireExecution        = errors.New("flora engine failed to resolve dependency graph")
	ErrRenameGeneratedFile  = errors.New("failed to rename generated container file")
)

type WireGenerator struct{}

func NewWireGenerator() *WireGenerator {
	return &WireGenerator{}
}

func isBuiltInType(name string) bool {
	switch name {
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64", "uintptr",
		"float32", "float64", "complex64", "complex128",
		"bool", "string", "error", "any", "byte", "rune":
		return true
	}
	return false
}

var wireTemplate = `//go:build wireinject
// +build wireinject

package {{.PackageName}}

import (
    "github.com/google/wire"
    {{range .Imports}}
    "{{.}}"
    {{end}}
)

{{range .ConfigWrappers}}
{{if .IsPrototype}}
func {{.WrapperName}}({{range $index, $param := .Params}}{{if $index}}, {{end}}{{$param.Name}} {{$param.Type}}{{end}}) func() ({{.ReturnType}}{{if .HasCleanup}}, func(){{end}}{{if .HasError}}, error{{end}}) {
    return func() ({{.ReturnType}}{{if .HasCleanup}}, func(){{end}}{{if .HasError}}, error{{end}}) {
        cfg := {{.ConfigPackagePrefix}}{{.ConfigStructName}}{}
        return cfg.{{.ConfigMethodName}}({{range $index, $param := .Params}}{{if $index}}, {{end}}{{$param.Name}}{{end}})
    }
}
{{else}}
func {{.WrapperName}}({{range $index, $param := .Params}}{{if $index}}, {{end}}{{$param.Name}} {{$param.Type}}{{end}}) ({{.ReturnType}}{{if .HasCleanup}}, func(){{end}}{{if .HasError}}, error{{end}}) {
    cfg := {{.ConfigPackagePrefix}}{{.ConfigStructName}}{}
    return cfg.{{.ConfigMethodName}}({{range $index, $param := .Params}}{{if $index}}, {{end}}{{$param.Name}}{{end}})
}
{{end}}
{{end}}

{{range .Prototypes}}
{{if not .IsConfig}}
func {{.WrapperName}}({{range $index, $param := .Params}}{{if $index}}, {{end}}{{$param.Name}} {{$param.Type}}{{end}}) func() ({{.ReturnType}}{{if .HasCleanup}}, func(){{end}}{{if .HasError}}, error{{end}}) {
    return func() ({{.ReturnType}}{{if .HasCleanup}}, func(){{end}}{{if .HasError}}, error{{end}}) {
        return {{.ConstructorCall}}({{range $index, $param := .Params}}{{if $index}}, {{end}}{{$param.Name}}{{end}})
    }
}
{{end}}
{{end}}

{{range .SliceBindings}}
func ProvideSliceOf{{.InterfaceName}}({{range .Implementations}}{{.ParamName}} {{.TypePrefix}}{{.StructName}}, {{end}}) []{{.InterfacePrefix}}{{.InterfaceName}} {
    return []{{.InterfacePrefix}}{{.InterfaceName}}{
        {{range .Implementations}}{{.ParamName}},{{end}}
    }
}
{{end}}

type FloraContainer struct {
    {{range .Providers}}
    {{.StructName}} {{if .IsPointer}}*{{end}}{{.TypePrefix}}{{.StructName}}
    {{end}}
    
    {{range .Prototypes}}
    {{.FieldName}} func() ({{.ReturnType}}{{if .HasCleanup}}, func(){{end}}{{if .HasError}}, error{{end}})
    {{end}}

    {{range .SliceBindings}}
    SliceOf{{.InterfaceName}} []{{.InterfacePrefix}}{{.InterfaceName}}
    {{end}}
}

func InitializeContainer() (*FloraContainer, func(), error) {
    wire.Build(
        {{range .Providers}}
        {{.CallPrefix}}{{.ConstructorName}},
        {{end}}
        {{range .Prototypes}}
        {{.WrapperName}},
        {{end}}
        {{range .Bindings}}
        wire.Bind(new({{.InterfacePrefix}}{{.InterfaceName}}), new({{if .IsPointer}}*{{end}}{{.ComponentPrefix}}{{.StructName}})),
        {{end}}
        {{range .SliceBindings}}
        ProvideSliceOf{{.InterfaceName}},
        {{end}}
        wire.Struct(new(FloraContainer), "*"),
    )
    return nil, nil, nil
}
`

type providerData struct {
	StructName      string
	CallPrefix      string
	TypePrefix      string
	ConstructorName string
	IsPointer       bool
}

type paramData struct {
	Name string
	Type string
}

type prototypeData struct {
	WrapperName     string
	FieldName       string
	ConstructorCall string
	ReturnType      string
	Params          []paramData
	HasCleanup      bool
	HasError        bool
	IsConfig        bool
}

type configWrapperData struct {
	WrapperName         string
	ConfigPackagePrefix string
	ConfigStructName    string
	ConfigMethodName    string
	ReturnType          string
	Params              []paramData
	HasCleanup          bool
	HasError            bool
	IsPrototype         bool
}

type bindingData struct {
	InterfacePrefix string
	InterfaceName   string
	ComponentPrefix string
	StructName      string
	IsPointer       bool
}

type sliceImplData struct {
	ParamName  string
	TypePrefix string
	StructName string
}

type sliceBindingData struct {
	InterfacePrefix string
	InterfaceName   string
	Implementations []sliceImplData
}

type templateData struct {
	PackageName    string
	Imports        []string
	Providers      []providerData
	Prototypes     []prototypeData
	ConfigWrappers []configWrapperData
	Bindings       []bindingData
	SliceBindings  []sliceBindingData
}

func (g *WireGenerator) Generate(outDir string, genCtx *engine.GeneratorContext) error {
	log := slog.With("pkg", "wiregen")

	if len(genCtx.Components) == 0 && len(genCtx.SliceBindings) == 0 {
		log.Debug("No components provided, skipping generation")
		return nil
	}

	absOutDir, err := filepath.Abs(outDir)
	if err != nil {
		chainErr := fmt.Errorf("%w: %w", ErrResolveOutputDir, err)
		return errs.Wrap(chainErr, "provided path: %s", outDir)
	}

	if err := os.MkdirAll(absOutDir, os.ModePerm); err != nil {
		chainErr := fmt.Errorf("%w: %w", ErrCreateOutputDir, err)
		return errs.Wrap(chainErr, "absolute path: %s", absOutDir)
	}

	pkgName := filepath.Base(absOutDir)
	pkgName = strings.ReplaceAll(pkgName, "-", "_")

	if buildPkg, err := build.Default.ImportDir(absOutDir, 0); err == nil {
		pkgName = buildPkg.Name
	} else if pkgName == "." || pkgName == "/" {
		pkgName = "main"
	}

	var generatedPkgPath string
	for _, comp := range genCtx.Components {
		if comp.PackageName == pkgName {
			generatedPkgPath = comp.PackagePath
			break
		}
	}

	data := templateData{
		PackageName: pkgName,
	}

	var providers []providerData
	var prototypes []prototypeData
	var configWrappers []configWrapperData
	var bindings []bindingData
	importSet := make(map[string]bool)

	for _, comp := range genCtx.Components {
		isConfig := comp.ConfigStructName != ""
		isBuiltIn := isBuiltInType(comp.StructName)

		compPrefix := ""
		if comp.PackageName != pkgName && !isBuiltIn {
			if comp.PackageName == "main" {
				return errs.Wrap(ErrMainComponentLeak, "cannot generate container in package '%s' because component '%s' belongs to package 'main'. Change output dir (-o) to your main directory or move the component.", pkgName, comp.StructName)
			}
			compPrefix = comp.PackageName + "."
			importSet[comp.PackagePath] = true
		}

		configPkgPrefix := ""
		if isConfig {
			if comp.ConfigPackageName != pkgName {
				if comp.ConfigPackageName == "main" {
					return errs.Wrap(ErrMainComponentLeak, "cannot generate container because config '%s' belongs to package 'main'.", comp.ConfigStructName)
				}
				configPkgPrefix = comp.ConfigPackageName + "."
				importSet[comp.ConfigPackagePath] = true
			}
		}

		var pData []paramData
		for _, p := range comp.Params {
			for _, imp := range p.Imports {
				importSet[imp] = true
			}
			pType := p.Type
			pType = strings.ReplaceAll(pType, "*"+pkgName+".", "*")
			pType = strings.ReplaceAll(pType, "[]"+pkgName+".", "[]")
			if after, ok := strings.CutPrefix(pType, pkgName+"."); ok {
				pType = after
			}
			pData = append(pData, paramData{Name: p.Name, Type: pType})
		}

		retType := compPrefix + comp.StructName
		if comp.IsPointer {
			retType = "*" + retType
		}

		if comp.Scope == "prototype" {
			wrapperName := "ProvidePrototype" + comp.StructName
			if isConfig {
				wrapperName = "ProvidePrototype_" + comp.ConfigStructName + "_" + comp.ConfigMethodName
				configWrappers = append(configWrappers, configWrapperData{
					WrapperName:         wrapperName,
					ConfigPackagePrefix: configPkgPrefix,
					ConfigStructName:    comp.ConfigStructName,
					ConfigMethodName:    comp.ConfigMethodName,
					ReturnType:          retType,
					Params:              pData,
					HasCleanup:          comp.HasCleanup,
					HasError:            comp.HasError,
					IsPrototype:         true,
				})
			}

			prototypes = append(prototypes, prototypeData{
				WrapperName:     wrapperName,
				FieldName:       comp.StructName + "Factory",
				ConstructorCall: compPrefix + comp.ConstructorName,
				ReturnType:      retType,
				Params:          pData,
				HasCleanup:      comp.HasCleanup,
				HasError:        comp.HasError,
				IsConfig:        isConfig,
			})

			for _, iface := range comp.Implements {
				ifacePrefix := ""
				if iface.PackageName != pkgName {
					if iface.PackageName == "main" {
						return errs.Wrap(ErrMainInterfaceLeak, "cannot generate container in package '%s' because interface '%s' belongs to package 'main'. Change output dir (-o) to your main directory or move the interface.", pkgName, iface.InterfaceName)
					}
					ifacePrefix = iface.PackageName + "."
					importSet[iface.PackagePath] = true
				}

				prototypes = append(prototypes, prototypeData{
					WrapperName:     "ProvidePrototype" + comp.StructName + "As" + iface.InterfaceName,
					FieldName:       iface.InterfaceName + "Factory",
					ConstructorCall: compPrefix + comp.ConstructorName,
					ReturnType:      ifacePrefix + iface.InterfaceName,
					Params:          pData,
					HasCleanup:      comp.HasCleanup,
					HasError:        comp.HasError,
					IsConfig:        isConfig,
				})
			}

		} else {
			wrapperName := comp.ConstructorName
			callPrefix := compPrefix

			if isConfig {
				callPrefix = ""
				configWrappers = append(configWrappers, configWrapperData{
					WrapperName:         wrapperName,
					ConfigPackagePrefix: configPkgPrefix,
					ConfigStructName:    comp.ConfigStructName,
					ConfigMethodName:    comp.ConfigMethodName,
					ReturnType:          retType,
					Params:              pData,
					HasCleanup:          comp.HasCleanup,
					HasError:            comp.HasError,
					IsPrototype:         false,
				})
			}

			providers = append(providers, providerData{
				StructName:      comp.StructName,
				CallPrefix:      callPrefix,
				TypePrefix:      compPrefix,
				ConstructorName: wrapperName,
				IsPointer:       comp.IsPointer,
			})

			for _, iface := range comp.Implements {
				ifacePrefix := ""
				if iface.PackageName != pkgName {
					if iface.PackageName == "main" {
						return errs.Wrap(ErrMainInterfaceLeak, "cannot generate container in package '%s' because interface '%s' belongs to package 'main'. Change output dir (-o) to your main directory or move the interface.", pkgName, iface.InterfaceName)
					}
					ifacePrefix = iface.PackageName + "."
					importSet[iface.PackagePath] = true
				}

				bindings = append(bindings, bindingData{
					InterfacePrefix: ifacePrefix,
					InterfaceName:   iface.InterfaceName,
					ComponentPrefix: compPrefix,
					StructName:      comp.StructName,
					IsPointer:       comp.IsPointer,
				})
			}
		}
	}

	var sliceBindingsData []sliceBindingData
	for _, sb := range genCtx.SliceBindings {
		slices.SortFunc(sb.Implementations, func(a, b *engine.ComponentMetadata) int {
			return cmp.Compare(a.Order, b.Order)
		})
		ifacePrefix := ""
		if sb.Interface.PackageName != pkgName && sb.Interface.PackageName != "main" {
			ifacePrefix = sb.Interface.PackageName + "."
			importSet[sb.Interface.PackagePath] = true
		}

		var impls []sliceImplData
		for i, impl := range sb.Implementations {
			typePrefix := ""
			if impl.IsPointer {
				typePrefix = "*"
			}
			if impl.PackageName != pkgName && impl.PackageName != "main" {
				typePrefix += impl.PackageName + "."
				importSet[impl.PackagePath] = true
			}
			impls = append(impls, sliceImplData{
				ParamName:  fmt.Sprintf("p%d", i),
				TypePrefix: typePrefix,
				StructName: impl.StructName,
			})
		}

		sliceBindingsData = append(sliceBindingsData, sliceBindingData{
			InterfacePrefix: ifacePrefix,
			InterfaceName:   sb.Interface.InterfaceName,
			Implementations: impls,
		})
	}

	data.Providers = providers
	data.Prototypes = prototypes
	data.ConfigWrappers = configWrappers
	data.Bindings = bindings
	data.SliceBindings = sliceBindingsData

	for imp := range importSet {
		if generatedPkgPath != "" && imp == generatedPkgPath {
			continue
		}
		data.Imports = append(data.Imports, imp)
	}

	tmpl, err := template.New("wire").Parse(wireTemplate)
	if err != nil {
		chainErr := fmt.Errorf("%w: %w", ErrParseTemplate, err)
		return errs.Wrap(chainErr, "template parsing failed")
	}

	tempFilePath := filepath.Join(absOutDir, "flora_injector.go")
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		chainErr := fmt.Errorf("%w: %w", ErrExecuteTemplate, err)
		return errs.Wrap(chainErr, "failed to apply data to template")
	}

	log.Debug("Writing temporary wire template", "path", tempFilePath)
	if err := os.WriteFile(tempFilePath, buf.Bytes(), 0644); err != nil {
		chainErr := fmt.Errorf("%w: %w", ErrWriteTempFile, err)
		return errs.Wrap(chainErr, "path: %s", tempFilePath)
	}

	defer func() {
		os.Remove(tempFilePath)
		tidyCmd := exec.Command("go", "mod", "tidy")
		tidyCmd.Dir = absOutDir
		_ = tidyCmd.Run()
	}()

	log.Debug("Ensuring google/wire dependency is present...")
	getCmd := exec.Command("go", "get", "github.com/google/wire@latest")
	getCmd.Dir = absOutDir
	if err := getCmd.Run(); err != nil {
		chainErr := fmt.Errorf("%w: %w", ErrEnsureWireDependency, err)
		return errs.Wrap(chainErr, "failed running 'go get github.com/google/wire@latest' in %s", absOutDir)
	}

	log.Debug("Running DI engine via Google Wire...")
	cmd := exec.Command("go", "run", "github.com/google/wire/cmd/wire@latest", "gen", ".")
	cmd.Dir = absOutDir

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		chainErr := fmt.Errorf("%w: %w", ErrWireExecution, err)
		return errs.Wrap(chainErr, "stderr:\n%s", stderr.String())
	}

	generatedWireFile := filepath.Join(absOutDir, "wire_gen.go")
	finalFloraFile := filepath.Join(absOutDir, "flora_container.go")

	log.Debug("Renaming generated file", "from", generatedWireFile, "to", finalFloraFile)
	if err := os.Rename(generatedWireFile, finalFloraFile); err != nil {
		chainErr := fmt.Errorf("%w: %w", ErrRenameGeneratedFile, err)
		return errs.Wrap(chainErr, "from %s to %s", generatedWireFile, finalFloraFile)
	}

	return nil
}
