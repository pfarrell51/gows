// this program reads MP3 or flac files and pulls out the meta data
// for title, artist and album
//
//
// this is not multi-processing safe

package main

import (
	"bytes"
	"fmt"
	"github.com/dhowden/tag"
	"github.com/dlclark/metaphone3"
	"io/fs"
	"os"
	"path"
	"regexp"
	"strings"
	"unicode"
)

type song struct {
	artist       string
	artistH      string
	artistHasThe bool
	album        string
	albumH       string
	title        string
	titleH       string
	path         string
}

var enc metaphone3.Encoder
var extRegex = regexp.MustCompile(".((M|m)(p|P)(3|4))|((F|f)(L|l)(A|a)(C|c))")

const nameP = "(([0-9A-Za-z]*)\\s*)*"
const divP = "-+"

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s DIRNAME", os.Args[0])
		os.Exit(1)
	}

	pathArg := path.Clean(os.Args[1])
	ProcessFiles(pathArg)
}

func ProcessFiles(pathArg string) {
	rmap := walkFiles(pathArg)
	processMap(rmap)
}
func justLetter(a string) string {
	buff := bytes.Buffer{}
	for _, c := range a {
		if unicode.IsLetter(c) {
			buff.WriteRune(c)
		} else if c == '_' {
			// ignore it
		} else if c == '-' {
			break
		}
	}
	return buff.String()
}

// this is the local  WalkDirFunc called by WalkDir for each file
// pathArg is the path to the base of our walk
// p is the current path/name
func processFile(pathArg string, sMap map[string]song, fsys fs.FS, p string, d fs.DirEntry, err error) error {
	if err != nil {
		fmt.Println("Error processing", p, " in ", d)
		fmt.Println("error is ", err)
		return nil
	}
	if d == nil || d.IsDir() || strings.HasPrefix(p, ".") {
		return nil
	}
	ext := extRegex.FindString(p)
	if len(ext) == 0 {
		return nil // not interesting extension
	}

	ps, pn := path.Split(p)
	aSong := whichSong(pn)
	aSong.path = path.Join(pathArg, ps, pn)

	sMap[aSong.titleH] = aSong
	return nil
}
func whichSong(sn string) song {
	rval := new(song)
	rval.titleH, _ = enc.Encode(justLetter(sn))
	m, err := tag.ReadFrom(sn)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(m.Format()) // The detected format.
	fmt.Print(m.Title())  // The title of the track (see Metadata interface for more details).
	return *rval
}

// walk all files, looking for nice GoPro created video files.
// fill in a map keyed by the desired new name order
func walkFiles(pathArg string) map[string]song {
	songMap := make(map[string]song)
	fsys := os.DirFS(pathArg)
	fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		err = processFile(pathArg, songMap, fsys, p, d, err)
		return nil
	})
	return songMap
}

// go thru the map, sort by key
// then create new ordering that makes sense to human
func processMap(m map[string]song) map[string]song {
	for _, aSong := range m {
		if aSong.artist == "" {
			fmt.Println(aSong.path)
		}
	}
	return m
}
