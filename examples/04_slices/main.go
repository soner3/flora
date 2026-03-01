package main

import "fmt"

func main() {
	container, cleanup, err := InitializeContainer()
	if err != nil {
		panic(err)
	}
	defer cleanup()

	fmt.Println("-------------------------------------------------")

	// Run the fully wired PluginManager
	container.PluginManager.RunAll()

	fmt.Println("-------------------------------------------------")
}
