// this program reads MP3 or flac files and pulls out the meta data
// for title, artist and album
//
//
// this is not multi-processing safe

package main

import (
	"fmt"
	"github.com/dhowden/tag"
	"io/fs"
	"os"
	"path"
	"regexp"
)

var extRegex = regexp.MustCompile(".((M|m)(p|P)(3|4))|((F|f)(L|l)(A|a)(C|c))$")

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s DIRNAME\n", os.Args[0])
		os.Exit(1)
	}
	pathArg := path.Clean(os.Args[1])
	walkFiles(pathArg)
}

// walk all files,
// fill in a map keyed by the desired new name order
func walkFiles(pathArg string) {
	c := 0
	fsys := os.DirFS(pathArg)
	fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		c++
		err = processFile(pathArg, fsys, p, d, err)
		return nil
	})
	fmt.Printf("processed %d songs\n", c)
	return
}

// this is the local  WalkDirFunc called by WalkDir for each file
// pathArg is the path to the base of our walk
// p is the current path/name
func processFile(pathArg string, fsys fs.FS, p string, d fs.DirEntry, err error) error {
	ext := extRegex.FindString(p)
	if len(ext) == 0 {
		return nil // not interesting extension
	}
	path := path.Join(pathArg, p)
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return fmt.Errorf("error in open: ", err)
	}

	m, err := tag.ReadFrom(file)
	if err != nil {
		return fmt.Errorf("%v %s", err, path)

	}
	if len(m.Title()) > 0 {
		// do nothing
	}
	return nil
}
