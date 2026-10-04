package generator

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// ensureGoModule makes dir a Go module for module (running `go mod init`
// if it isn't one yet), then runs `go mod tidy` so every import the
// generated source uses is resolved into go.mod/go.sum. It shells out to
// the Go toolchain already required to run iluvatar itself, rather than
// hand-templating version-pinned go.mod/go.sum content that would drift.
//
// It must run after generateProtos: every grpc resource imports its own
// <resource>/pb package, which only has .go source once protoc has
// compiled the resource's .proto.
func ensureGoModule(dir, module string) error {
	_, err := os.Stat(filepath.Join(dir, "go.mod"))
	switch {
	case err == nil:
		// already a module
	case os.IsNotExist(err):
		if err := runGo(dir, "mod", "init", module); err != nil {
			return err
		}
	default:
		return fmt.Errorf("checking for existing go.mod: %w", err)
	}

	return runGo(dir, "mod", "tidy")
}

func runGo(dir string, args ...string) error {
	cmd := exec.Command("go", args...)
	cmd.Dir = dir

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("go %v: %w\n%s", args, err, out)
	}
	return nil
}
