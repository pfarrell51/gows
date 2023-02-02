// this program helps clean up the mp3 and flac files in my hits playlist
// its main task is to normalize the file names to relect the artist and Song title
//
// this is not multi-processing safe

// Bugs

package musictools

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path"
	"regexp"
	"runtime"
	"strings"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func init() {
	GetEncoder().MaxLength = 24
	LoadArtistMap()
}

// top level entry point, takes the path to the directory to walk/process
func ProcessFiles(pathArg string) {
	if GetFlags().ZDumpArtist {
		DumpGptree()
		return
	}
	rmap := WalkFiles(pathArg)
	ProcessMap(pathArg, rmap)
}

var regAnd = regexp.MustCompile("(?i) (and|the) ")

// takes a string and returns just the letters. Also removes the words "and" and "the" from the string
// since they are essentially noise words.
func JustLetter(a string) string {
	buff := bytes.Buffer{}
	loc := []int{0, 0}
	for j := 0; j < 4; j++ { // 4 allows the space before the keyword (and/the), as we back up
		loc = regAnd.FindStringIndex(a)
		if len(loc) < 1 {
			break
		}
		a = a[:loc[0]] + a[loc[1]-1:]
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

var sortKeyExp = regexp.MustCompile("^[A-Z](-|_)")
var extRegex = regexp.MustCompile("((M|m)(p|P)(3|4))|((F|f)(L|l)(A|a)(C|c))$")
var underToSpace = regexp.MustCompile("_")

var cReg = regexp.MustCompile(",\\s")
var dReg = regexp.MustCompile("-\\s")
var commaExp = regexp.MustCompile(",\\s")
var slashReg = regexp.MustCompile("/")

func pathLastTwo(s string) (artist, album string) {
	parts := slashReg.FindAllStringIndex(s, -1)
	switch len(parts) {
	case 1:
		artist = s[:parts[0][0]]
		album = ""
	case 2:
		artist = s[:parts[0][0]]
		album = s[parts[0][1]:parts[1][0]]
	default:
		fmt.Printf("PIB, no directory found in %s\n", s)
	}
	return artist, album
}

// parse the file info to find artist and Song title
// most of my music files have file names with the artist name, a hyphen and then the track title
// so this pulls out the information and fills in the "Song" object.
func parseFilename(pathArg, p string) *Song {
	if GetFlags().Debug {
		fmt.Printf("pf: %s\n", p)
	}
	var rval = new(Song)
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

	var groupN, SongN string
	ps, _ := path.Split(string(nameB))
	if len(ps) > 0 {
		rval.artistInDirectory = true
		rval.Artist, rval.Album = pathLastTwo(ps)
		rval.outPath = path.Join(rval.outPath, ps)

		nameB = nameB[len(ps):]
		SongN = string(nameB)
	}
	words := cReg.FindAllIndex(nameB, -1)
	dash := dReg.FindIndex(nameB)

	var semiExp = regexp.MustCompile("; ")
	if semiExp.Match(nameB) {
		semiLoc := semiExp.FindIndex(nameB)
		SongN = strings.TrimSpace(string(nameB[:semiLoc[0]]))
		groupN = strings.TrimSpace(string(nameB[semiLoc[1]:]))
		rval.alreadyNew = true
	} else {
		//  no ;
		if dash != nil {
			sa := string(nameB[:dash[0]])
			sb := string(nameB[dash[1]:])
			if len(words) == 0 {
				groupN = sa
				SongN = sb
			} else {
				sa := sa
				sb := sb
				if strings.HasPrefix(sa, "The ") {
					sa = sa[4:]
				}
				if strings.HasPrefix(sb, "The ") {
					sb = sb[4:]
				}
				ta, _ := GetEncoder().Encode(JustLetter(sa))
				tb, _ := GetEncoder().Encode(JustLetter(sb))
				_, OKa := gptree.Get(ta)
				if OKa {
					SongN = sb
					groupN = sa
				}
				_, OKb := gptree.Get(tb)
				if OKb {
					groupN = sb
					SongN = sa
				}
			}
		} else {
			// fmt.Println("no dash and no ; try , ")
			commasS := commaExp.FindIndex(nameB)
			if commasS == nil || len(commasS) == 0 {
				SongN = string(nameB)
			} else {
				SongN = string(nameB[:commasS[0]])
				groupN = string(nameB[commasS[1]:])
			}
		}
	}
	groupN = cases.Title(language.English, cases.NoLower).String(strings.TrimSpace(groupN))
	rval.Title = strings.TrimSpace(SongN)
	rval.titleH, _ = GetEncoder().Encode(JustLetter(SongN))
	rval.Artist = strings.TrimSpace(groupN)
	//fmt.Printf("main %s by %s\n", rval.title, rval.artist)
	if strings.HasPrefix(rval.Artist, "The ") {
		rval.artistHasThe = true
		rval.Artist = rval.Artist[4:]
	}
	rval.artistH, _ = GetEncoder().Encode(JustLetter(rval.Artist))
	_, ok := gptree.Get(rval.artistH)
	rval.artistKnown = ok
	return rval
}

const divP = " -+" // want space for names like Led Zeppelin - Bron-Yr-Aur
var dashRegex = regexp.MustCompile(divP)

// this is the local  WalkDirFunc called by WalkDir for each file
// pathArg is the path to the base of our walk
// p is the current path/name
func processFile(pathArg string, sMap map[string]Song, fsys fs.FS, p string, d fs.DirEntry, err error) error {
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
	aSong := new(Song)
	if GetFlags().JsonOutput {
		aSong = GetMetaData(pathArg, p)
		key, _ := GetEncoder().Encode(JustLetter(aSong.Title))
		aSong.titleH = key
		sMap[key] = *aSong
		return nil
	}
	aSong = parseFilename(pathArg, p)
	if aSong == nil {
		return nil
	}
	v := sMap[aSong.titleH]
	if len(v.titleH) > 0 {
		if aSong.artistH == v.artistH {
			fmt.Printf("#existing duplicate Song for %s %s == %s\n", aSong.inPath, aSong.Title, v.Title)
		} else {
			fmt.Printf("#possible dup Song for %s %s == %s %s\n", aSong.inPath, aSong.Title,
				v.Title, v.Artist)
			aSong.titleH += "1"
		}
		return nil
	}

	sMap[aSong.titleH] = *aSong
	return nil
}

// walk all files, looking for nice GoPro created video files.
// fill in a map keyed by the desired new name order
func WalkFiles(pathArg string) map[string]Song {
	songMap := make(map[string]Song)
	fsys := os.DirFS(pathArg)
	fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		err = processFile(pathArg, songMap, fsys, p, d, err)
		return nil
	})
	return songMap
}

// this is the output routine. it goes thru the map and produces output
// appropriate for the specified flag
func ProcessMap(pathArg string, m map[string]Song) map[string]Song {
	if GetFlags().JsonOutput {
		PrintJson(m)
		return m
	}
	uniqueArtists := make(map[string]string) // we just need a set, but use a map

	for _, aSong := range m {
		switch {
		case GetFlags().DoRename:
			cmd := "mv"
			if runtime.GOOS == "windows" {
				cmd = "ren "
			}
			switch {
			case aSong.alreadyNew:
				continue
				//fmt.Printf("pM aNew %s \"%s\" \"%s/%s; %s%s\"\n", cmd, aSong.inPath,
				//	pathArg, aSong.title, aSong.artist, aSong.ext)
			case aSong.Artist == "":
				fmt.Printf("#rename artist is blank %s\n", aSong.inPath)
				cmd = "#" + cmd
				continue
			case aSong.artistInDirectory:
				fmt.Printf("%s \"%s\" \"%s/%s%s\"\n", cmd, aSong.inPath,
					aSong.outPath, aSong.Title, aSong.ext)
				continue
			}
			if aSong.Artist == "" {
				fmt.Printf("%s \"%s\" \"%s/%s%s\"\n", cmd, aSong.inPath,
					aSong.outPath, aSong.Title, aSong.ext)
			} else {
				fmt.Printf("%s \"%s\" \"%s/%s; %s%s\"\n", cmd, aSong.inPath,
					aSong.outPath, aSong.Title, aSong.Artist, aSong.ext)
			}
		case GetFlags().JustList:
			the := ""
			if aSong.artistHasThe {
				the = "The "
			}
			fmt.Printf("%s by %s%s\n", aSong.Title, the, aSong.Artist)
		case GetFlags().ShowArtistNotInMap && !aSong.artistKnown:
			prim, _ := GetEncoder().Encode(JustLetter(aSong.Artist))
			if len(prim) == 0 && len(aSong.Artist) == 0 {
				// fmt.Printf("prim: %s, a: %v\n", prim, aSong)
				continue
			}
			uniqueArtists[prim] = aSong.Artist
		case GetFlags().NoGroup:
			if aSong.Artist == "" {
				fmt.Printf("nogroup %s\n", aSong.inPath)
			}
		default:
		}
	}
	if GetFlags().ShowArtistNotInMap {
		for k, v := range uniqueArtists {
			fmt.Printf("addto map k: %s v: %s\n", k, v)
		}
	}
	return m
}

// specialized function, dumps the Artist map
func DumpGptree() {
	if GetFlags().ZDumpArtist {
		gptree.Each(func(key string, v string) {
			fmt.Printf("\"%s\", \n", v)
		})
	}
}
