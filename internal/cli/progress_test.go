package cli

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProgressBarDraws(t *testing.T) {
	var out bytes.Buffer
	bar := &progressBar{w: &out, enabled: true}

	bar.Step(1, 2, "Compiling protos")
	require.Equal(t, "\r["+strings.Repeat("█", 15)+strings.Repeat("░", 15)+"]  50% (1/2) Compiling protos", out.String())

	out.Reset()
	bar.Logf("installing %s", "protoc")
	require.Contains(t, out.String(), "\rinstalling protoc\n\r[", "message should replace the bar's line, then the bar is redrawn")

	out.Reset()
	bar.Done()
	require.Contains(t, out.String(), strings.Repeat("█", 30)+"] 100% (2/2) Done")
	require.True(t, strings.HasSuffix(out.String(), "\n"))
}

func TestProgressBarDisabledOnlyLogs(t *testing.T) {
	var out bytes.Buffer
	bar := newProgressBar(&out)

	bar.Step(1, 2, "Compiling protos")
	bar.Logf("installing %s", "protoc")
	bar.Done()

	require.Equal(t, "installing protoc\n", out.String())
}

func TestIsTerminalReturnsFalseForNonFileWriter(t *testing.T) {
	// bytes.Buffer is not an *os.File, so it can never be a terminal.
	require.False(t, isTerminal(&bytes.Buffer{}))
}

func TestIsTerminalReturnsFalseForRegularFile(t *testing.T) {
	// A regular file is not a character device, so not a terminal.
	f, err := os.CreateTemp(t.TempDir(), "not-a-terminal-*")
	require.NoError(t, err)
	defer func() {
		_ = f.Close()
	}()
	require.False(t, isTerminal(f), "a regular file should not be considered a terminal")
}

func TestProgressBarDrawSkipsWhenTotalIsZero(t *testing.T) {
	var out bytes.Buffer
	bar := &progressBar{w: &out, enabled: true, total: 0}

	bar.Step(0, 0, "Nothing")
	require.Empty(t, out.String(), "draw should not write anything when total is 0")
}

func TestProgressBarDrawSkipsWhenDisabled(t *testing.T) {
	var out bytes.Buffer
	bar := &progressBar{w: &out, enabled: false}

	bar.Step(1, 2, "Something")
	require.Empty(t, out.String(), "draw should not write anything when disabled")
}

func TestProgressBarStopClearsLine(t *testing.T) {
	var out bytes.Buffer
	bar := &progressBar{w: &out, enabled: true, lineLen: 10}

	bar.Stop()
	// Stop() ends the bar's line with a newline when lineLen > 0.
	require.Equal(t, "\n", out.String())
}

func TestProgressBarClearWritesSpaces(t *testing.T) {
	var out bytes.Buffer
	bar := &progressBar{w: &out, enabled: true, lineLen: 10}

	bar.clear()
	require.Equal(t, "\r          \r", out.String())
	require.Zero(t, bar.lineLen)
}
