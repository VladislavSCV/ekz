package theme

import (
	"embed"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:embed themes/*.yaml
var embeddedFS embed.FS

// Theme — настраиваемая предметная область (билет ДЭ).
// Меняйте YAML на экзамене, не трогая код шаблона.
type Theme struct {
	ID          string `yaml:"id"`
	Name        string `yaml:"name"`
	Tagline     string `yaml:"tagline"`
	ExamVariant string `yaml:"exam_variant"`
	ExamCode    string `yaml:"exam_code"`
	Domain      string `yaml:"domain_description"`

	Brand struct {
		Primary string `yaml:"primary"`
		Dark    string `yaml:"dark"`
		Accent  string `yaml:"accent"`
	} `yaml:"brand"`

	Admin struct {
		Login    string `yaml:"login"`
		Password string `yaml:"password"`
	} `yaml:"admin"`

	Labels struct {
		BookingSingular string `yaml:"booking_singular"`
		BookingPlural   string `yaml:"booking_plural"`
		OptionField     string `yaml:"option_field"`
		DateField       string `yaml:"date_field"`
		PaymentField    string `yaml:"payment_field"`
		DatePlaceholder string `yaml:"date_placeholder"`
		SliderSubtitle  string `yaml:"slider_subtitle"`
	} `yaml:"labels"`

	StatusNew       string `yaml:"status_new"`
	StatusAssigned  string `yaml:"status_assigned"`
	StatusCompleted string `yaml:"status_completed"`

	Options  []Option `yaml:"options"`
	Payments []Option `yaml:"payments"`
}

type Option struct {
	Value string `yaml:"value"`
	Label string `yaml:"label"`
}

func (t Theme) AllStatuses() []string {
	return []string{t.StatusNew, t.StatusAssigned, t.StatusCompleted}
}

func (t Theme) AdminStatuses() []string {
	return []string{t.StatusAssigned, t.StatusCompleted}
}

func ListEmbedded() ([]Theme, error) {
	entries, err := embeddedFS.ReadDir("themes")
	if err != nil {
		return nil, err
	}
	var out []Theme
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".yaml") {
			continue
		}
		if strings.HasPrefix(e.Name(), "_") {
			continue
		}
		th, err := LoadEmbedded(strings.TrimSuffix(e.Name(), ".yaml"))
		if err != nil {
			return nil, err
		}
		out = append(out, th)
	}
	return out, nil
}

func LoadEmbedded(id string) (Theme, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return Theme{}, fmt.Errorf("id темы не указан")
	}
	data, err := embeddedFS.ReadFile("themes/" + id + ".yaml")
	if err != nil {
		return Theme{}, fmt.Errorf("тема %q не найдена: %w", id, err)
	}
	return Parse(data)
}

func LoadFile(path string) (Theme, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Theme{}, err
	}
	return Parse(data)
}

func Parse(data []byte) (Theme, error) {
	var t Theme
	if err := yaml.Unmarshal(data, &t); err != nil {
		return Theme{}, err
	}
	if err := t.validate(); err != nil {
		return Theme{}, err
	}
	return t, nil
}

func (t *Theme) validate() error {
	if strings.TrimSpace(t.Name) == "" {
		return fmt.Errorf("theme: поле name обязательно")
	}
	if strings.TrimSpace(t.Admin.Login) == "" || strings.TrimSpace(t.Admin.Password) == "" {
		return fmt.Errorf("theme: admin.login и admin.password обязательны")
	}
	if len(t.Options) < 1 || len(t.Payments) < 1 {
		return fmt.Errorf("theme: нужны options и payments (минимум по одному)")
	}
	for _, s := range []string{t.StatusNew, t.StatusAssigned, t.StatusCompleted} {
		if strings.TrimSpace(s) == "" {
			return fmt.Errorf("theme: все три статуса обязательны")
		}
	}
	if t.Brand.Primary == "" {
		t.Brand.Primary = "#1a5fb4"
	}
	if t.Brand.Dark == "" {
		t.Brand.Dark = "#0d3d7a"
	}
	if t.Brand.Accent == "" {
		t.Brand.Accent = "#ff6b35"
	}
	if t.Labels.BookingSingular == "" {
		t.Labels.BookingSingular = "Заявка"
	}
	if t.Labels.BookingPlural == "" {
		t.Labels.BookingPlural = "Заявки"
	}
	if t.Labels.OptionField == "" {
		t.Labels.OptionField = "Тип услуги"
	}
	if t.Labels.DateField == "" {
		t.Labels.DateField = "Дата"
	}
	if t.Labels.PaymentField == "" {
		t.Labels.PaymentField = "Способ оплаты"
	}
	if t.Labels.DatePlaceholder == "" {
		t.Labels.DatePlaceholder = "01.06.2026"
	}
	if t.Labels.SliderSubtitle == "" {
		t.Labels.SliderSubtitle = t.Tagline
	}
	return nil
}

func WriteTemplate(dst string) error {
	data, err := embeddedFS.ReadFile("themes/_template.yaml")
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0o644)
}
