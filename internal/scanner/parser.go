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
	"fmt"
	"go/types"
	"log/slog"
	"math"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/soner3/flora/internal/engine"
	"github.com/soner3/flora/internal/errs"
	"golang.org/x/tools/go/packages"
)

var (
	ErrProviderFuncNotFound = errors.New("provider not found")
	ErrInvalidProviderFunc  = errors.New("invalid provider func")
	ErrUnknownMarker        = errors.New("unknown marker")
	ErrInterfaceCollision   = errors.New("interface collision")
	ErrInvalidInterface     = errors.New("invalid interface")
	ErrNoImplementation     = errors.New("no component implements interface")
	ErrInvalidSlice         = errors.New("invalid slice")
	ErrInvalidScope         = errors.New("invalid scope")
	ErrInvalidOrder         = errors.New("invalid order")
)

const (
	ScopeSingleton = "singleton"
	ScopePrototype = "prototype"

	ComponentMarker = "github.com/soner3/flora.Component"
)

var scopes = []string{
	ScopeSingleton,
	ScopePrototype,
}

var markers = []string{
	ComponentMarker,
}

type scannedComponent struct {
	Metadata *engine.ComponentMetadata
	PtrType  *types.Pointer
}

type componentInfo struct {
	Pkg        *packages.Package
	Name       string
	TypeName   *types.TypeName
	StructType *types.Struct
	Marker     string
	Tag        string
}

var log = slog.With("pkg", "scanner")

// ParsePackages parses the given packages and returns a GeneratorContext
// containing the parsed components and slice bindings
func ParsePackages(pkgs []*packages.Package) (*engine.GeneratorContext, error) {
	log.Debug("Parsing components from packages", "package_count", len(pkgs))

	compInfos := parseMarkedComponents(pkgs)

	log.Debug("Marked components found", "count", len(*compInfos))

	neededInterfaces := make(map[string]types.Type)
	neededSlices := make(map[string]types.Type)
	scannedComponents := make([]*scannedComponent, 0)

	for _, compInfo := range *compInfos {
		scannedComp, err := processComponent(&compInfo, &neededInterfaces, &neededSlices)
		if err != nil {
			return nil, err
		}
		scannedComponents = append(scannedComponents, scannedComp)
	}

	log.Debug("Resolving interface implementations", "interfaces_needed", len(neededInterfaces))
	if err := bindInterfacesToComponents(scannedComponents, neededInterfaces); err != nil {
		return nil, err
	}

	log.Debug("Resolving slice bindings", "slices_needed", len(neededSlices))

	sliceBindings, err := bindSlicesToComponents(scannedComponents, neededSlices)
	if err != nil {
		return nil, err
	}

	var finalMetadata []*engine.ComponentMetadata
	for _, comp := range scannedComponents {
		finalMetadata = append(finalMetadata, comp.Metadata)
	}

	log.Debug("Successfully parsed all components", "total", len(finalMetadata), "slices", len(sliceBindings))

	return &engine.GeneratorContext{
		Components:    finalMetadata,
		SliceBindings: sliceBindings,
	}, nil

}

// parseMarkedComponents finds all components in the given packages that are marked with the flora markers
func parseMarkedComponents(pkgs []*packages.Package) *[]componentInfo {

	components := make([]componentInfo, 0)

	for _, pkg := range pkgs {
		scope := pkg.Types.Scope()
		for _, name := range scope.Names() {
			obj := scope.Lookup(name)
			if typeName, ok := obj.(*types.TypeName); ok {
				if structType, ok := typeName.Type().Underlying().(*types.Struct); ok {
					isComponent, marker, tag := isMarkedWith(structType)

					if isComponent {
						components = append(components, componentInfo{
							Pkg:        pkg,
							Name:       name,
							TypeName:   typeName,
							StructType: structType,
							Marker:     marker,
							Tag:        tag,
						})
					}
				}
			}
		}
	}

	return &components
}

// isMarkedWith checks if the struct is marked with any of the flora markers
// and returns the marker and tag
func isMarkedWith(structType *types.Struct) (bool, string, string) {
	for i := 0; i < structType.NumFields(); i++ {
		field := structType.Field(i)
		for _, marker := range markers {
			if field.Anonymous() && field.Type().String() == marker {
				return true, marker, structType.Tag(i)
			}
		}

	}
	return false, "", ""
}

// parseFloraTag parses the flora tag and sets the metadata accordingly
func parseFloraTag(rawTag string, metadata *engine.ComponentMetadata) error {
	metadata.ConstructorName = "New" + metadata.StructName
	metadata.IsPrimary = false
	metadata.Scope = ScopeSingleton
	metadata.Order = math.MaxInt32

	if rawTag == "" {
		return nil
	}

	structTag := reflect.StructTag(rawTag)
	val := structTag.Get("flora")

	if val == "" {
		return nil
	}

	parts := strings.SplitSeq(val, ",")
	for part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		switch {
		case part == "primary":
			metadata.IsPrimary = true
		case strings.HasPrefix(part, "constructor="):
			metadata.ConstructorName = strings.TrimPrefix(part, "constructor=")
		case strings.HasPrefix(part, "scope="):
			scope := strings.TrimPrefix(part, "scope=")
			if !slices.Contains(scopes, scope) {
				return errs.Wrap(ErrInvalidScope, "invalid scope '%s' for component '%s' in package '%s'", scope, metadata.StructName, metadata.PackageName)
			}
			metadata.Scope = scope
		case strings.HasPrefix(part, "order="):
			orderStr := strings.TrimPrefix(part, "order=")
			order, err := strconv.Atoi(orderStr)
			if err != nil {
				return errs.Wrap(ErrInvalidOrder, "invalid order '%s' for component '%s' in package '%s' (must be an integer)", orderStr, metadata.StructName, metadata.PackageName)
			}
			metadata.Order = order
		default:
			metadata.ConstructorName = part
		}
	}

	return nil
}

// processComponent processes a component and returns a scannedComponent
func processComponent(compInfo *componentInfo, neededInterfaces, neededSlices *map[string]types.Type) (*scannedComponent, error) {
	metadata := &engine.ComponentMetadata{
		StructName:  compInfo.Name,
		PackageName: compInfo.Pkg.Name,
		PackagePath: compInfo.Pkg.PkgPath,
	}

	switch compInfo.Marker {
	case ComponentMarker:
		if err := parseFloraTag(compInfo.Tag, metadata); err != nil {
			return nil, err
		}

		obj := compInfo.Pkg.Types.Scope().Lookup(metadata.ConstructorName)

		err := processProviderFunc(compInfo, metadata, obj, neededInterfaces, neededSlices)
		if err != nil {
			return nil, err
		}

		return &scannedComponent{
			Metadata: metadata,
			PtrType:  types.NewPointer(compInfo.TypeName.Type()),
		}, nil
	default:
		chainErr := fmt.Errorf("%w: %v", ErrUnknownMarker, compInfo.Marker)
		return nil, errs.Wrap(chainErr, "unknown marker '%s' for component '%s' in package '%s'",
			compInfo.Marker, compInfo.Name, compInfo.Pkg.Name)
	}

}

// processProviderFunc validates the provider function and populates
// the needed interfaces and slices in compInfo
func processProviderFunc(compInfo *componentInfo, metadata *engine.ComponentMetadata, obj types.Object, neededInterfaces, neededSlices *map[string]types.Type) error {

	sig, err := validateProviderFunc(compInfo, metadata, obj)
	if err != nil {
		return err
	}

	for v := range sig.Params().Variables() {
		paramType := v.Type()

		if iface, isInterface := paramType.Underlying().(*types.Interface); isInterface {
			if !iface.Empty() {
				(*neededInterfaces)[paramType.String()] = paramType
			}
		}

		if sliceType, isSlice := paramType.(*types.Slice); isSlice {
			elemType := sliceType.Elem()
			if iface, isInterface := elemType.Underlying().(*types.Interface); isInterface {
				if !iface.Empty() {
					(*neededSlices)[elemType.String()] = elemType
				}
			}
		}

		if sigParam, isFunc := paramType.(*types.Signature); isFunc {

			if sigParam.Params().Len() > 0 {
				chainErr := fmt.Errorf("%w: %v", ErrInvalidProviderFunc, sigParam)
				return errs.Wrap(chainErr, "invalid prototype provider func: '%s' for component '%s': prototype provider func must not have parameters",
					metadata.ConstructorName, metadata.StructName)
			}

			if _, _, err := validateReturnValues(sigParam, metadata.ConstructorName, metadata.StructName, metadata.PackageName); err != nil {
				return err
			}

			retType := sigParam.Results().At(0).Type()
			if iface, isInterface := retType.Underlying().(*types.Interface); isInterface {
				if !iface.Empty() {
					(*neededInterfaces)[retType.String()] = retType
				}
			}
		}

	}

	return nil
}

// validateProviderFunc validates the object is a provider function and
// returns the signature if valid
func validateProviderFunc(compInfo *componentInfo, metadata *engine.ComponentMetadata, obj types.Object) (*types.Signature, error) {

	if obj == nil {
		chainErr := fmt.Errorf("%w: %v", ErrProviderFuncNotFound, obj)
		return nil, errs.Wrap(chainErr, "provider '%s' not found for component '%s' in package '%s'",
			metadata.ConstructorName, metadata.StructName, metadata.PackageName)
	}

	funcObj, ok := obj.(*types.Func)
	if !ok {
		chainErr := fmt.Errorf("%w: %v", ErrInvalidProviderFunc, funcObj)
		return nil, errs.Wrap(chainErr, "expected '%s' to be a function for component '%s', but it is a %T",
			metadata.ConstructorName, metadata.StructName, obj)
	}

	sig := funcObj.Type().(*types.Signature)

	hasCleanup, hasErr, err := validateReturnValues(sig, metadata.ConstructorName, metadata.StructName, metadata.PackageName)
	if err != nil {
		return nil, err
	}

	metadata.HasCleanup = hasCleanup
	metadata.HasError = hasErr

	results := sig.Results()
	firstType := results.At(0).Type()

	var baseRetType types.Type

	if ptr, isPtr := firstType.(*types.Pointer); isPtr {
		metadata.IsPointer = true
		baseRetType = ptr.Elem()
	} else {
		metadata.IsPointer = false
		baseRetType = firstType
	}

	if !types.Identical(baseRetType, compInfo.TypeName.Type()) {
		chainErr := fmt.Errorf("%w: %v", ErrInvalidProviderFunc, firstType)
		return nil, errs.Wrap(chainErr, "invalid provider func: '%s' returns '%s', but must return '%s' or '*%s'",
			metadata.ConstructorName, firstType.String(), compInfo.TypeName.Name(), compInfo.TypeName.Name())
	}

	params := sig.Params()

	for i := 0; i < params.Len(); i++ {
		param := params.At(i)
		paramType := param.Type()
		var baseParamType types.Type

		if ptr, isPtr := paramType.(*types.Pointer); isPtr {
			baseParamType = ptr.Elem()
		} else {
			baseParamType = paramType
		}

		if types.Identical(baseParamType, compInfo.TypeName.Type()) {
			chainErr := fmt.Errorf("%w: %v", ErrInvalidProviderFunc, paramType)
			return nil, errs.Wrap(chainErr, "circular dependency: provider func '%s' for component '%s' cannot require its own type as a parameter",
				metadata.ConstructorName, metadata.StructName)
		}

		var imports []string
		qualifier := func(p *types.Package) string {
			if p.Path() != metadata.PackagePath {
				imports = append(imports, p.Path())
			}
			return p.Name()
		}

		paramTypeStr := types.TypeString(paramType, qualifier)

		metadata.Params = append(metadata.Params, engine.ParamMetadata{
			Name:    fmt.Sprintf("p%d", i),
			Type:    paramTypeStr,
			Imports: imports,
		})
	}

	return sig, nil
}

// isCleanupFunc checks if the type is a cleanup function (func())
func isCleanupFunc(t types.Type) bool {
	sig, ok := t.(*types.Signature)
	if !ok {
		return false
	}
	return sig.Params().Len() == 0 && sig.Results().Len() == 0
}

// bindInterfacesToComponents binds the needed interfaces to the components that implement them
func bindInterfacesToComponents(components []*scannedComponent, neededInterfaces map[string]types.Type) error {
	for neededName, neededType := range neededInterfaces {
		iface := neededType.Underlying().(*types.Interface)

		var implementers []*scannedComponent

		for _, comp := range components {
			if types.Implements(comp.PtrType, iface) {
				implementers = append(implementers, comp)
			}
		}

		bindToComp := func(comp *scannedComponent, ifaceType types.Type) error {
			if named, ok := ifaceType.(*types.Named); ok {
				comp.Metadata.Implements = append(comp.Metadata.Implements, engine.InterfaceMetadata{
					PackageName:   named.Obj().Pkg().Name(),
					PackagePath:   named.Obj().Pkg().Path(),
					InterfaceName: named.Obj().Name(),
				})
				log.Debug("Bound interface to component", "interface", neededName, "component", implementers[0].Metadata.StructName)

			} else {
				chainErr := fmt.Errorf("%w: %v", ErrInvalidInterface, ifaceType)
				return errs.Wrap(chainErr, "cannot bind anonymous interface '%s' to component '%s' in package '%s': only named interfaces are supported",
					ifaceType.String(), comp.Metadata.StructName, comp.Metadata.PackageName)
			}
			return nil
		}

		if len(implementers) == 1 {
			if err := bindToComp(implementers[0], neededType); err != nil {
				return err
			}
		} else if len(implementers) > 1 {
			var primaryComp *scannedComponent
			primaryCount := 0

			for i, impl := range implementers {
				if impl.Metadata.IsPrimary {
					primaryCount++
					primaryComp = implementers[i]
				}
			}

			switch primaryCount {
			case 1:
				if err := bindToComp(primaryComp, neededType); err != nil {
					return err
				}
			case 0:
				chainErr := fmt.Errorf("%w: %v", ErrInterfaceCollision, implementers)
				return errs.Wrap(chainErr, "interface collision: %d components implement injected interface '%s', but none is marked 'primary'", len(implementers), neededName)
			default:
				chainErr := fmt.Errorf("%w: %v", ErrInterfaceCollision, implementers)
				return errs.Wrap(chainErr, "interface collision: multiple components implementing '%s' are marked as 'primary'", neededName)
			}

		} else {
			chainErr := fmt.Errorf("%w: %v", ErrNoImplementation, neededName)
			return errs.Wrap(chainErr, "no component found that implements interface '%s'", neededName)
		}
	}
	return nil
}

// bindSlicesToComponents binds the needed slices to the components that implement them
func bindSlicesToComponents(components []*scannedComponent, neededSlices map[string]types.Type) ([]*engine.SliceBindingMetadata, error) {
	var sliceBindings []*engine.SliceBindingMetadata

	for neededName, neededType := range neededSlices {
		iface := neededType.Underlying().(*types.Interface)
		var implementers []*engine.ComponentMetadata

		for _, comp := range components {
			if types.Implements(comp.PtrType, iface) {
				implementers = append(implementers, comp.Metadata)
			}
		}

		if named, ok := neededType.(*types.Named); ok {
			sliceBindings = append(sliceBindings, &engine.SliceBindingMetadata{
				Interface: engine.InterfaceMetadata{
					PackageName:   named.Obj().Pkg().Name(),
					PackagePath:   named.Obj().Pkg().Path(),
					InterfaceName: named.Obj().Name(),
				},
				Implementations: implementers,
			})
			log.Debug("Resolved slice binding", "interface", neededName, "implementations_count", len(implementers))
		} else {
			chainErr := fmt.Errorf("%w: %v", ErrInvalidSlice, neededType)
			return nil, errs.Wrap(chainErr, "cannot bind anonymous slice '%s': only named slices and interfaces are supported", neededType.String())
		}
	}
	return sliceBindings, nil
}

// validateReturnValues validates the return values of a provider function
// Returns: (hasCleanup, hasError, error)
func validateReturnValues(sig *types.Signature, constructorName, structName, pkgName string) (bool, bool, error) {
	results := sig.Results()
	numResults := results.Len()

	if numResults == 0 || numResults > 3 {
		chainErr := fmt.Errorf("%w: %v", ErrInvalidProviderFunc, results)
		return false, false, errs.Wrap(chainErr, "invalid provider func: '%s' for component '%s' in package '%s' must return 1, 2, or 3 values",
			constructorName, structName, pkgName)
	}

	firstType := results.At(0).Type()
	if firstType.String() == "error" || isCleanupFunc(firstType) {
		chainErr := fmt.Errorf("%w: %v", ErrInvalidProviderFunc, firstType)
		return false, false, errs.Wrap(chainErr, "invalid provider func: '%s' for component '%s': 1st return value must neither be 'error' nor 'func()'", constructorName, structName)
	}

	hasCleanup, hasErr := false, false

	if numResults == 2 {
		secondType := results.At(1).Type()
		hasErr = secondType.String() == "error"
		hasCleanup = isCleanupFunc(secondType)
		if !hasErr && !hasCleanup {
			chainErr := fmt.Errorf("%w: %v", ErrInvalidProviderFunc, secondType)
			return false, false, errs.Wrap(chainErr, "invalid provider func: '%s' for component '%s': 2nd return value must be 'error' or 'func()'", constructorName, structName)
		}
	}

	if numResults == 3 {
		secondType := results.At(1).Type()
		if !isCleanupFunc(secondType) {
			chainErr := fmt.Errorf("%w: %v", ErrInvalidProviderFunc, secondType)
			return false, false, errs.Wrap(chainErr, "invalid provider func: '%s' for component '%s': 2nd return value must be 'func()' when returning 3 values", constructorName, structName)
		}
		thirdType := results.At(2).Type()
		if thirdType.String() != "error" {
			chainErr := fmt.Errorf("%w: %v", ErrInvalidProviderFunc, thirdType)
			return false, false, errs.Wrap(chainErr, "invalid provider func: '%s' for component '%s': 3rd return value must be 'error' when returning 3 values", constructorName, structName)
		}
		hasCleanup, hasErr = true, true
	}

	return hasCleanup, hasErr, nil
}
