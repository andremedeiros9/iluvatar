package generator

import "fmt"

// writeCompose generates docker-compose.yml and .env.example for driver.
// sqlite has no separate server to orchestrate, so nothing is written for
// it.
func writeCompose(dir, driver string, data templateData) error {
	switch driver {
	case "sqlite":
		return nil
	case "", "postgres":
		if err := render(dir, "docker-compose.yml", "templates/compose/postgres.yml.tmpl", data); err != nil {
			return fmt.Errorf("could not render docker-compose.yml for postgres: %w", err)
		}
		return render(dir, ".env.example", "templates/compose/postgres.env.tmpl", data)
	case "mysql":
		if err := render(dir, "docker-compose.yml", "templates/compose/mysql.yml.tmpl", data); err != nil {
			return fmt.Errorf("could not render docker-compose.yml for mysql: %w", err)
		}
		return render(dir, ".env.example", "templates/compose/mysql.env.tmpl", data)
	default:
		return fmt.Errorf("unsupported database driver %q (want %q, %q, or %q)", driver, "postgres", "mysql", "sqlite")
	}
}
