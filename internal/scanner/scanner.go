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
package scanner

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/soner3/mint/internal/errs"
	"golang.org/x/tools/go/packages"
)

var (
	ErrLoadPackages = errors.New("failed to load packages")
	ErrCompile      = errors.New("compile error in package")
)

func ScanPackages(rootDir string) ([]*packages.Package, error) {
	log := slog.With("pkg", "scanner")

	log.Debug("Scanning packages", "rootDir", rootDir)

	cfg := &packages.Config{
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedSyntax |
			packages.NeedTypes |
			packages.NeedTypesInfo,
		Dir: rootDir,
	}

	log.Debug("Loading packages via packages.Load...")
	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		chainErr := fmt.Errorf("%w: %w", ErrLoadPackages, err)
		return nil, errs.Wrap(chainErr, "directory: %s", rootDir)
	}

	var validPkgs []*packages.Package
	for _, pkg := range pkgs {
		if len(pkg.Errors) > 0 {
			chainErr := fmt.Errorf("%w: %w", ErrCompile, pkg.Errors[0])
			return nil, errs.Wrap(chainErr, "package ID: %s", pkg.ID)
		}

		validPkgs = append(validPkgs, pkg)
	}

	log.Debug("Successfully filtered packages", "total_loaded", len(pkgs), "valid_count", len(validPkgs))

	return validPkgs, nil
}
