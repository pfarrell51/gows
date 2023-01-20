// this program helps clean up the mp3 and flac files in my hits playlist
// its main task is to normalize the file names to relect the artist and song title
//
// this is not multi-processing safe

package main

import (
	"bytes"
	"fmt"
	"github.com/dlclark/metaphone3"
	g "github.com/zyedidia/generic"
	"github.com/zyedidia/generic/btree"
	"io/fs"
	"os"
	"path"
	"regexp"
	"strings"
	"unicode"
)

var gptree = btree.New[string, string](g.Less[string])
var enc metaphone3.Encoder
var extRegex = regexp.MustCompile(".((M|m)(p|P)(3|4))|((F|f)(L|l)(A|a)(C|c))")

const nameP = "(([0-9A-Za-z]*)\\s*)*"
const divP = "-+"

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

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s DIRNAME", os.Args[0])
		os.Exit(1)
	}

	loadMetaPhone()
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
func loadMetaPhone() {
	groupNames := [...]string{"ABBA", "Alison_Krauss", "AllmanBrothers", "Almanac_Singers", "Animals",
		"Arlo_Guthrie", "Band", "Basia", "BeachBoys", "Beatles", "BlindFaith", "BloodSweatTears", "Boston",
		"BrewerAndShipley", "BuffaloSpringfield", "Byrds", "CensorBeep.mp4", "Chesapeake",
		"Cream", "Crosby_Stills_Nash",
		"David_Bromberg", "Derek_Dominos", "Dire_Straits", "Doobie_Brothers", "Doors", "Dylan", "Elton_John",
		"Emmylou_Harris", "Fleetwood_Mac", "Heart", "James_Taylor", "Jefferson_Airplane", "Jethro_Tull",
		"John_Denver", "John_Hartford", "John_Starling", "Joni_Mitchell", "Judy_Collins", "Kingston_Trio",
		"Led_Zepplin", "Linda_Ronstadt", "Lynyrd_Skynyrd", "Mamas_Popas", "Meatloaf", "Mike_Auldridge",
		"New_Riders_Purple_Sage", "Pablo_Cruise", "Paul_Simon", "Peter_Paul_Mary", "Rolling_Stones",
		"Roy_Orbison", "Santana", "Seals_Croft", "Seldom_Scene", "Simon_Garfunkel", "Steely_Dan",
		"5th_Dimension", "TonyRice", "Traveling_Wilburys", "Who", "Yes",
	}
	for _, n := range groupNames {
		prim, sec := enc.Encode(justLetter(n))
		gptree.Put(prim, n)
		if len(sec) > 0 {
			gptree.Put(sec, n)
		}
	}
}

// most of my music files have file names with the artist name, a hyphen and then the track title
// so this pulls out the information and fills in the "song" object.
func splitFilename(name string) *song {
	var regMulti = regexp.MustCompile(nameP + divP)
	var regPunch = regexp.MustCompile(divP)
	var rval = new(song)
	nameB := []byte(strings.Trim(name, " 	_"))
	punchS := regPunch.Find(nameB)
	if punchS == nil {
		// no group
		rval.title = name
		rval.titleH, _ = enc.Encode(justLetter(name))
		return rval
	}
	groupS := regMulti.Find(nameB)
	if groupS == nil {
		fmt.Println("PIB, group empty ", groupS)
		return rval
	}
	songN := strings.Trim(name[len(groupS):], " 	_")
	rval.title = songN
	rval.titleH, _ = enc.Encode(justLetter(songN))
	rval.artist = string(groupS[0 : len(groupS)-2])
	if strings.HasPrefix(rval.artist, "The ") {
		rval.artistHasThe = true
		rval.artist = rval.artist[4:]
	}
	rval.artistH, _ = enc.Encode(justLetter(rval.artist))
	return rval
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
	aSong := splitFilename(pn)
	aSong.path = path.Join(pathArg, ps, pn)

	sMap[aSong.titleH] = *aSong
	return nil
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
