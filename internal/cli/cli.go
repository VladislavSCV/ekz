package cli

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/VladislavSCV/ekz/internal/codegen"
	"github.com/VladislavSCV/ekz/internal/config"
	"github.com/VladislavSCV/ekz/internal/theme"
	"github.com/VladislavSCV/ekz/internal/wizard"
)

func Run(args []string) error {
	if len(args) > 0 {
		switch args[0] {
		case "themes":
			return listThemes()
		case "presets":
			return listPresets()
		case "theme":
			return themeCmd(args[1:])
		case "help", "-h", "--help":
			printHelp()
			return nil
		}
	}
	return runGenerate()
}

func printHelp() {
	fmt.Println(`ekz — универсальный генератор проектов ДЭ (09.02.07)

  ekz                              мастер: столбцы БД, страницы, статусы
  ekz presets                      список встроенных шаблонов (project.yaml)
  ekz -name proj -config project.yaml
  ekz -name proj -quick -preset food-delivery
  ekz -name proj -quick -theme conferences   (то же, устаревший флаг -theme)
  ekz themes | ekz theme init

Схема: вы описываете project.yaml → ekz генерирует SQLite/GORM, API и Vite-страницы.
Не «любые 20 таблиц», а надёжный каркас типового билета.`)
}

func listPresets() error {
	fmt.Println("Встроенные шаблоны (-quick -preset <id>):")
	for _, p := range config.ListBuiltinPresets() {
		fmt.Printf("  • %-16s %s — %s\n", p.ID, p.Title, p.Description)
	}
	fmt.Println("\nУниверсально: ekz  →  «Свой билет с нуля» или шаблон + правка полей")
	return nil
}

func listThemes() error {
	list, err := theme.ListEmbedded()
	if err != nil {
		return err
	}
	fmt.Println("Темы оформления для -quick -theme (устар., предпочтите -preset):")
	for _, t := range list {
		fmt.Printf("  • %-14s %s\n", t.ID, t.Name)
	}
	_ = listPresets()
	return nil
}

func themeCmd(args []string) error {
	if len(args) == 0 || args[0] == "help" {
		fmt.Println("  ekz theme init [файл]  — шаблон theme.yaml (устар., используйте мастер ekz)")
		return nil
	}
	if args[0] == "init" {
		dst := "theme.yaml"
		if len(args) > 1 {
			dst = args[1]
		}
		return theme.WriteTemplate(dst)
	}
	return fmt.Errorf("неизвестная команда: ekz theme %s", args[0])
}

func runGenerate() error {
	flagName := flag.String("name", "", "название папки проекта")
	flagConfig := flag.String("config", "", "project.yaml (без мастера)")
	flagQuick := flag.Bool("quick", false, "быстрый режим без мастера")
	flagPreset := flag.String("preset", "", "шаблон project.yaml для -quick (food-delivery, conferences)")
	flagTheme := flag.String("theme", "", "устар.: то же что -preset для conferences")
	flagThemeFile := flag.String("theme-file", "", "theme.yaml для -quick (legacy)")
	flag.Parse()

	fmt.Println("╔══════════════════════════════════════════════════╗")
	fmt.Println("║  ekz — генератор ДЭ (09.02.07)                  ║")
	fmt.Println("╚══════════════════════════════════════════════════╝")
	fmt.Println()

	var schema config.ProjectSchema
	var err error

	switch {
	case *flagConfig != "":
		schema, err = config.LoadFile(*flagConfig)
		if *flagName != "" {
			schema.ProjectName = *flagName
			schema.ProjectSlug = config.ModuleSlug(*flagName)
		}
	case *flagQuick:
		if *flagName == "" {
			return fmt.Errorf("укажите -name для быстрого режима")
		}
		presetID := strings.TrimSpace(*flagPreset)
		if presetID == "" {
			presetID = strings.TrimSpace(*flagTheme)
		}
		if presetID == "" {
			presetID = "conferences"
		}
		if *flagThemeFile != "" {
			var th theme.Theme
			th, err = theme.LoadFile(*flagThemeFile)
			if err != nil {
				return err
			}
			schema = config.FromConferencesTheme(*flagName, th)
		} else if _, err := config.LoadBuiltinPreset(presetID); err == nil {
			schema, err = config.LoadBuiltinPreset(presetID)
			if err != nil {
				return err
			}
			schema.ProjectName = *flagName
			schema.ProjectSlug = config.ModuleSlug(*flagName)
		} else if presetID == "conferences" {
			var th theme.Theme
			th, err = theme.LoadEmbedded("conferences")
			if err != nil {
				return err
			}
			schema = config.FromConferencesTheme(*flagName, th)
		} else {
			return fmt.Errorf("неизвестный шаблон %q (ekz presets)", presetID)
		}
	default:
		schema, err = wizard.Run()
	}
	if err != nil {
		return err
	}
	if err := schema.Validate(); err != nil {
		return err
	}

	cwd, _ := os.Getwd()
	target := filepath.Join(cwd, schema.ProjectName)
	if info, err := os.Stat(target); err == nil && info.IsDir() {
		if ents, _ := os.ReadDir(target); len(ents) > 0 {
			return fmt.Errorf("папка %s не пуста", target)
		}
	}

	fmt.Printf("\nПортал: %s\n", schema.Portal.Name)
	fmt.Printf("Сущность: %s → таблица %s (%d полей)\n", schema.Main.Label, schema.Main.Table, len(schema.Main.Fields))
	if len(schema.Main.Statuses) > 0 {
		fmt.Printf("Статусы: %s\n", strings.Join(schema.Main.Statuses, " → "))
	}
	fmt.Printf("Папка:   %s\n\n", target)

	if err := codegen.Generate(target, schema); err != nil {
		return err
	}

	fmt.Println("✓ Проект создан (project.yaml сохранён в корне)")
	fmt.Println("  backend:  cd", schema.ProjectName, "/backend && go mod tidy && go run .")
	fmt.Println("  frontend: cd", schema.ProjectName, "/frontend && npm install && npm run dev")
	return nil
}
