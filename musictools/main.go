// this program helps clean up the mp3 and flac files in my hits playlist
// its main task is to normalize the file names to relect the artist and song title
//
// this is not multi-processing safe

package main

import (
	"bytes"
	"flag"
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
var doRename bool
var justList bool
var noGroup bool

type song struct {
	artist            string
	artistH           string
	artistHasThe      bool
	inArtistDirectory bool
	album             string
	albumH            string
	title             string
	titleH            string
	path              string
	ext               string
}

func init() {
	enc.MaxLength = 24
	loadMetaPhone()
}
func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s [flags] directory-spec\n", os.Args[0])
		os.Exit(1)
	}
	flag.Usage = func() {
		w := flag.CommandLine.Output() // may be os.Stderr - but not necessarily
		fmt.Fprintf(w, "Usage of %s: [flags] directory-spec \n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(w, "default is to list files that need love.\n")

	}
	flag.BoolVar(&doRename, "r", false, "rename - perform rename function on needed files")
	flag.BoolVar(&justList, "l", false, "list - `list files")
	flag.BoolVar(&noGroup, "n", false, "nogroup - `list files that do not have an artist/group in the title")
	flag.Parse()

	pathArg := path.Clean(flag.Arg(0))
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

const nameP = "(([0-9A-Za-z]*)\\s*)*"
const divP = " -+" // want space for names like Led Zeppelin - Bron-Yr-Aur
var regMulti = regexp.MustCompile(nameP + divP)
var regPunch = regexp.MustCompile(divP)
var newStyle = regexp.MustCompile(":\\s")

// most of my music files have file names with the artist name, a hyphen and then the track title
// so this pulls out the information and fills in the "song" object.
func splitFilename(ps, pn string) *song {
	var rval = new(song)
	nameB := []byte(strings.TrimSpace(pn))
	punctS := regPunch.FindIndex(nameB)
	newStyleS := newStyle.FindIndex(nameB)
	var groupN, songN string
	switch {
	case newStyleS != nil:
		groupN = string(nameB[:newStyleS[1]])
		songN = strings.TrimSpace(string(nameB[newStyleS[1]:]))
	case punctS == nil:
		// no punct => no group. Use what you have as song title
		songN = strings.TrimSpace(pn)
		if len(ps) > 1 {
			fmt.Printf("no punct g: %s t: %s\n", ps, rval.title)
			rval.inArtistDirectory = true
			groupN = ps
		}
	default:
		// fall thru, old style
		groupS := regMulti.Find(nameB)
		if groupS == nil {
			fmt.Println("PIB, group empty ", groupS)
			return rval
		}
		groupN = strings.TrimSpace(string(groupS[:len(groupS)-2]))
		songN = strings.TrimSpace(pn[punctS[1]:])
	}
	rval.title = songN
	rval.titleH, _ = enc.Encode(justLetter(songN))
	rval.artist = groupN
	if strings.HasPrefix(rval.artist, "The ") {
		rval.artistHasThe = true
		rval.artist = rval.artist[4:]
	}
	rval.artistH, _ = enc.Encode(justLetter(rval.artist))
	return rval
}

var extRegex = regexp.MustCompile("((M|m)(p|P)(3|4))|((F|f)(L|l)(A|a)(C|c))$")

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
	ext := path.Ext(p)
	if len(ext) == 0 {
		return nil // not interesting extension
	}
	extR := extRegex.FindStringIndex(p)
	if extR == nil {
		return nil // not interesting extension
	}
	shortp := p[0 : extR[0]-1]
	ps, pn := path.Split(shortp)
	aSong := splitFilename(ps, pn)
	aSong.ext = ext
	aSong.path = path.Join(pathArg, ps, pn) + ext
	v := sMap[aSong.titleH]
	if len(v.titleH) > 0 {
		fmt.Printf("existing song for %s %s == %s\n", aSong.path, aSong.title, v.title)
		return nil
	}
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
		switch {
		case doRename:
			if aSong.artist == "" {
				fmt.Printf("rename artist is blank %s\n", aSong.path)
			}
		case justList:
			fmt.Printf("%s by %s\n", aSong.title, aSong.artist)
		case noGroup:
			if aSong.artist == "" {
				fmt.Printf("no artist/group %s\n", aSong.path)
			}
		default:
			fmt.Printf("default %s\n", aSong.path)
		}
	}
	return m
}
