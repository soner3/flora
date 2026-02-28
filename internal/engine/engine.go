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
package engine

type InterfaceMetadata struct {
	PackageName   string
	PackagePath   string
	InterfaceName string
}

type ParamMetadata struct {
	Name    string
	Type    string
	Imports []string
}

type ComponentMetadata struct {
	PackageName       string
	PackagePath       string
	StructName        string
	ConstructorName   string
	IsPrimary         bool
	Scope             string
	IsPointer         bool
	HasCleanup        bool
	HasError          bool
	Order             int
	ConfigStructName  string
	ConfigMethodName  string
	ConfigPackageName string
	ConfigPackagePath string
	Implements        []InterfaceMetadata
	Params            []ParamMetadata
}

type SliceBindingMetadata struct {
	Interface       InterfaceMetadata
	Implementations []*ComponentMetadata
}

type GeneratorContext struct {
	Components    []*ComponentMetadata
	SliceBindings []*SliceBindingMetadata
}

type Generator interface {
	Generate(targetDir string, genCtx *GeneratorContext) error
}
