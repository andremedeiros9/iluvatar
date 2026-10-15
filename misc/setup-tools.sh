#!/usr/bin/env bash
set -uo pipefail

# Installs the development tools used by `make lint` and `make vulns`:
# golangci-lint and govulncheck.
#
# Expects go to be on PATH (or GO_INSTALL_DIR/bin prepended by the caller).
#
# If the network is unavailable, each tool that can't be downloaded is
# reported individually and the script exits non-zero at the end.

GO_INSTALL_DIR="${GO_INSTALL_DIR:-/usr/local/go}"

TOOLS=(
    "github.com/golangci/golangci-lint/v2/cmd/golangci-lint"
    "golang.org/x/vuln/cmd/govulncheck"
)

FAILED=0

for tool in "${TOOLS[@]}"; do
    name="$(basename "$tool")"
    echo "Installing $name..."
    if PATH="${GO_INSTALL_DIR}/bin:${PATH}" GOPROXY="${GOPROXY}" go install "$tool@latest" 2>/dev/null; then
        echo "  $name installed."
        continue
    fi

    # Fall back to direct fetch (useful when the proxy is unreachable,
    # e.g. IPv6-only routing in WSL).
    echo "  $name: proxy failed, trying direct..."
    if PATH="${GO_INSTALL_DIR}/bin:${PATH}" GOPROXY=direct go install "$tool@latest" 2>/dev/null; then
        echo "  $name installed (direct)."
        continue
    fi

    echo "  $name: failed to install (network/proxy issue?)." >&2
    FAILED=1
done

if [ "$FAILED" -ne 0 ]; then
    echo >&2
    echo "Some tools could not be installed. Make sure you have network access" >&2
    echo "to proxy.golang.org (or github.com for direct), or override GOPROXY:" >&2
    echo "  make setup GOPROXY=https://goproxy.io" >&2
    echo "Then re-run: make setup-tools" >&2
    exit 1
fi

echo "Development tools installed."
