// this program reads MP3 and guesses where the souirce flac files are

//
//
// this is not multi-processing safe

package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/dhowden/tag"
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
	artistHasThe bool
	album        string
	title        string
	path         string
}

var noTheRegex = regexp.MustCompile("^((T|t)(H|h)(E|e)) ")
var csvW *csv.Writer

const divP = " -+" // want space for names like Led Zeppelin - Bron-Yr-Aur

var dashRegex = regexp.MustCompile(divP)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s DIRNAME\n", os.Args[0])
		os.Exit(1)
	}
	pathArg := path.Clean(os.Args[1])
	csvW = csv.NewWriter(os.Stdout)
	defer csvW.Flush()
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
