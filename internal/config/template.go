package config

import "fmt"

// ExampleSchema — заготовка project.yaml для правки в редакторе.
func ExampleSchema() ProjectSchema {
	p := ProjectSchema{
		ProjectName: "my-ticket",
		ProjectSlug: "myticket",
	}
	p.Portal.Name = "Мой портал"
	p.Portal.Tagline = "Краткое описание билета"
	p.Admin.Login = "Admin26"
	p.Admin.Password = "Demo20"
	p.Brand.Primary = "#1a5fb4"
	p.Brand.Dark = "#0d3d7a"
	p.Brand.Accent = "#ff6b35"
	p.Main = EntitySchema{
		Label: "Заявка",
		Code:  "Booking",
		Table: "bookings",
		Fields: []FieldSchema{
			{Column: "title", Label: "Название", Type: "string", Widget: "text", Required: true},
			{Column: "event_date", Label: "Дата", Type: "date", Widget: "date", Required: true},
		},
		Statuses: []string{"Новая", "В работе", "Завершена"},
	}
	p.Reviews = true
	p.Pages = PagesSchema{
		Login: true, Register: true, CreateForm: true,
		Cabinet: true, Admin: true, Slider: true,
	}
	return p
}

// WriteSchemaTemplate сохраняет пример project.yaml.
func WriteSchemaTemplate(path string) error {
	p := ExampleSchema()
	if err := p.Validate(); err != nil {
		return err
	}
	if err := SaveFile(path, p); err != nil {
		return err
	}
	fmt.Printf("✓ Шаблон project.yaml: %s\n", path)
	fmt.Println("  Правьте поля, затем: ekz -name <папка> -config", path)
	return nil
}

// ExportBuiltinPreset записывает встроенный шаблон в файл (без генерации кода).
func ExportBuiltinPreset(presetID, path string) error {
	p, err := LoadBuiltinPreset(presetID)
	if err != nil {
		return err
	}
	if path == "" {
		path = presetID + ".yaml"
	}
	if err := SaveFile(path, p); err != nil {
		return err
	}
	fmt.Printf("✓ %s → %s\n", presetID, path)
	fmt.Println("  Отредактируйте файл, затем: ekz -name <папка> -config", path)
	return nil
}
