// Copyright © 2026 Soner Astan astansoner@gmail.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import "fmt"

func main() {
	container, cleanup, err := InitializeContainer()
	if err != nil {
		panic(err)
	}
	defer cleanup()

	fmt.Println("-------------------------------------------------")

	// 1. Run the Interface-based Slice
	container.PluginManager.RunAll()

	fmt.Println("-------------------------------------------------")

	// 2. Run the Pointer-based Slice
	container.Router.PrintRoutes()

	fmt.Println("-------------------------------------------------")
}
