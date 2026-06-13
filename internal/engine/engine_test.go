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
package engine

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolvePackageName(t *testing.T) {
	testcases := []struct {
		name     string
		setupDir func(t *testing.T) string
		want     string
	}{
		{
			name: "TestBaseName",
			setupDir: func(t *testing.T) string {
				return t.TempDir()
			},
			want: "",
		},
		{
			name: "TestReadsGoBuildPackageName",
			setupDir: func(t *testing.T) string {
				dir := t.TempDir()
				if err := os.WriteFile(filepath.Join(dir, "dummy.go"), []byte("package custompkg\n"), 0644); err != nil {
					t.Fatalf("failed to write dummy.go: %v", err)
				}
				return dir
			},
			want: "custompkg",
		},
		{
			name: "TestRootDirFallbackToMain",
			setupDir: func(t *testing.T) string {
				return "/"
			},
			want: "main",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			dir := tc.setupDir(t)
			got := ResolvePackageName(dir)

			want := tc.want
			if tc.name == "TestBaseName" {
				want = filepath.Base(dir)
			}

			if got != want {
				t.Errorf("ResolvePackageName(%q) = %q, want %q", dir, got, want)
			}
		})
	}
}
