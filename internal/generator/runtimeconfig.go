package generator

import (
	"path/filepath"

	"github.com/andremedeiros9/iluvatar/internal/config"
)

// writeRuntimeConfig snapshots cfg into config/<project-name>.toml: the
// file the generated application reads at runtime, as opposed to
// iluvatar.toml, which only drives generation.
func writeRuntimeConfig(dir string, cfg *config.Config) error {
	path := filepath.Join(dir, "config", cfg.Project.Name+".toml")
	return config.Save(path, cfg)
}
