#!/bin/sh
# Installs iluvatar on Linux or macOS.
#
# Run straight from the web, it downloads the latest release:
#
#   curl -fsSL https://github.com/andremedeiros9/iluvatar/releases/latest/download/install.sh | sh
#
# Run from an extracted release archive (./install.sh), it installs the
# iluvatar binary sitting next to it instead, without downloading anything.
#
# The binary goes into ~/.local/bin; set ILUVATAR_INSTALL_DIR to put it
# somewhere else.
set -eu

repo="andremedeiros9/iluvatar"
install_dir="${ILUVATAR_INSTALL_DIR:-$HOME/.local/bin}"

# bundled is the binary shipped alongside this script in a release
# archive. When the script is piped into sh, $0 is the shell itself, so
# there is no "alongside" to look in.
bundled=""
case "$0" in
    *install.sh)
        script_dir="$(cd "$(dirname "$0")" && pwd)"
        if [ -f "$script_dir/iluvatar" ]; then
            bundled="$script_dir/iluvatar"
        fi
        ;;
esac

if [ -n "$bundled" ]; then
    mkdir -p "$install_dir"
    cp "$bundled" "$install_dir/iluvatar"
else
    case "$(uname -s)" in
        Linux) os="linux" ;;
        Darwin) os="darwin" ;;
        *)
            echo "iluvatar: unsupported OS $(uname -s); on Windows use install.ps1" >&2
            exit 1
            ;;
    esac

    case "$(uname -m)" in
        x86_64 | amd64) arch="amd64" ;;
        arm64 | aarch64) arch="arm64" ;;
        *)
            echo "iluvatar: unsupported architecture $(uname -m)" >&2
            exit 1
            ;;
    esac

    archive="iluvatar_${os}_${arch}.tar.gz"
    base_url="https://github.com/${repo}/releases/latest/download"

    tmp="$(mktemp -d)"
    trap 'rm -rf "$tmp"' EXIT

    echo "Downloading ${archive}..."
    curl -fsSL -o "$tmp/$archive" "$base_url/$archive"
    curl -fsSL -o "$tmp/checksums.txt" "$base_url/checksums.txt"

    if command -v sha256sum >/dev/null 2>&1; then
        actual="$(sha256sum "$tmp/$archive" | cut -d ' ' -f 1)"
    elif command -v shasum >/dev/null 2>&1; then
        actual="$(shasum -a 256 "$tmp/$archive" | cut -d ' ' -f 1)"
    else
        actual=""
        echo "iluvatar: no sha256sum or shasum found, skipping checksum verification" >&2
    fi

    if [ -n "$actual" ]; then
        expected="$(awk -v file="$archive" '$2 == file { print $1 }' "$tmp/checksums.txt")"
        if [ "$actual" != "$expected" ]; then
            echo "iluvatar: checksum mismatch for ${archive}" >&2
            exit 1
        fi
    fi

    mkdir -p "$install_dir"
    tar -xzf "$tmp/$archive" -C "$install_dir" iluvatar
fi

chmod +x "$install_dir/iluvatar"

# A browser download is quarantined by macOS, which would refuse to run
# the binary; the user chose to install it, so lift that.
if [ "$(uname -s)" = "Darwin" ]; then
    xattr -d com.apple.quarantine "$install_dir/iluvatar" 2>/dev/null || true
fi

echo "Installed iluvatar to $install_dir/iluvatar"

case ":$PATH:" in
    *":$install_dir:"*) ;;
    *) echo "Add $install_dir to your PATH to run it as 'iluvatar'." ;;
esac
