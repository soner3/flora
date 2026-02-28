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

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"gen"},
	Short:   "Generates the type-safe Flora DI container",
	Long: `Scans the specified input directory for 'flora.Component' and 'flora.Configuration' tags.
It resolves the dependency graph, validates missing or duplicate providers, 
and uses Google Wire under the hood to generate a reflection-free, type-safe DI container.

The resulting 'flora_container.go' will be placed in your specified output directory.`,
	Example: `  # Scan current directory and generate container in the 'flora' folder (defaults)
  flora generate

  # Scan specific directory and output to the 'cmd/server' package
  flora generate -i ./internal -o ./cmd/server
  
  # Using the alias
  flora gen -i ./pkg/services`,
	SilenceUsage: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		log := slog.With("pkg", "cmd")
		log.Debug("Validating flags", "input", inputDir, "output", outputDir)

		info, err := os.Stat(inputDir)
		if err != nil {
			return errs.Wrap(err, "invalid directory provided for flag 'input': %s (directory does not exist)", inputDir)
		}
		if !info.IsDir() {
			return errs.Wrap(err, "invalid path provided for flag 'input': %s is a file, but must be a directory", inputDir)
		}

		log.Debug("Flags are valid")
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return app.RunGenerate(inputDir, outputDir)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringVarP(&inputDir, "input", "i", ".", "Input directory to scan")
	generateCmd.Flags().StringVarP(&outputDir, "output", "o", "flora", "Output directory for the generated container")
}
