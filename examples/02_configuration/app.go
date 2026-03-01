package main

import (
	"fmt"

	"github.com/soner3/flora"
)

// ---------------------------------------------------------
// 1. Third-Party Struct (Simulation)
// ---------------------------------------------------------
// Imagine this is a struct from an external library (e.g., "database/sql").
// You cannot modify this struct to add 'flora.Component' to it.
type ExternalLogger struct {
	Prefix string
}

func (l *ExternalLogger) PrintInfo(msg string) {
	fmt.Printf("[%s] INFO: %s\n", l.Prefix, msg)
}

// ---------------------------------------------------------
// 2. The Configuration
// ---------------------------------------------------------
// We use a Configuration struct to tell Flora how to create external types.
type AppConfig struct {
	flora.Configuration // Tells Flora to scan this struct for provider methods
}

// Magic comments tell Flora how to register this method in the container.
// flora:primary
func (c *AppConfig) ProvideExternalLogger() (*ExternalLogger, func(), error) {
	fmt.Println("-> [Config] Initializing ExternalLogger...")
	logger := &ExternalLogger{Prefix: "SYSTEM"}

	// Flora natively supports cleanup functions for graceful shutdowns
	cleanup := func() {
		fmt.Println("-> [Config] Cleaning up ExternalLogger resources...")
	}

	// We return the instance, the cleanup function, and nil (no error)
	return logger, cleanup, nil
}

// ---------------------------------------------------------
// 3. The Consumer
// ---------------------------------------------------------
type Worker struct {
	flora.Component // This is our own code, so we use Component
	logger          *ExternalLogger
}

// Flora automatically sees that Worker needs an ExternalLogger
// and uses our AppConfig to provide it!
func NewWorker(logger *ExternalLogger) *Worker {
	return &Worker{logger: logger}
}

func (w *Worker) DoWork() {
	w.logger.PrintInfo("Worker has successfully started its job!")
}
