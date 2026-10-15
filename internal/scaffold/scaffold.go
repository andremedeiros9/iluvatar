package scaffold

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	"github.com/andremedeiros9/iluvatar/internal/config"
)

//go:embed templates/iluvatar.toml.tmpl
var templatesFS embed.FS

// Options customizes the config rendered by Init.
type Options struct {
	// Name is the project name recorded in iluvatar.toml.
	Name string
	// Module is the Go module path recorded in iluvatar.toml.
	Module string
	// Architecture is the code layout to generate: "vertical" (default)
	// or "horizontal".
	Architecture string
	// Port is the gRPC server's listen port.
	Port int
	// DatabaseDriver is one of: postgres, mysql, sqlite.
	DatabaseDriver string
	// CIProvider is one of: github, gitlab.
	CIProvider string
	// CloudProvider is one of: none, aws.
	CloudProvider string
	// CloudRegion is the provider region (only used with aws).
	CloudRegion string
}

// Init creates path (if needed) and writes a starter iluvatar.toml into it.
// It fails if path already contains an iluvatar.toml, so it never
// overwrites an existing project.
func Init(path string, opts Options) error {
	if err := os.MkdirAll(path, 0o755); err != nil {
		return fmt.Errorf("error while creating project directory: %w", err)
	}

	configPath := filepath.Join(path, config.FileName)
	switch _, err := os.Stat(configPath); {
	case err == nil:
		return fmt.Errorf("error: %s already exists", configPath)
	case !os.IsNotExist(err):
		return fmt.Errorf("error while checking for existing %s: %w", config.FileName, err)
	}

	// Apply sensible defaults for any zero-valued options.
	architecture := opts.Architecture
	if architecture == "" {
		architecture = config.ArchitectureVertical
	}
	port := opts.Port
	if port == 0 {
		port = 8080
	}
	databaseDriver := opts.DatabaseDriver
	if databaseDriver == "" {
		databaseDriver = "postgres"
	}
	ciProvider := opts.CIProvider
	if ciProvider == "" {
		ciProvider = "github"
	}
	cloudProvider := opts.CloudProvider
	if cloudProvider == "" {
		cloudProvider = "none"
	}
	cloudRegion := opts.CloudRegion
	if cloudRegion == "" {
		cloudRegion = "us-east-1"
	}

	tmpl, err := template.ParseFS(templatesFS, "templates/iluvatar.toml.tmpl")
	if err != nil {
		return fmt.Errorf("error while parsing config template: %w", err)
	}

	f, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("error while creating %s: %w", config.FileName, err)
	}
	defer func() {
		_ = f.Close()
	}()

	data := struct {
		Name           string
		Module         string
		GoVersion      string
		Architecture   string
		Port           int
		DatabaseDriver string
		CIProvider     string
		CloudProvider  string
		CloudRegion    string
	}{
		Name:           opts.Name,
		Module:         opts.Module,
		GoVersion:      strings.TrimPrefix(runtime.Version(), "go"),
		Architecture:   architecture,
		Port:           port,
		DatabaseDriver: databaseDriver,
		CIProvider:     ciProvider,
		CloudProvider:  cloudProvider,
		CloudRegion:    cloudRegion,
	}

	if err := tmpl.Execute(f, data); err != nil {
		return fmt.Errorf("error while writing %s: %w", config.FileName, err)
	}

	return nil
}
