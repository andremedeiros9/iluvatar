// Package cli wires up the iluvatar command line interface.
//
// It defines the root "iluvatar" command and its subcommands, starting
// with "init" (write a fresh iluvatar.toml to a new project directory)
// and later "generate" (turn an edited iluvatar.toml into a working Go
// REST API). This package parses flags and args and delegates all real
// work to the config, scaffold, and generator packages — it should never
// contain file I/O or code generation logic itself.
package cli
