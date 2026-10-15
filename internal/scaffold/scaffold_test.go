package scaffold_test

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/andremedeiros9/iluvatar/internal/config"
	"github.com/andremedeiros9/iluvatar/internal/scaffold"
	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "widgets")

	err := scaffold.Init(path, scaffold.Options{
		Name:   "widgets",
		Module: "github.com/example/widgets",
	})
	require.NoError(t, err, "initializing project")

	cfg, err := config.Load(filepath.Join(path, config.FileName))
	require.NoError(t, err, "loading generated config")

	require.Equal(t, "widgets", cfg.Project.Name)
	require.Equal(t, "github.com/example/widgets", cfg.Project.Module)
}

func TestInitCreatesDirectory(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "does", "not", "exist", "widgets")

	err := scaffold.Init(path, scaffold.Options{
		Name:   "widgets",
		Module: "github.com/example/widgets",
	})
	require.NoError(t, err, "initializing project should create missing directories")

	_, err = os.Stat(path)
	require.NoError(t, err, "project directory should exist")

	_, err = os.Stat(filepath.Join(path, config.FileName))
	require.NoError(t, err, "config file should exist")
}

func TestInitGoVersion(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "widgets")

	err := scaffold.Init(path, scaffold.Options{
		Name:   "widgets",
		Module: "github.com/example/widgets",
	})
	require.NoError(t, err)

	data, err := os.ReadFile(filepath.Join(path, config.FileName))
	require.NoError(t, err)

	expected := "go_version = \"" + strings.TrimPrefix(runtime.Version(), "go") + "\""
	require.Contains(t, string(data), expected, "go_version should match runtime version")
}

func TestInitTemplateContent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "widgets")

	err := scaffold.Init(path, scaffold.Options{
		Name:   "widgets",
		Module: "github.com/example/widgets",
	})
	require.NoError(t, err)

	data, err := os.ReadFile(filepath.Join(path, config.FileName))
	require.NoError(t, err)

	content := string(data)
	require.Contains(t, content, `name       = "widgets"`)
	require.Contains(t, content, `module     = "github.com/example/widgets"`)
	require.Contains(t, content, `[project]`)
	require.Contains(t, content, `[server]`)
	require.Contains(t, content, `[database]`)
	require.Contains(t, content, `[ci]`)
	require.Contains(t, content, `[cloud]`)
	require.Contains(t, content, `provider = "github"`)
	require.Contains(t, content, `provider = "none"`)
}

func TestInitEmptyOptions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "widgets")

	err := scaffold.Init(path, scaffold.Options{})
	require.NoError(t, err, "initializing with empty options should succeed")

	cfg, err := config.Load(filepath.Join(path, config.FileName))
	require.NoError(t, err)
	require.Empty(t, cfg.Project.Name)
	require.Empty(t, cfg.Project.Module)
}

func TestInitNestedPath(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "a", "b", "c", "widgets")

	err := scaffold.Init(path, scaffold.Options{
		Name:   "widgets",
		Module: "github.com/example/widgets",
	})
	require.NoError(t, err, "initializing nested path should create all parent directories")

	cfg, err := config.Load(filepath.Join(path, config.FileName))
	require.NoError(t, err)
	require.Equal(t, "widgets", cfg.Project.Name)
}

func TestInitFailsWhenPathIsAFile(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "a-file")
	require.NoError(t, os.WriteFile(filePath, []byte("not a directory"), 0o644))

	err := scaffold.Init(filePath, scaffold.Options{
		Name:   "widgets",
		Module: "github.com/example/widgets",
	})
	require.Error(t, err, "initializing where path is a file should fail")
	require.Contains(t, err.Error(), "error while creating project directory")
}

func TestInitRefusesToOverwrite(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "widgets")
	opts := scaffold.Options{Name: "widgets", Module: "github.com/example/widgets"}

	err := scaffold.Init(path, opts)
	require.NoError(t, err, "initializing project")

	err = scaffold.Init(path, opts)
	require.Error(t, err, "initializing project in existing directory should fail")
	require.Contains(t, err.Error(), "already exists", "error should mention the config file already exists")
}
