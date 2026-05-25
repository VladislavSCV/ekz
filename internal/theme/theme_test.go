package theme

import (
	"testing"
)

func TestLoadEmbeddedConferences(t *testing.T) {
	th, err := LoadEmbedded("conferences")
	if err != nil {
		t.Fatal(err)
	}
	if th.ID != "conferences" {
		t.Errorf("id: got %q", th.ID)
	}
	if len(th.Options) < 1 || len(th.Payments) < 1 {
		t.Fatal("options and payments required")
	}
	for _, s := range th.AllStatuses() {
		if s == "" {
			t.Fatal("empty status")
		}
	}
}

func TestListEmbedded(t *testing.T) {
	list, err := ListEmbedded()
	if err != nil {
		t.Fatal(err)
	}
	if len(list) < 1 {
		t.Fatal("expected at least one theme")
	}
}

func TestParseInvalid(t *testing.T) {
	_, err := Parse([]byte("name: only-name"))
	if err == nil {
		t.Fatal("expected validation error")
	}
}
