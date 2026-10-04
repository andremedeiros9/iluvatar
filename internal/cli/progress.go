package cli

import (
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

// progressBarWidth is how many cells the bar itself spans.
const progressBarWidth = 30

// progressBar draws a single-line progress bar that redraws itself in
// place. It only draws when writing to a terminal; piped or redirected
// output gets no bar, just the messages passed to Logf.
type progressBar struct {
	w       io.Writer
	enabled bool

	done, total int
	name        string
	// lineLen is the width of the bar currently on screen, so the next
	// draw can blank out whatever it doesn't overwrite.
	lineLen int
}

func newProgressBar(w io.Writer) *progressBar {
	return &progressBar{w: w, enabled: isTerminal(w)}
}

func isTerminal(w io.Writer) bool {
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	info, err := f.Stat()
	return err == nil && info.Mode()&os.ModeCharDevice != 0
}

// Step redraws the bar with done of total steps finished and name as the
// step now running.
func (b *progressBar) Step(done, total int, name string) {
	b.done, b.total, b.name = done, total, name
	b.draw()
}

// Logf prints a message on its own line above the bar.
func (b *progressBar) Logf(format string, args ...any) {
	b.clear()
	_, _ = fmt.Fprintf(b.w, format+"\n", args...)
	b.draw()
}

// Done fills the bar and ends its line.
func (b *progressBar) Done() {
	b.Step(b.total, b.total, "Done")
	b.Stop()
}

// Stop ends the bar's line as it stands, so whatever is printed next
// (an error, say) starts on a fresh one.
func (b *progressBar) Stop() {
	if b.lineLen > 0 {
		_, _ = fmt.Fprintln(b.w)
		b.lineLen = 0
	}
}

func (b *progressBar) draw() {
	if !b.enabled || b.total == 0 {
		return
	}

	filled := progressBarWidth * b.done / b.total
	line := fmt.Sprintf("[%s%s] %3d%% (%d/%d) %s",
		strings.Repeat("█", filled),
		strings.Repeat("░", progressBarWidth-filled),
		100*b.done/b.total, b.done, b.total, b.name,
	)

	lineLen := utf8.RuneCountInString(line)
	padding := strings.Repeat(" ", max(b.lineLen-lineLen, 0))
	_, _ = fmt.Fprint(b.w, "\r"+line+padding)
	b.lineLen = lineLen
}

func (b *progressBar) clear() {
	if b.lineLen > 0 {
		_, _ = fmt.Fprint(b.w, "\r"+strings.Repeat(" ", b.lineLen)+"\r")
		b.lineLen = 0
	}
}
