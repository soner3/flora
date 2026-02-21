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
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/soner3/weld/internal/errs"
	"github.com/spf13/cobra"
)

var logLevel string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "weld",
	Short: "Compile-time Dependency Injection for Go",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		setupLogger()
	},
	SilenceErrors: true,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		if weldErr, ok := err.(*errs.WeldError); ok {
			slog.Error(weldErr.Error(), "id", weldErr.ID)
			slog.Debug("Error Stacktrace", "trace", weldErr.StackTrace)
		} else {
			slog.Error(err.Error())
		}

		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.weld.yaml)")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "info", "Log level (debug, info, warn, error)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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
