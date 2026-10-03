// Command iluvatar is the CLI entry point for the Iluvatar code generator.
// It is responsible for parsing the command line arguments and invoking the appropriate subcommand,
// which will then call the internal scaffold package to generate the code.
package main

import (
	"fmt"
	"os"

	"github.com/andremedeiros9/iluvatar/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
