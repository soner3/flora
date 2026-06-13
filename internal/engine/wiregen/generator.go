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

const (
	ContainerFileName = "flora_container.go"
	WireFileName      = "wire_gen.go"
	InjectorFileName  = "flora_injector.go"
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

func sanitize(name string) string {
	var sb strings.Builder
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
			sb.WriteRune(r)
		} else {
			sb.WriteRune('_')
		}
	}
	return sb.String()
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

{{range .Aliases}}
type {{.AliasName}} {{.OriginalType}}
{{end}}

{{range .Wrappers}}
{{if .IsPrototype}}
func {{.WrapperName}}({{range $index, $param := .Params}}{{if $index}}, {{end}}{{$param.Name}} {{$param.Type}}{{end}}) func() ({{.ReturnType}}{{if .HasCleanup}}, func(){{end}}{{if .HasError}}, error{{end}}) {
    return func() ({{.ReturnType}}{{if .HasCleanup}}, func(){{end}}{{if .HasError}}, error{{end}}) {
        {{if .IsConfig}}
        flora_cfg_struct := {{.ConfigPrefix}}{{.ConfigStruct}}{}
        return flora_cfg_struct.{{.ConfigMethod}}({{range $index, $param := .Params}}{{if $index}}, {{end}}{{$param.CastCall}}{{end}})
        {{else}}
        return {{.DirectCall}}({{range $index, $param := .Params}}{{if $index}}, {{end}}{{$param.CastCall}}{{end}})
        {{end}}
    }
}
{{else}}
func {{.WrapperName}}({{range $index, $param := .Params}}{{if $index}}, {{end}}{{$param.Name}} {{$param.Type}}{{end}}) ({{.ReturnType}}{{if .HasCleanup}}, func(){{end}}{{if .HasError}}, error{{end}}) {
    {{if .IsConfig}}
    flora_cfg_struct := {{.ConfigPrefix}}{{.ConfigStruct}}{}
    {{if and .HasCleanup .HasError}}
    val, cleanup, err := flora_cfg_struct.{{.ConfigMethod}}({{range $index, $param := .Params}}{{if $index}}, {{end}}{{$param.CastCall}}{{end}})
    {{else if .HasCleanup}}
    val, cleanup := flora_cfg_struct.{{.ConfigMethod}}({{range $index, $param := .Params}}{{if $index}}, {{end}}{{$param.CastCall}}{{end}})
    {{else if .HasError}}
    val, err := flora_cfg_struct.{{.ConfigMethod}}({{range $index, $param := .Params}}{{if $index}}, {{end}}{{$param.CastCall}}{{end}})
    {{else}}
    val := flora_cfg_struct.{{.ConfigMethod}}({{range $index, $param := .Params}}{{if $index}}, {{end}}{{$param.CastCall}}{{end}})
    {{end}}
    {{else}}
    {{if and .HasCleanup .HasError}}
    val, cleanup, err := {{.DirectCall}}({{range $index, $param := .Params}}{{if $index}}, {{end}}{{$param.CastCall}}{{end}})
    {{else if .HasCleanup}}
    val, cleanup := {{.DirectCall}}({{range $index, $param := .Params}}{{if $index}}, {{end}}{{$param.CastCall}}{{end}})
    {{else if .HasError}}
    val, err := {{.DirectCall}}({{range $index, $param := .Params}}{{if $index}}, {{end}}{{$param.CastCall}}{{end}})
    {{else}}
    val := {{.DirectCall}}({{range $index, $param := .Params}}{{if $index}}, {{end}}{{$param.CastCall}}{{end}})
    {{end}}
    {{end}}

    {{if .NeedsProviderAlias}}
    return {{.ReturnType}}(val){{if .HasCleanup}}, cleanup{{end}}{{if .HasError}}, err{{end}}
    {{else}}
    return val{{if .HasCleanup}}, cleanup{{end}}{{if .HasError}}, err{{end}}
    {{end}}
}
{{end}}
{{end}}

{{range .Bindings}}
func {{.WrapperName}}(val {{.ParamType}}) {{.InterfaceType}} {
    {{if .NeedsCast}}
    return ({{.OriginalType}})(val)
    {{else}}
    return val
    {{end}}
}
{{end}}

{{range .PrototypeBindings}}
func {{.WrapperName}}(factory {{.ParamType}}) func() ({{.InterfaceType}}{{if .HasCleanup}}, func(){{end}}{{if .HasError}}, error{{end}}) {
    return func() ({{.InterfaceType}}{{if .HasCleanup}}, func(){{end}}{{if .HasError}}, error{{end}}) {
        {{if and .HasCleanup .HasError}}
        val, cleanup, err := factory()
        return val, cleanup, err
        {{else if .HasCleanup}}
        val, cleanup := factory()
        return val, cleanup
        {{else if .HasError}}
        val, err := factory()
        return val, err
        {{else}}
        return factory()
        {{end}}
    }
}
{{end}}

{{range .SliceBindings}}
func ProvideSliceOf{{.InterfaceName}}({{range $index, $impl := .Implementations}}{{if $index}}, {{end}}{{$impl.ParamName}} {{$impl.ParamType}}{{end}}) []{{.InterfaceType}} {
    return []{{.InterfaceType}}{
        {{range .Implementations}}{{if .NeedsCast}}({{.OriginalType}})({{.ParamName}}),{{else}}{{.ParamName}},{{end}}{{end}}
    }
}
{{end}}

{{range .PrimaryAliases}}
func ProvidePrimary_{{.AliasName}}(val {{.AliasName}}) {{.OriginalType}} {
    return ({{.OriginalType}})(val)
}
{{end}}

type FloraContainer struct {
    {{range .Fields}}
    {{.Name}} {{.Type}}
    {{end}}
}

func InitializeContainer() (*FloraContainer, func(), error) {
    wire.Build(
        {{range .Providers}}
        {{.Call}},
        {{end}}
        {{range .PrimaryAliases}}
        ProvidePrimary_{{.AliasName}},
        {{end}}
        {{range .Bindings}}
        {{.WrapperName}},
        {{end}}
        {{range .PrototypeBindings}}
        {{.WrapperName}},
        {{end}}
        {{range .SliceBindings}}
        ProvideSliceOf{{.InterfaceName}},
        {{end}}
        wire.Struct(new(FloraContainer), "*"),
    )
    return nil, nil, nil
}
`

type aliasData struct {
	AliasName    string
	OriginalType string
}

type paramData struct {
	Name     string
	Type     string
	CastCall string
}

type wrapperData struct {
	WrapperName        string
	IsConfig           bool
	ConfigPrefix       string
	ConfigStruct       string
	ConfigMethod       string
	DirectCall         string
	Params             []paramData
	ReturnType         string
	HasCleanup         bool
	HasError           bool
	IsPrototype        bool
	NeedsProviderAlias bool
}

type bindingData struct {
	WrapperName   string
	ParamType     string
	OriginalType  string
	InterfaceType string
	NeedsCast     bool
}

type prototypeBindingData struct {
	WrapperName   string
	ParamType     string
	InterfaceType string
	HasCleanup    bool
	HasError      bool
}

type sliceImplData struct {
	ParamName    string
	ParamType    string
	OriginalType string
	NeedsCast    bool
}

type sliceBindingData struct {
	InterfaceName   string
	InterfaceType   string
	Implementations []sliceImplData
}

type primaryAliasData struct {
	AliasName    string
	OriginalType string
}

type containerField struct {
	Name string
	Type string
}

type providerData struct {
	Call string
}

type templateData struct {
	PackageName       string
	Imports           []string
	Aliases           map[string]aliasData
	Wrappers          []wrapperData
	Bindings          []bindingData
	PrototypeBindings []prototypeBindingData
	SliceBindings     []sliceBindingData
	PrimaryAliases    []primaryAliasData
	Providers         []providerData
	Fields            []containerField
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

	pkgName := engine.ResolvePackageName(absOutDir)

	var generatedPkgPath string
	for _, comp := range genCtx.Components {
		if comp.PackageName == pkgName {
			generatedPkgPath = comp.PackagePath
			break
		}
	}

	importSet := make(map[string]bool)
	typeCount := make(map[string]int)
	reqQualifiers := make(map[string]bool)

	formatParamType := func(pType string) string {
		pType = strings.ReplaceAll(pType, "*"+pkgName+".", "*")
		pType = strings.ReplaceAll(pType, "[]"+pkgName+".", "[]")
		if after, ok := strings.CutPrefix(pType, pkgName+"."); ok {
			pType = after
		}
		return pType
	}

	for _, comp := range genCtx.Components {
		compPrefix := ""
		if comp.PackageName != pkgName && !isBuiltInType(comp.StructName) {
			compPrefix = comp.PackageName + "."
			importSet[comp.PackagePath] = true
		}

		retType := compPrefix + comp.StructName
		if comp.IsPointer {
			retType = "*" + retType
		}

		typeCount[retType]++

		for _, p := range comp.Params {
			for _, imp := range p.Imports {
				importSet[imp] = true
			}
			if p.RequestedQualifier != "" {
				reqQualifiers[p.RequestedQualifier] = true
			}
		}
	}

	var wrappers []wrapperData
	var bindings []bindingData
	var prototypeBindings []prototypeBindingData
	var primaryAliases []primaryAliasData
	var providers []providerData
	var fields []containerField
	aliases := make(map[string]aliasData)

	for _, comp := range genCtx.Components {
		isBuiltIn := isBuiltInType(comp.StructName)
		isPrototype := comp.Scope == "prototype"
		isConfig := comp.ConfigStructName != ""

		compPrefix := ""
		if comp.PackageName != pkgName && !isBuiltIn {
			if comp.PackageName == "main" {
				return errs.Wrap(ErrMainComponentLeak, "cannot generate container in package '%s' because component '%s' belongs to package 'main'", pkgName, comp.StructName)
			}
			compPrefix = comp.PackageName + "."
			importSet[comp.PackagePath] = true
		}

		configPkgPrefix := ""
		if isConfig && comp.ConfigPackageName != pkgName {
			if comp.ConfigPackageName == "main" {
				return errs.Wrap(ErrMainComponentLeak, "cannot generate container because config '%s' belongs to package 'main'.", comp.ConfigStructName)
			}
			configPkgPrefix = comp.ConfigPackageName + "."
			importSet[comp.ConfigPackagePath] = true
		}

		retType := compPrefix + comp.StructName
		if comp.IsPointer {
			retType = "*" + retType
		}

		needsProviderAlias := typeCount[retType] > 1 || reqQualifiers[comp.QualifierName]

		providerAliasName := "FloraQualifier_" + sanitize(comp.QualifierName)

		if needsProviderAlias {
			aliases[providerAliasName] = aliasData{
				AliasName:    providerAliasName,
				OriginalType: retType,
			}
		}

		var pData []paramData
		needsWrapper := isConfig || isPrototype || needsProviderAlias

		for _, p := range comp.Params {
			pType := formatParamType(p.Type)
			castCall := p.Name

			if p.RequestedQualifier != "" {
				needsWrapper = true
				aliasName := "FloraQualifier_" + sanitize(p.RequestedQualifier)
				pType = aliasName
				castCall = fmt.Sprintf("(%s)(%s)", formatParamType(p.Type), p.Name)
			}
			pData = append(pData, paramData{Name: p.Name, Type: pType, CastCall: castCall})
		}

		finalReturnType := retType
		if needsProviderAlias {
			finalReturnType = providerAliasName
		}

		prefix := "ProvideWrapper_"
		if isPrototype {
			prefix = "ProvidePrototype_"
		}

		wrapperName := prefix + sanitize(comp.ConstructorName)
		if isConfig {
			wrapperName = prefix + sanitize(comp.ConfigStructName) + "_" + sanitize(comp.ConfigMethodName)
		}

		if needsWrapper {
			wrappers = append(wrappers, wrapperData{
				WrapperName:        wrapperName,
				IsConfig:           isConfig,
				ConfigPrefix:       configPkgPrefix,
				ConfigStruct:       comp.ConfigStructName,
				ConfigMethod:       comp.ConfigMethodName,
				DirectCall:         compPrefix + comp.ConstructorName,
				Params:             pData,
				ReturnType:         finalReturnType,
				HasCleanup:         comp.HasCleanup,
				HasError:           comp.HasError,
				IsPrototype:        isPrototype,
				NeedsProviderAlias: needsProviderAlias,
			})
			providers = append(providers, providerData{Call: wrapperName})
		} else {
			providers = append(providers, providerData{Call: compPrefix + comp.ConstructorName})
		}

		for _, iface := range comp.Implements {
			if iface.PackageName == "main" {
				if comp.PackageName != "main" {
					return errs.Wrap(ErrMainInterfaceLeak, "interface '%s' belongs to package 'main', but component '%s' is in package '%s'. Go forbids subpackages from importing main. Please move the interface to an internal package.",
						iface.TypeName, comp.StructName, comp.PackageName)
				}

			}
			ifacePrefix := ""
			if iface.PackageName != pkgName {
				ifacePrefix = iface.PackageName + "."
				importSet[iface.PackagePath] = true
			}

			if isPrototype {
				retSig := finalReturnType
				if comp.HasCleanup && comp.HasError {
					retSig += ", func(), error"
				} else if comp.HasCleanup {
					retSig += ", func()"
				} else if comp.HasError {
					retSig += ", error"
				}
				paramType := "func() (" + retSig + ")"

				pbWrapperName := fmt.Sprintf("ProvidePrototypeBinding_%s_As_%s", sanitize(comp.QualifierName), sanitize(iface.TypeName))
				prototypeBindings = append(prototypeBindings, prototypeBindingData{
					WrapperName:   pbWrapperName,
					ParamType:     paramType,
					InterfaceType: ifacePrefix + iface.TypeName,
					HasCleanup:    comp.HasCleanup,
					HasError:      comp.HasError,
				})

			} else {
				bWrapperName := fmt.Sprintf("ProvideBinding_%s_As_%s", sanitize(comp.QualifierName), sanitize(iface.TypeName))
				bindings = append(bindings, bindingData{
					WrapperName:   bWrapperName,
					ParamType:     finalReturnType,
					OriginalType:  retType,
					InterfaceType: ifacePrefix + iface.TypeName,
					NeedsCast:     needsProviderAlias,
				})
			}
		}

		if needsProviderAlias && comp.IsPrimary {
			primaryAliases = append(primaryAliases, primaryAliasData{
				AliasName:    providerAliasName,
				OriginalType: retType,
			})
		}

		fieldName := sanitize(comp.QualifierName)

		if len(fieldName) > 0 {
			r := []rune(fieldName)
			if r[0] >= 'a' && r[0] <= 'z' {
				r[0] -= 32
			}
			fieldName = string(r)
		}

		if isPrototype {
			retSig := finalReturnType
			if comp.HasCleanup && comp.HasError {
				retSig += ", func(), error"
			} else if comp.HasCleanup {
				retSig += ", func()"
			} else if comp.HasError {
				retSig += ", error"
			}
			fields = append(fields, containerField{Name: fieldName + "Factory", Type: "func() (" + retSig + ")"})
		} else {
			fields = append(fields, containerField{Name: fieldName, Type: finalReturnType})
		}
	}

	var sliceBindingsData []sliceBindingData
	for _, sb := range genCtx.SliceBindings {
		slices.SortFunc(sb.Implementations, func(a, b *engine.ComponentMetadata) int {
			return cmp.Compare(a.Order, b.Order)
		})

		if sb.Type.PackageName == "main" {
			for _, impl := range sb.Implementations {
				if impl.PackageName != "main" {
					return errs.Wrap(ErrMainInterfaceLeak, "slice interface '%s' belongs to package 'main', but component '%s' is in package '%s'. Go forbids subpackages from importing main. Please move the interface to an internal package.",
						sb.Type.TypeName, impl.StructName, impl.PackageName)
				}
			}
			if pkgName != "main" {
				return errs.Wrap(ErrMainInterfaceLeak, "slice interface belongs to package 'main' (Go forbids importing main)")
			}
		}

		if sb.Type.PackagePath != "" && sb.Type.PackageName != pkgName {
			importSet[sb.Type.PackagePath] = true
		}

		formattedType := sb.Type.TypeName

		formattedType = strings.ReplaceAll(formattedType, "*"+pkgName+".", "*")
		if after, ok := strings.CutPrefix(formattedType, pkgName+"."); ok {
			formattedType = after
		}

		var impls []sliceImplData
		for i, impl := range sb.Implementations {
			implPrefix := ""
			if impl.PackageName != pkgName && !isBuiltInType(impl.StructName) {
				implPrefix = impl.PackageName + "."
			}
			implRetType := implPrefix + impl.StructName
			if impl.IsPointer {
				implRetType = "*" + implRetType
			}

			implNeedsAlias := typeCount[implRetType] > 1 || reqQualifiers[impl.QualifierName]
			paramType := implRetType
			if implNeedsAlias {
				paramType = "FloraQualifier_" + sanitize(impl.QualifierName)
			}

			impls = append(impls, sliceImplData{
				ParamName:    fmt.Sprintf("p%d", i),
				ParamType:    paramType,
				OriginalType: implRetType,
				NeedsCast:    implNeedsAlias,
			})
		}

		cleanName := sb.Type.TypeName
		cleanName = strings.ReplaceAll(cleanName, "*", "Ptr_")
		cleanName = strings.ReplaceAll(cleanName, "[]", "SliceOf_")
		cleanName = strings.ReplaceAll(cleanName, ".", "_")

		funcNameSuffix := sanitize(cleanName)

		if len(funcNameSuffix) > 0 {
			r := []rune(funcNameSuffix)
			if r[0] >= 'a' && r[0] <= 'z' {
				r[0] -= 32
			}
			funcNameSuffix = string(r)
		}

		sliceBindingsData = append(sliceBindingsData, sliceBindingData{
			InterfaceName:   funcNameSuffix,
			InterfaceType:   formattedType,
			Implementations: impls,
		})
	}

	data := templateData{
		PackageName:       pkgName,
		Aliases:           aliases,
		Wrappers:          wrappers,
		Bindings:          bindings,
		PrototypeBindings: prototypeBindings,
		SliceBindings:     sliceBindingsData,
		PrimaryAliases:    primaryAliases,
		Providers:         providers,
		Fields:            fields,
	}

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

	tempFilePath := filepath.Join(absOutDir, InjectorFileName)
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

	generatedWireFile := filepath.Join(absOutDir, WireFileName)
	finalFloraFile := filepath.Join(absOutDir, ContainerFileName)

	log.Debug("Renaming generated file", "from", generatedWireFile, "to", finalFloraFile)
	if err := os.Rename(generatedWireFile, finalFloraFile); err != nil {
		chainErr := fmt.Errorf("%w: %w", ErrRenameGeneratedFile, err)
		return errs.Wrap(chainErr, "from %s to %s", generatedWireFile, finalFloraFile)
	}

	return nil
}
