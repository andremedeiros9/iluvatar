package cli

import (
	"errors"
	"log"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/andremedeiros9/iluvatar/internal/config"
	"github.com/andremedeiros9/iluvatar/internal/scaffold"
)

func newInitCommand() *cobra.Command {
	var module string

	cmd := &cobra.Command{
		Use:   "init <path>",
		Short: "Create a new iluvatar project at <path>",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]
			name := filepath.Base(filepath.Clean(path))

			mod := module
			if mod == "" {
				mod = name
			}

			if err := scaffold.Init(path, scaffold.Options{
				Name:   name,
				Module: mod,
			}); err != nil {
				return errors.New("failed to initialize project: " + err.Error())
			}

			log.New(cmd.OutOrStdout(), "", 0).Printf("Created %s\n", filepath.Join(path, config.FileName)) 

			return nil
		},
	}

	cmd.Flags().StringVar(&module, "module", "", "Go module path for the generated project (defaults to the directory name)")

	return cmd
}
