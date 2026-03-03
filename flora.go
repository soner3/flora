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
package flora

// Component is a marker struct designed to be embedded into your own domain structs
// (e.g., Services, Repositories, Handlers). Embedding this tells Flora's AST scanner
// to automatically discover the struct, find its constructor, and wire its dependencies.
//
// By default, Flora looks for an exported function named 'New<StructName>' (e.g., NewUserService),
// treats the component as a Singleton, and automatically binds it to any interfaces it implements.
//
// You can customize this behavior using the `flora` struct tag:
//
//	type MyService struct {
//		flora.Component `flora:"primary,scope=prototype,constructor=BuildService,order=1"`
//	}
//
// Available Tag Options:
//   - constructor=<Name>: Overrides the default constructor lookup (e.g., constructor=BuildMyService).
//   - primary: Marks this component as the primary implementation of an interface to resolve collisions.
//   - scope=prototype: Changes the lifecycle from Singleton to a Factory (fresh instance per injection).
//   - order=<Number>: Defines the exact sorting position when injected as part of a Slice ([]Interface).
type Component struct{}

// Configuration is a marker struct designed to be embedded into configuration provider structs.
// It is strictly meant for integrating third-party types (like *sql.DB), functional adapters
// (like HTTP middlewares), or any types you do not own and cannot embed flora.Component into.
//
// Flora scans all exported methods of a Configuration struct and registers their return types
// as providers in the DI container.
//
// Because Go does not allow struct tags on functions, you configure these provider methods
// using "Magic Comments" placed immediately above the method signature.
//
// Example:
//
//	type DatabaseConfig struct { flora.Configuration }
//
//	// flora:primary,scope=prototype
//	func (c *DatabaseConfig) ProvideDatabase() *sql.DB { ... }
//
// Available Magic Comments:
//   - // flora:primary : Marks the returned type as the primary implementation to resolve collisions.
//   - // flora:scope=prototype : Changes the lifecycle to a Factory (fresh instance per injection).
//   - // flora:order=<Number> : Defines the exact sorting position when injected via Slice ([]Interface).
type Configuration struct{}
