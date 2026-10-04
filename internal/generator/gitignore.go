package generator

import (
	"fmt"
	"os"
	"path/filepath"
)

// writeGitignore copies the static gitignore template to .gitignore. It
// has no placeholders, so it's written directly rather than through
// render.
func writeGitignore(dir string) error {
	data, err := templatesFS.ReadFile("templates/gitignore.tmpl")
	if err != nil {
		return fmt.Errorf("reading gitignore template: %w", err)
	}

	if err := os.WriteFile(filepath.Join(dir, ".gitignore"), data, 0o644); err != nil {
		return fmt.Errorf("writing .gitignore: %w", err)
	}

	return nil
}
