// this program reads MP3 and flac files
// and writes out commands to fix unicode characters in the
// artist, title and album metadata tags
//
// this is not multi-processing safe

package main

import (
	_ "embed"
	"fmt"
	"github.com/dhowden/tag"
	"io/fs"
	"os"
	"path"
	"regexp"
	"strings"
)

type song struct {
	artist string
	album  string
	title  string
	path   string
}

var repuni = map[string]string{
	"é":      "e",
	"ó":      "o",
	"ś":      "s",
	"É":      "E",
	"…":      "",
	"\u2013": "-", //En dash
	"\u2014": "-", //Em dash
	"\u2015": "―", // Horizontal bar
}
var noTheRegex = regexp.MustCompile("^((T|t)(H|h)(E|e)) ")

const divP = " -+" // want space for names like Led Zeppelin - Bron-Yr-Aur

var dashRegex = regexp.MustCompile(divP)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s DIRNAME\n", os.Args[0])
		os.Exit(1)
	}

	pathArg := path.Clean(os.Args[1])

	fsys := os.DirFS(pathArg)
	fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		processFile(pathArg, fsys, p, d, err)
		return err
	})
	return
}

var extRegex = regexp.MustCompile(".((M|m)(p|P)(3|4))|((F|f)(L|l)(A|a)(C|c))$")

// this is the local  WalkDirFunc called by WalkDir for each file
// pathArg is the path to the base of our walk
// p is the current path/name
func processFile(pathArg string, fsys fs.FS, p string, d fs.DirEntry, err error) {
	if err != nil {
		fmt.Println("Error processing", p, " in ", d)
		fmt.Println("error is ", err)
		return
	}
	if d == nil || d.IsDir() || strings.HasPrefix(p, ".") {
		return
	}
	ext := extRegex.FindString(p)
	if len(ext) == 0 {
		return // not interesting extension
	}
	aSong := new(song)
	aSong.path = path.Join(pathArg, p)
	getMetadata(aSong)
	replaceUnicode(aSong)
	printFixCommand(aSong)
}
func getMetadata(sn *song) *song {
	file, err := os.Open(sn.path)
	defer file.Close()
	if err != nil {
		fmt.Println(err)
		return sn
	}
	m, err := tag.ReadFrom(file)
	if err != nil {
		fmt.Printf("%v %s", err, sn.title)
		return sn
	}
	sn.title = m.Title() // The title of the track (see Metadata interface for more details).
	if sn.title == "" {
		_, filename := path.Split(sn.path)
		sn.title = strings.TrimSpace(filename)
	}
	sn.artist = m.Artist()
	sn.artist = strings.ReplaceAll(sn.artist, ",", "")

	if noTheRegex.MatchString(sn.artist) {
		sn.artist = sn.artist[4:]
		//fmt.Printf("%s %s %s\n", sn.title, sn.artist, sn.path)
	}
	sn.album = m.Album()
	if noTheRegex.MatchString(sn.album) {
		sn.album = sn.album[4:]
	}
	return sn
}

func locateUnicode(aSong *song) {
	sn.title = replaceUnicode(sn.title)
	sn.artist = replaceUnicode(sn.artist)
	sn.album = replaceUnicode(sn.album)
}
func replaceUnicode(s string) string {
	return ""
}

func printFixCommand(aSong *song) {
	fmt.Printf("%s %s %s\n", aSong.title, aSong.artist, aSong.path)
}
