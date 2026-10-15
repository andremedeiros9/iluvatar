package generator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProtocInstallDir(t *testing.T) {
	home, err := os.UserHomeDir()
	require.NoError(t, err)

	got, err := protocInstallDir()
	require.NoError(t, err)
	require.Equal(t, filepath.Join(home, ".iluvatar", "protoc"), got)
}
