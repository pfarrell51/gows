// test driver for gopro rename utility
package main

import (
	"fstest"
	"testing"
)

func TestProcessFile(t *testing.T) {

	fsys := fstest.MapFS{
		"file.go":                {},
		"subfolder/subfolder.go": {},
		"subfolder2/another.go":  {},
		"subfolder2/file.go":     {},
	}
	ProcessFiles(fsys, ".")

}
