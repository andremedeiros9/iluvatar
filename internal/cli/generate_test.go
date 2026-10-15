package cli

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewGenerateCommand(t *testing.T) {
	cmd := newGenerateCommand()

	require.Equal(t, "generate [path]", cmd.Use)
	require.Equal(t, "Generate a project's build tooling from its iluvatar.toml", cmd.Short)
	require.NotNil(t, cmd.Args, "generate should validate its arguments")

	// generate takes zero or one positional argument ([path])
	err := cmd.Args(cmd, []string{})
	require.NoError(t, err, "generate with no args should succeed")

	err = cmd.Args(cmd, []string{"."})
	require.NoError(t, err, "generate with one arg should succeed")

	err = cmd.Args(cmd, []string{"path1", "path2"})
	require.Error(t, err, "generate with two args should fail")
}

func TestNewGenerateCommandHasRunE(t *testing.T) {
	cmd := newGenerateCommand()
	require.NotNil(t, cmd.RunE, "generate should have a RunE handler")
}
