package codegen

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/VladislavSCV/ekz/internal/config"
)

func TestGeneratePresetsBackendStarts(t *testing.T) {
	presets := []struct {
		id   string
		name string
	}{
		{"food-delivery", "test-food"},
		{"conferences", "test-conf"},
	}
	for _, tc := range presets {
		t.Run(tc.id, func(t *testing.T) {
			schema, err := config.LoadBuiltinPreset(tc.id)
			if err != nil {
				t.Fatal(err)
			}
			schema.ProjectName = tc.name
			schema.ProjectSlug = config.ModuleSlug(tc.name)

			dir := t.TempDir()
			if err := Generate(dir, schema); err != nil {
				t.Fatalf("generate: %v", err)
			}
			assertModuleImportsMatch(t, dir, schema.ProjectSlug)
			assertProjectYAMLRoundTrip(t, dir, schema)
			testBackendBuildAndListen(t, filepath.Join(dir, "backend"), schema.Admin.Login, schema.Admin.Password)
		})
	}
}

func assertModuleImportsMatch(t *testing.T, projectDir, slug string) {
	t.Helper()
	modPath := filepath.Join(projectDir, "backend", "go.mod")
	data, err := os.ReadFile(modPath)
	if err != nil {
		t.Fatal(err)
	}
	modLine := strings.TrimSpace(strings.Split(string(data), "\n")[0])
	wantMod := "module " + slug + "/backend"
	if modLine != wantMod {
		t.Fatalf("go.mod: %q, want %q", modLine, wantMod)
	}
	mainPath := filepath.Join(projectDir, "backend", "main.go")
	main, err := os.ReadFile(mainPath)
	if err != nil {
		t.Fatal(err)
	}
	wantImport := `"` + slug + `/backend/internal/database"`
	if !strings.Contains(string(main), wantImport) {
		t.Fatalf("main.go missing import %s", wantImport)
	}
}

func assertProjectYAMLRoundTrip(t *testing.T, projectDir string, want config.ProjectSchema) {
	t.Helper()
	path := filepath.Join(projectDir, "project.yaml")
	got, err := config.LoadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if got.Portal.Name != want.Portal.Name {
		t.Errorf("portal.name: %q vs %q", got.Portal.Name, want.Portal.Name)
	}
	if got.Main.Code != want.Main.Code {
		t.Errorf("main.code: %q vs %q", got.Main.Code, want.Main.Code)
	}
}

func testBackendBuildAndListen(t *testing.T, backendDir, adminLogin, adminPassword string) {
	t.Helper()
	exeName := "srv"
	if runtime.GOOS == "windows" {
		exeName += ".exe"
	}
	exe := filepath.Join(t.TempDir(), exeName)

	tidy := exec.Command("go", "mod", "tidy")
	tidy.Dir = backendDir
	if out, err := tidy.CombinedOutput(); err != nil {
		t.Fatalf("go mod tidy: %v\n%s", err, out)
	}

	build := exec.Command("go", "build", "-o", exe, ".")
	build.Dir = backendDir
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("go build: %v\n%s", err, out)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	run := exec.CommandContext(ctx, exe)
	run.Dir = backendDir
	if err := run.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if run.Process != nil {
			_ = run.Process.Kill()
		}
		_, _ = run.Process.Wait()
	}()

	deadline := time.Now().Add(8 * time.Second)
	var lastErr error
	for time.Now().Before(deadline) {
		if run.ProcessState != nil && run.ProcessState.Exited() {
			out, _ := run.CombinedOutput()
			t.Fatalf("backend exited early: %s", out)
		}
		resp, err := http.Post(
			"http://127.0.0.1:8080/api/login",
			"application/json",
			bytes.NewReader(mustJSON(t, map[string]string{
				"login":    adminLogin,
				"password": adminPassword,
			})),
		)
		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return
			}
			body, _ := io.ReadAll(resp.Body)
			lastErr = fmt.Errorf("login status %d: %s", resp.StatusCode, body)
		} else {
			lastErr = err
		}
		time.Sleep(200 * time.Millisecond)
	}
	t.Fatalf("backend did not accept login: %v", lastErr)
}

func mustJSON(t *testing.T, v any) []byte {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return b
}
