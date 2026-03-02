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

	// Let's simulate two different users making a request to our server.
	// Since the server uses a prototype factory, each request will get a UNIQUE session!
	container.Server.HandleRequest("Alice")

	// A small sleep so the random number generator creates a different ID
	container.Server.HandleRequest("Bob")
}
