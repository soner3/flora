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
