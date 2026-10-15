package cli

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewRootCommand(t *testing.T) {
	root := NewRootCommand()

	require.Equal(t, "iluvatar", root.Use)
	require.Equal(t, "Generate RESTful Go APIs from a config file", root.Short)
	require.True(t, root.SilenceUsage, "SilenceUsage should be true so the CLI controls error output")
	require.True(t, root.SilenceErrors, "SilenceErrors should be true so the CLI controls error output")

	// init subcommand must be registered
	initCmd, _, err := root.Find([]string{"init"})
	require.NoError(t, err)
	require.NotNil(t, initCmd, "init subcommand should be registered")
	require.Equal(t, "init <path>", initCmd.Use)

	// generate subcommand must be registered
	genCmd, _, err := root.Find([]string{"generate"})
	require.NoError(t, err)
	require.NotNil(t, genCmd, "generate subcommand should be registered")
	require.Equal(t, "generate [path]", genCmd.Use)
}

func TestRootCommandNoArgsPrintsHelp(t *testing.T) {
	root := NewRootCommand()
	root.SetArgs([]string{})

	var out bytes.Buffer
	root.SetOut(&out)

	_, err := root.ExecuteC()
	require.NoError(t, err, "running iluvatar with no subcommand should not error (SilenceErrors is true)")
	require.Contains(t, out.String(), "Available Commands", "help should be printed when no subcommand is given")
}

func TestRootCommandHelpDoesNotError(t *testing.T) {
	root := NewRootCommand()
	root.SetArgs([]string{"--help"})

	// --help should print help and exit with code 0, not return an error
	_, err := root.ExecuteC()
	require.NoError(t, err, "--help should not return an error")
}
