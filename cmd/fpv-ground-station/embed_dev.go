//go:build dev

package main

import "io/fs"

func webDistFS() (fs.FS, error) {
	return nil, nil
}
