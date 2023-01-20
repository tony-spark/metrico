// Package db contains database assets (e.g. migrations)
package db

import "embed"

//go:embed "migrations"
var EmbeddedDBFiles embed.FS
