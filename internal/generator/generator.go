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
	Module              string
	GoVersion           string
	GolangciLintVersion string
	Port                int
	CloudProvider       string
	CloudRegion         string
}

// Generate writes the generated project's build tooling, deployment
// scaffolding, and gRPC CRUD implementation into dir, based on cfg: a
// Makefile, a CI pipeline, a Dockerfile, a docker-compose setup (when the
// database needs one), a .gitignore, the runtime config the generated app
// itself will load (config/<project-name>.toml, unlike iluvatar.toml,
// which only drives generation), and, for every configured resource, a
// model/repository/service/.proto plus a migration. If any resources are
// configured, it also ensures dir is a Go module with every dependency
// the generated source needs resolved (`go mod init`/`go mod tidy`).
func Generate(dir string, cfg *config.Config) error {
	cloudProvider, err := normalizeCloudProvider(cfg.Cloud.Provider)
	if err != nil {
		return err
	}

	data := templateData{
		BinaryName:          cfg.Project.Name,
		Module:              cfg.Project.Module,
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
	if err := writeGitignore(dir); err != nil {
		return err
	}
	if err := writeCRUD(dir, cfg, data); err != nil {
		return err
	}

	if len(cfg.Resources) == 0 {
		return nil
	}
	return ensureGoModule(dir, cfg.Project.Module)
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
