package generator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRenderGoProducesValidGo(t *testing.T) {
	dir := t.TempDir()
	data := templateData{
		BinaryName:          "widgets",
		Module:              "github.com/example/widgets",
		GoVersion:           "1.27.1",
		GolangciLintVersion: "v2.14.0",
		Port:                8080,
	}

	err := renderGo(dir, filepath.Join("cmd", "widgets", "main.go"), "templates/cmd/main.go.tmpl", data)
	require.NoError(t, err)

	_, err = os.Stat(filepath.Join(dir, "cmd", "widgets", "main.go"))
	require.NoError(t, err, "rendered Go file should exist")

	// The file should be valid Go source (gofmt-formatted)
	content, err := os.ReadFile(filepath.Join(dir, "cmd", "widgets", "main.go"))
	require.NoError(t, err)
	require.Contains(t, string(content), "package main")
}

func TestRenderGoInvalidTemplateFails(t *testing.T) {
	dir := t.TempDir()
	data := templateData{}

	err := renderGo(dir, "bad.go", "templates/does-not-exist.tmpl", data)
	require.Error(t, err)
	require.Contains(t, err.Error(), "parsing")
}

func TestRenderGoCreatesParentDirectories(t *testing.T) {
	dir := t.TempDir()
	data := templateData{BinaryName: "widgets"}

	err := renderGo(dir, filepath.Join("deep", "nested", "main.go"), "templates/cmd/main.go.tmpl", data)
	require.NoError(t, err)

	_, err = os.Stat(filepath.Join(dir, "deep", "nested", "main.go"))
	require.NoError(t, err)
}

func TestRender(t *testing.T) {
	dir := t.TempDir()
	data := templateData{BinaryName: "widgets"}

	err := render(dir, "Makefile", "templates/Makefile.tmpl", data)
	require.NoError(t, err)

	content, err := os.ReadFile(filepath.Join(dir, "Makefile"))
	require.NoError(t, err)
	require.Contains(t, string(content), "BINARY := widgets")
	require.Contains(t, string(content), "golangci-lint run ./...")
}

func TestRenderCreatesParentDirectories(t *testing.T) {
	dir := t.TempDir()
	data := templateData{}

	err := render(dir, filepath.Join("a", "b", "c", "file.txt"), "templates/Makefile.tmpl", data)
	require.NoError(t, err)

	_, err = os.Stat(filepath.Join(dir, "a", "b", "c", "file.txt"))
	require.NoError(t, err)
}
