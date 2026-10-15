package generator

import (
	"fmt"
	"path/filepath"

	"github.com/andremedeiros9/iluvatar/internal/config"
)

// projectData is the set of values grpcserver/server.go renders from: it
// needs every resource at once, to register all their services on the
// same server.
type projectData struct {
	Module    string
	Resources []resourceData
}

// writeCRUD generates the gRPC CRUD implementation for every resource in
// cfg: per-resource model/repository/service/.proto files, the shared
// server bootstrap, the database bootstrap, cmd/<name>/main.go, and one
// migration pair per resource. It's a no-op if cfg has no resources.
//
// Every resource must declare communication = "grpc" (rest/rpc aren't
// implemented yet) and every field must use a supported type; both are
// validated here.
//
// When data.Architecture is "horizontal", resources are laid out in
// internal/domain, internal/services, internal/rpc, and internal/api
// instead of the vertical internal/<resource> + internal/grpcserver +
// internal/db layout.
func writeCRUD(dir string, cfg *config.Config, data templateData) error {
	if len(cfg.Resources) == 0 {
		return nil
	}

	resources := make([]resourceData, 0, len(cfg.Resources))
	for i, res := range cfg.Resources {
		if res.Communication != "grpc" {
			return fmt.Errorf("resource %q: communication %q is not implemented yet (only %q is supported)", res.Name, res.Communication, "grpc")
		}

		rd, err := buildResourceData(cfg.Project.Module, cfg.Database.Driver, res)
		if err != nil {
			return fmt.Errorf("could not build resource data: %w", err)
		}
		resources = append(resources, rd)

		if err := writeResourceFiles(dir, rd, data.Architecture); err != nil {
			return fmt.Errorf("could not write resource files: %w", err)
		}
		if err := writeMigration(dir, i+1, rd); err != nil {
			return fmt.Errorf("could not write migration: %w", err)
		}
	}

	if data.Architecture == config.ArchitectureHorizontal {
		if err := writeAPIServer(dir, cfg.Project.Module, resources); err != nil {
			return fmt.Errorf("could not write API server: %w", err)
		}
		if err := writeAPIDB(dir, cfg.Database.Driver, data); err != nil {
			return fmt.Errorf("could not write database files: %w", err)
		}
	} else {
		if err := writeGRPCServer(dir, cfg.Project.Module, resources); err != nil {
			return fmt.Errorf("could not write gRPC server: %w", err)
		}
		if err := writeDB(dir, cfg.Database.Driver, data); err != nil {
			return fmt.Errorf("could not write database files: %w", err)
		}
	}
	return writeMain(dir, data)
}

func writeResourceFiles(dir string, rd resourceData, architecture string) error {
	if architecture == config.ArchitectureHorizontal {
		return writeHorizontalResourceFiles(dir, rd)
	}
	base := filepath.Join("internal", rd.PackageName)

	if err := renderGo(dir, filepath.Join(base, "model.go"), "templates/resource/model.go.tmpl", rd); err != nil {
		return fmt.Errorf("could not render model.go: %w", err)
	}
	if err := renderGo(dir, filepath.Join(base, "repository.go"), "templates/resource/repository.go.tmpl", rd); err != nil {
		return fmt.Errorf("could not render repository.go: %w", err)
	}
	if err := renderGo(dir, filepath.Join(base, "service.go"), "templates/resource/service.go.tmpl", rd); err != nil {
		return fmt.Errorf("could not render service.go: %w", err)
	}
	return render(dir, filepath.Join(base, "pb", rd.Name+".proto"), "templates/resource/resource.proto.tmpl", rd)
}

func writeHorizontalResourceFiles(dir string, rd resourceData) error {
	// internal/domain/<resource>/model.go  — entity
	if err := renderGo(dir, filepath.Join("internal", "domain", rd.PackageName, "model.go"), "templates/horizontal/domain/model.go.tmpl", rd); err != nil {
		return fmt.Errorf("could not render domain model.go: %w", err)
	}
	// internal/domain/<resource>/repository.go — repository interface
	if err := renderGo(dir, filepath.Join("internal", "domain", rd.PackageName, "repository.go"), "templates/horizontal/domain/repository.go.tmpl", rd); err != nil {
		return fmt.Errorf("could not render domain repository.go: %w", err)
	}
	// internal/services/<resource>/repository.go — concrete sqlx repository
	if err := renderGo(dir, filepath.Join("internal", "services", rd.PackageName, "repository.go"), "templates/horizontal/services/repository.go.tmpl", rd); err != nil {
		return fmt.Errorf("could not render services repository.go: %w", err)
	}
	// internal/services/<resource>/service.go — application service
	if err := renderGo(dir, filepath.Join("internal", "services", rd.PackageName, "service.go"), "templates/horizontal/services/service.go.tmpl", rd); err != nil {
		return fmt.Errorf("could not render services service.go: %w", err)
	}
	// internal/rpc/<resource>/pb/<resource>.proto — proto definition
	return render(dir, filepath.Join("internal", "rpc", rd.PackageName, "pb", rd.Name+".proto"), "templates/horizontal/rpc/resource.proto.tmpl", rd)
}

func writeAPIServer(dir, module string, resources []resourceData) error {
	projData := projectData{Module: module, Resources: resources}
	return renderGo(dir, filepath.Join("internal", "api", "server.go"), "templates/horizontal/api/server.go.tmpl", projData)
}

func writeAPIDB(dir, driver string, data templateData) error {
	switch driver {
	case "", "postgres":
		return renderGo(dir, filepath.Join("internal", "api", "db.go"), "templates/horizontal/api/db.go.tmpl", data)
	case "mysql":
		return renderGo(dir, filepath.Join("internal", "api", "db.go"), "templates/horizontal/api/mysql.go.tmpl", data)
	case "sqlite":
		return renderGo(dir, filepath.Join("internal", "api", "db.go"), "templates/horizontal/api/sqlite.go.tmpl", data)
	default:
		return fmt.Errorf("unsupported database driver %q (want %q, %q, or %q)", driver, "postgres", "mysql", "sqlite")
	}
}

func writeMigration(dir string, seq int, rd resourceData) error {
	name := fmt.Sprintf("%04d_create_%s", seq, rd.TableName)

	if err := render(dir, filepath.Join("migrations", name+".up.sql"), "templates/migrations/create_resource.up.sql.tmpl", rd); err != nil {
		return fmt.Errorf("could not render migration up.sql: %w", err)
	}
	return render(dir, filepath.Join("migrations", name+".down.sql"), "templates/migrations/create_resource.down.sql.tmpl", rd)
}

func writeGRPCServer(dir, module string, resources []resourceData) error {
	data := projectData{Module: module, Resources: resources}
	return renderGo(dir, filepath.Join("internal", "grpcserver", "server.go"), "templates/grpcserver/server.go.tmpl", data)
}

func writeDB(dir, driver string, data templateData) error {
	switch driver {
	case "", "postgres":
		return renderGo(dir, filepath.Join("internal", "db", "db.go"), "templates/db/postgres.go.tmpl", data)
	case "mysql":
		return renderGo(dir, filepath.Join("internal", "db", "db.go"), "templates/db/mysql.go.tmpl", data)
	case "sqlite":
		return renderGo(dir, filepath.Join("internal", "db", "db.go"), "templates/db/sqlite.go.tmpl", data)
	default:
		return fmt.Errorf("unsupported database driver %q (want %q, %q, or %q)", driver, "postgres", "mysql", "sqlite")
	}
}

func writeMain(dir string, data templateData) error {
	tmpl := "templates/cmd/main.go.tmpl"
	if data.Architecture == config.ArchitectureHorizontal {
		tmpl = "templates/horizontal/cmd/main.go.tmpl"
	}
	return renderGo(dir, filepath.Join("cmd", data.BinaryName, "main.go"), tmpl, data)
}
