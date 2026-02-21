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
	"go/types"
	"log/slog"
	"reflect"
	"strings"

	"github.com/soner3/weld/internal/engine"
	"github.com/soner3/weld/internal/errs"
	"golang.org/x/tools/go/packages"
)

type scannedComponent struct {
	Metadata *engine.ComponentMetadata
	PtrType  *types.Pointer
}

func ParseComponents(pkgs []*packages.Package) ([]*engine.ComponentMetadata, error) {
	log := slog.With("pkg", "scanner")
	log.Debug("Parsing components from packages", "package_count", len(pkgs))

	var components []scannedComponent
	neededInterfaces := make(map[string]types.Type)

	for _, pkg := range pkgs {
		scope := pkg.Types.Scope()
		for _, name := range scope.Names() {
			obj := scope.Lookup(name)
			if obj == nil {
				continue
			}
			if typeName, ok := obj.(*types.TypeName); ok {
				if structType, ok := typeName.Type().Underlying().(*types.Struct); ok {
					isComponent, rawTag := getWeldComponentInfo(structType)

					if isComponent {
						metadata := &engine.ComponentMetadata{
							PackageName: pkg.Name,
							PackagePath: pkg.PkgPath,
							StructName:  name,
						}

						parseWeldTag(rawTag, metadata)

						log.Debug("Found component", "component", name, "package", pkg.Name, "constructor", metadata.ConstructorName)

						constructorObj := scope.Lookup(metadata.ConstructorName)

						if constructorObj == nil {
							return nil, errs.Wrap(nil, "constructor '%s' not found for component '%s' in package '%s'",
								metadata.ConstructorName, metadata.StructName, pkg.Name)
						}

						funcObj, ok := constructorObj.(*types.Func)
						if !ok {
							return nil, errs.Wrap(nil, "expected '%s' to be a function for component '%s', but it is a %T",
								metadata.ConstructorName, metadata.StructName, constructorObj)
						}

						sig := funcObj.Type().(*types.Signature)
						params := sig.Params()

						for v := range params.Variables() {
							paramType := v.Type()

							if iface, isInterface := paramType.Underlying().(*types.Interface); isInterface {
								if !iface.Empty() {
									neededInterfaces[paramType.String()] = paramType
								}
							}
						}

						components = append(components, scannedComponent{
							Metadata: metadata,
							PtrType:  types.NewPointer(typeName.Type()),
						})
					}
				}
			}
		}
	}

	log.Debug("Resolving interface implementations", "interfaces_needed", len(neededInterfaces))

	for neededName, neededType := range neededInterfaces {
		iface := neededType.Underlying().(*types.Interface)

		var implementers []scannedComponent

		for _, comp := range components {
			if types.Implements(comp.PtrType, iface) {
				implementers = append(implementers, comp)
			}
		}

		if len(implementers) == 1 {
			bindInterfaceToComponent(&implementers[0], neededType)
			log.Debug("Bound interface to component", "interface", neededName, "component", implementers[0].Metadata.StructName)
		} else if len(implementers) > 1 {
			var primaryComp *scannedComponent
			primaryCount := 0

			for i, impl := range implementers {
				if impl.Metadata.IsPrimary {
					primaryCount++
					primaryComp = &implementers[i]
				}
			}

			switch primaryCount {
			case 1:
				bindInterfaceToComponent(primaryComp, neededType)
				log.Debug("Bound interface to primary component", "interface", neededName, "component", primaryComp.Metadata.StructName)
			case 0:
				return nil, errs.Wrap(nil, "interface collision: %d components implement injected interface '%s', but none is marked 'primary'", len(implementers), neededName)
			default:
				return nil, errs.Wrap(nil, "interface collision: multiple components implementing '%s' are marked as 'primary'", neededName)
			}

		} else {
			return nil, errs.Wrap(nil, "no component found that implements interface '%s'", neededName)
		}
	}

	var finalMetadata []*engine.ComponentMetadata
	for _, comp := range components {
		finalMetadata = append(finalMetadata, comp.Metadata)
	}

	log.Debug("Successfully parsed all components", "total", len(finalMetadata))
	return finalMetadata, nil
}

func bindInterfaceToComponent(comp *scannedComponent, ifaceType types.Type) {
	if named, ok := ifaceType.(*types.Named); ok {
		comp.Metadata.Implements = append(comp.Metadata.Implements, engine.InterfaceMetadata{
			PackageName:   named.Obj().Pkg().Name(),
			PackagePath:   named.Obj().Pkg().Path(),
			InterfaceName: named.Obj().Name(),
		})
	}
}

func getWeldComponentInfo(structType *types.Struct) (bool, string) {
	for i := 0; i < structType.NumFields(); i++ {
		field := structType.Field(i)
		if field.Anonymous() && field.Type().String() == "github.com/soner3/weld.Component" {
			return true, structType.Tag(i)
		}
	}
	return false, ""
}

func parseWeldTag(rawTag string, metadata *engine.ComponentMetadata) {
	metadata.ConstructorName = "New" + metadata.StructName
	metadata.IsPrimary = false
	metadata.Scope = "singleton"

	if rawTag == "" {
		return
	}

	structTag := reflect.StructTag(rawTag)
	val := structTag.Get("weld")

	if val == "" {
		return
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
			metadata.Scope = strings.TrimPrefix(part, "scope=")
		default:
			metadata.ConstructorName = part
		}
	}

}
