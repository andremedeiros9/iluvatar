package cli

import "github.com/spf13/cobra"

// NewRootCommand builds the "iluvatar" command tree.
func NewRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:           "iluvatar",
		Short:         "Generate RESTful Go APIs from a config file",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	root.AddCommand(newInitCommand())

	return root
}

// Execute runs the CLI with the process's actual args.
func Execute() error {
	return NewRootCommand().Execute()
}
