package config_test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/andremedeiros9/iluvatar/internal/config"
)

func TestLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, config.FileName)

	err := os.WriteFile(path, []byte(`
[project]
name = "widgets"
module = "github.com/example/widgets"
go_version = "1.25.6"

[server]
framework = "net/http"
port = 8080

[database]
driver = "postgres"
dsn = ""

[ci]
provider = "gitlab"

[cloud]
provider = "aws"
region = "eu-west-1"

[[resources]]
name = "widget"
communication = "grpc"

  [[resources.fields]]
  name = "sku"
  type = "string"
`), 0o644)
	if err != nil {
		t.Fatalf("writing fixture: %v", err)
	}

	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if cfg.Project.Name != "widgets" {
		t.Errorf("Project.Name = %q, want %q", cfg.Project.Name, "widgets")
	}
	if cfg.Project.Module != "github.com/example/widgets" {
		t.Errorf("Project.Module = %q, want %q", cfg.Project.Module, "github.com/example/widgets")
	}
	if cfg.Server.Port != 8080 {
		t.Errorf("Server.Port = %d, want %d", cfg.Server.Port, 8080)
	}
	if cfg.Database.Driver != "postgres" {
		t.Errorf("Database.Driver = %q, want %q", cfg.Database.Driver, "postgres")
	}
	if cfg.CI.Provider != "gitlab" {
		t.Errorf("CI.Provider = %q, want %q", cfg.CI.Provider, "gitlab")
	}
	if cfg.Cloud.Provider != "aws" {
		t.Errorf("Cloud.Provider = %q, want %q", cfg.Cloud.Provider, "aws")
	}
	if cfg.Cloud.Region != "eu-west-1" {
		t.Errorf("Cloud.Region = %q, want %q", cfg.Cloud.Region, "eu-west-1")
	}
	if len(cfg.Resources) != 1 || cfg.Resources[0].Name != "widget" {
		t.Fatalf("Resources = %+v, want one resource named %q", cfg.Resources, "widget")
	}
	if cfg.Resources[0].Communication != "grpc" {
		t.Errorf("Resources[0].Communication = %q, want %q", cfg.Resources[0].Communication, "grpc")
	}
	if len(cfg.Resources[0].Fields) != 1 || cfg.Resources[0].Fields[0].Name != "sku" {
		t.Errorf("Resources[0].Fields = %+v, want one field named %q", cfg.Resources[0].Fields, "sku")
	}
}

func TestSave(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config", "widgets.toml")

	want := &config.Config{
		Project:  config.Project{Name: "widgets", Module: "github.com/example/widgets", GoVersion: "1.27.1"},
		Server:   config.Server{Framework: "net/http", Port: 8080},
		Database: config.Database{Driver: "postgres", DSN: ""},
		CI:       config.CI{Provider: "github"},
		Cloud:    config.Cloud{Provider: "aws", Region: "eu-west-1"},
		Resources: []config.Resource{
			{Name: "widget", Communication: "rest", Fields: []config.Field{{Name: "sku", Type: "string"}}},
		},
	}

	if err := config.Save(path, want); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := config.Load(path)
	if err != nil {
		t.Fatalf("Load after Save: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("round-tripped config mismatch:\ngot  %+v\nwant %+v", got, want)
	}
}
