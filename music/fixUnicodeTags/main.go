// this program reads MP3 and flac files
// and writes out commands to fix(remove or replace) unicode characters in the
// artist, title and album metadata tags
//
// this is not multi-processing safe
// this is not general, it uses a specific look up table rather
// than the official "PRECIS"
// (Preparation, Enforcement, and Comparison of Internationalized Strings in Application Protocols)
// and is documented in RFC7564.
//
// for example Joshua Bell's music often uses Cyrillic or Polish, which this does not handle.

package main

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/dhowden/tag"
	"github.com/pfarrell51/gows/music/util"
)

const fixer = "mp3tag"

type song struct {
	artist string
	album  string
	title  string
	path   string
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
	locateUnicode(aSong)
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
	}
	sn.album = m.Album()
	if noTheRegex.MatchString(sn.album) {
		sn.album = sn.album[4:]
	}
	return sn
}

func locateUnicode(sn *song) {
	fmt.Printf("working on %s\n", sn.title)
	var changed bool
	sn.title = util.CleanUni(sn.title, &changed)
	sn.artist = util.CleanUni(sn.artist, &changed)
	sn.album = util.CleanUni(sn.album, &changed)
	if changed {
		fmt.Println("it changed\b\n")
	}
}

func printFixCommand(aSong *song) {
	//  fixer [options] <file>
	// --str-title              Set the metadata title.
	//    --str-artist             Set the metadata artist.
	//    --str-album              Set the metadata album.

	fmt.Printf("%s --str-title \"%s\" --str-album \"%s\" --str-artist \"%s\" \"%s\" \n",
		fixer, aSong.title, aSong.album, aSong.artist, aSong.path)
}
