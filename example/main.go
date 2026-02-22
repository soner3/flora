package main

import (
	"fmt"
	"os"
)

func main() {
	container, err := InitializeContainer()
	if err != nil {
		fmt.Printf("Failed to initialize DI container: %v\n", err)
		os.Exit(1)
	}

	container.UserService.PrintUser()
}
