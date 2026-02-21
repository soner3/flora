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
	"log/slog"

	"github.com/soner3/weld/internal/engine/wiregen"
	"github.com/soner3/weld/internal/scanner"
)

func RunGenerate(dir string) error {
	log := slog.With("pkg", "app")

	log.Info("Starting weld generation...", "dir", dir)

	log.Debug("Scanning packages for weld components...")
	pkgs, err := scanner.ScanPackages(dir)
	if err != nil {
		return err
	}

	components, err := scanner.ParseComponents(pkgs)
	if err != nil {
		return err
	}

	if len(components) == 0 {
		log.Warn("No weld components found. Nothing to generate.")
		return nil
	}

	log.Info("Scan complete", "components_found", len(components))

	log.Debug("Generating DI container...")
	gen := wiregen.New()
	if err := gen.Generate(dir, components); err != nil {
		return err
	}

	log.Info("Successfully generated weld container!")
	return nil
}
