# iluvatar

```
                      \       |       /
                 '.    \      |      /    .'
                   '.   .-"""""""""-.   .'
             ---      .'/////|||\\\\\'.      ---
                     /////.-"""""-.\\\\\
                    |////'  _   _  '\\\\|
                    |///   (o) (o)   \\\|
                   (|//       |       \\|)
                    |/\      /_\      /\|
                    '/\\\  .-===-.  ///\'
                     \\\\\\\'---'///////
                      \\\\\\\\|////////
                       '\\\\\\|//////'
                         '\\\\|////'
                           '\\|//'
                             \|/
                              '

 ___  _      _   _ __     __    _     _____     _     ____
|_ _|| |    | | | |\ \   / /   / \   |_   _|   / \   |  _ \
 | | | |    | | | | \ \ / /   / _ \    | |    / _ \  | |_) |
 | | | |___ | |_| |  \ V /   / ___ \   | |   / ___ \ |  _ <
|___||_____| \___/    \_/   /_/   \_\  |_|  /_/   \_\|_| \_\

            The One who sings your services into being.
```

## Install
Pick the instructions for your system.

### Linux and macOS

```sh
curl -fsSL https://github.com/andremedeiros9/iluvatar/releases/latest/download/install.sh | sh
```

This puts `iluvatar` in `~/.local/bin`. If that directory isn't on your
`PATH`, the installer tells you so. Set `ILUVATAR_INSTALL_DIR` to install
somewhere else.

### Windows

In PowerShell:

```powershell
irm https://github.com/andremedeiros9/iluvatar/releases/latest/download/install.ps1 | iex
```

This puts `iluvatar.exe` in `%LOCALAPPDATA%\Programs\iluvatar` and adds that
directory to your user `PATH`. Open a new terminal afterwards. Set
`ILUVATAR_INSTALL_DIR` to install somewhere else.

### Manual download

Download the archive for your system from the
[latest release](https://github.com/andremedeiros9/iluvatar/releases/latest)
and extract it. Each archive holds the `iluvatar` binary and an installer
that puts it in the same place as the commands above, without downloading
anything else:

- **Linux and macOS**: run `./install.sh` from the extracted folder.
- **Windows**: double-click `install.cmd` in the extracted folder.

| System                | Archive                         |
| --------------------- | ------------------------------- |
| Linux (x86-64)        | `iluvatar_linux_amd64.tar.gz`   |
| Linux (ARM64)         | `iluvatar_linux_arm64.tar.gz`   |
| macOS (Intel)         | `iluvatar_darwin_amd64.tar.gz`  |
| macOS (Apple Silicon) | `iluvatar_darwin_arm64.tar.gz`  |
| Windows (x86-64)      | `iluvatar_windows_amd64.zip`    |
| Windows (ARM64)       | `iluvatar_windows_arm64.zip`    |

### With Go

```sh
go install github.com/andremedeiros9/iluvatar/cmd/iluvatar@latest
```

## Uninstall

Uninstalling removes the `iluvatar` binary and the `protoc` copy iluvatar
may have downloaded into `~/.iluvatar`. Projects you generated are not
touched.

### Linux and macOS

```sh
rm -f ~/.local/bin/iluvatar
rm -rf ~/.iluvatar
```

### Windows

In PowerShell:

```powershell
$dir = "$env:LOCALAPPDATA\Programs\iluvatar"
Remove-Item -Recurse -Force $dir, "$env:USERPROFILE\.iluvatar" -ErrorAction SilentlyContinue

# Take the install directory back out of your user PATH.
$path = ([Environment]::GetEnvironmentVariable('Path', 'User') -split ';') -ne $dir
[Environment]::SetEnvironmentVariable('Path', ($path -join ';'), 'User')
```

### Other installs

- If you set `ILUVATAR_INSTALL_DIR`, or moved the binary yourself after a
  manual download, delete it from that directory instead.
- If you installed with `go install`, delete `iluvatar` (`iluvatar.exe` on
  Windows) from the directory `go env GOPATH` prints, under `bin`.

`iluvatar generate` also installs the `protoc-gen-go` and
`protoc-gen-go-grpc` plugins with `go install` when they are missing. They
live in the same Go `bin` directory and can be deleted from there if you
don't use them elsewhere.

A code generator for RESTful APIs with basic CRUD actions. Generated code is
always Go, for now.

## Usage (planned)

```
iluvatar init <path>
```

Creates `<path>` and writes a starter `iluvatar.toml` into it, pre-filled
with defaults (database driver, HTTP framework, etc.) for you to edit.
Running `iluvatar generate` against that file will then scaffold the actual
Go REST API.

## Architecture

```
cmd/iluvatar/            Entry point: parses args, hands off to internal/cli.

internal/
  cli/                   Cobra command tree (root, init, generate). Thin:
                         flag/arg parsing only, delegates real work below.

  config/                The Config struct mirroring iluvatar.toml's schema
                         (project, server, database, resources). Shared by
                         scaffold (to render defaults) and generator (to
                         read the user's edited file).

  scaffold/              Backs `iluvatar init`. Creates the target directory
                         and renders templates/iluvatar.toml.tmpl into a
                         starter config file.
    templates/           Templates for files `init` writes out.

  generator/             Backs `iluvatar generate` (not yet implemented).
                         Will turn a parsed Config into a working Go REST
                         API — handlers, models, routing, DB access — using
                         its own templates/ directory.
```

Each `internal/*` package keeps a `doc.go` with its package-level
documentation, so the layout above stays in sync with the code.
