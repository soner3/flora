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
package app

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/soner3/flora/internal/engine/wiregen"
)

func TestRunGenerate(t *testing.T) {
	testcases := []struct {
		name    string
		dir     string
		outDir  string
		wantErr bool
	}{
		{
			name:    "TestScanError",
			dir:     "./testdata/scan_err",
			outDir:  t.TempDir(),
			wantErr: true,
		},
		{
			name:    "TestParseError",
			dir:     "./testdata/parse_err",
			outDir:  t.TempDir(),
			wantErr: true,
		},
		{
			name:    "TestZeroComponents",
			dir:     "./testdata/empty",
			outDir:  t.TempDir(),
			wantErr: false,
		},
		{
			name:    "TestGenerateError",
			dir:     "./testdata/happy",
			outDir:  "invalid\x00path",
			wantErr: true,
		},
		{
			name:    "TestSuccess",
			dir:     "./testdata/happy",
			outDir:  "",
			wantErr: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			outDir := tc.outDir

			if outDir == "" {
				tmpDir, err := os.MkdirTemp(".", "flora_app_test_*")
				if err != nil {
					t.Fatal(err)
				}
				t.Cleanup(func() {
					if removeErr := os.RemoveAll(tmpDir); removeErr != nil {
						t.Logf("warning: failed to clean up temp dir %s: %v", tmpDir, removeErr)
					}
				})
				outDir = tmpDir
			}

			err := RunGenerate(tc.dir, outDir)

			if tc.wantErr {
				if err == nil {
					t.Errorf("expected an error, but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("did not expect an error, but got: %v", err)
				}
			}
		})
	}
}

func TestRunWatch_Success(t *testing.T) {

	watchDir := t.TempDir()
	outDir := t.TempDir()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)

	_ = os.WriteFile(filepath.Join(watchDir, "dummy_file.txt"), []byte(""), 0644)
	_ = os.Mkdir(filepath.Join(watchDir, ".hidden_folder"), 0755)

	go func() {
		errCh <- RunWatch(ctx, "./testdata/happy", outDir, watchDir)
	}()

	time.Sleep(100 * time.Millisecond)

	_ = os.WriteFile(filepath.Join(watchDir, "ignored.txt"), []byte("text"), 0644)
	_ = os.WriteFile(filepath.Join(watchDir, "flora_container.go"), []byte("package main"), 0644)

	dummyFile := filepath.Join(watchDir, "trigger.go")
	_ = os.WriteFile(dummyFile, []byte("package trigger\n"), 0644)

	time.Sleep(400 * time.Millisecond)

	cancel()

	select {
	case err := <-errCh:
		if err != nil {
			t.Errorf("RunWatch returned unexpected error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("RunWatch did not shut down in time. Goroutine or WaitGroup deadlock!")
	}
}

func TestRunWatch_InvalidDir(t *testing.T) {
	ctx := context.Background()
	err := RunWatch(ctx, "./testdata/happy", t.TempDir(), "invalid\x00path")
	if err == nil {
		t.Error("Expected an error for invalid watch directory, but got nil")
	}
}

func TestRunWatch_WalkError(t *testing.T) {
	ctx := context.Background()
	watchDir := t.TempDir()

	noPermDir := filepath.Join(watchDir, "noperm")
	os.Mkdir(noPermDir, 0755)
	os.Chmod(noPermDir, 0000)

	defer os.Chmod(noPermDir, 0755)

	err := RunWatch(ctx, ".", ".", watchDir)
	if err == nil {
		t.Error("Expected an error because directory is unreadable")
	}
}

func TestRunWatch_AbsError(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()

	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)

	os.Chdir(tmpDir)
	os.RemoveAll(tmpDir)

	err := RunWatch(ctx, ".", ".", ".")
	if err == nil {
		t.Error("Expected an error for filepath.Abs")
	}
}

func TestRunWatch_WatcherError(t *testing.T) {
	originalWatcherFunc := newWatcherFunc
	newWatcherFunc = func() (FileWatcher, error) {
		return nil, errors.New("mocked watcher init error")
	}
	defer func() { newWatcherFunc = originalWatcherFunc }()

	err := RunWatch(context.Background(), ".", ".", t.TempDir())
	if err == nil {
		t.Error("Expected an error from NewWatcher")
	}
}

func TestNewWatcherFunc_Error(t *testing.T) {
	originalNewWatcher := newWatcher

	defer func() { newWatcher = originalNewWatcher }()

	newWatcher = func() (*fsnotify.Watcher, error) {
		return nil, errors.New("simulated OS out of file descriptors")
	}
	watcher, err := newWatcherFunc()
	if err == nil {
		t.Error("Expected an error from newWatcherFunc because osNewWatcher failed, but got nil")
	}
	if watcher != nil {
		t.Error("Expected watcher to be nil when an error occurs")
	}
}

type mockWatcher struct {
	eventsChan chan fsnotify.Event
	errorsChan chan error
}

func (m *mockWatcher) Add(name string) error         { return nil }
func (m *mockWatcher) Close() error                  { return nil }
func (m *mockWatcher) Events() <-chan fsnotify.Event { return m.eventsChan }
func (m *mockWatcher) Errors() <-chan error          { return m.errorsChan }

func TestRunWatch_MockedChannels(t *testing.T) {

	testcases := []struct {
		name       string
		actionFunc func(events chan fsnotify.Event, errs chan error)
	}{
		{
			name: "TestEventsChannelClosed",
			actionFunc: func(events chan fsnotify.Event, errs chan error) {
				close(events)
			},
		},
		{
			name: "TestErrorsChannelClosed",
			actionFunc: func(events chan fsnotify.Event, errs chan error) {
				close(errs)
			},
		},
		{
			name: "TestWatcherErrorReceived",
			actionFunc: func(events chan fsnotify.Event, errs chan error) {
				err := errors.New("simulated watcher error")
				errs <- err
				close(errs)
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			eventsCh := make(chan fsnotify.Event)
			errorsCh := make(chan error)

			originalWatcherFunc := newWatcherFunc
			newWatcherFunc = func() (FileWatcher, error) {
				return &mockWatcher{
					eventsChan: eventsCh,
					errorsChan: errorsCh,
				}, nil
			}
			defer func() { newWatcherFunc = originalWatcherFunc }()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			errCh := make(chan error, 1)

			outDir := t.TempDir()
			go func() {
				errCh <- RunWatch(ctx, ".", outDir, t.TempDir())
			}()

			time.Sleep(50 * time.Millisecond)

			tc.actionFunc(eventsCh, errorsCh)

			select {
			case err := <-errCh:
				if err != nil {
					t.Errorf("expected RunWatch to exit with nil, got: %v", err)
				}
			case <-time.After(1 * time.Second):
				t.Fatalf("RunWatch deadlock! Did not exit after channel action.")
			}
		})
	}
}


func TestWithContainerStub(t *testing.T) {
	containerFile := wiregen.ContainerFileName

	testcases := []struct {
		name           string
		initialContent string
		operation      func(dir string) func() error
		wantErr        bool
		wantContent    string
	}{
		{
			name:           "TestSuccessNoRollback",
			initialContent: "original content",
			operation: func(dir string) func() error {
				return func() error {
					return os.WriteFile(filepath.Join(dir, containerFile), []byte("new generated content"), 0644)
				}
			},
			wantErr:     false,
			wantContent: "new generated content",
		},
		{
			name:           "TestErrorTriggersRollback",
			initialContent: "original backup",
			operation: func(dir string) func() error {
				return func() error {
					_ = os.WriteFile(filepath.Join(dir, containerFile), []byte("broken temporary code"), 0644)
					return errors.New("simulated crash")
				}
			},
			wantErr:     true,
			wantContent: "original backup",
		},
		{
			name:           "TestNoPriorFile",
			initialContent: "",
			operation: func(dir string) func() error {
				return func() error { return nil }
			},
			wantErr:     false,
			wantContent: "",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			containerPath := filepath.Join(dir, containerFile)

			if tc.initialContent != "" {
				if err := os.WriteFile(containerPath, []byte(tc.initialContent), 0644); err != nil {
					t.Fatalf("failed to write initial container file: %v", err)
				}
			}

			err := WithContainerStub(dir, tc.operation(dir))

			if tc.wantErr {
				if err == nil {
					t.Error("expected an error, but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("did not expect an error, but got: %v", err)
				}
			}

			if tc.wantContent != "" {
				got, readErr := os.ReadFile(containerPath)
				if readErr != nil {
					t.Fatalf("failed to read container file after operation: %v", readErr)
				}
				if string(got) != tc.wantContent {
					t.Errorf("file content mismatch:\n got:  %q\n want: %q", string(got), tc.wantContent)
				}
			}
		})
	}
}

