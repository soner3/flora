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

	"github.com/soner3/weld/internal/errs"
	"github.com/spf13/cobra"
)

var dir string

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:          "generate",
	Aliases:      []string{"gen"},
	Short:        "Generate weld files",
	Long:         `Generate weld files from the given directory.`,
	SilenceUsage: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		slog.Debug("Validating flag 'dir'", "dir", dir)

		if _, err := os.Stat(dir); err != nil {
			return errs.Wrap(err, "invalid directory provided for flag 'dir': %s", dir)
		}

		slog.Debug("Flag 'dir' is valid", "dir", dir)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringVarP(&dir, "dir", "d", ".", "Directory to scan")
}
