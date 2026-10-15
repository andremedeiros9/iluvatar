package cli

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/andremedeiros9/iluvatar/internal/config"
	"github.com/andremedeiros9/iluvatar/internal/scaffold"
)

func newInitCommand() *cobra.Command {
	var module string
	var architecture string
	var nonInteractive bool

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

			opts := scaffold.Options{
				Name:             name,
				Module:           mod,
				Architecture:     architecture,
				Port:             8080,
				DatabaseDriver:   "postgres",
				CIProvider:       "github",
				CloudProvider:    "none",
				CloudRegion:      "us-east-1",
			}

			interactive := !nonInteractive && isStdinTerminal()
			if interactive {
				if err := promptOptions(cmd.InOrStdin(), &opts, name, mod); err != nil {
					return errors.New("failed to initialize project: " + err.Error())
				}
			}

			bar := newProgressBar(cmd.ErrOrStderr())
			bar.Step(0, 1, "Writing "+config.FileName)
			if err := scaffold.Init(path, opts); err != nil {
				bar.Stop()
				return errors.New("failed to initialize project: " + err.Error())
			}
			bar.Done()

			configPath := filepath.Join(path, config.FileName)
			if abs, err := filepath.Abs(configPath); err == nil {
				configPath = abs
			}
			log.New(cmd.OutOrStdout(), "", 0).Printf("Created %s\n", configPath)

			return nil
		},
	}

	cmd.Flags().StringVar(&module, "module", "", "Go module path for the generated project (defaults to the directory name)")
	cmd.Flags().StringVar(&architecture, "architecture", "vertical", "Code layout: \"vertical\" (default) or \"horizontal\"")
	cmd.Flags().BoolVar(&nonInteractive, "non-interactive", false, "Skip interactive prompts and use defaults/flags")

	return cmd
}

// promptOptions interactively asks the user for every iluvatar.toml value,
// filling opts with the answers. Defaults come from opts itself.
func promptOptions(in io.Reader, opts *scaffold.Options, name, moduleDefault string) error {
	reader := bufio.NewReader(in)

	fmt.Println("Let's set up your iluvatar project. Press enter to keep the value in brackets.")

	n, err := prompt(reader, "Project name", opts.Name)
	if err != nil {
		return err
	}
	opts.Name = n

	m, err := prompt(reader, "Go module path", moduleDefault)
	if err != nil {
		return err
	}
	opts.Module = m

	arch, err := promptSelect(reader, "Code architecture", architectureOptions, 0)
	if err != nil {
		return err
	}
	opts.Architecture = arch

	portStr, err := prompt(reader, "Server port", strconv.Itoa(opts.Port))
	if err != nil {
		return err
	}
	if portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil && p > 0 {
			opts.Port = p
		}
	}

	driver, err := promptSelect(reader, "Database driver", databaseDriverOptions, 0)
	if err != nil {
		return err
	}
	opts.DatabaseDriver = driver

	ci, err := promptSelect(reader, "CI provider", ciProviderOptions, 0)
	if err != nil {
		return err
	}
	opts.CIProvider = ci

	cloud, err := promptSelect(reader, "Cloud provider", cloudProviderOptions, 0)
	if err != nil {
		return err
	}
	opts.CloudProvider = cloud

	if opts.CloudProvider == "aws" {
		region, err := prompt(reader, "AWS region", opts.CloudRegion)
		if err != nil {
			return err
		}
		opts.CloudRegion = region
	}

	addResources, err := promptBool(reader, "Add a CRUD resource now?", false)
	if err != nil {
		return err
	}
	if addResources {
		if err := promptResource(reader, opts.Module); err != nil {
			return err
		}
	}

	return nil
}

func promptResource(reader *bufio.Reader, module string) error {
	rname, err := promptNonEmpty(reader, "Resource name (e.g. user)", "")
	if err != nil {
		return err
	}

	comm, err := promptSelect(reader, "Communication protocol", communicationOptions, 0)
	if err != nil {
		return err
	}

	fmt.Printf("Resource %q (%s) noted. You can add more resources and fields in iluvatar.toml.\n", rname, comm)
	_ = module
	return nil
}
