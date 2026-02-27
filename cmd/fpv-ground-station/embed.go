//go:build !dev

package main

import (
	"embed"
	"io/fs"
)

//go:embed all:dist
var distFS embed.FS

func webDistFS() (fs.FS, error) {
	return fs.Sub(distFS, "dist")
}
