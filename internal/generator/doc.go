// Package generator turns a parsed, user-edited Config into a working Go
// REST API with CRUD handlers for each configured resource.
//
// This is the engine behind the future "iluvatar generate" command. It is
// not implemented yet; for now this package only reserves the package's
// place in the architecture, and will hold its own Go source templates
// (handlers, models, routing, database access) under ./templates once
// code generation lands.
package generator
