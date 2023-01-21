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
var extRegex = regexp.MustCompile(".((M|m)(p|P)(3|4))|((F|f)(L|l)(A|a)(C|c))$")

const divP = "-+"

var dashRegex = regexp.MustCompile(divP)

var songMap map[string]song

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s DIRNAME", os.Args[0])
		os.Exit(1)
	}

	songMap = make(map[string]song)
	pathArg := path.Clean(os.Args[1])
	ProcessFiles(pathArg)
}

func ProcessFiles(pathArg string) {
	walkFiles(pathArg)
	processMap()
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
func processFile(pathArg string, fsys fs.FS, p string, d fs.DirEntry, err error) error {
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
	aSong := new(song)
	aSong.path = path.Join(pathArg, p)
	whichSong(aSong)
	songMap[aSong.titleH] = *aSong
	return nil
}
func whichSong(sn *song) *song {
	file, err := os.Open(sn.path)
	defer file.Close()
	if err != nil {
		fmt.Println(err)
		return sn
	}

	m, err := tag.ReadFrom(file)
	if err != nil {
		fmt.Println(err)
	}
	sn.title = m.Title() // The title of the track (see Metadata interface for more details).
	if sn.title == "" {
		_, filename := path.Split(sn.path)
		punchIdx := dashRegex.FindStringIndex(filename)
		if punchIdx != nil {
			sn.title = strings.Trim(filename[punchIdx[1]:], " 	")
		}
	}
	sn.titleH, _ = enc.Encode(justLetter(sn.title))
	sn.artist = m.Artist()
	sn.album = m.Album()
	return sn
}

// walk all files,
// fill in a map keyed by the desired new name order
func walkFiles(pathArg string) {
	fsys := os.DirFS(pathArg)
	fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		err = processFile(pathArg, fsys, p, d, err)
		return nil
	})
	return
}

// go thru the map, sort by key
// then create new ordering that makes sense to human
func processMap() {
	for _, aSong := range songMap {
		//if aSong.artist == "" {
		fmt.Printf("%s %s %s\n", aSong.title, aSong.artist, aSong.path)
		//}
	}
	return
}
