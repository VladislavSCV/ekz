package config

import (
	"embed"
	"fmt"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:embed builtin/*.yaml
var builtinFS embed.FS

// BuiltinPreset — готовый project.yaml для типового билета.
type BuiltinPreset struct {
	ID          string
	Title       string
	Description string
}

var builtinCatalog = []BuiltinPreset{
	{ID: "food-delivery", Title: "Доставка еды", Description: "Заказы, адрес, блюдо, статусы доставки"},
	{ID: "conferences", Title: "Конференции.РФ", Description: "Бронь помещений, оплата, статусы мероприятия"},
}

// ListBuiltinPresets возвращает встроенные шаблоны (полный project.yaml).
func ListBuiltinPresets() []BuiltinPreset {
	out := make([]BuiltinPreset, len(builtinCatalog))
	copy(out, builtinCatalog)
	return out
}

// LoadBuiltinPreset загружает шаблон; project_name в YAML — заглушка.
func LoadBuiltinPreset(id string) (ProjectSchema, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return ProjectSchema{}, fmt.Errorf("id шаблона не указан")
	}
	data, err := builtinFS.ReadFile("builtin/" + id + ".yaml")
	if err != nil {
		return ProjectSchema{}, fmt.Errorf("шаблон %q не найден (доступны: %s)", id, strings.Join(builtinIDs(), ", "))
	}
	var p ProjectSchema
	if err := yaml.Unmarshal(data, &p); err != nil {
		return ProjectSchema{}, err
	}
	return p, nil
}

func builtinIDs() []string {
	ids := make([]string, len(builtinCatalog))
	for i, b := range builtinCatalog {
		ids[i] = b.ID
	}
	sort.Strings(ids)
	return ids
}
