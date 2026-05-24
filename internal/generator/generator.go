package generator

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"ekz/internal/presets"
)

type fileData struct {
	ProjectName string
	ProjectSlug string
}

func Generate(p presets.Preset, targetDir, projectName string) error {
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return err
	}

	data := fileData{
		ProjectName: projectName,
		ProjectSlug: slug(projectName),
	}

	root := p.Root
	err := fs.WalkDir(p.Template, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}

		destRel := rel
		if strings.HasSuffix(destRel, ".tmpl") {
			destRel = strings.TrimSuffix(destRel, ".tmpl")
		}
		destPath := filepath.Join(targetDir, filepath.FromSlash(destRel))

		if d.IsDir() {
			return os.MkdirAll(destPath, 0o755)
		}

		content, err := fs.ReadFile(p.Template, path)
		if err != nil {
			return err
		}

		if strings.HasSuffix(path, ".tmpl") {
			tmpl, err := template.New(filepath.Base(path)).Parse(string(content))
			if err != nil {
				return fmt.Errorf("шаблон %s: %w", path, err)
			}
			var buf bytes.Buffer
			if err := tmpl.Execute(&buf, data); err != nil {
				return fmt.Errorf("выполнение %s: %w", path, err)
			}
			content = buf.Bytes()
		}

		if err := os.MkdirAll(filepath.Dir(destPath), 0o755); err != nil {
			return err
		}
		return os.WriteFile(destPath, content, 0o644)
	})
	if err != nil {
		return err
	}

	if err := initGit(targetDir); err != nil {
		fmt.Fprintf(os.Stderr, "предупреждение: %v\n", err)
	}
	return nil
}

func slug(name string) string {
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

func initGit(dir string) error {
	cmd := exec.Command("git", "init")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git init: %w (%s)", err, strings.TrimSpace(string(out)))
	}
	return nil
}
