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
// configured, it also compiles their .proto files with protoc (installing
// the latest protoc and its Go plugins first if the machine doesn't have
// them) and ensures dir is a Go module with every dependency the
// generated source needs resolved (`go mod init`/`go mod tidy`).
//
// Pass WithProgress to be told about each step as it starts.
func Generate(dir string, cfg *config.Config, opts ...Option) error {
	var progress Progress = stderrProgress{}
	for _, opt := range opts {
		opt(&progress)
	}

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

	type step struct {
		name string
		run  func() error
	}
	steps := []step{
		{"Writing Makefile", func() error { return writeMakefile(dir, data) }},
		{"Writing CI pipeline", func() error { return writeCI(dir, cfg.CI.Provider, data) }},
		{"Writing Dockerfile", func() error { return writeDockerfile(dir, data) }},
		{"Writing docker-compose", func() error { return writeCompose(dir, cfg.Database.Driver, data) }},
		{"Writing runtime config", func() error { return writeRuntimeConfig(dir, cfg) }},
		{"Writing .gitignore", func() error { return writeGitignore(dir) }},
	}
	if len(cfg.Resources) > 0 {
		steps = append(steps,
			step{"Writing CRUD code", func() error { return writeCRUD(dir, cfg, data) }},
			step{"Compiling protos", func() error { return generateProtos(dir, progress.Logf) }},
			step{"Resolving Go dependencies", func() error { return ensureGoModule(dir, cfg.Project.Module) }},
		)
	}

	for i, s := range steps {
		progress.Step(i, len(steps), s.name)
		if err := s.run(); err != nil {
			return err
		}
	}
	return nil
}

func normalizeCloudProvider(raw string) (string, error) {
	switch raw {
	case "", "none":
		return "none", nil
	case "aws":
		return raw, nil
	default:
		return "", fmt.Errorf("unsupported cloud provider %q (want %q or %q)", raw, "none", "aws")
	}
}
