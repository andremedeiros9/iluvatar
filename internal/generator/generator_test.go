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

func TestGenerateCRUDRejectsNonGRPCCommunication(t *testing.T) {
	dir := t.TempDir()
	cfg := &config.Config{
		Project:  config.Project{Name: "widgets", Module: "github.com/example/widgets"},
		Database: config.Database{Driver: "sqlite"},
		Resources: []config.Resource{
			{Name: "user", Communication: "rest"},
		},
	}

	err := generator.Generate(dir, cfg)
	if err == nil {
		t.Fatal("Generate: expected error for communication = \"rest\", got nil")
	}
	if !strings.Contains(err.Error(), "grpc") {
		t.Errorf("error %q does not mention grpc", err)
	}
}

func TestGenerateCRUDRejectsUnsupportedFieldType(t *testing.T) {
	dir := t.TempDir()
	cfg := &config.Config{
		Project:  config.Project{Name: "widgets", Module: "github.com/example/widgets"},
		Database: config.Database{Driver: "sqlite"},
		Resources: []config.Resource{
			{
				Name:          "user",
				Communication: "grpc",
				Fields:        []config.Field{{Name: "avatar", Type: "blob"}},
			},
		},
	}

	if err := generator.Generate(dir, cfg); err == nil {
		t.Fatal("Generate: expected error for unsupported field type, got nil")
	}
}

// TestGenerateCRUDResourceFiles exercises the full gRPC CRUD generation
// path for one resource, including ensureGoModule's real `go mod tidy
// -e` call, so it needs network access and is slower than the rest of
// this package's tests.
func TestGenerateCRUDResourceFiles(t *testing.T) {
	dir := t.TempDir()
	cfg := &config.Config{
		Project:  config.Project{Name: "widgets", Module: "github.com/example/widgets", GoVersion: "1.27.1"},
		Server:   config.Server{Port: 9090},
		Database: config.Database{Driver: "postgres"},
		Resources: []config.Resource{
			{
				Name:          "user",
				Communication: "grpc",
				Fields: []config.Field{
					{Name: "email", Type: "string"},
					{Name: "signup_at", Type: "time"},
					{Name: "referrer_id", Type: "uuid"},
				},
			},
		},
	}

	if err := generator.Generate(dir, cfg); err != nil {
		t.Fatalf("Generate: %v", err)
	}

	assertContains(t, filepath.Join(dir, "internal", "user", "pb", "user.proto"), "google.protobuf.Timestamp signup_at")
	assertContains(t, filepath.Join(dir, "internal", "user", "pb", "user.proto"), "string referrer_id")
	assertContains(t, filepath.Join(dir, "internal", "user", "model.go"), `db:"email"`)
	assertContains(t, filepath.Join(dir, "internal", "user", "model.go"), "ReferrerId uuid.UUID")
	assertContains(t, filepath.Join(dir, "internal", "user", "repository.go"), "INSERT INTO user")
	assertContains(t, filepath.Join(dir, "internal", "user", "service.go"), "pb.UserServiceServer")
	assertContains(t, filepath.Join(dir, "internal", "grpcserver", "server.go"), "RegisterUserServiceServer")
	assertContains(t, filepath.Join(dir, "internal", "db", "db.go"), "jackc/pgx")
	assertContains(t, filepath.Join(dir, "cmd", "widgets", "main.go"), `fmt.Sprintf(":%d", 9090)`)
	assertContains(t, filepath.Join(dir, "migrations", "0001_create_user.up.sql"), "CREATE TABLE user")
	assertContains(t, filepath.Join(dir, "migrations", "0001_create_user.down.sql"), "DROP TABLE user")

	if _, err := os.Stat(filepath.Join(dir, "go.mod")); err != nil {
		t.Errorf("go.mod should have been created: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "go.sum")); err != nil {
		t.Errorf("go.sum should have been created: %v", err)
	}
}

func TestGenerateNoCRUDFilesWithoutResources(t *testing.T) {
	dir := t.TempDir()
	cfg := &config.Config{
		Project:  config.Project{Name: "widgets", Module: "github.com/example/widgets"},
		Database: config.Database{Driver: "sqlite"},
	}

	if err := generator.Generate(dir, cfg); err != nil {
		t.Fatalf("Generate: %v", err)
	}

	for _, p := range []string{"go.mod", "cmd", "internal/db", "internal/grpcserver", "migrations"} {
		if _, err := os.Stat(filepath.Join(dir, p)); !os.IsNotExist(err) {
			t.Errorf("%s should not exist when no resources are configured", p)
		}
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
