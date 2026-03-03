package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/soner3/flora/internal/errs"
	"github.com/spf13/cobra"
)

func TestExecute_Success(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	rootCmd.SetArgs([]string{"--help"})

	Execute()

	if !strings.Contains(buf.String(), "Usage:") {
		t.Errorf("expected output to contain 'Usage:', got %s", buf.String())
	}
}

func TestExecute_StandardError(t *testing.T) {
	var exitCode int
	originalOsExit := osExit
	osExit = func(code int) { exitCode = code }
	defer func() { osExit = originalOsExit }()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	rootCmd.SetArgs([]string{"--invalid-flag"})

	Execute()

	if exitCode != 1 {
		t.Errorf("expected osExit(1) to be called for standard error, got %d", exitCode)
	}
}

func TestExecute_FloraError(t *testing.T) {
	var exitCode int
	originalOsExit := osExit
	osExit = func(code int) { exitCode = code }
	defer func() { osExit = originalOsExit }()

	dummyCmd := &cobra.Command{
		Use: "trigger-flora-error",
		RunE: func(cmd *cobra.Command, args []string) error {
			return errs.Wrap(nil, "this is a test flora error")
		},
	}
	rootCmd.AddCommand(dummyCmd)

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"trigger-flora-error"})

	Execute()

	if exitCode != 1 {
		t.Errorf("expected osExit(1) to be called for FloraError, got %d", exitCode)
	}
}

func TestSetupLogger(t *testing.T) {
	levels := []string{"debug", "info", "warn", "error", "invalid_level"}

	for _, level := range levels {
		t.Run("Level_"+level, func(t *testing.T) {
			logLevel = level
			setupLogger(rootCmd)
		})
	}
}
