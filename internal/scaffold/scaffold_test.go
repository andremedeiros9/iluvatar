package scaffold_test

import (
	"path/filepath"
	"testing"

	"github.com/andremedeiros9/iluvatar/internal/config"
	"github.com/andremedeiros9/iluvatar/internal/scaffold"
)

func TestInit(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "widgets")

	err := scaffold.Init(path, scaffold.Options{
		Name:   "widgets",
		Module: "github.com/example/widgets",
	})
	if err != nil {
		t.Fatalf("Init: %v", err)
	}

	cfg, err := config.Load(filepath.Join(path, config.FileName))
	if err != nil {
		t.Fatalf("Load generated config: %v", err)
	}

	if cfg.Project.Name != "widgets" {
		t.Errorf("Project.Name = %q, want %q", cfg.Project.Name, "widgets")
	}
	if cfg.Project.Module != "github.com/example/widgets" {
		t.Errorf("Project.Module = %q, want %q", cfg.Project.Module, "github.com/example/widgets")
	}
}

func TestInitRefusesToOverwrite(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "widgets")
	opts := scaffold.Options{Name: "widgets", Module: "github.com/example/widgets"}

	if err := scaffold.Init(path, opts); err != nil {
		t.Fatalf("first Init: %v", err)
	}

	if err := scaffold.Init(path, opts); err == nil {
		t.Fatal("second Init: expected error, got nil")
	}
}
