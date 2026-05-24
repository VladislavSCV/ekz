package config

import "github.com/VladislavSCV/ekz/internal/theme"

// FromConferencesTheme — быстрый пресет «Конференции.РФ» (билет 2026, вариант 2).
func FromConferencesTheme(projectName string, th theme.Theme) ProjectSchema {
	fields := []FieldSchema{
		{Column: "room_type", Label: th.Labels.OptionField, Type: "enum", Widget: "select", Required: true,
			Options: copyOpts(th.Options)},
		{Column: "start_date", Label: th.Labels.DateField, Type: "date", Widget: "date", Required: true},
		{Column: "payment_method", Label: th.Labels.PaymentField, Type: "enum", Widget: "select", Required: true,
			Options: copyOpts(th.Payments)},
	}
	p := ProjectSchema{
		ProjectName: projectName,
		ProjectSlug: ModuleSlug(projectName),
		Main: EntitySchema{
			Label:    th.Labels.BookingSingular,
			Code:     "Booking",
			Table:    "bookings",
			Fields:   fields,
			Statuses: []string{th.StatusNew, th.StatusAssigned, th.StatusCompleted},
		},
		Reviews: true,
		Pages: PagesSchema{
			Login: true, Register: true, CreateForm: true,
			Cabinet: true, Admin: true, Slider: true,
		},
	}
	p.Portal.Name = th.Name
	p.Portal.Tagline = th.Tagline
	p.Admin.Login = th.Admin.Login
	p.Admin.Password = th.Admin.Password
	p.Brand.Primary = th.Brand.Primary
	p.Brand.Dark = th.Brand.Dark
	p.Brand.Accent = th.Brand.Accent
	_ = p.Validate()
	return p
}

func copyOpts(in []theme.Option) []Option {
	out := make([]Option, len(in))
	for i, o := range in {
		out[i] = Option{Value: o.Value, Label: o.Label}
	}
	return out
}
