package presets

import "embed"

// Preset — технический каркас (Go + Vite). Предметная область задаётся theme YAML.
type Preset struct {
	ID       string
	Name     string
	Label    string
	Template embed.FS
	Root     string
}

//go:embed deweb/*
var dewebFS embed.FS

var registry = []Preset{
	{
		ID:       "de",
		Name:     "Универсальный шаблон ДЭ",
		Label:    "Веб-приложение 09.02.07 (настраивается темой YAML)",
		Template: dewebFS,
		Root:     "deweb",
	},
}

func List() []Preset {
	out := make([]Preset, len(registry))
	copy(out, registry)
	return out
}

func ByID(id string) (Preset, bool) {
	// Обратная совместимость: старый -preset conferences
	if id == "conferences" {
		id = "de"
	}
	for _, p := range registry {
		if p.ID == id {
			return p, true
		}
	}
	return Preset{}, false
}
