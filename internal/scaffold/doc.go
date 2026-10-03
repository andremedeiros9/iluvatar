// Package scaffold creates a brand new iluvatar project directory.
//
// It backs the "iluvatar init <path>" command: given a destination path,
// it creates the directory (if needed) and renders the templates in
// ./templates into a starter iluvatar.toml, pre-filled with sensible
// defaults and ready for the user to edit before running
// "iluvatar generate".
package scaffold
