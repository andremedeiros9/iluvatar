package generator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/andremedeiros9/iluvatar/internal/config"
	"github.com/stretchr/testify/require"
)

func TestWriteDBUnsupportedDriver(t *testing.T) {
	dir := t.TempDir()
	data := templateData{BinaryName: "widgets"}

	err := writeDB(dir, "oracle", data)
	require.Error(t, err)
	require.Contains(t, err.Error(), "unsupported database driver")
}

func TestWriteDBDrivers(t *testing.T) {
	drivers := []string{"", "postgres", "mysql", "sqlite"}
	data := templateData{BinaryName: "widgets"}

	for _, driver := range drivers {
		t.Run(driver, func(t *testing.T) {
			dir := t.TempDir()
			err := writeDB(dir, driver, data)
			require.NoError(t, err, "driver %q should not error", driver)

			_, err = os.Stat(filepath.Join(dir, "internal", "db", "db.go"))
			require.NoError(t, err, "db.go should exist for driver %q", driver)
		})
	}
}

func TestWriteCRUDNoResourcesIsNoOp(t *testing.T) {
	dir := t.TempDir()
	cfg := &config.Config{
		Project: config.Project{Name: "widgets"},
	}

	err := writeCRUD(dir, cfg, templateData{})
	require.NoError(t, err, "writeCRUD with no resources should be a no-op")
}
