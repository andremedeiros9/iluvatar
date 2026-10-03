# iluvatar

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