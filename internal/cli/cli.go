package cli

import (
	"errors"
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
		case "schema":
			return schemaCmd(args[1:])
		case "theme":
			return themeCmd(args[1:])
		case "help", "-h", "--help":
			printHelp()
			return nil
		}
	}
	return runGenerate()
}

func listPresets() error {
	fmt.Println("Встроенные шаблоны (id для -preset и schema export):")
	for _, p := range config.ListBuiltinPresets() {
		fmt.Printf("  • %-16s %s — %s\n", p.ID, p.Title, p.Description)
	}
	fmt.Println("\n  ekz -name proj -preset conferences")
	fmt.Println("  ekz schema export conferences")
	fmt.Println("  ekz help")
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
	flagPreset := flag.String("preset", "", "шаблон: food-delivery, conferences (с -name = без мастера)")
	flagTheme := flag.String("theme", "", "устар.: то же что -preset для conferences")
	flagThemeFile := flag.String("theme-file", "", "theme.yaml для -quick (legacy)")
	flag.Parse()

	fmt.Println("╔══════════════════════════════════════════════════╗")
	fmt.Println("║  ekz — генератор ДЭ (09.02.07)                  ║")
	fmt.Println("╚══════════════════════════════════════════════════╝")
	fmt.Println()

	var schema config.ProjectSchema
	var err error

	quick := *flagQuick || (strings.TrimSpace(*flagName) != "" && strings.TrimSpace(*flagPreset) != "")

	switch {
	case *flagConfig != "":
		schema, err = config.LoadFile(*flagConfig)
		if *flagName != "" {
			schema.ProjectName = *flagName
			schema.ProjectSlug = config.ModuleSlug(*flagName)
		}
	case quick:
		if strings.TrimSpace(*flagName) == "" {
			return fmt.Errorf("укажите -name (например: -name food-delivery -preset food-delivery)")
		}
		schema, err = loadQuickSchema(*flagName, *flagPreset, *flagTheme, *flagThemeFile)
	default:
		schema, err = wizard.Run()
	}
	if errors.Is(err, wizard.ErrSchemaOnly) {
		return nil
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

func loadQuickSchema(projectName, preset, themeID, themeFile string) (config.ProjectSchema, error) {
	presetID := strings.TrimSpace(preset)
	if presetID == "" {
		presetID = strings.TrimSpace(themeID)
	}
	if presetID == "" {
		presetID = "conferences"
	}
	if strings.TrimSpace(themeFile) != "" {
		th, err := theme.LoadFile(themeFile)
		if err != nil {
			return config.ProjectSchema{}, err
		}
		return config.FromConferencesTheme(projectName, th), nil
	}
	if _, err := config.LoadBuiltinPreset(presetID); err == nil {
		schema, err := config.LoadBuiltinPreset(presetID)
		if err != nil {
			return config.ProjectSchema{}, err
		}
		schema.ProjectName = projectName
		schema.ProjectSlug = config.ModuleSlug(projectName)
		return schema, nil
	}
	if presetID == "conferences" {
		th, err := theme.LoadEmbedded("conferences")
		if err != nil {
			return config.ProjectSchema{}, err
		}
		return config.FromConferencesTheme(projectName, th), nil
	}
	return config.ProjectSchema{}, fmt.Errorf("неизвестный шаблон %q (ekz presets)", presetID)
}
