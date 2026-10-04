package cli

import (
	"errors"
	"log"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/andremedeiros9/iluvatar/internal/config"
	"github.com/andremedeiros9/iluvatar/internal/generator"
)

func newGenerateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate [path]",
		Short: "Generate a project's build tooling from its iluvatar.toml",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := "."
			if len(args) == 1 {
				dir = args[0]
			}

			cfg, err := config.Load(filepath.Join(dir, config.FileName))
			if err != nil {
				return errors.New("failed to load config: " + err.Error())
			}

			if err := generator.Generate(dir, cfg); err != nil {
				return errors.New("failed to generate project: " + err.Error())
			}

			// escrever o diretorio desde da root
			log.New(cmd.OutOrStdout(), "", 0).Printf("Generated project in %s\n", dir)
			return nil
		},
	}

	return cmd
}
