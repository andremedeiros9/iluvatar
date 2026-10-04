package generator

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProtocPlatform(t *testing.T) {
	type testCase struct {
		goos, goarch string
		want         string
	}
	tC := []testCase{
		{goos: "windows", goarch: "amd64", want: "win64"},
		{goos: "linux", goarch: "amd64", want: "linux-x86_64"},
		{goos: "linux", goarch: "arm64", want: "linux-aarch_64"},
		{goos: "darwin", goarch: "arm64", want: "osx-aarch_64"},
	}

	for _, tc := range tC {
		t.Run(tc.goos+"/"+tc.goarch, func(t *testing.T) {
			got, err := protocPlatform(tc.goos, tc.goarch)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	_, err := protocPlatform("plan9", "amd64")
	require.Error(t, err)
	require.Contains(t, err.Error(), "no prebuilt protoc release")
}

func TestExtractZip(t *testing.T) {
	archive := writeZip(t, map[string]string{
		"bin/protoc": "binary",
		"include/google/protobuf/timestamp.proto": "syntax",
	})
	dest := filepath.Join(t.TempDir(), "protoc")

	err := extractZip(archive, dest)
	require.NoError(t, err)

	data, err := os.ReadFile(filepath.Join(dest, "bin", "protoc"))
	require.NoError(t, err)
	require.Equal(t, "binary", string(data))

	data, err = os.ReadFile(filepath.Join(dest, "include", "google", "protobuf", "timestamp.proto"))
	require.NoError(t, err)
	require.Equal(t, "syntax", string(data))
}

func TestExtractZipRejectsPathsOutsideDest(t *testing.T) {
	archive := writeZip(t, map[string]string{"../escaped": "nope"})
	parent := t.TempDir()

	err := extractZip(archive, filepath.Join(parent, "protoc"))
	require.Error(t, err)
	require.Contains(t, err.Error(), "unsafe path")

	_, err = os.Stat(filepath.Join(parent, "escaped"))
	require.True(t, os.IsNotExist(err), "entry should not have been written outside dest")
}

func writeZip(t *testing.T, files map[string]string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "archive.zip")
	f, err := os.Create(path)
	require.NoError(t, err)

	w := zip.NewWriter(f)
	for name, content := range files {
		entry, err := w.Create(name)
		require.NoError(t, err)
		_, err = entry.Write([]byte(content))
		require.NoError(t, err)
	}
	require.NoError(t, w.Close())
	require.NoError(t, f.Close())

	return path
}
