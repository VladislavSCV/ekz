package presets

import "embed"

type Preset struct {
	ID       string
	Name     string
	Label    string
	Template embed.FS
	Root     string // subdirectory inside embed FS
}

//go:embed conferences/*
var conferencesFS embed.FS

var registry = []Preset{
	{
		ID:       "conferences",
		Name:     "Конференции.РФ",
		Label:    "Конференции.РФ (Вариант №2)",
		Template: conferencesFS,
		Root:     "conferences",
	},
}

func List() []Preset {
	out := make([]Preset, len(registry))
	copy(out, registry)
	return out
}

func ByID(id string) (Preset, bool) {
	for _, p := range registry {
		if p.ID == id {
			return p, true
		}
	}
	return Preset{}, false
}
