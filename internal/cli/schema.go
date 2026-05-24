package cli

import (
	"fmt"

	"github.com/VladislavSCV/ekz/internal/config"
)

func schemaCmd(args []string) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "-h" {
		printSchemaHelp()
		return nil
	}
	switch args[0] {
	case "init":
		path := "project.yaml"
		if len(args) > 1 {
			path = args[1]
		}
		return config.WriteSchemaTemplate(path)
	case "export":
		if len(args) < 2 {
			return fmt.Errorf("укажите id шаблона: ekz schema export conferences [файл]")
		}
		out := ""
		if len(args) > 2 {
			out = args[2]
		}
		return config.ExportBuiltinPreset(args[1], out)
	default:
		return fmt.Errorf("неизвестно: ekz schema %s (ekz schema help)", args[0])
	}
}

func printSchemaHelp() {
	fmt.Print(`project.yaml — описание билета (поля БД, страницы, админ). Не код.

КОМАНДЫ
  ekz schema init [файл]           пример с полями (по умолчанию project.yaml)
  ekz schema export <id> [файл]    шаблон: conferences | food-delivery
  ekz presets                      список id

ПОЛЯ YAML
  project_name, portal.name, portal.tagline
  main.label, main.code, main.table
  main.fields[]     column, label, type (string|text|int|date|enum), options, required
  main.statuses[]   для админки
  pages             login, register, create_form, cabinet, admin, slider
  reviews           true | false
  admin.login, admin.password
  brand.primary, brand.dark, brand.accent

ПОСЛЕ ПРАВКИ
  ekz -name myproj -config project.yaml

ПОЛУЧИТЬ YAML БЕЗ ФАЙЛОВ
  ekz → «Только project.yaml — без генерации кода»

Справка: ekz help
`)
}
