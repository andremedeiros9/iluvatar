package cli

import (
	"bytes"
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
