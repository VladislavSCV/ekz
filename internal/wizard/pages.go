package wizard

import (
	"github.com/AlecAivazis/survey/v2"

	"github.com/VladislavSCV/ekz/internal/config"
)

var pageOpts = []string{
	"Вход (login)",
	"Регистрация",
	"Форма создания записи",
	"Личный кабинет (список)",
	"Панель администратора",
	"Слайдер в кабинете (модуль 2 ДЭ)",
}

// AskPages — выбор генерируемых HTML-страниц.
func AskPages(defaults config.PagesSchema) (config.PagesSchema, error) {
	var def []string
	if defaults.Login {
		def = append(def, pageOpts[0])
	}
	if defaults.Register {
		def = append(def, pageOpts[1])
	}
	if defaults.CreateForm {
		def = append(def, pageOpts[2])
	}
	if defaults.Cabinet {
		def = append(def, pageOpts[3])
	}
	if defaults.Admin {
		def = append(def, pageOpts[4])
	}
	if defaults.Slider {
		def = append(def, pageOpts[5])
	}
	if len(def) == 0 {
		def = pageOpts
	}

	var picked []string
	if err := survey.AskOne(&survey.MultiSelect{
		Message: "Какие страницы генерировать?",
		Options: pageOpts,
		Default: def,
	}, &picked); err != nil {
		return config.PagesSchema{}, err
	}

	var p config.PagesSchema
	for _, s := range picked {
		switch s {
		case pageOpts[0]:
			p.Login = true
		case pageOpts[1]:
			p.Register = true
		case pageOpts[2]:
			p.CreateForm = true
		case pageOpts[3]:
			p.Cabinet = true
		case pageOpts[4]:
			p.Admin = true
		case pageOpts[5]:
			p.Slider = true
		}
	}
	if !p.Login && !p.Register {
		p.Login = true
		p.Register = true
	}
	return p, nil
}
