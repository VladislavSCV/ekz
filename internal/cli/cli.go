package cli

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"

	"github.com/VladislavSCV/ekz/internal/generator"
	"github.com/VladislavSCV/ekz/internal/presets"
)

func Run() error {
	flagName := flag.String("name", "", "название проекта (папка)")
	flagPreset := flag.String("preset", "", "id пресета (conferences)")
	flag.Parse()
	fmt.Println("╔══════════════════════════════════════════════════╗")
	fmt.Println("║  Генератор решений ДЭ (09.02.07)                 ║")
	fmt.Println("║  Разработчик веб и мультимедийных приложений     ║")
	fmt.Println("╚══════════════════════════════════════════════════╝")
	fmt.Println()

	var preset presets.Preset
	var projectName string

	if *flagName != "" && *flagPreset != "" {
		projectName = strings.TrimSpace(*flagName)
		p, ok := presets.ByID(*flagPreset)
		if !ok {
			return fmt.Errorf("неизвестный пресет: %s (доступно: conferences)", *flagPreset)
		}
		preset = p
	} else {
		if err := survey.AskOne(&survey.Input{
			Message: "Название проекта (папка):",
			Help:    "Будет создана в текущей директории",
		}, &projectName, survey.WithValidator(survey.Required)); err != nil {
			return err
		}
		projectName = strings.TrimSpace(projectName)
		if projectName == "" {
			return fmt.Errorf("название проекта не может быть пустым")
		}

		list := presets.List()
		options := make([]string, len(list))
		byLabel := make(map[string]presets.Preset, len(list))
		for i, p := range list {
			options[i] = p.Label
			byLabel[p.Label] = p
		}

		var chosen string
		if err := survey.AskOne(&survey.Select{
			Message: "Выберите пресет (тему билета):",
			Options: options,
		}, &chosen); err != nil {
			return err
		}

		var ok bool
		preset, ok = byLabel[chosen]
		if !ok {
			return fmt.Errorf("неизвестный пресет: %s", chosen)
		}
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	target := filepath.Join(cwd, projectName)

	info, err := os.Stat(target)
	if err == nil {
		if !info.IsDir() {
			return fmt.Errorf("%s существует и не является папкой", target)
		}
		entries, _ := os.ReadDir(target)
		if len(entries) > 0 {
			return fmt.Errorf("папка %s не пуста — выберите другое имя или очистите её", target)
		}
	}

	fmt.Printf("\nГенерация «%s» → %s\n\n", preset.Name, target)

	if err := generator.Generate(preset, target, projectName); err != nil {
		return err
	}

	fmt.Println("✓ Проект успешно создан!")
	fmt.Println()
	fmt.Println("Дальнейшие шаги:")
	fmt.Println("  1. cd", projectName)
	fmt.Println("  2. cd backend && go run .")
	fmt.Println("  3. cd frontend && npm install && npm run dev")
	fmt.Println()
	fmt.Println("Подробности — в README.md проекта.")
	return nil
}
