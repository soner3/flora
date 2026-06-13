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
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/soner3/flora/internal/engine/wiregen"
	"github.com/soner3/flora/internal/errs"
	"github.com/soner3/flora/internal/scanner"
)

type FileWatcher interface {
	Add(name string) error
	Close() error
	Events() <-chan fsnotify.Event
	Errors() <-chan error
}

type realWatcher struct {
	w *fsnotify.Watcher
}

func (rw *realWatcher) Add(name string) error         { return rw.w.Add(name) }
func (rw *realWatcher) Close() error                  { return rw.w.Close() }
func (rw *realWatcher) Events() <-chan fsnotify.Event { return rw.w.Events }
func (rw *realWatcher) Errors() <-chan error          { return rw.w.Errors }

var newWatcher = fsnotify.NewWatcher

var newWatcherFunc = func() (FileWatcher, error) {
	w, err := newWatcher()
	if err != nil {
		return nil, err
	}
	return &realWatcher{w: w}, nil
}

func RunGenerate(inputDir, outputDir string) error {
	log := slog.With("pkg", "app")

	log.Info("Starting flora generation...", "dir", inputDir, "out", outputDir)

	return WithContainerStub(outputDir, func() error {
		log.Debug("Scanning packages for flora components...")
		pkgs, err := scanner.ScanPackages(inputDir)
		if err != nil {
			return err
		}

		genCtx, err := scanner.ParsePackages(pkgs)
		if err != nil {
			return err
		}

		if len(genCtx.Components) == 0 && len(genCtx.SliceBindings) == 0 {
			log.Warn("No flora components found. Nothing to generate.")
			return nil
		}

		log.Info("Scan complete", "components_found", len(genCtx.Components), "slice_bindings_found", len(genCtx.SliceBindings))

		log.Debug("Generating DI container...")
		gen := wiregen.NewWireGenerator()
		if err := gen.Generate(outputDir, genCtx); err != nil {
			return err
		}

		log.Info("Successfully generated flora container!")
		return nil
	})
}

// WithContainerStub backs up flora_container.go, overwrites it with a minimal stub
// so the AST scan does not fail on stale types, runs operation, and restores the
// original file if operation returns an error.
func WithContainerStub(outputDir string, operation func() error) (err error) {
	absOutDir, _ := filepath.Abs(outputDir)

	containerPath := filepath.Join(absOutDir, wiregen.ContainerFileName)

	pkgName := strings.ReplaceAll(filepath.Base(absOutDir), "-", "_")

	var backupContent []byte
	backupContent, _ = os.ReadFile(containerPath)

	defer func() {
		if err != nil && len(backupContent) > 0 {
			_ = os.WriteFile(containerPath, backupContent, 0644)
		}
	}()

	stub := fmt.Sprintf("//go:build !wireinject\n// +build !wireinject\n\npackage %s\ntype FloraContainer struct{}\nfunc InitializeContainer() (*FloraContainer, func(), error) { return nil, nil, nil }\n", pkgName)
	_ = os.WriteFile(containerPath, []byte(stub), 0644)

	err = operation()
	return err
}

// RunWatch starts a file watcher in the specified directory and triggers RunGenerate on changes.
func RunWatch(ctx context.Context, inputDir, outputDir, watchDir string) error {
	log := slog.With("pkg", "app")

	absWatchDir, err := filepath.Abs(watchDir)
	if err != nil {
		return errs.Wrap(err, "failed to resolve absolute path for watch directory")
	}

	watcher, err := newWatcherFunc()
	if err != nil {
		return errs.Wrap(err, "failed to initialize file watcher")
	}
	defer watcher.Close()

	err = filepath.Walk(absWatchDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if strings.HasPrefix(info.Name(), ".") && info.Name() != "." {
				return filepath.SkipDir
			}
			return watcher.Add(path)
		}
		return nil
	})

	if err != nil {
		return errs.Wrap(err, "failed to scan directory for watching: %s", absWatchDir)
	}

	log.Info("Flora Watch Mode started", "watch_dir", absWatchDir)

	if err := RunGenerate(inputDir, outputDir); err != nil {
		log.Error("Initial generation failed", "error", err)
	}

	var timer *time.Timer
	debounceDuration := 300 * time.Millisecond

	var genWg sync.WaitGroup

	for {
		select {
		case <-ctx.Done():
			if timer != nil {
				timer.Stop()
			}
			genWg.Wait()
			return nil

		case event, ok := <-watcher.Events():
			log.Debug("File event received", "event", event, "ok", ok)
			if !ok {
				log.Debug("Events channel closed, exiting")
				return nil
			}

			baseName := filepath.Base(event.Name)

			if event.Op == fsnotify.Chmod ||
				!strings.HasSuffix(event.Name, ".go") ||
				baseName == "flora_container.go" ||
				baseName == "flora_injector.go" ||
				baseName == "wire_gen.go" ||
				baseName == "wire.go" {
				log.Debug("File event not relevant, skipping", "event", event)
				continue
			}

			if timer != nil {
				log.Debug("Debounce timer already active, stopping it", "event", event)
				timer.Stop()
			}

			timer = time.AfterFunc(debounceDuration, func() {
				log.Debug("Debounce timer expired, starting generation")
				genWg.Add(1)
				defer genWg.Done()

				log.Info("File change detected. Regenerating...", "file", filepath.Base(event.Name))
				if err := RunGenerate(inputDir, outputDir); err != nil {
					log.Error("Generation failed", "error", err)
				}
			})

		case err, ok := <-watcher.Errors():
			if !ok {
				log.Debug("Error channel closed, exiting")
				return nil
			}
			log.Error("Watcher encountered an error", "error", err)
		}
	}
}
