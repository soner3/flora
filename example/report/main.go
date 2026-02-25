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
package report

import (
	"fmt"

	"github.com/soner3/flora"
	"github.com/soner3/flora/example/config"
)

type DocumentGenerator interface {
	Generate(title string) string
}

type PdfGenerator struct {
	flora.Component `flora:"scope=prototype,primary"`
	cfg             config.Config
}

func NewPdfGenerator(cfg config.Config) (*PdfGenerator, func(), error) {
	fmt.Println("   -> [PdfGenerator] New instance created! (Connecting to port", cfg.Port, ")")

	if cfg.Port == 0 {
		return nil, nil, fmt.Errorf("invalid configuration")
	}

	cleanup := func() {
		fmt.Println("   -> [PdfGenerator] Cleanup: Temporary PDF files deleted.")
	}

	return &PdfGenerator{
		cfg: cfg,
	}, cleanup, nil
}

func (p *PdfGenerator) Generate(title string) string {
	return fmt.Sprintf("📄 PDF document '%s' generated successfully.", title)
}

type NonPrimaryDocumentGenerator struct {
	flora.Component `flora:"scope=prototype"`
	cfg             config.Config
}

func NewNonPrimaryDocumentGenerator(cfg config.Config) (*NonPrimaryDocumentGenerator, func(), error) {
	fmt.Println("   -> [NonPrimaryDocumentGenerator] New instance created! (Connecting to port", cfg.Port, ")")

	if cfg.Port == 0 {
		return nil, nil, fmt.Errorf("invalid configuration")
	}

	cleanup := func() {
		fmt.Println("   -> [NonPrimaryDocumentGenerator] Cleanup: Temporary PDF files deleted.")
	}

	return &NonPrimaryDocumentGenerator{
		cfg: cfg,
	}, cleanup, nil
}

func (p *NonPrimaryDocumentGenerator) Generate(title string) string {
	return fmt.Sprintf("📄 Non-primary document '%s' generated successfully.", title)
}
