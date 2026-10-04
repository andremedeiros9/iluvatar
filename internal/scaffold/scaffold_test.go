package scaffold_test

import (
	"path/filepath"
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

func TestInitRefusesToOverwrite(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "widgets")
	opts := scaffold.Options{Name: "widgets", Module: "github.com/example/widgets"}

	err := scaffold.Init(path, opts)
	require.NoError(t, err, "initializing project")

	err = scaffold.Init(path, opts)
	require.Error(t, err, "initializing project in existing directory should fail")
}
