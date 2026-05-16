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
	"log/slog"
	"os"

	"github.com/soner3/flora/internal/app"
	"github.com/soner3/flora/internal/errs"
	"github.com/spf13/cobra"
)

var inputDir string
var outputDir string
var watch bool
var watchDir string

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"gen"},
	Short:   "Generates the type-safe Flora DI container",
	Long: `The 'flora generate' command is the core engine of the Flora Dependency Injection framework.

It statically analyzes your Go source code (AST) looking for 'flora.Component' 
and 'flora.Configuration' markers to automatically discover your application's architecture.

Features:
  - Scans directories recursively to build a complete dependency graph.
  - Generates a clean, human-readable 'flora_container.go' file.
  - Supports an interactive '--watch' mode for rapid local development.`,
	Example: `  # 1. Default: Scan the current directory and generate the container in the './flora' folder
  flora generate

  # 2. Custom Paths: Scan the './internal' folder and output to a specific package
  flora generate --input ./internal --output ./cmd/server
  
  # 3. Watch Mode: Run in the background and auto-regenerate on any file change in the current directory
  flora gen --watch

  # 4. Advanced Watch: Watch a specific source directory while scanning another
  flora gen -i ./pkg/services -o ./internal/di --watch --watch-dir ./pkg`,
	SilenceUsage: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		log := slog.With("pkg", "cmd")
		log.Debug("Validating flags", "input", inputDir, "output", outputDir, "watch", watch, "watchDir", watchDir)

		if err := validateDir(inputDir, "input"); err != nil {
			return err
		}

		if err := validateDir(outputDir, "output"); err != nil {
			return err
		}

		if watch {
			log.Debug("Watch flag is enabled")
			if err := validateDir(watchDir, "watch-dir"); err != nil {
				return err
			}
		}

		log.Debug("Flags are valid")
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if watch {
			return app.RunWatch(cmd.Context(), inputDir, outputDir, watchDir)
		}

		return app.RunGenerate(inputDir, outputDir)
	},
}

func validateDir(dir string, flagName string) error {
	info, err := os.Stat(dir)
	if err != nil {
		return errs.Wrap(err, "invalid directory provided for flag '%s': %s (directory does not exist)", flagName, dir)
	}
	if !info.IsDir() {
		return errs.Wrap(err, "invalid path provided for flag '%s': %s is a file, but must be a directory", flagName, dir)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringVarP(&inputDir, "input", "i", ".", "Input directory to scan")
	generateCmd.Flags().StringVarP(&outputDir, "output", "o", "flora", "Output directory for the generated container")
	generateCmd.Flags().BoolVarP(&watch, "watch", "w", false, "Watch for file changes and regenerate the container automatically")
	generateCmd.Flags().StringVarP(&watchDir, "watch-dir", "d", ".", "Directory to watch for file changes")
}
