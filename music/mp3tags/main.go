// this program reads MP3 or flac files and pulls out the meta data
// for title, artist, album and path
//
//
// this is not multi-processing safe

package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/dhowden/tag"
	"github.com/dlclark/metaphone3"
	"io/fs"
	"log"
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
var noTheRegex = regexp.MustCompile("^((T|t)(H|h)(E|e)) ")
var csvW *csv.Writer

const divP = " -+" // want space for names like Led Zeppelin - Bron-Yr-Aur

var dashRegex = regexp.MustCompile(divP)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s DIRNAME\n", os.Args[0])
		os.Exit(1)
	}
	flag.Parse()
	pathArg := path.Clean(flag.Arg(0))
	csvW = csv.NewWriter(os.Stdout)
	defer csvW.Flush()
	fsys := os.DirFS(pathArg)
	fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		processFile(pathArg, fsys, p, d, err)
		return err
	})
	return
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
	printCSV(aSong)
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
		punchIdx := dashRegex.FindStringIndex(filename)
		if punchIdx != nil {
			sn.title = strings.TrimSpace(filename[punchIdx[1]:])
		}
	}
	sn.titleH, _ = enc.Encode(justLetter(sn.title))
	sn.artist = m.Artist()

	if noTheRegex.MatchString(sn.artist) {
		sn.artist = sn.artist[4:]
		//fmt.Printf("%s %s %s\n", sn.title, sn.artist, sn.path)
	}
	sn.album = m.Album()
	return sn
}

// go thru the map, sort by key
// then create new ordering that makes sense to human
func printCSV(aSong *song) {
	var record []string
	record = make([]string, 4)

	record[0] = aSong.title
	record[1] = aSong.artist
	record[2] = aSong.album
	record[3] = aSong.path

	csvW.Write(record)
	//			fmt.Printf("%s %s %s\n", aSong.title, aSong.artist, aSong.path)

	if err := csvW.Error(); err != nil {
		log.Fatalln("error writing csv:", err)
	}
	return
}
