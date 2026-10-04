package generator

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// protocReleasesURL is where protoc's prebuilt release archives live.
const protocReleasesURL = "https://github.com/protocolbuffers/protobuf/releases"

// protocPlugins are the protoc plugins needed to turn a resource's .proto
// into its *.pb.go/*_grpc.pb.go files, and the package `go install` builds
// each of them from.
var protocPlugins = []struct{ name, pkg string }{
	{"protoc-gen-go", "google.golang.org/protobuf/cmd/protoc-gen-go@latest"},
	{"protoc-gen-go-grpc", "google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest"},
}

// generateProtos compiles every resource's .proto under dir into its
// *.pb.go/*_grpc.pb.go files, the same way the generated Makefile's
// `make proto` does. protoc and its Go plugins are installed first if the
// machine doesn't have them yet.
func generateProtos(dir string, logf func(format string, args ...any)) error {
	protos, err := filepath.Glob(filepath.Join(dir, "internal", "*", "pb", "*.proto"))
	if err != nil {
		return fmt.Errorf("looking for .proto files: %w", err)
	}
	if len(protos) == 0 {
		return nil
	}

	protoc, err := ensureProtoc(logf)
	if err != nil {
		return err
	}

	args := []string{
		"--go_out=.", "--go_opt=paths=source_relative",
		"--go-grpc_out=.", "--go-grpc_opt=paths=source_relative",
	}
	for _, plugin := range protocPlugins {
		pluginPath, err := ensureProtocPlugin(plugin.name, plugin.pkg, logf)
		if err != nil {
			return err
		}
		args = append(args, "--plugin="+plugin.name+"="+pluginPath)
	}
	for _, proto := range protos {
		rel, err := filepath.Rel(dir, proto)
		if err != nil {
			return fmt.Errorf("resolving %s: %w", proto, err)
		}
		args = append(args, filepath.ToSlash(rel))
	}

	cmd := exec.Command(protoc, args...)
	cmd.Dir = dir

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("protoc %v: %w\n%s", args, err, out)
	}
	return nil
}

// ensureProtoc returns the path to a protoc binary: the one on PATH if
// there is one, otherwise the one iluvatar installed on a previous run,
// otherwise a freshly installed latest release.
func ensureProtoc(logf func(format string, args ...any)) (string, error) {
	if p, err := exec.LookPath("protoc"); err == nil {
		return p, nil
	}

	root, err := protocInstallDir()
	if err != nil {
		return "", err
	}
	if p, err := exec.LookPath(filepath.Join(root, "bin", "protoc")); err == nil {
		return p, nil
	}

	if err := installProtoc(root, logf); err != nil {
		return "", fmt.Errorf("installing protoc: %w", err)
	}

	p, err := exec.LookPath(filepath.Join(root, "bin", "protoc"))
	if err != nil {
		return "", fmt.Errorf("protoc not found after installing it into %s: %w", root, err)
	}
	return p, nil
}

// protocInstallDir is where iluvatar installs protoc when the machine
// doesn't have one: ~/.iluvatar/protoc, holding the release's bin/ and
// include/ side by side, which is the layout protoc needs to find the
// well-known types (google/protobuf/*.proto) on its own.
func protocInstallDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("locating home directory: %w", err)
	}
	return filepath.Join(home, ".iluvatar", "protoc"), nil
}

// installProtoc downloads the latest protoc release for this OS/arch and
// unpacks it into root.
func installProtoc(root string, logf func(format string, args ...any)) error {
	platform, err := protocPlatform(runtime.GOOS, runtime.GOARCH)
	if err != nil {
		return err
	}

	version, err := latestProtocVersion()
	if err != nil {
		return err
	}

	logf("protoc not found, installing v%s into %s", version, root)

	url := fmt.Sprintf("%s/download/v%s/protoc-%s-%s.zip", protocReleasesURL, version, version, platform)
	archive, err := download(url)
	if err != nil {
		return err
	}
	defer func() {
		_ = os.Remove(archive)
	}()

	if err := os.RemoveAll(root); err != nil {
		return fmt.Errorf("clearing %s: %w", root, err)
	}
	if err := extractZip(archive, root); err != nil {
		return err
	}

	logf("installed protoc v%s; add %s to your PATH to use it outside iluvatar (e.g. `make proto`)", version, filepath.Join(root, "bin"))
	return nil
}

// protocPlatform maps a GOOS/GOARCH pair to the platform suffix protoc's
// release archives are named with (protoc-<version>-<platform>.zip).
func protocPlatform(goos, goarch string) (string, error) {
	switch goos + "/" + goarch {
	case "windows/amd64":
		return "win64", nil
	case "windows/386":
		return "win32", nil
	case "linux/amd64":
		return "linux-x86_64", nil
	case "linux/386":
		return "linux-x86_32", nil
	case "linux/arm64":
		return "linux-aarch_64", nil
	case "darwin/amd64":
		return "osx-x86_64", nil
	case "darwin/arm64":
		return "osx-aarch_64", nil
	default:
		return "", fmt.Errorf("no prebuilt protoc release for %s/%s; install protoc manually", goos, goarch)
	}
}

// latestProtocVersion resolves protoc's latest release version (without
// the leading "v") from where GitHub redirects releases/latest to, which,
// unlike the GitHub API, isn't rate limited for anonymous callers.
func latestProtocVersion() (string, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Get(protocReleasesURL + "/latest")
	if err != nil {
		return "", fmt.Errorf("resolving latest protoc release: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	location := resp.Header.Get("Location")
	version := strings.TrimPrefix(path.Base(location), "v")
	if location == "" || version == "" || !strings.Contains(location, "/releases/tag/") {
		return "", fmt.Errorf("resolving latest protoc release: unexpected response %s (Location %q)", resp.Status, location)
	}
	return version, nil
}

// download fetches url into a temporary file and returns its path; the
// caller removes it.
func download(url string) (string, error) {
	client := &http.Client{Timeout: 10 * time.Minute}

	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("downloading %s: %w", url, err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("downloading %s: %s", url, resp.Status)
	}

	f, err := os.CreateTemp("", "iluvatar-protoc-*.zip")
	if err != nil {
		return "", fmt.Errorf("creating temporary file: %w", err)
	}

	_, err = io.Copy(f, resp.Body)
	if closeErr := f.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		_ = os.Remove(f.Name())
		return "", fmt.Errorf("downloading %s: %w", url, err)
	}
	return f.Name(), nil
}

// extractZip unpacks the zip archive at src into dest, refusing entries
// that would land outside dest.
func extractZip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return fmt.Errorf("opening %s: %w", src, err)
	}
	defer func() {
		_ = r.Close()
	}()

	for _, f := range r.File {
		if !filepath.IsLocal(f.Name) {
			return fmt.Errorf("extracting %s: unsafe path %q", src, f.Name)
		}
		target := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(target, 0o755); err != nil {
				return fmt.Errorf("creating %s: %w", target, err)
			}
			continue
		}
		if err := extractZipFile(f, target); err != nil {
			return fmt.Errorf("extracting %s: %w", f.Name, err)
		}
	}
	return nil
}

func extractZipFile(f *zip.File, target string) error {
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return err
	}

	in, err := f.Open()
	if err != nil {
		return err
	}
	defer func() {
		_ = in.Close()
	}()

	// Keep the archive's mode bits so bin/protoc stays executable.
	out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.Mode().Perm()|0o644)
	if err != nil {
		return err
	}

	_, err = io.Copy(out, in)
	if closeErr := out.Close(); err == nil {
		err = closeErr
	}
	return err
}

// ensureProtocPlugin returns the path to the protoc plugin called name,
// running `go install pkg` first if it isn't installed yet. The path is
// handed to protoc explicitly, so the plugin doesn't need to be on PATH.
func ensureProtocPlugin(name, pkg string, logf func(format string, args ...any)) (string, error) {
	if p, err := exec.LookPath(name); err == nil {
		return p, nil
	}

	bin, err := goBinDir()
	if err != nil {
		return "", err
	}
	if p, err := exec.LookPath(filepath.Join(bin, name)); err == nil {
		return p, nil
	}

	logf("%s not found, installing %s", name, pkg)
	if err := runGo("", "install", pkg); err != nil {
		return "", err
	}

	p, err := exec.LookPath(filepath.Join(bin, name))
	if err != nil {
		return "", fmt.Errorf("%s not found after installing it into %s: %w", name, bin, err)
	}
	return p, nil
}

// goBinDir is the directory `go install` writes binaries to: GOBIN, or
// the first GOPATH entry's bin/ when GOBIN isn't set.
func goBinDir() (string, error) {
	out, err := exec.Command("go", "env", "GOBIN", "GOPATH").Output()
	if err != nil {
		return "", fmt.Errorf("go env GOBIN GOPATH: %w", err)
	}

	lines := strings.Split(strings.ReplaceAll(string(out), "\r\n", "\n"), "\n")
	if gobin := strings.TrimSpace(lines[0]); gobin != "" {
		return gobin, nil
	}
	if len(lines) > 1 {
		if gopaths := filepath.SplitList(strings.TrimSpace(lines[1])); len(gopaths) > 0 && gopaths[0] != "" {
			return filepath.Join(gopaths[0], "bin"), nil
		}
	}
	return "", errors.New("could not determine where `go install` puts binaries (GOBIN and GOPATH are both empty)")
}
