package codegen

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/VladislavSCV/ekz/internal/config"
)

//go:embed templates/*
var templatesFS embed.FS

//go:embed static/**/*
var staticFS embed.FS

func Generate(target string, schema config.ProjectSchema) error {
	if err := os.MkdirAll(target, 0o755); err != nil {
		return err
	}
	data := schema

	funcMap := template.FuncMap{"lower": strings.ToLower, "add": func(a, b int) int { return a + b }}
	err := fs.WalkDir(templatesFS, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".tmpl") {
			return nil
		}
		rel := strings.TrimPrefix(path, "templates/")
		rel = strings.TrimPrefix(rel, "templates\\")
		destRel := strings.TrimSuffix(rel, ".tmpl")
		if skipTemplate(destRel, data) {
			return nil
		}
		dest := filepath.Join(target, filepath.FromSlash(destRel))
		content, err := fs.ReadFile(templatesFS, path)
		if err != nil {
			return err
		}
		tmpl, err := template.New(filepath.Base(path)).Funcs(funcMap).Parse(string(content))
		if err != nil {
			return fmt.Errorf("parse %s: %w", path, err)
		}
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, data); err != nil {
			return fmt.Errorf("execute %s: %w", path, err)
		}
		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(dest, buf.Bytes(), 0o644); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	if err := fs.WalkDir(staticFS, "static", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, _ := filepath.Rel("static", path)
		dest := filepath.Join(target, filepath.FromSlash(rel))
		content, err := fs.ReadFile(staticFS, path)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			return err
		}
		return os.WriteFile(dest, content, 0o644)
	}); err != nil {
		return err
	}

	if err := config.SaveFile(filepath.Join(target, "project.yaml"), schema); err != nil {
		return err
	}

	if err := writeGitignore(target); err != nil {
		return err
	}
	if err := initGit(target); err != nil {
		fmt.Fprintf(os.Stderr, "предупреждение: %v\n", err)
	}
	return nil
}

func writeGitignore(target string) error {
	return os.WriteFile(filepath.Join(target, ".gitignore"), []byte(`backend/data/
*.db
node_modules/
dist/
`), 0o644)
}

func skipTemplate(destRel string, s config.ProjectSchema) bool {
	switch filepath.ToSlash(destRel) {
	case "frontend/index.html":
		return !s.Pages.Login
	case "frontend/register.html":
		return !s.Pages.Register
	case "frontend/create.html", "frontend/src/pages/create.js":
		return !s.Pages.CreateForm
	case "frontend/cabinet.html", "frontend/src/pages/cabinet.js":
		return !s.Pages.Cabinet
	case "frontend/admin.html", "frontend/src/pages/admin.js":
		return !s.Pages.Admin
	case "frontend/public/slider/slide.svg":
		return !s.Pages.Slider
	}
	return false
}

func initGit(dir string) error {
	cmd := exec.Command("git", "init")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git init: %w (%s)", err, strings.TrimSpace(string(out)))
	}
	return nil
}
