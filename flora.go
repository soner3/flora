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

// Component is a marker struct that is embedded in components.
// It allows Flora to auto-discover and wire the struct using tags.
type Component struct{}

// Configuration is a marker struct that is embedded in configuration classes.
// Flora will scan all methods of a Configuration struct and register them as
// providers. Methods can be configured using magic comments (e.g., // flora:primary).
type Configuration struct{}
