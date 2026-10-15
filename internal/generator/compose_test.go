package generator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWriteComposeUnsupportedDriver(t *testing.T) {
	dir := t.TempDir()
	data := templateData{}

	err := writeCompose(dir, "oracle", data)
	require.Error(t, err)
	require.Contains(t, err.Error(), "unsupported database driver")
}

func TestWriteComposeSkippedForSQLite(t *testing.T) {
	dir := t.TempDir()
	data := templateData{}

	err := writeCompose(dir, "sqlite", data)
	require.NoError(t, err)

	_, err = os.Stat(filepath.Join(dir, "docker-compose.yml"))
	require.True(t, os.IsNotExist(err), "docker-compose.yml should not exist for sqlite")

	_, err = os.Stat(filepath.Join(dir, ".env.example"))
	require.True(t, os.IsNotExist(err), ".env.example should not exist for sqlite")
}
