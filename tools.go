//go:build tools
// +build tools

// Official way to track tool dependencies with go modules:
// https://go.dev/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module

package main

import (
	_ "github.com/go-swagger/go-swagger/cmd/swagger"
	_ "github.com/golang-migrate/migrate/v4/cmd/migrate"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/sqlc-dev/sqlc/cmd/sqlc"
	_ "github.com/vektra/mockery/v2"
)
