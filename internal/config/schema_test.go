package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/VladislavSCV/ekz/internal/theme"
)

func TestLoadTestdataFoodDelivery(t *testing.T) {
	path := filepath.Join("..", "..", "testdata", "food-delivery.yaml")
	p, err := LoadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if p.Main.Code != "FoodOrder" {
		t.Errorf("main.code: got %q", p.Main.Code)
	}
	if len(p.Main.Fields) != 4 {
		t.Fatalf("fields: got %d", len(p.Main.Fields))
	}
}

func TestValidateBuiltinPresets(t *testing.T) {
	for _, id := range []string{"food-delivery", "conferences"} {
		p, err := LoadBuiltinPreset(id)
		if err != nil {
			t.Fatalf("%s: %v", id, err)
		}
		p.ProjectName = "exam-" + id
		p.ProjectSlug = ModuleSlug(p.ProjectName)
		if err := p.Validate(); err != nil {
			t.Fatalf("%s validate: %v", id, err)
		}
	}
}

func TestFromConferencesThemeMatchesBuiltin(t *testing.T) {
	builtin, err := LoadBuiltinPreset("conferences")
	if err != nil {
		t.Fatal(err)
	}
	th, err := theme.LoadEmbedded("conferences")
	if err != nil {
		t.Fatal(err)
	}
	fromTheme := FromConferencesTheme("conf", th)
	fromTheme.ProjectName = builtin.ProjectName
	fromTheme.ProjectSlug = builtin.ProjectSlug
	if fromTheme.Main.Table != builtin.Main.Table {
		t.Errorf("table: theme %q builtin %q", fromTheme.Main.Table, builtin.Main.Table)
	}
	if len(fromTheme.Main.Statuses) != len(builtin.Main.Statuses) {
		t.Fatalf("statuses: %v vs %v", fromTheme.Main.Statuses, builtin.Main.Statuses)
	}
}

func TestModuleSlug(t *testing.T) {
	cases := map[string]string{
		"food-delivery": "fooddelivery",
		"My Project":    "myproject",
		"":              "project",
	}
	for in, want := range cases {
		if got := ModuleSlug(in); got != want {
			t.Errorf("ModuleSlug(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestLoadFileMissing(t *testing.T) {
	_, err := LoadFile(filepath.Join(t.TempDir(), "missing.yaml"))
	if err == nil {
		t.Fatal("expected error")
	}
	if !os.IsNotExist(err) {
		t.Fatalf("expected IsNotExist, got %v", err)
	}
}
