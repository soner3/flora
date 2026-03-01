package main

import "net/http"

func main() {
	container, cleanup, err := InitializeContainer()
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// Run the fully wired Webserver
	if err = container.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}

}
