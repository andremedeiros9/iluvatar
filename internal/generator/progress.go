package generator

import (
	"fmt"
	"os"
)

// Progress is told how Generate is getting on, so a caller can show it.
type Progress interface {
	// Step is called right before each step starts: done of total steps
	// have finished, and name describes the one starting now.
	Step(done, total int, name string)
	// Logf reports something the user should know about mid-step, such
	// as a missing tool being installed.
	Logf(format string, args ...any)
}

// Option customizes a Generate call.
type Option func(*Progress)

// WithProgress makes Generate report its progress to p.
func WithProgress(p Progress) Option {
	return func(progress *Progress) {
		*progress = p
	}
}

// stderrProgress is the Progress used when the caller doesn't pass one:
// steps go unreported, and messages are written to stderr.
type stderrProgress struct{}

func (stderrProgress) Step(int, int, string) {}

func (stderrProgress) Logf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
}
