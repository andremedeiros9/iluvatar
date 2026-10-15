package generator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWriteCIUnsupportedProvider(t *testing.T) {
	dir := t.TempDir()
	data := templateData{}

	err := writeCI(dir, "bitbucket", data)
	require.Error(t, err)
	require.Contains(t, err.Error(), "unsupported ci provider")
}
