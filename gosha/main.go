// this program helps clean up the mp3 and flac files in my hits playlist
// its main task is to normalize the file names to relect the artist and song title
//
// this is not multi-processing safe

// Bugs

package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

type song struct {
	alreadyNew bool
	artist     string
	artistH    string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s [flags] directory-spec\n", os.Args[0])
		os.Exit(1)
	}
	start := time.Now()

	if false {
		ch := make(chan song)
		v := <-ch
		fmt.Println(v)
	}

	pathArg := path.Clean(os.Args[1])
	ProcessFiles(pathArg)
	duration := time.Since(start)
	fmt.Printf("# %v\n", duration)
}

func ProcessFiles(pathArg string) {
	walkFiles(pathArg)
}

// parse the file info to find artist and song title
// most of my music files have file names with the artist name, a hyphen and then the track title
// so this pulls out the information and fills in the "song" object.
func parseFilename(pathArg, p string) {
	fmt.Printf("sf: %s\n", p)
	fn := path.Join(pathArg, p)
	f, err := os.Open(fn)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%x", h.Sum(nil))
}

// this is the local  WalkDirFunc called by WalkDir for each file
// pathArg is the path to the base of our walk
// p is the current path/name
func processFile(pathArg string, fsys fs.FS, p string, d fs.DirEntry, err error) error {
	if err != nil {
		fmt.Println("Error processing", p, " in ", d)
		fmt.Println("error is ", err)
		return nil
	}
	if d == nil || d.IsDir() || strings.HasPrefix(p, ".") {
		return nil
	}

	parseFilename(pathArg, p)
	return nil
}

// walk all files, looking for nice GoPro created video files.
// fill in a map keyed by the desired new name order
func walkFiles(pathArg string) {
	fsys := os.DirFS(pathArg)
	fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		err = processFile(pathArg, fsys, p, d, err)
		return nil
	})
}

// go thru the map, sort by key
// then create new ordering that makes sense to human
func processMap(pathArg string, m map[string]song) map[string]song {
	rval := make(map[string]song) // we just need a set, but use a map
	return rval
}
