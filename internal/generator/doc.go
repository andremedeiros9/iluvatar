// Package generator turns a parsed, user-edited Config into a working Go
// REST API with CRUD handlers for each configured resource.
//
// This is the engine behind the "iluvatar generate" command. So far it
// writes the generated project's build tooling: a Makefile and a CI
// pipeline (GitHub Actions or GitLab CI, chosen by cfg.CI.Provider). CRUD
// handlers, models, routing, and database access are not implemented yet.
package generator
