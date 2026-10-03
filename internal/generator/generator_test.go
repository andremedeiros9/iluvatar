package generator_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/andremedeiros9/iluvatar/internal/config"
	"github.com/andremedeiros9/iluvatar/internal/generator"
)

func TestGenerateGitHub(t *testing.T) {
	dir := t.TempDir()
	cfg := &config.Config{
		Project:  config.Project{Name: "widgets", GoVersion: "1.27.1"},
		Server:   config.Server{Port: 8080},
		Database: config.Database{Driver: "postgres"},
		CI:       config.CI{Provider: "github"},
	}

	if err := generator.Generate(dir, cfg); err != nil {
		t.Fatalf("Generate: %v", err)
	}

	assertContains(t, filepath.Join(dir, "Makefile"), "BINARY := widgets")
	assertContains(t, filepath.Join(dir, ".github", "workflows", "ci.yml"), "golangci-lint-action@v7")

	if _, err := os.Stat(filepath.Join(dir, ".gitlab-ci.yml")); !os.IsNotExist(err) {
		t.Errorf(".gitlab-ci.yml should not be written for the github provider")
	}
}

func TestGenerateGitLab(t *testing.T) {
	dir := t.TempDir()
	cfg := &config.Config{
		Project:  config.Project{Name: "widgets", GoVersion: "1.27.1"},
		Server:   config.Server{Port: 8080},
		Database: config.Database{Driver: "postgres"},
		CI:       config.CI{Provider: "gitlab"},
	}

	if err := generator.Generate(dir, cfg); err != nil {
		t.Fatalf("Generate: %v", err)
	}

	assertContains(t, filepath.Join(dir, ".gitlab-ci.yml"), "image: golang:1.27.1")

	if _, err := os.Stat(filepath.Join(dir, ".github")); !os.IsNotExist(err) {
		t.Errorf(".github should not be written for the gitlab provider")
	}
}

func TestGenerateUnsupportedCIProvider(t *testing.T) {
	dir := t.TempDir()
	cfg := &config.Config{
		Project: config.Project{Name: "widgets"},
		CI:      config.CI{Provider: "bitbucket"},
	}

	if err := generator.Generate(dir, cfg); err == nil {
		t.Fatal("Generate: expected error for unsupported CI provider, got nil")
	}
}

func TestGenerateDockerfile(t *testing.T) {
	dir := t.TempDir()
	cfg := &config.Config{
		Project:  config.Project{Name: "widgets", GoVersion: "1.27.1"},
		Server:   config.Server{Port: 9090},
		Database: config.Database{Driver: "sqlite"},
	}

	if err := generator.Generate(dir, cfg); err != nil {
		t.Fatalf("Generate: %v", err)
	}

	assertContains(t, filepath.Join(dir, "Dockerfile"), "FROM golang:1.27.1-alpine AS builder")
	assertContains(t, filepath.Join(dir, "Dockerfile"), "EXPOSE 9090")
	assertContains(t, filepath.Join(dir, "Dockerfile"), "./cmd/widgets")
}

func TestGenerateComposePostgres(t *testing.T) {
	dir := t.TempDir()
	cfg := &config.Config{
		Project:  config.Project{Name: "widgets"},
		Server:   config.Server{Port: 8080},
		Database: config.Database{Driver: "postgres"},
	}

	if err := generator.Generate(dir, cfg); err != nil {
		t.Fatalf("Generate: %v", err)
	}

	assertContains(t, filepath.Join(dir, "docker-compose.yml"), "image: postgres:16-alpine")
	assertContains(t, filepath.Join(dir, ".env.example"), "POSTGRES_DB=widgets")
}

func TestGenerateComposeMySQL(t *testing.T) {
	dir := t.TempDir()
	cfg := &config.Config{
		Project:  config.Project{Name: "widgets"},
		Server:   config.Server{Port: 8080},
		Database: config.Database{Driver: "mysql"},
	}

	if err := generator.Generate(dir, cfg); err != nil {
		t.Fatalf("Generate: %v", err)
	}

	assertContains(t, filepath.Join(dir, "docker-compose.yml"), "image: mysql:8")
	assertContains(t, filepath.Join(dir, ".env.example"), "MYSQL_DATABASE=widgets")
}

func TestGenerateComposeSkippedForSQLite(t *testing.T) {
	dir := t.TempDir()
	cfg := &config.Config{
		Project:  config.Project{Name: "widgets"},
		Server:   config.Server{Port: 8080},
		Database: config.Database{Driver: "sqlite"},
	}

	if err := generator.Generate(dir, cfg); err != nil {
		t.Fatalf("Generate: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dir, "docker-compose.yml")); !os.IsNotExist(err) {
		t.Errorf("docker-compose.yml should not be written for the sqlite driver")
	}
	if _, err := os.Stat(filepath.Join(dir, ".env.example")); !os.IsNotExist(err) {
		t.Errorf(".env.example should not be written for the sqlite driver")
	}
}

func TestGenerateUnsupportedDatabaseDriver(t *testing.T) {
	dir := t.TempDir()
	cfg := &config.Config{
		Project:  config.Project{Name: "widgets"},
		Database: config.Database{Driver: "oracle"},
	}

	if err := generator.Generate(dir, cfg); err == nil {
		t.Fatal("Generate: expected error for unsupported database driver, got nil")
	}
}

func TestGenerateRuntimeConfig(t *testing.T) {
	dir := t.TempDir()
	cfg := &config.Config{
		Project:  config.Project{Name: "widgets", Module: "github.com/example/widgets"},
		Server:   config.Server{Port: 8080},
		Database: config.Database{Driver: "postgres", DSN: "postgres://example"},
	}

	if err := generator.Generate(dir, cfg); err != nil {
		t.Fatalf("Generate: %v", err)
	}

	path := filepath.Join(dir, "config", "widgets.toml")
	loaded, err := config.Load(path)
	if err != nil {
		t.Fatalf("Load generated runtime config: %v", err)
	}

	if loaded.Project.Module != "github.com/example/widgets" {
		t.Errorf("Project.Module = %q, want %q", loaded.Project.Module, "github.com/example/widgets")
	}
	if loaded.Database.DSN != "postgres://example" {
		t.Errorf("Database.DSN = %q, want %q", loaded.Database.DSN, "postgres://example")
	}

	if _, err := os.Stat(filepath.Join(dir, config.FileName)); !os.IsNotExist(err) {
		t.Errorf("iluvatar.toml should not be written by Generate")
	}
}

func TestGenerateGitignore(t *testing.T) {
	dir := t.TempDir()
	cfg := &config.Config{
		Project:  config.Project{Name: "widgets"},
		Database: config.Database{Driver: "sqlite"},
	}

	if err := generator.Generate(dir, cfg); err != nil {
		t.Fatalf("Generate: %v", err)
	}

	assertContains(t, filepath.Join(dir, ".gitignore"), ".env")
}

func TestGenerateAWSCloudTargets(t *testing.T) {
	dir := t.TempDir()
	cfg := &config.Config{
		Project:  config.Project{Name: "widgets"},
		Database: config.Database{Driver: "sqlite"},
		Cloud:    config.Cloud{Provider: "aws", Region: "eu-west-1"},
	}

	if err := generator.Generate(dir, cfg); err != nil {
		t.Fatalf("Generate: %v", err)
	}

	assertContains(t, filepath.Join(dir, "Makefile"), "AWS_REGION     ?= eu-west-1")
	assertContains(t, filepath.Join(dir, "Makefile"), "docker-push:")
}

func TestGenerateNoCloudProviderOmitsAWSTargets(t *testing.T) {
	dir := t.TempDir()
	cfg := &config.Config{
		Project:  config.Project{Name: "widgets"},
		Database: config.Database{Driver: "sqlite"},
	}

	if err := generator.Generate(dir, cfg); err != nil {
		t.Fatalf("Generate: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dir, "Makefile"))
	if err != nil {
		t.Fatalf("reading Makefile: %v", err)
	}
	if strings.Contains(string(data), "docker-push") {
		t.Errorf("Makefile should not contain AWS targets when cloud.provider is unset:\n%s", data)
	}
}

func TestGenerateUnsupportedCloudProvider(t *testing.T) {
	dir := t.TempDir()
	cfg := &config.Config{
		Project: config.Project{Name: "widgets"},
		Cloud:   config.Cloud{Provider: "gcp"},
	}

	if err := generator.Generate(dir, cfg); err == nil {
		t.Fatal("Generate: expected error for unsupported cloud provider, got nil")
	}
}

func assertContains(t *testing.T, path, substr string) {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading %s: %v", path, err)
	}

	if !strings.Contains(string(data), substr) {
		t.Errorf("%s does not contain %q:\n%s", path, substr, data)
	}
}
