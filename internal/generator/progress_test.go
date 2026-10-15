package generator

import "testing"

func TestStderrProgress(t *testing.T) {
	var p Progress = stderrProgress{}

	// stderrProgress should not panic and should be a no-op
	p.Step(0, 1, "Writing Makefile")
	p.Logf("installing %s", "protoc")
}
