#!/usr/bin/env bash
set -euo pipefail

# Installs the Go toolchain if it isn't already on PATH.
# Honors GO_VERSION and GO_INSTALL_DIR (defaults: 1.27.1, /usr/local/go).

GO_VERSION="${GO_VERSION:-1.27.1}"
GO_INSTALL_DIR="${GO_INSTALL_DIR:-/usr/local/go}"

if command -v go >/dev/null 2>&1; then
    echo "Go already installed: $(go version)"
    exit 0
fi

GO_OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
GO_ARCH="$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/')"

case "$GO_OS" in
    windows)
        echo "Please install Go manually from https://go.dev/dl/ and re-run make setup." >&2
        exit 1
        ;;
esac

GO_TARBALL="go${GO_VERSION}.${GO_OS}-${GO_ARCH}.tar.gz"
GO_URL="https://go.dev/dl/${GO_TARBALL}"

echo "Installing Go ${GO_VERSION} for ${GO_OS}/${GO_ARCH}..."
curl -fsSL "$GO_URL" -o /tmp/go.tar.gz
sudo tar -C /usr/local -xzf /tmp/go.tar.gz
rm -f /tmp/go.tar.gz

echo "Go installed to ${GO_INSTALL_DIR}."
echo
echo "To use Go in this and future terminals, add this to your shell config:"
echo
case "${SHELL:-}" in
    */zsh)
        echo "  echo 'export PATH=\${PATH}:${GO_INSTALL_DIR}/bin' >> ~/.zshrc"
        echo "  source ~/.zshrc"
        ;;
    */bash|*/sh)
        echo "  echo 'export PATH=\${PATH}:${GO_INSTALL_DIR}/bin' >> ~/.bashrc"
        echo "  source ~/.bashrc"
        ;;
    *)
        echo "  export PATH=\${PATH}:${GO_INSTALL_DIR}/bin"
        echo "(add the line above to your shell's rc file, e.g. ~/.bashrc or ~/.zshrc)"
        ;;
esac
echo
echo "Then re-run: make setup"
