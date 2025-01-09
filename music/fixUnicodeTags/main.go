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
	"github.com/dhowden/tag"
	"io/fs"
	"os"
	"path"
	"regexp"
	"strings"
	"unicode"
)

type song struct {
	artist string
	album  string
	title  string
	path   string
}

var repuni = map[string]string{
	"à":      "a",
	"á":      "a",
	"ã":      "a",
	"ç":      "c",
	"é":      "e",
	"è":      "e",
	"ë":      "e",
	"ì":      "i",
	"í":      "i",
	"ñ":      "n",
	"ò":      "o",
	"ó":      "o",
	"ô":      "o",
	"ö":      "o",
	"ś":      "s",
	"ù":      "u",
	"ú":      "u",
	"È":      "E",
	"É":      "E",
	"Ù":      "U",
	"Ú":      "U",
	"’":      "'",
	"´":      "'",
	"`":      "'",
	"“":      "\"",
	"”":      "\"",
	"«":      "\"",
	"»":      "\"",
	"…":      "",
	"⁄":      " ",
	"\u00A0": " ", // non-breaking space
	"\u2010": "-", // hyphen
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
		//fmt.Printf("%s %s %s\n", sn.title, sn.artist, sn.path)
	}
	sn.album = m.Album()
	if noTheRegex.MatchString(sn.album) {
		sn.album = sn.album[4:]
	}
	return sn
}

func locateUnicode(sn *song) {
	//	fmt.Printf("working on %s\n", sn.title)
	sn.title = replaceUnicode(sn.title)
	sn.artist = replaceUnicode(sn.artist)
	sn.album = replaceUnicode(sn.album)
}
func replaceUnicode(s string) string {
	var sb, ub strings.Builder
	for _, runeValue := range s {
		if runeValue > unicode.MaxASCII { // Check if the rune is not an ASCII character
			fmt.Printf("%c (U+%04X)\n", runeValue, runeValue)
			ub.WriteRune(runeValue)
		} else {
			if ub.Len() == 0 {
				sb.WriteRune(runeValue)
			} else {
				replace(&sb, &ub)
			}
		}
		if ub.Len() > 0 {
			replace(&sb, &ub)
		}
	}

	//fmt.Printf("processed %s\n", sb.String())
	return sb.String()
}
func replace(sb, ub *strings.Builder) string {
	// lookup
	k := ub.String()
	v, ok := repuni[k]
	if !ok {
		fmt.Printf("**** error: lookup failed for %s\n", k)
	}
	//fmt.Printf("k: %s f?:%t v: %s (U+%04X)\n", k, ok, v, v)
	sb.WriteString(v)
	ub.Reset()
	return v
}
func printFixCommand(aSong *song) {
	fmt.Printf("t:%s al:%s sr: %s p: %s\n", aSong.title, aSong.album, aSong.artist, aSong.path)
}
