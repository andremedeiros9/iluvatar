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
// grpcserver bootstrap, the database bootstrap, cmd/<name>/main.go, and
// one migration pair per resource. It's a no-op if cfg has no resources.
//
// Every resource must declare communication = "grpc" (rest/rpc aren't
// implemented yet) and every field must use a supported type; both are
// validated here.
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
			return err
		}
		resources = append(resources, rd)

		if err := writeResourceFiles(dir, rd); err != nil {
			return err
		}
		if err := writeMigration(dir, i+1, rd); err != nil {
			return err
		}
	}

	if err := writeGRPCServer(dir, cfg.Project.Module, resources); err != nil {
		return err
	}
	if err := writeDB(dir, cfg.Database.Driver, data); err != nil {
		return err
	}
	return writeMain(dir, data)
}

func writeResourceFiles(dir string, rd resourceData) error {
	base := filepath.Join("internal", rd.PackageName)

	if err := renderGo(dir, filepath.Join(base, "model.go"), "templates/resource/model.go.tmpl", rd); err != nil {
		return err
	}
	if err := renderGo(dir, filepath.Join(base, "repository.go"), "templates/resource/repository.go.tmpl", rd); err != nil {
		return err
	}
	if err := renderGo(dir, filepath.Join(base, "service.go"), "templates/resource/service.go.tmpl", rd); err != nil {
		return err
	}
	return render(dir, filepath.Join(base, "pb", rd.Name+".proto"), "templates/resource/resource.proto.tmpl", rd)
}

func writeMigration(dir string, seq int, rd resourceData) error {
	name := fmt.Sprintf("%04d_create_%s", seq, rd.TableName)

	if err := render(dir, filepath.Join("migrations", name+".up.sql"), "templates/migrations/create_resource.up.sql.tmpl", rd); err != nil {
		return err
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
	return renderGo(dir, filepath.Join("cmd", data.BinaryName, "main.go"), "templates/cmd/main.go.tmpl", data)
}
