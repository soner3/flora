package main

import (
	"fmt"
	"os"

	"github.com/soner3/weld/example/weld"
)

func main() {
	container, err := weld.InitializeContainer()
	if err != nil {
		fmt.Printf("Failed to initialize DI container: %v\n", err)
		os.Exit(1)
	}

	container.UserService.PrintUser()
}
