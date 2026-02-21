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
	"log/slog"

	"github.com/soner3/weld/internal/errs"
	"golang.org/x/tools/go/packages"
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
		return nil, errs.Wrap(err, "failed to load packages in dir %s", rootDir)
	}

	var validPkgs []*packages.Package
	for _, pkg := range pkgs {
		if len(pkg.Errors) > 0 {
			return nil, errs.Wrap(pkg.Errors[0], "compile error in package %s", pkg.ID)
		}

		if pkg.Name != "main" {
			validPkgs = append(validPkgs, pkg)
		}
	}

	log.Debug("Successfully filtered packages", "total_loaded", len(pkgs), "valid_count", len(validPkgs))

	return validPkgs, nil
}
