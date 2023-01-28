// this program helps clean up the mp3 and flac files in my hits playlist
// its main task is to normalize the file names to relect the artist and song title
//
// this is not multi-processing safe

// Bugs

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path"
	"regexp"
	"runtime"
	"strings"
	"time"
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
	alreadyNew        bool
	artist            string
	artistH           string
	artistHasThe      bool
	artistInDirectory bool
	artistKnown       bool
	album             string
	albumH            string
	title             string
	titleH            string
	inPath            string
	outPath           string
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
	start := time.Now()
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

	if false {
		ch := make(chan song)
		v := <-ch
		fmt.Println(v)
	}

	pathArg := path.Clean(flag.Arg(0))
	ProcessFiles(pathArg)
	duration := time.Since(start)
	fmt.Printf("# %v\n", duration)
}

func ProcessFiles(pathArg string) {
	if !zDumpArtist {
		rmap := walkFiles(pathArg)
		processMap(pathArg, rmap)
	} else {
		dumpGptree()
	}
}

var regAnd = regexp.MustCompile("(?i) (and|the) ")

func justLetter(a string) string {
	buff := bytes.Buffer{}
	loc := []int{0, 0}
	for j := 0; j < 4; j++ {	// 4 allows to and and two the, but it will nearly alwaysbreak before that
		loc = regAnd.FindStringIndex(a)
		if len(loc) < 1 {
			break
		}
		a = a[0:loc[0]] + a[loc[1]-1:]
	}
	for _, c := range a {
		if unicode.IsLetter(c) {
			buff.WriteRune(c)
		} else if c == '_' || c == '&' || unicode.IsSpace(c) {
			// ignore it
		} else if c == '-' {
			break
		}
	}
	return buff.String()
}
func loadMetaPhone() {
	groupNames := [...]string{
		"5th_Dimension", "ABBA", "Alice Cooper", "Alison_Krauss", "AllmanBrothers", "Almanac_Singers",
		"Animals", "Aquarius", "Aretha Franklin", "Arlo_Guthrie", "Association", "Average White Band",
		"Band", "Basia", "BeachBoys", "Beatles", "Bee Gees", "Billy Joel", "BlindFaith",
		"BloodSweatTears", "Blue Oyster Cult", "Blues Brothers", "Bob Dylan", "Boston", "Box Tops", "Bread",
		"Brewer and Shipley", "Brewer & Shipley", "BuffaloSpringfield", "Byrds",
		"Carole King", "Carpenters", "Cheap Trick", "Chesapeake", "Cream", "Crosby & Nash",
		"Crosby and Nash", "Crosby Stills & Nash", "Crosby Stills And Nash", "CSN&Y",
		"Crosby Stills Nash Young", "Crosby Stills Nash & Young", "David Allan Coe",
		"David Bowie", "David_Bromberg", "Deep Purple", "Derek and the Dominos",
		"Derek_Dominos", "Detroit Wheels",
		"Dire_Straits", "Doc Watson", "Don McLean", "Doobie_Brothers", "Doors", "Dylan",
		"Elton_John", "Emerson, Lake & Palmer", "Emmylou_Harris", "Fifth Dimension",
		"Fleetwood_Mac", "Genesis",
		"George Harrison", "Graham Nash", "Gram Parsons", "Hall and Oates", "Hall & Oates",
		"Heart", "Isley Brothers", "Jackie Wilson", "Jackson Browne",
		"James_Taylor", "Jefferson_Airplane", "Jethro_Tull", "Jimmy Buffett", "John_Denver",
		"John_Hartford", "John_Starling", "Joni_Mitchell", "Judy_Collins", "Kansas",
		"KC The Sunshine Band", "Kingston_Trio", "Led_Zepplin", "Linda_Ronstadt",
		"Lovin Spoonful", "Lynyrd_Skynyrd",
		"Mamas And Papas", "Mamas & The Papas", "Maria Muldaur",
		"Meatloaf", "Mike_Auldridge", "Mith Ryder & Detroit Wheels", "Moody Blues",
		"Neal Young", "Neil Diamond",
		"New Riders of the Purple Sage", "New_Riders_Purple_Sage",
		"Nitty Gritty Dirt Band", "Oates", "Otis Redding", "Pablo_Cruise", "Palmer",
		"Paul_Simon", "Peter_Paul_Mary", "Rascals", "Ringo Starr", "Roberta Flack", "Rolling_Stones",
		"Roy_Orbison", "Sam And Dave", "Santana", "Seals and Crofts", "Seals_Croft", "Seldom_Scene",
		"Shadows Of Knight", "Simon and Garfunkel", "Simon_Garfunkel", "Sonny And Cher",
		"Spoonful", "Seals & Crofts", "Steely_Dan", "Steppenwolf", "Steven_Stills",
		"Stevie Ray Vaughan and Double Trouble", "Sting", "Sunshine Band",
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

var sortKeyExp = regexp.MustCompile("^[A-Z](-|_)")
var extRegex = regexp.MustCompile("((M|m)(p|P)(3|4))|((F|f)(L|l)(A|a)(C|c))$")
var underToSpace = regexp.MustCompile("_")

var cReg = regexp.MustCompile(",\\s")
var dReg = regexp.MustCompile("-\\s")
var commaExp = regexp.MustCompile(",\\s")

// parse the file info to find artist and song title
// most of my music files have file names with the artist name, a hyphen and then the track title
// so this pulls out the information and fills in the "song" object.
func parseFilename(pathArg, p string) *song {
	// fmt.Printf("sf: %s\n", p)
	var rval = new(song)
	rval.inPath = path.Join(pathArg, p)
	rval.outPath = pathArg
	nameB := []byte(strings.TrimSpace(p))
	if sortKeyExp.Match(nameB) {
		nameB = nameB[2:]
	}
	nameB = underToSpace.ReplaceAll(nameB, []byte(" "))
	extR := extRegex.FindIndex(nameB)
	if extR == nil || len(extR) == 0 {
		return rval
	}
	ext := path.Ext(p)
	rval.ext = ext
	nameB = nameB[0 : extR[0]-1]

	var groupN, songN string
	ps, _ := path.Split(string(nameB))
	if len(ps) > 0 {
		rval.artistInDirectory = true
		rval.outPath = path.Join(rval.outPath, ps)
		nameB = nameB[len(ps):]
		songN = string(nameB)
		groupN = ps[0 : len(ps)-1] // cut off trailing slash
	}
	words := cReg.FindAllIndex(nameB, -1)
	dash := dReg.FindIndex(nameB)

	var semiExp = regexp.MustCompile("; ")
	if semiExp.Match(nameB) {
		semiLoc := semiExp.FindIndex(nameB)
		songN = strings.TrimSpace(string(nameB[:semiLoc[0]]))
		groupN = strings.TrimSpace(string(nameB[semiLoc[1]:]))
		rval.alreadyNew = true
	} else {
		//  no ;
		if dash != nil {
			sa := string(nameB[:dash[0]])
			sb := string(nameB[dash[1]:])
			if len(words) == 0 {
				groupN = sa
				songN = sb
			} else {
				sa := sa
				sb := sb
				if strings.HasPrefix(sa, "The ") {
					sa = sa[4:]
				}
				if strings.HasPrefix(sb, "The ") {
					sb = sb[4:]
				}
				ta, _ := enc.Encode(justLetter(sa))
				tb, _ := enc.Encode(justLetter(sb))
				_, OKa := gptree.Get(ta)
				if OKa {
					songN = sb
					groupN = sa
				}
				_, OKb := gptree.Get(tb)
				if OKb {
					groupN = sb
					songN = sa
				}
			}
		} else {
			// fmt.Println("no dash and no ; try , ")
			commasS := commaExp.FindIndex(nameB)
			if commasS == nil || len(commasS) == 0 {
				songN = string(nameB)
			} else {
				songN = string(nameB[:commasS[0]])
				groupN = string(nameB[commasS[1]:])
			}
		}
	}
	groupN = cases.Title(language.English, cases.NoLower).String(strings.TrimSpace(groupN))
	rval.title = strings.TrimSpace(songN)
	rval.titleH, _ = enc.Encode(justLetter(songN))
	rval.artist = strings.TrimSpace(groupN)
	//fmt.Printf("main %s by %s\n", rval.title, rval.artist)
	if strings.HasPrefix(rval.artist, "The ") {
		rval.artistHasThe = true
		rval.artist = rval.artist[4:]
	}
	rval.artistH, _ = enc.Encode(justLetter(rval.artist))
	_, ok := gptree.Get(rval.artistH)
	rval.artistKnown = ok
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
	extR := extRegex.FindStringIndex(p)
	if extR == nil {
		return nil // not interesting extension
	}

	aSong := parseFilename(pathArg, p)
	if aSong == nil {
		return nil
	}
	v := sMap[aSong.titleH]
	if len(v.titleH) > 0 {
		if aSong.artistH == v.artistH {
			fmt.Printf("#existing duplicate song for %s %s == %s\n", aSong.inPath, aSong.title, v.title)
		} else {
			fmt.Printf("#possible dup song for %s %s == %s %s\n", aSong.inPath, aSong.title,
				v.title, v.artist)
			aSong.titleH += "1"
		}
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
	uniqueArtists := make(map[string]string) // we just need a set, but use a map

	for _, aSong := range m {
		switch {
		case doRename:
			cmd := "mv"
			if runtime.GOOS == "windows" {
				cmd = "ren "
			}
			switch {
			case aSong.alreadyNew:
				continue
				//fmt.Printf("pM aNew %s \"%s\" \"%s/%s; %s%s\"\n", cmd, aSong.inPath,
				//	pathArg, aSong.title, aSong.artist, aSong.ext)
			case aSong.artist == "":
				fmt.Printf("#rename artist is blank %s\n", aSong.inPath)
				cmd = "#" + cmd
				continue
			case aSong.artistInDirectory:
				fmt.Printf("%s \"%s\" \"%s/%s%s\"\n", cmd, aSong.inPath,
					aSong.outPath, aSong.title, aSong.ext)
				continue
			}
			if aSong.artist == "" {
				fmt.Printf("%s \"%s\" \"%s/%s%s\"\n", cmd, aSong.inPath,
					aSong.outPath, aSong.title, aSong.ext)
			} else {
				fmt.Printf("%s \"%s\" \"%s/%s; %s%s\"\n", cmd, aSong.inPath,
					aSong.outPath, aSong.title, aSong.artist, aSong.ext)
			}
		case justList:
			the := ""
			if aSong.artistHasThe {
				the = "The "
			}
			fmt.Printf("%s by %s%s\n", aSong.title, the, aSong.artist)
		case showArtistNotInMap && !aSong.artistKnown:
			prim, _ := enc.Encode(justLetter(aSong.artist))
			if len(prim) == 0 && len(aSong.artist) == 0 {
				// fmt.Printf("prim: %s, a: %v\n", prim, aSong)
				continue
			}
			uniqueArtists[prim] = aSong.artist
		case noGroup:
			if aSong.artist == "" {
				fmt.Printf("nogroup %s\n", aSong.inPath)
			}
		default:
		}
	}
	if showArtistNotInMap {
		for k, v := range uniqueArtists {
			fmt.Printf("addto map k: %s v: %s\n", k, v)
		}
	}
	return m
}
func dumpGptree() {
	if zDumpArtist {
		gptree.Each(func(key string, v string) {
			fmt.Printf("\"%s\", \n", v)
		})
	}
}
