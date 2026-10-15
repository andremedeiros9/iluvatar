package cli

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewInitCommand(t *testing.T) {
	cmd := newInitCommand()

	require.Equal(t, "init <path>", cmd.Use)
	require.Equal(t, "Create a new iluvatar project at <path>", cmd.Short)
	require.NotNil(t, cmd.Args, "init should validate its arguments")

	// init takes exactly one positional argument (<path>)
	err := cmd.Args(cmd, []string{})
	require.Error(t, err, "init with no args should fail")

	err = cmd.Args(cmd, []string{"path1", "path2"})
	require.Error(t, err, "init with two args should fail")

	err = cmd.Args(cmd, []string{"my-project"})
	require.NoError(t, err, "init with exactly one arg should succeed")

	// --module flag should exist with a string default of ""
	flag := cmd.Flags().Lookup("module")
	require.NotNil(t, flag, "--module flag should be registered")
	require.Equal(t, "string", flag.Value.Type())
	require.Equal(t, "", flag.DefValue)
	require.Equal(t,
		"Go module path for the generated project (defaults to the directory name)",
		flag.Usage,
	)
}

func TestNewInitCommandHasRunE(t *testing.T) {
	cmd := newInitCommand()
	require.NotNil(t, cmd.RunE, "init should have a RunE handler")
}
