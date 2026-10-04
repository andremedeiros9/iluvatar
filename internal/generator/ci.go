package generator

import (
	"fmt"
	"path/filepath"
)

func writeCI(dir, provider string, data templateData) error {
	switch provider {
	case "", "github":
		return render(dir, filepath.Join(".github", "workflows", "ci.yml"), "templates/ci/github.yml.tmpl", data)
	case "gitlab":
		return render(dir, ".gitlab-ci.yml", "templates/ci/gitlab.yml.tmpl", data)
	default:
		return fmt.Errorf("unsupported ci provider %q (want %q or %q)", provider, "github", "gitlab")
	}
}
