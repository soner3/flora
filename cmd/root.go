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
package cmd

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/soner3/flora/internal/errs"
	"github.com/spf13/cobra"
)

var logLevel string

var (
	Version = "0.1.0"
	Build   = "unknown"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "flora",
	Version: fmt.Sprintf("%s, build %s", Version, Build),
	Short:   "Compile-time Dependency Injection for Go",
	Long: `Flora is a powerful, reflection-free Dependency Injection framework for Go.

By analyzing your codebase for '@Component' and '@Configuration' tags, 
Flora automatically resolves your dependency graph and uses Google Wire 
under the hood to generate safe, readable, and highly performant 
initialization code at compile time.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		setupLogger()
	},
	SilenceErrors: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		if floraErr, ok := err.(*errs.FloraError); ok {
			slog.Error(floraErr.Error(), "id", floraErr.ID)
			slog.Debug("Error Stacktrace", "trace", floraErr.StackTrace)
		} else {
			slog.Error(err.Error())
		}

		os.Exit(1)
	}
}

func init() {
	rootCmd.SetVersionTemplate("Flora version {{.Version}}\n")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "info", "Log level (debug, info, warn, error)")
}

func setupLogger() {
	var level slog.Level
	var writer io.Writer = os.Stdout
	invalidLevel := false

	switch strings.ToLower(logLevel) {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
		writer = os.Stderr
	default:
		invalidLevel = true
		level = slog.LevelInfo
	}

	logger := slog.New(slog.NewTextHandler(writer, &slog.HandlerOptions{
		Level: level,
	}))

	slog.SetDefault(logger)

	if invalidLevel {
		slog.Warn("Invalid log level provided. Defaulting to 'info'.", "provided", logLevel)
	}
}
