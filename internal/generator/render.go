package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

// render executes the template at templatePath and writes the result to
// filepath.Join(dir, relPath), creating any missing parent directories.
func render(dir, relPath, templatePath string, data templateData) error {
	tmpl, err := template.ParseFS(templatesFS, templatePath)
	if err != nil {
		return fmt.Errorf("parsing %s: %w", templatePath, err)
	}

	outPath := filepath.Join(dir, relPath)
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return fmt.Errorf("creating directory for %s: %w", relPath, err)
	}

	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("creating %s: %w", relPath, err)
	}
	defer func() {
		_ = f.Close()
	}()

	if err := tmpl.Execute(f, data); err != nil {
		return fmt.Errorf("writing %s: %w", relPath, err)
	}

	return nil
}
