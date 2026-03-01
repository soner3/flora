package main

import "fmt"

func main() {
	// Initialize the generated container
	container, cleanup, err := InitializeContainer()
	if err != nil {
		panic(err)
	}

	// Deferring cleanup ensures the Logger's cleanup function is actually called
	defer cleanup()

	fmt.Println("-------------------------------------------------")

	// Use the fully wired Worker
	container.Worker.DoWork()

	fmt.Println("-------------------------------------------------")
	fmt.Println("Application is shutting down. Watch the cleanup happen:")
}
