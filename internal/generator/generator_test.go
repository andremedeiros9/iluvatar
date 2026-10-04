package generator_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/andremedeiros9/iluvatar/internal/config"
	"github.com/andremedeiros9/iluvatar/internal/generator"
	"github.com/stretchr/testify/require"
)

func TestGenerateCIProvider(t *testing.T) {
	dirGitHub := t.TempDir()
	cfgGitHub := &config.Config{
		Project:  config.Project{Name: "widgets", GoVersion: "1.27.1"},
		Server:   config.Server{Port: 8080},
		Database: config.Database{Driver: "postgres"},
		CI:       config.CI{Provider: "github"},
	}

	dirGitLab := t.TempDir()
	cfgGitLab := &config.Config{
		Project:  config.Project{Name: "widgets", GoVersion: "1.27.1"},
		Server:   config.Server{Port: 8080},
		Database: config.Database{Driver: "postgres"},
		CI:       config.CI{Provider: "gitlab"},
	}

	unsupportedDir := t.TempDir()
	cfgUnsupported := &config.Config{
		Project: config.Project{Name: "widgets"},
		CI:      config.CI{Provider: "bitbucket"},
	}

	type testCase struct {
		name       string
		dir        string
		cfg        *config.Config
		assertions func(error)
	}
	tC := []testCase{
		{
			name: "GenerateGitHub",
			dir:  dirGitHub,
			cfg:  cfgGitHub,
			assertions: func(err error) {
				require.NoError(t, err)
				assertContains(t, filepath.Join(dirGitHub, "Makefile"), "BINARY := widgets")
				assertContains(t, filepath.Join(dirGitHub, ".github", "workflows", "ci.yml"), "golangci-lint-action@v7")
				if _, err := os.Stat(filepath.Join(dirGitHub, ".gitlab-ci.yml")); !os.IsNotExist(err) {
					require.Error(t, err)
					require.Equal(t, err.Error(), ".gitlab-ci.yml should not be written for the github provider")
				}
			},
		},
		{
			name: "GenerateGitLab",
			dir:  dirGitLab,
			cfg:  cfgGitLab,
			assertions: func(err error) {
				require.NoError(t, err)
				assertContains(t, filepath.Join(dirGitLab, ".gitlab-ci.yml"), "image: golang:1.27.1")
				if _, err := os.Stat(filepath.Join(dirGitLab, ".github")); !os.IsNotExist(err) {
					require.Error(t, err)
					require.Equal(t, err.Error(), ".github should not be written for the gitlab provider")
				}
			},
		},
		{
			name: "GenerateUnsupportedCIProvider",
			dir:  unsupportedDir,
			cfg:  cfgUnsupported,
			assertions: func(err error) {
				require.Error(t, err)
				require.Equal(t, err.Error(), "unsupported ci provider \"bitbucket\" (want \"github\" or \"gitlab\")")
			},
		},
	}

	for _, tc := range tC {
		t.Run(tc.name, func(t *testing.T) {
			err := generator.Generate(tc.dir, tc.cfg)
			tc.assertions(err)
		})
	}
}

func TestGenerateDockerfile(t *testing.T) {
	dir := t.TempDir()
	cfg := &config.Config{
		Project:  config.Project{Name: "widgets", GoVersion: "1.27.1"},
		Server:   config.Server{Port: 9090},
		Database: config.Database{Driver: "sqlite"},
	}

	err := generator.Generate(dir, cfg)
	require.NoError(t, err)

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

	err := generator.Generate(dir, cfg)
	require.NoError(t, err)

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

	err := generator.Generate(dir, cfg)
	require.NoError(t, err)

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

	err := generator.Generate(dir, cfg)
	require.NoError(t, err)

	_, err = os.Stat(filepath.Join(dir, "docker-compose.yml"))
	require.True(t, os.IsNotExist(err), "docker-compose.yml should not be written for the sqlite driver")
	_, err = os.Stat(filepath.Join(dir, ".env.example"))
	require.True(t, os.IsNotExist(err), ".env.example should not be written for the sqlite driver")
}

func TestGenerateUnsupportedDatabaseDriver(t *testing.T) {
	dir := t.TempDir()
	cfg := &config.Config{
		Project:  config.Project{Name: "widgets"},
		Database: config.Database{Driver: "oracle"},
	}

	err := generator.Generate(dir, cfg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "unsupported database driver")
}

func TestGenerateRuntimeConfig(t *testing.T) {
	dir := t.TempDir()
	cfg := &config.Config{
		Project:  config.Project{Name: "widgets", Module: "github.com/example/widgets"},
		Server:   config.Server{Port: 8080},
		Database: config.Database{Driver: "postgres", DSN: "postgres://example"},
	}

	err := generator.Generate(dir, cfg)
	require.NoError(t, err)

	path := filepath.Join(dir, "config", "widgets.toml")
	loaded, err := config.Load(path)
	require.NoError(t, err, "loading generated config")

	require.Equal(t, cfg.Project.Name, loaded.Project.Name)
	require.Equal(t, cfg.Project.Module, loaded.Project.Module)
	require.Equal(t, cfg.Server.Port, loaded.Server.Port)
	require.Equal(t, cfg.Database.Driver, loaded.Database.Driver)
	require.Equal(t, cfg.Database.DSN, loaded.Database.DSN)

	_, err = os.Stat(filepath.Join(dir, config.FileName))
	require.True(t, os.IsNotExist(err), "iluvatar.toml should not be written by Generate")
}

func TestGenerateGitignore(t *testing.T) {
	dir := t.TempDir()
	cfg := &config.Config{
		Project:  config.Project{Name: "widgets"},
		Database: config.Database{Driver: "sqlite"},
	}

	err := generator.Generate(dir, cfg)
	require.NoError(t, err)

	assertContains(t, filepath.Join(dir, ".gitignore"), ".env")
}

func TestGenerateAWSCloudTargets(t *testing.T) {
	dir := t.TempDir()
	cfg := &config.Config{
		Project:  config.Project{Name: "widgets"},
		Database: config.Database{Driver: "sqlite"},
		Cloud:    config.Cloud{Provider: "aws", Region: "eu-west-1"},
	}

	err := generator.Generate(dir, cfg)
	require.NoError(t, err)

	assertContains(t, filepath.Join(dir, "Makefile"), "AWS_REGION     ?= eu-west-1")
	assertContains(t, filepath.Join(dir, "Makefile"), "docker-push:")
}

func TestGenerateNoCloudProviderOmitsAWSTargets(t *testing.T) {
	dir := t.TempDir()
	cfg := &config.Config{
		Project:  config.Project{Name: "widgets"},
		Database: config.Database{Driver: "sqlite"},
	}

	err := generator.Generate(dir, cfg)
	require.NoError(t, err)

	data, err := os.ReadFile(filepath.Join(dir, "Makefile"))
	require.NoError(t, err, "reading Makefile")
	require.False(t, strings.Contains(string(data), "docker-push:"), "Makefile should not contain docker-push target when no cloud provider is configured")
}

func TestGenerateUnsupportedCloudProvider(t *testing.T) {
	dir := t.TempDir()
	cfg := &config.Config{
		Project: config.Project{Name: "widgets"},
		Cloud:   config.Cloud{Provider: "gcp"},
	}

	err := generator.Generate(dir, cfg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "unsupported cloud provider")
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
	require.Error(t, err)
	require.Contains(t, err.Error(), "grpc")
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

	err := generator.Generate(dir, cfg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "unsupported field type")
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

	err := generator.Generate(dir, cfg)
	require.NoError(t, err)

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

	_, err = os.Stat(filepath.Join(dir, "go.mod"))
	require.NoError(t, err, "go.mod should have been created")

	_, err = os.Stat(filepath.Join(dir, "go.sum"))
	require.NoError(t, err, "go.sum should have been created")

}

func TestGenerateNoCRUDFilesWithoutResources(t *testing.T) {
	dir := t.TempDir()
	cfg := &config.Config{
		Project:  config.Project{Name: "widgets", Module: "github.com/example/widgets"},
		Database: config.Database{Driver: "sqlite"},
	}

	err := generator.Generate(dir, cfg)
	require.NoError(t, err)

	for _, p := range []string{"go.mod", "cmd", "internal/db", "internal/grpcserver", "migrations"} {
		_, err := os.Stat(filepath.Join(dir, p))
		require.Error(t, err, "%s should not exist when no resources are configured", p)
	}
}

func assertContains(t *testing.T, path, substr string) {
	t.Helper()

	data, err := os.ReadFile(path)
	require.NoError(t, err, "reading %s", path)

	require.Contains(t, string(data), substr)
}
