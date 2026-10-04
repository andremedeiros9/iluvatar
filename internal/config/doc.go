// Package config defines the schema of iluvatar.toml, the file that
// describes a generated project: its module path, HTTP framework,
// database choice, and the list of resources to scaffold CRUD endpoints
// for.
//
// It will provide the Config struct that both scaffold (to render
// defaults) and generator (to read the user's edited file) depend on,
// plus a Load function to parse and validate iluvatar.toml from disk.
package config
