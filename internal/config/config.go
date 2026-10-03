package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

// FileName is the name iluvatar looks for in a project directory.
const FileName = "iluvatar.toml"

// Config is the parsed contents of a project's iluvatar.toml.
type Config struct {
	Project   Project    `toml:"project"`
	Server    Server     `toml:"server"`
	Database  Database   `toml:"database"`
	CI        CI         `toml:"ci"`
	Cloud     Cloud      `toml:"cloud"`
	Resources []Resource `toml:"resources"`
}

// Project identifies the generated module.
type Project struct {
	Name      string `toml:"name"`
	Module    string `toml:"module"`
	GoVersion string `toml:"go_version"`
}

// Server configures the generated HTTP server.
type Server struct {
	// Framework is one of: net/http, chi, gin, echo.
	Framework string `toml:"framework"`
	Port      int    `toml:"port"`
}

// Database configures the generated data layer.
type Database struct {
	// Driver is one of: postgres, mysql, sqlite.
	Driver string `toml:"driver"`
	DSN    string `toml:"dsn"`
}

// CI configures the generated project's CI pipeline.
type CI struct {
	// Provider is one of: github, gitlab. Defaults to github.
	Provider string `toml:"provider"`
}

// Cloud configures the cloud provider the generated image targets.
type Cloud struct {
	// Provider is one of: none, aws. Defaults to none.
	Provider string `toml:"provider"`
	// Region is the provider region (e.g. an AWS region). Only read when
	// Provider is a real cloud provider.
	Region string `toml:"region"`
}

// Field is a single column/attribute on a Resource.
type Field struct {
	Name string `toml:"name"`
	Type string `toml:"type"`
}

// Resource describes one set of CRUD endpoints to generate.
type Resource struct {
	Name   string  `toml:"name"`
	Fields []Field `toml:"fields"`
	// Communication is the protocol this resource's endpoints are exposed
	// over: rest, grpc, or rpc. Defaults to rest.
	Communication string `toml:"communication"`
}

// Load reads and parses iluvatar.toml at path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}

	var cfg Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	return &cfg, nil
}

// Save encodes cfg as TOML and writes it to path, creating any missing
// parent directories.
func Save(path string, cfg *Config) error {
	data, err := toml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("encoding config: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("creating directory for config: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("writing config: %w", err)
	}

	return nil
}
