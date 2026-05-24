package config

import (
	"fmt"
	"strings"
	"unicode"
)

// ProjectSchema — единый источник правды для генерации (сохраняется как project.yaml).
type ProjectSchema struct {
	ProjectName string `yaml:"project_name"`
	ProjectSlug string `yaml:"project_slug"`

	Portal struct {
		Name    string `yaml:"name"`
		Tagline string `yaml:"tagline"`
	} `yaml:"portal"`

	Admin struct {
		Login    string `yaml:"login"`
		Password string `yaml:"password"`
	} `yaml:"admin"`

	Brand struct {
		Primary string `yaml:"primary"`
		Dark    string `yaml:"dark"`
		Accent  string `yaml:"accent"`
	} `yaml:"brand"`

	// Main — основная бизнес-сущность билета (заявка, заказ доставки, бронь…).
	Main EntitySchema `yaml:"main"`

	Reviews bool       `yaml:"reviews"`
	Pages   PagesSchema `yaml:"pages"`
}

type EntitySchema struct {
	Label    string        `yaml:"label"`     // «Заказ доставки»
	Code     string        `yaml:"code"`      // DeliveryOrder (Go)
	Table    string        `yaml:"table"`     // delivery_orders
	Fields   []FieldSchema `yaml:"fields"`
	Statuses []string      `yaml:"statuses"` // пусто = без статусов
}

type FieldSchema struct {
	Column   string   `yaml:"column"`          // snake_case в БД
	GoName   string   `yaml:"go_name,omitempty"` // PascalCase для Go
	Label    string   `yaml:"label"`           // подпись в форме
	Type     string   `yaml:"type"`            // string, text, int, date, enum
	Widget   string   `yaml:"widget"`          // text, textarea, number, date, select
	Options  []Option `yaml:"options"`         // для enum/select
	Required bool     `yaml:"required"`
}

type Option struct {
	Value string `yaml:"value"`
	Label string `yaml:"label"`
}

type PagesSchema struct {
	Login      bool `yaml:"login"`
	Register   bool `yaml:"register"`
	CreateForm bool `yaml:"create_form"`
	Cabinet    bool `yaml:"cabinet"`
	Admin      bool `yaml:"admin"`
	Slider     bool `yaml:"slider"`
}

func (p *ProjectSchema) Validate() error {
	if strings.TrimSpace(p.ProjectName) == "" {
		return fmt.Errorf("project_name обязателен")
	}
	if strings.TrimSpace(p.Portal.Name) == "" {
		return fmt.Errorf("portal.name обязателен")
	}
	if strings.TrimSpace(p.Admin.Login) == "" || strings.TrimSpace(p.Admin.Password) == "" {
		return fmt.Errorf("admin.login и admin.password обязательны")
	}
	if strings.TrimSpace(p.Main.Code) == "" {
		return fmt.Errorf("main.code обязателен (имя сущности для Go)")
	}
	if len(p.Main.Fields) == 0 {
		return fmt.Errorf("добавьте хотя бы одно поле в main.fields")
	}
	if p.Pages.Admin && len(p.Main.Statuses) == 0 {
		return fmt.Errorf("панель администратора требует статусы (main.statuses)")
	}
	if p.Reviews && len(p.Main.Statuses) == 0 {
		return fmt.Errorf("отзывы требуют статусы заявок")
	}
	for i, f := range p.Main.Fields {
		if err := f.validate(); err != nil {
			return fmt.Errorf("поле #%d: %w", i+1, err)
		}
	}
	if p.ProjectSlug == "" {
		p.ProjectSlug = ModuleSlug(p.ProjectName)
	}
	if !p.Pages.Login && !p.Pages.Register {
		p.Pages.Login = true
		p.Pages.Register = true
	}
	if (p.Pages.Cabinet || p.Pages.CreateForm || p.Pages.Admin) && !p.Pages.Login {
		return fmt.Errorf("для кабинета/формы/админки нужна страница входа")
	}
	p.Main.normalize()
	return nil
}

func (f *FieldSchema) validate() error {
	if strings.TrimSpace(f.Column) == "" {
		return fmt.Errorf("column обязателен")
	}
	if strings.TrimSpace(f.Label) == "" {
		f.Label = f.Column
	}
	switch f.Type {
	case "string", "text", "int", "date", "enum":
	default:
		return fmt.Errorf("неподдерживаемый type: %s", f.Type)
	}
	if f.Type == "enum" && len(f.Options) == 0 {
		return fmt.Errorf("для enum нужны options")
	}
	if f.Widget == "" {
		f.Widget = defaultWidget(f.Type)
	}
	return nil
}

func defaultWidget(t string) string {
	switch t {
	case "text":
		return "textarea"
	case "int":
		return "number"
	case "date":
		return "date"
	case "enum":
		return "select"
	default:
		return "text"
	}
}

func (e *EntitySchema) normalize() {
	if e.Code == "" {
		e.Code = ToPascal(e.Label)
	} else if strings.ContainsAny(e.Code, "_- ") || !isPascalIdent(e.Code) {
		e.Code = ToPascal(e.Code)
	} else {
		// FoodOrder уже в camelCase — не ломать в Foodorder
		e.Code = ToPascal(splitCamelWords(e.Code))
	}
	if e.Table == "" {
		e.Table = ToSnake(e.Code)
	}
	for i := range e.Fields {
		e.Fields[i].Column = ToSnake(e.Fields[i].Column)
		if e.Fields[i].GoName == "" {
			e.Fields[i].GoName = ToPascal(e.Fields[i].Column)
		}
	}
}

func (f FieldSchema) GoType() string {
	switch f.Type {
	case "int":
		return "int"
	default:
		return "string"
	}
}

func (p ProjectSchema) MainHasStatus() bool {
	return len(p.Main.Statuses) > 0
}

func (p ProjectSchema) StatusNew() string {
	if len(p.Main.Statuses) == 0 {
		return ""
	}
	return p.Main.Statuses[0]
}

func (e EntitySchema) RouteName() string {
	return ToSnake(e.Code) + "s"
}

func ToSnake(s string) string {
	var b strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				b.WriteByte('_')
			}
			b.WriteRune(unicode.ToLower(r))
		} else if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		} else if r == '-' || r == ' ' || r == '_' {
			b.WriteRune('_')
		}
	}
	return strings.Trim(b.String(), "_")
}

func isPascalIdent(s string) bool {
	if s == "" {
		return false
	}
	for i, r := range s {
		if i == 0 && !unicode.IsUpper(r) {
			return false
		}
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// splitCamelWords вставляет пробелы перед заглавными: FoodOrder → Food Order.
func splitCamelWords(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	var b strings.Builder
	for i, r := range runes {
		if i > 0 && unicode.IsUpper(r) {
			prev := runes[i-1]
			nextLower := i+1 < len(runes) && unicode.IsLower(runes[i+1])
			if unicode.IsLower(prev) || (unicode.IsUpper(prev) && nextLower) {
				b.WriteRune(' ')
			}
		}
		b.WriteRune(r)
	}
	return b.String()
}

func ToPascal(s string) string {
	s = splitCamelWords(s)
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == '_' || r == '-' || r == ' ' || r == '.'
	})
	var b strings.Builder
	for _, p := range parts {
		if p == "" {
			continue
		}
		runes := []rune(strings.ToLower(p))
		runes[0] = unicode.ToUpper(runes[0])
		b.WriteString(string(runes))
	}
	out := b.String()
	if out == "" {
		return "Record"
	}
	return out
}

// Slug — имя папки (может содержать дефисы).
func Slug(name string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(name) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		} else if r == ' ' || r == '-' || r == '_' {
			b.WriteRune('-')
		}
	}
	s := strings.Trim(b.String(), "-")
	if s == "" {
		return "project"
	}
	return s
}

// ModuleSlug — валидный сегмент go.mod (без дефисов).
func ModuleSlug(name string) string {
	s := strings.ReplaceAll(Slug(name), "-", "")
	if s == "" {
		return "project"
	}
	return s
}
