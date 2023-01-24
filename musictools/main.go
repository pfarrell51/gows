// this program helps clean up the mp3 and flac files in my hits playlist
// its main task is to normalize the file names to relect the artist and song title
//
// this is not multi-processing safe

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path"
	"regexp"
	"strings"
	"unicode"

	"github.com/dlclark/metaphone3"
	g "github.com/zyedidia/generic"
	"github.com/zyedidia/generic/btree"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var gptree = btree.New[string, string](g.Less[string])
var enc metaphone3.Encoder
var showArtistNotInMap bool
var doRename bool
var justList bool
var noGroup bool
var zDumpArtist bool

type song struct {
	artist            string
	artistH           string
	artistHasThe      bool
	artistInDirectory bool
	artistKnown       bool
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
	flag.BoolVar(&showArtistNotInMap, "a", false, "artist map  list artist not in gpmap")
	flag.BoolVar(&justList, "l", false, "list - `list files")
	flag.BoolVar(&noGroup, "n", false, "nogroup - `list files that do not have an artist/group in the title")
	flag.BoolVar(&doRename, "r", false, "rename - perform rename function on needed files")
	flag.BoolVar(&zDumpArtist, "z", false, "list artist names one per line")
	flag.Parse()

	pathArg := path.Clean(flag.Arg(0))
	ProcessFiles(pathArg)
}

func ProcessFiles(pathArg string) {
	rmap := walkFiles(pathArg)
	processMap(pathArg, rmap)
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
	groupNames := [...]string{
		"5th_Dimension", "ABBA", "Alison_Krauss", "AllmanBrothers", "Almanac_Singers",
		"Animals", "Aretha Franklin", "Arlo_Guthrie", "Association", "Band", "Basia",
		"BeachBoys", "Beatles", "BlindFaith", "BloodSweatTears", "Boston", "Box Tops",
		"Brewer and Shipley", "Brewer & Shipley", "BuffaloSpringfield", "Byrds",
		"Carole King", "Carpenters", "Cheap Trick", "Chesapeake", "Cream", "Crosby & Nash",
		"Crosby and Nash", "Crosby Stills And Nash", "Crosby_Stills_Nash", "David Allan Coe", 
		"David Bowie", "David_Bromberg", "Deep Purple", "Derek and the Dominos", 
		"Derek_Dominos", "Detroit Wheels",
		"Dire_Straits", "Doc Watson", "Don McLean", "Doobie_Brothers", "Doors", "Dylan",
		"Elton_John", "Emmylou_Harris", "Fifth Dimension", "Fleetwood_Mac", "Genesis",
		"George Harrison", "Graham Nash", "Hall and Oates", "Heart", "Isley Brothers",
		"Jackie Wilson", "Jackson Browne",
		"James_Taylor", "Jefferson_Airplane", "Jethro_Tull", "Jimmy Buffett", "John_Denver",
		"John_Hartford", "John_Starling", "Joni_Mitchell", "Judy_Collins", "Kansas",
		"Kingston_Trio", "Led_Zepplin", "Linda_Ronstadt", "Lynyrd_Skynyrd",
		"Mamas And Papas", "Maria Muldaur",
		"Meatloaf", "Mike_Auldridge", "Moody Blues", "Neal Young", "Neil Diamond",
		"New Riders of the Purple Sage", "New_Riders_Purple_Sage",
		"Nitty Gritty Dirt Band", "Oates", "Otis Redding", "Pablo_Cruise", "Palmer",
		"Paul_Simon", "Peter_Paul_Mary", "Rascals", "Ringo Starr", "Roberta Flack", "Rolling_Stones",
		"Roy_Orbison", "Sam And Dave", "Santana", "Seals and Crofts", "Seals_Croft", "Seldom_Scene",
		"Shadows Of Knight", "Simon and Garfunkel", "Simon_Garfunkel", "Sonny And Cher",
		"Spoonful", "Steely_Dan", "Steppenwolf", "Steven_Stills", "Sting", "Sunshine Band",
		"Three Dog Night", "TonyRice", "Traveling_Wilburys", "Turtles", "Warren Zevon",
		"Who", "Wilson Pickett", "Yes",
	}
	for _, n := range groupNames {
		prim, sec := enc.Encode(justLetter(n))
		gptree.Put(prim, n)
		if len(sec) > 0 {
			gptree.Put(sec, n)
		}
	}
}

const nameP = "(([0-9A-Za-z&]*)\\s*)*"
const divP = " -+" // want space for names like Led Zeppelin - Bron-Yr-Aur
var regMulti = regexp.MustCompile(nameP + divP)
var regDash = regexp.MustCompile(divP)
var newStyle = regexp.MustCompile(":\\s")
var sortKeyExp = regexp.MustCompile("^[A-Z](-|_)")

// most of my music files have file names with the artist name, a hyphen and then the track title
// so this pulls out the information and fills in the "song" object.
func splitFilename(ps, pn string) *song {
	var rval = new(song)
	nameB := []byte(strings.TrimSpace(pn))
	dashS := regDash.FindIndex(nameB)
	newStyleS := newStyle.FindIndex(nameB)
	var groupN, songN string
	switch {
	case newStyleS != nil:
		groupN = string(nameB[newStyleS[1]:])
		songN = strings.TrimSpace(string(nameB[:newStyleS[0]]))
	case dashS == nil:
		// no punct => no group. Use what you have as song title
		songN = strings.TrimSpace(pn)
		if len(ps) > 1 {
			rval.artistInDirectory = true
			groupN = ps[:len(ps)-1]
			groupN = cases.Title(language.English, cases.NoLower).String(groupN)
		}
	default:
		// fall thru, old style
		sortKey := sortKeyExp.Find(nameB) // cut out leading "X_"
		if len(sortKey) > 0 {
			nameB = nameB[2:]
		}
		groupS := regMulti.Find(nameB)
		if groupS == nil {
			fmt.Println("PIB, group empty ", groupS)
			return rval
		}
		groupN = strings.TrimSpace(string(groupS[:len(groupS)-2]))
		songN = strings.TrimSpace(pn[dashS[1]:])
	}
	rval.title = songN
	rval.titleH, _ = enc.Encode(justLetter(songN))
	rval.artist = strings.TrimSpace(groupN)
	if strings.HasPrefix(rval.artist, "The ") {
		rval.artistHasThe = true
		rval.artist = rval.artist[4:]
	}
	rval.artistH, _ = enc.Encode(justLetter(rval.artist))
	_, ok := gptree.Get(rval.artistH)
	rval.artistKnown = ok
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
func processMap(pathArg string, m map[string]song) map[string]song {
	uniqueArtists := make(map[string]bool)

	for _, aSong := range m {
		switch {
		case doRename:
			if aSong.artist == "" {
				fmt.Printf("#rename artist is blank %s\n", aSong.path)
				continue
			}
			if aSong.artistInDirectory {
				continue
			}
			fmt.Printf("mv '%s' '%s/%s: %s%s'\n", aSong.path,
				pathArg, aSong.artist, aSong.title, aSong.ext)
		case justList:
			the := ""
			if aSong.artistHasThe {
				the = "The "
			}
			fmt.Printf("%s by %s%s\n", aSong.title, the, aSong.artist)
		case showArtistNotInMap && !aSong.artistKnown:
			uniqueArtists[aSong.artist] = true
		case noGroup:
			if aSong.artist == "" {
				fmt.Printf("nogroup %s\n", aSong.path)
			}
		default:
		}
	}
	if showArtistNotInMap {
		for k, _ := range uniqueArtists {
			fmt.Printf("addto map k: %s\n", k)
		}
	}
	if zDumpArtist {
		gptree.Each(func(key string, v string) {
			fmt.Printf("\"%s\", \n", v)
		})
	}
	return m
}
