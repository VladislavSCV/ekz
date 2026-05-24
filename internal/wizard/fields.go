package wizard

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"

	"github.com/VladislavSCV/ekz/internal/config"
)

// AskFields — интерактивное описание столбцов главной таблицы.
func AskFields(existing []config.FieldSchema) ([]config.FieldSchema, error) {
	fields := append([]config.FieldSchema(nil), existing...)
	if len(fields) > 0 {
		fmt.Println("\nТекущие поля:")
		for i, f := range fields {
			fmt.Printf("  %d) %s — %s (%s)\n", i+1, f.Column, f.Label, f.Type)
		}
		var action string
		opts := []string{
			"Оставить как есть",
			"Добавить поля к списку",
			"Заменить все (ввести заново)",
		}
		if err := survey.AskOne(&survey.Select{
			Message: "Поля таблицы:",
			Options: opts,
			Default: opts[0],
		}, &action); err != nil {
			return nil, err
		}
		switch action {
		case opts[1]:
			// fall through to add loop
		case opts[2]:
			fields = nil
		default:
			return fields, nil
		}
	}

	fmt.Println("\n─── Поля таблицы (столбцы) ───")
	fmt.Println("Системные id, user_id, created_at и таблица users добавятся автоматически.")
	for {
		f, stop, err := askField()
		if err != nil {
			return nil, err
		}
		if stop && len(fields) == 0 {
			return nil, fmt.Errorf("нужно хотя бы одно поле")
		}
		if !stop {
			fields = append(fields, f)
		}
		if stop {
			break
		}
		var again bool
		_ = survey.AskOne(&survey.Confirm{Message: "Добавить ещё поле?", Default: true}, &again)
		if !again {
			break
		}
	}
	return fields, nil
}

func askField() (config.FieldSchema, bool, error) {
	var f config.FieldSchema
	if err := survey.AskOne(&survey.Input{
		Message: "Имя столбца в БД (snake_case, Enter без ввода = закончить):",
	}, &f.Column); err != nil {
		return f, false, err
	}
	f.Column = strings.TrimSpace(f.Column)
	if f.Column == "" {
		return f, true, nil
	}
	if err := survey.AskOne(&survey.Input{Message: "Подпись в форме:", Default: f.Column}, &f.Label); err != nil {
		return f, false, err
	}
	typeOpt := []string{
		"string — строка",
		"text — длинный текст",
		"int — число",
		"date — дата ДД.ММ.ГГГГ",
		"enum — выбор из списка",
	}
	var typeChoice string
	if err := survey.AskOne(&survey.Select{Message: "Тип данных:", Options: typeOpt}, &typeChoice); err != nil {
		return f, false, err
	}
	switch {
	case strings.HasPrefix(typeChoice, "string"):
		f.Type = "string"
	case strings.HasPrefix(typeChoice, "text"):
		f.Type = "text"
	case strings.HasPrefix(typeChoice, "int"):
		f.Type = "int"
	case strings.HasPrefix(typeChoice, "date"):
		f.Type = "date"
	default:
		f.Type = "enum"
	}
	if f.Type == "enum" {
		for {
			var val, label string
			if err := survey.AskOne(&survey.Input{Message: "Значение option (пусто = стоп):"}, &val); err != nil {
				return f, false, err
			}
			if strings.TrimSpace(val) == "" {
				break
			}
			_ = survey.AskOne(&survey.Input{Message: "Подпись в списке:", Default: val}, &label)
			f.Options = append(f.Options, config.Option{Value: val, Label: label})
		}
		if len(f.Options) == 0 {
			return f, false, fmt.Errorf("для enum добавьте хотя бы один вариант")
		}
	}
	_ = survey.AskOne(&survey.Confirm{Message: "Обязательное поле?", Default: true}, &f.Required)
	return f, false, nil
}

func askStatuses() ([]string, error) {
	fmt.Println("Введите статусы по порядку. Первый — при создании записи.")
	var statuses []string
	for i := 1; ; i++ {
		var s string
		msg := fmt.Sprintf("Статус #%d (пусто = закончить):", i)
		if i == 1 {
			msg = "Статус при создании (например: Новая):"
		}
		def := ""
		if i == 1 {
			def = "Новая"
		}
		if err := survey.AskOne(&survey.Input{Message: msg, Default: def}, &s); err != nil {
			return nil, err
		}
		s = strings.TrimSpace(s)
		if s == "" {
			if len(statuses) == 0 {
				return nil, fmt.Errorf("нужен минимум один статус")
			}
			break
		}
		statuses = append(statuses, s)
		if i >= 5 {
			break
		}
		var more bool
		_ = survey.AskOne(&survey.Confirm{Message: "Ещё статус?", Default: len(statuses) < 3}, &more)
		if !more {
			break
		}
	}
	return statuses, nil
}
