package generator

import (
	"embed"
	"fmt"

	"github.com/andremedeiros9/iluvatar/internal/config"
)

//go:embed templates
var templatesFS embed.FS

// golangciLintVersion pins the golangci-lint release used by generated CI
// pipelines, matching the version iluvatar itself lints with.
const golangciLintVersion = "v2.14.0"

// templateData is the set of values available to every generator template.
type templateData struct {
	BinaryName          string
	GoVersion           string
	GolangciLintVersion string
	Port                int
	CloudProvider       string
	CloudRegion         string
}

// Generate writes the generated project's build tooling and deployment
// scaffolding into dir, based on cfg: a Makefile, a CI pipeline, a
// Dockerfile, a docker-compose setup (when the database needs one), a
// .gitignore, and the runtime config the generated app itself will load
// (config/<project-name>.toml). Unlike iluvatar.toml, which only exists to
// drive generation, that file travels with the generated project.
func Generate(dir string, cfg *config.Config) error {
	cloudProvider, err := normalizeCloudProvider(cfg.Cloud.Provider)
	if err != nil {
		return err
	}

	data := templateData{
		BinaryName:          cfg.Project.Name,
		GoVersion:           cfg.Project.GoVersion,
		GolangciLintVersion: golangciLintVersion,
		Port:                cfg.Server.Port,
		CloudProvider:       cloudProvider,
		CloudRegion:         cfg.Cloud.Region,
	}

	if err := writeMakefile(dir, data); err != nil {
		return err
	}
	if err := writeCI(dir, cfg.CI.Provider, data); err != nil {
		return err
	}
	if err := writeDockerfile(dir, data); err != nil {
		return err
	}
	if err := writeCompose(dir, cfg.Database.Driver, data); err != nil {
		return err
	}
	if err := writeRuntimeConfig(dir, cfg); err != nil {
		return err
	}

	return writeGitignore(dir)
}

func normalizeCloudProvider(raw string) (string, error) {
	switch raw {
	case "":
		return "none", nil
	case "none", "aws":
		return raw, nil
	default:
		return "", fmt.Errorf("unsupported cloud provider %q (want %q or %q)", raw, "none", "aws")
	}
}
