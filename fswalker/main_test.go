package main

import (
	"archive/zip"
	"main"
	"os"
	"testing"
	"testing/fstest"
)

func TestFilesInMemory(t *testing.T) {
	t.Parallel()
	fsys := fstest.MapFS{
		"file.go":                {},
		"subfolder/subfolder.go": {},
		"subfolder2/another.go":  {},
		"subfolder2/file.go":     {},
	}
	want := 4
	got := main.Files(fsys)
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}
