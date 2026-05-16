package cmd

import (
	"bytes"
	"runtime/debug"
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

func TestSetupBuildInfo(t *testing.T) {
	// 1. Globale Variablen am Ende des Tests wiederherstellen (damit andere Tests nicht kaputt gehen)
	origReadFunc := readBuildInfoFunc
	origVersion := Version
	origBuild := Build
	defer func() {
		readBuildInfoFunc = origReadFunc
		Version = origVersion
		Build = origBuild
	}()

	testcases := []struct {
		name            string
		initialVersion  string
		initialBuild    string
		mockReadFunc    func() (*debug.BuildInfo, bool)
		expectedVersion string
		expectedBuild   string
	}{
		{
			name:           "GoReleaser Already Set",
			initialVersion: "v1.0.0",
			initialBuild:   "abcdef1",
			mockReadFunc: func() (*debug.BuildInfo, bool) {
				return nil, false
			},
			expectedVersion: "v1.0.0",
			expectedBuild:   "abcdef1",
		},
		{
			name:           "ReadBuildInfo Fails",
			initialVersion: "dev",
			initialBuild:   "unknown",
			mockReadFunc: func() (*debug.BuildInfo, bool) {
				return nil, false
			},
			expectedVersion: "dev",
			expectedBuild:   "unknown",
		},
		{
			name:           "Version is (devel)",
			initialVersion: "dev",
			initialBuild:   "unknown",
			mockReadFunc: func() (*debug.BuildInfo, bool) {
				return &debug.BuildInfo{
					Main: debug.Module{Version: "(devel)"},
				}, true
			},
			expectedVersion: "dev",
			expectedBuild:   "unknown",
		},
		{
			name:           "Happy Path - Long Revision & Clean",
			initialVersion: "dev",
			initialBuild:   "unknown",
			mockReadFunc: func() (*debug.BuildInfo, bool) {
				return &debug.BuildInfo{
					Main: debug.Module{Version: "v1.2.3"},
					Settings: []debug.BuildSetting{
						{Key: "vcs.revision", Value: "1234567890abcdef"},
					},
				}, true
			},
			expectedVersion: "v1.2.3",
			expectedBuild:   "1234567",
		},
		{
			name:           "Happy Path - Short Revision & Dirty",
			initialVersion: "dev",
			initialBuild:   "unknown",
			mockReadFunc: func() (*debug.BuildInfo, bool) {
				return &debug.BuildInfo{
					Main: debug.Module{Version: "v2.0.0"},
					Settings: []debug.BuildSetting{
						{Key: "vcs.revision", Value: "abc"},
						{Key: "vcs.modified", Value: "true"},
					},
				}, true
			},
			expectedVersion: "v2.0.0",
			expectedBuild:   "abc-dirty",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			Version = tc.initialVersion
			Build = tc.initialBuild
			readBuildInfoFunc = tc.mockReadFunc
			setupBuildInfo()

			if Version != tc.expectedVersion {
				t.Errorf("Expected Version %q, but got %q", tc.expectedVersion, Version)
			}
			if Build != tc.expectedBuild {
				t.Errorf("Expected Build %q, but got %q", tc.expectedBuild, Build)
			}
		})
	}
}
