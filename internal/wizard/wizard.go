package wizard

import (
	"errors"
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"

	"github.com/VladislavSCV/ekz/internal/config"
)

// ErrSchemaOnly — мастер выгрузил YAML, генерация кода не нужна.
var ErrSchemaOnly = errors.New("project.yaml сохранён")

// Run — интерактивный мастер: шаблон, свой билет или project.yaml.
func Run() (config.ProjectSchema, error) {
	fmt.Println("\n─── Как собрать проект? ───")
	fmt.Println("На экзамене почти всегда: users + одна бизнес-таблица + опционально отзывы.")
	fmt.Println("Вы задаёте столбцы и страницы — ekz генерирует БД, API и фронт из project.yaml.")
	fmt.Println("Управление: ↑↓ и Enter (как в меню).\n")

	opts := []string{
		"Свой билет с нуля (любая тема: доставка, отель, курсы…)",
		"Шаблон: доставка еды",
		"Шаблон: Конференции.РФ",
		"Продолжить из project.yaml",
		"Только project.yaml — без генерации кода",
	}
	var choice string
	if err := survey.AskOne(&survey.Select{
		Message: "Сценарий:",
		Options: opts,
		Default: opts[0],
	}, &choice); err != nil {
		return config.ProjectSchema{}, err
	}

	switch choice {
	case opts[1]:
		return runFromBuiltin("food-delivery")
	case opts[2]:
		return runFromBuiltin("conferences")
	case opts[3]:
		return runFromYAMLFile()
	case opts[4]:
		return runSchemaOnly()
	default:
		return runCustom()
	}
}

func runSchemaOnly() (config.ProjectSchema, error) {
	sub := []string{
		"Пустой пример (ekz schema init)",
		"Шаблон: conferences",
		"Шаблон: food-delivery",
	}
	var pick string
	if err := survey.AskOne(&survey.Select{Message: "Что выгрузить в YAML?", Options: sub}, &pick); err != nil {
		return config.ProjectSchema{}, err
	}
	var path string
	_ = survey.AskOne(&survey.Input{Message: "Имя файла:", Default: "project.yaml"}, &path)
	var err error
	switch pick {
	case sub[0]:
		err = config.WriteSchemaTemplate(path)
	case sub[1]:
		err = config.ExportBuiltinPreset("conferences", path)
	default:
		err = config.ExportBuiltinPreset("food-delivery", path)
	}
	if err != nil {
		return config.ProjectSchema{}, err
	}
	return config.ProjectSchema{}, ErrSchemaOnly
}

func runFromYAMLFile() (config.ProjectSchema, error) {
	var path string
	if err := survey.AskOne(&survey.Input{
		Message: "Путь к project.yaml:",
		Default: "project.yaml",
	}, &path, survey.WithValidator(survey.Required)); err != nil {
		return config.ProjectSchema{}, err
	}
	p, err := config.LoadFile(path)
	if err != nil {
		return config.ProjectSchema{}, err
	}
	fmt.Printf("\nЗагружен: %s, таблица %s, %d полей\n", p.Portal.Name, p.Main.Table, len(p.Main.Fields))
	return finalizeProject(p, false)
}

func runFromBuiltin(id string) (config.ProjectSchema, error) {
	p, err := config.LoadBuiltinPreset(id)
	if err != nil {
		return config.ProjectSchema{}, err
	}
	meta, _ := findBuiltinMeta(id)
	if meta.ID != "" {
		fmt.Printf("\nШаблон «%s»: %s\n", meta.Title, meta.Description)
	}
	return finalizeProject(p, true)
}

func findBuiltinMeta(id string) (config.BuiltinPreset, bool) {
	for _, m := range config.ListBuiltinPresets() {
		if m.ID == id {
			return m, true
		}
	}
	return config.BuiltinPreset{}, false
}

// finalizeProject — имя папки, портал, поля/страницы (опционально), админ.
func finalizeProject(p config.ProjectSchema, fromTemplate bool) (config.ProjectSchema, error) {
	fmt.Println("\n─── Проект ───")
	defName := p.ProjectName
	if fromTemplate || defName == "" || defName == "food-delivery" || defName == "conferences" {
		defName = ""
	}
	if err := survey.AskOne(&survey.Input{
		Message: "Название папки проекта:",
		Default: defName,
	}, &p.ProjectName, survey.WithValidator(survey.Required)); err != nil {
		return p, err
	}
	p.ProjectSlug = config.ModuleSlug(p.ProjectName)

	fmt.Println("\n─── Портал ───")
	_ = survey.AskOne(&survey.Input{
		Message: "Название портала (шапка):",
		Default: p.Portal.Name,
	}, &p.Portal.Name, survey.WithValidator(survey.Required))
	_ = survey.AskOne(&survey.Input{
		Message: "Слоган:",
		Default: p.Portal.Tagline,
	}, &p.Portal.Tagline)

	var tuneFields bool
	if fromTemplate {
		_ = survey.AskOne(&survey.Confirm{
			Message: "Изменить столбцы БД из шаблона?",
			Default: false,
		}, &tuneFields)
	} else {
		tuneFields = true
	}
	if tuneFields {
		fields, err := AskFields(p.Main.Fields)
		if err != nil {
			return p, err
		}
		p.Main.Fields = fields
	}

	if !fromTemplate {
		return runEntityTail(p)
	}

	var tunePages bool
	_ = survey.AskOne(&survey.Confirm{
		Message: "Изменить набор страниц из шаблона?",
		Default: false,
	}, &tunePages)
	if tunePages {
		pages, err := AskPages(p.Pages)
		if err != nil {
			return p, err
		}
		p.Pages = pages
	}

	return askAdminAndBrand(p)
}

func runCustom() (config.ProjectSchema, error) {
	var p config.ProjectSchema

	fmt.Println("\n─── 1. Проект ───")
	if err := survey.AskOne(&survey.Input{Message: "Название папки проекта:"}, &p.ProjectName, survey.WithValidator(survey.Required)); err != nil {
		return p, err
	}
	p.ProjectSlug = config.ModuleSlug(p.ProjectName)

	fmt.Println("\n─── 2. Портал ───")
	if err := survey.AskOne(&survey.Input{Message: "Название портала (шапка сайта):", Default: "Мой портал"}, &p.Portal.Name, survey.WithValidator(survey.Required)); err != nil {
		return p, err
	}
	_ = survey.AskOne(&survey.Input{Message: "Краткий слоган:"}, &p.Portal.Tagline)

	fmt.Println("\n─── 3. Главная сущность (таблица в БД) ───")
	if err := survey.AskOne(&survey.Input{
		Message: "Название сущности (для человека):",
		Default: "Заявка",
		Help:    "Например: Заказ доставки, Заявка, Бронь",
	}, &p.Main.Label, survey.WithValidator(survey.Required)); err != nil {
		return p, err
	}
	if err := survey.AskOne(&survey.Input{
		Message: "Имя для кода (латиница, PascalCase):",
		Default: config.ToPascal(p.Main.Label),
		Help:    "Например: FoodOrder, Booking",
	}, &p.Main.Code, survey.WithValidator(survey.Required)); err != nil {
		return p, err
	}
	defaultTable := config.ToSnake(p.Main.Code)
	if !strings.HasSuffix(defaultTable, "s") {
		defaultTable += "s"
	}
	_ = survey.AskOne(&survey.Input{Message: "Имя таблицы в БД:", Default: defaultTable}, &p.Main.Table)

	fields, err := AskFields(nil)
	if err != nil {
		return p, err
	}
	p.Main.Fields = fields

	return runEntityTail(p)
}

func runEntityTail(p config.ProjectSchema) (config.ProjectSchema, error) {
	fmt.Println("\n─── Статусы (для админки и отзывов) ───")
	var useStatus bool
	if len(p.Main.Statuses) > 0 {
		def := true
		_ = survey.AskOne(&survey.Confirm{
			Message: fmt.Sprintf("Использовать статусы (%s)?", strings.Join(p.Main.Statuses, " → ")),
			Default: def,
		}, &useStatus)
		if !useStatus {
			p.Main.Statuses = nil
		}
	} else {
		if err := survey.AskOne(&survey.Confirm{
			Message: "Нужны статусы записи? (админ меняет, пользователь видит)",
			Default: true,
		}, &useStatus); err != nil {
			return p, err
		}
	}
	if useStatus && len(p.Main.Statuses) == 0 {
		var statErr error
		p.Main.Statuses, statErr = askStatuses()
		if statErr != nil {
			return p, statErr
		}
	} else if useStatus {
		var redo bool
		_ = survey.AskOne(&survey.Confirm{Message: "Перезадать статусы вручную?", Default: false}, &redo)
		if redo {
			var statErr error
			p.Main.Statuses, statErr = askStatuses()
			if statErr != nil {
				return p, statErr
			}
		}
	}

	fmt.Println("\n─── Отзывы ───")
	defReviews := len(p.Main.Statuses) > 0 && p.Reviews
	if err := survey.AskOne(&survey.Confirm{
		Message: "Таблица отзывов (после смены статуса админом)?",
		Default: defReviews,
	}, &p.Reviews); err != nil {
		return p, err
	}
	if p.Reviews && len(p.Main.Statuses) == 0 {
		return p, fmt.Errorf("отзывы возможны только со статусами")
	}

	fmt.Println("\n─── Страницы ───")
	pages, err := AskPages(p.Pages)
	if err != nil {
		return p, err
	}
	p.Pages = pages

	return askAdminAndBrand(p)
}

func askAdminAndBrand(p config.ProjectSchema) (config.ProjectSchema, error) {
	fmt.Println("\n─── Администратор ───")
	_ = survey.AskOne(&survey.Input{Message: "Логин админа:", Default: p.Admin.Login}, &p.Admin.Login)
	_ = survey.AskOne(&survey.Input{Message: "Пароль админа:", Default: p.Admin.Password}, &p.Admin.Password)

	if p.Brand.Primary == "" {
		p.Brand.Primary = "#1a5fb4"
		p.Brand.Dark = "#0d3d7a"
		p.Brand.Accent = "#ff6b35"
	}
	var tuneBrand bool
	_ = survey.AskOne(&survey.Confirm{
		Message: "Задать цвета бренда (CSS)?",
		Default: false,
	}, &tuneBrand)
	if tuneBrand {
		_ = survey.AskOne(&survey.Input{Message: "Основной цвет (#hex):", Default: p.Brand.Primary}, &p.Brand.Primary)
		_ = survey.AskOne(&survey.Input{Message: "Тёмный цвет:", Default: p.Brand.Dark}, &p.Brand.Dark)
		_ = survey.AskOne(&survey.Input{Message: "Акцент:", Default: p.Brand.Accent}, &p.Brand.Accent)
	}

	if err := p.Validate(); err != nil {
		return p, err
	}
	return p, nil
}
