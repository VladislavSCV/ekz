package config

import (
	"bytes"
	"os"

	"gopkg.in/yaml.v3"
)

func LoadFile(path string) (ProjectSchema, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return ProjectSchema{}, err
	}
	var p ProjectSchema
	if err := yaml.Unmarshal(data, &p); err != nil {
		return ProjectSchema{}, err
	}
	if err := p.Validate(); err != nil {
		return ProjectSchema{}, err
	}
	return p, nil
}

func SaveFile(path string, p ProjectSchema) error {
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(&p); err != nil {
		return err
	}
	_ = enc.Close()
	return os.WriteFile(path, buf.Bytes(), 0o644)
}
