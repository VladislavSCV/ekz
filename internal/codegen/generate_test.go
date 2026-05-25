package codegen

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/VladislavSCV/ekz/internal/config"
)

func TestGenerateFrontendArtifacts(t *testing.T) {
	schema, err := config.LoadBuiltinPreset("food-delivery")
	if err != nil {
		t.Fatal(err)
	}
	schema.ProjectName = "fe-check"
	schema.ProjectSlug = config.ModuleSlug(schema.ProjectName)

	dir := t.TempDir()
	if err := Generate(dir, schema); err != nil {
		t.Fatal(err)
	}
	for _, rel := range []string{
		"frontend/package.json",
		"frontend/vite.config.js",
		"frontend/src/theme-config.js",
		"frontend/index.html",
	} {
		if _, err := os.Stat(filepath.Join(dir, rel)); err != nil {
			t.Errorf("missing %s: %v", rel, err)
		}
	}
	themeJS, err := os.ReadFile(filepath.Join(dir, "frontend/src/theme-config.js"))
	if err != nil {
		t.Fatal(err)
	}
	body := string(themeJS)
	for _, want := range []string{`export const PORTAL_NAME`, `export const STATUS_NEW`, "Новый"} {
		if !strings.Contains(body, want) {
			t.Fatalf("theme-config.js missing %q:\n%s", want, themeJS)
		}
	}
}
