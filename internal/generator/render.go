package generator

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"text/template"
)

// render executes the template at templatePath and writes the result to
// filepath.Join(dir, relPath), creating any missing parent directories.
func render(dir, relPath, templatePath string, data any) error {
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

// renderGo works like render, but gofmt-formats the executed template
// before writing it. That means template whitespace never matters, and a
// template that produces invalid Go fails loudly at generation time
// instead of silently shipping broken source.
func renderGo(dir, relPath, templatePath string, data any) error {
	tmpl, err := template.ParseFS(templatesFS, templatePath)
	if err != nil {
		return fmt.Errorf("parsing %s: %w", templatePath, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("executing %s: %w", templatePath, err)
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("formatting %s (generated from %s): %w\n\n%s", relPath, templatePath, err, buf.String())
	}

	outPath := filepath.Join(dir, relPath)
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return fmt.Errorf("creating directory for %s: %w", relPath, err)
	}

	if err := os.WriteFile(outPath, formatted, 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", relPath, err)
	}

	return nil
}
