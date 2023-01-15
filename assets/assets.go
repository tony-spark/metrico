// Package assets contains various assets (e.g. web templates)
package assets

import "embed"

//go:embed "templates"
var EmbeddedAssets embed.FS
