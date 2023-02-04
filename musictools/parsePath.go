// parsePath section parses artist/album/song title from filename
//
// this is not multi-processing safe

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
var ExtRegex = regexp.MustCompile("((M|m)(p|P)(3|4))|((F|f)(L|l)(A|a)(C|c))$")
var underToSpace = regexp.MustCompile("_")

var cReg = regexp.MustCompile(",\\s")
var dReg = regexp.MustCompile("-\\s")
var commaExp = regexp.MustCompile(",\\s")
var slashReg = regexp.MustCompile("/")

func pathLastTwo(s string) (artist, album string) {
	if matched, _ := regexp.MatchString("^[A-Z]( |-|_)", s); matched {
		s = s[2:]
	}
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
	if sortKeyExp.MatchString(artist) {
		fmt.Printf("removing sort key from %s\n", artist)
		artist = artist[2:]
	}
	return StandardizeArtist(artist), album
}

// parse the file info to find artist and Song title
// most of my music files have file names with the artist name, a hyphen and then the track title
// so this pulls out the information and fills in the "Song" object.
func parseFilename(pathArg, p string) *Song {
	if GetFlags().Debug {
		fmt.Printf("\npf: %s\n", p)
	}
	var rval = new(Song)
	rval.BasicPathSetup(pathArg, p)
	nameB := []byte(strings.TrimSpace(p))
	if sortKeyExp.Match(nameB) {
		nameB = nameB[2:]
	}
	nameB = underToSpace.ReplaceAll(nameB, []byte(" "))
	extR := ExtRegex.FindIndex(nameB)
	if extR == nil || len(extR) == 0 {
		return rval
	}
	ext := path.Ext(p)
	rval.ext = ext
	nameB = nameB[0 : extR[0]-1]

	var groupN, SongN string
	groupN = rval.Artist // default to using the name from the subdirectory
	ps, _ := path.Split(string(nameB))
	if len(ps) > 0 {
		rval.outPath = path.Join(rval.outPath, ps)
		rval.artistInDirectory = true
		rval.Artist, rval.Album = pathLastTwo(ps)
		nameB = nameB[len(ps):]
		SongN = StandardizeTitle(string(nameB))
	}
	words := cReg.FindAllIndex(nameB, -1)
	dash := dReg.FindIndex(nameB)
	var semiExp = regexp.MustCompile("; ")
	if semiExp.Match(nameB) {
		semiLoc := semiExp.FindIndex(nameB)
		SongN = StandardizeTitle(string(nameB[:semiLoc[0]]))
		groupN = StandardizeArtist(string(nameB[semiLoc[1]:]))
		rval.alreadyNew = true
	} else {
		//  no ;
		if dash != nil {
			sa := string(nameB[:dash[0]]) // have a dash, so split by it
			sb := string(nameB[dash[1]:])
			if len(words) == 0 {
				groupN = StandardizeArtist(sa)
				SongN = StandardizeTitle(sb)
			} else {
				ta, ta2 := EncodeTitle(sa)
				_, OKa := Gptree.Get(ta)
				_, OKa2 := Gptree.Get(ta2)
				if OKa || OKa2 {
					SongN = StandardizeTitle(sb)
					groupN = StandardizeArtist(sa)
				}
				tb, tb2 := EncodeArtist(sb)
				_, OKb := Gptree.Get(tb)
				_, OKb2 := Gptree.Get(tb2)
				if OKb || OKb2 {
					groupN = StandardizeArtist(sb)
					SongN = StandardizeTitle(sa)
				}
				if !OKa && !OKa2 && !OKb && !OKb2 {
					fmt.Printf("not OK in hash %s\n", string(nameB))
					fmt.Printf("ta %s ta2 %s, tb %s Tb2 %s\n", ta, ta2, tb, tb2)
				}
			}
		} else {
			commasS := commaExp.FindIndex(nameB)
			if commasS == nil || len(commasS) == 0 {
				SongN = string(nameB)
			} else {
				SongN = string(nameB[:commasS[0]])
				groupN = string(nameB[commasS[1]:])
			}
		}
	}
	groupN = cases.Title(language.English, cases.NoLower).String(StandardizeArtist(groupN))
	rval.Title = StandardizeTitle(SongN)
	rval.titleH, _ = EncodeTitle(SongN)
	rval.Artist = StandardizeArtist(groupN)
	rval.artistH, _ = EncodeArtist(rval.Artist)
	_, ok := Gptree.Get(rval.artistH)
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
	extR := ExtRegex.FindStringIndex(p)
	if extR == nil {
		return nil // not interesting extension
	}
	var aSong *Song
	if GetFlags().JsonOutput || GetFlags().DoRenameMetadata {
		aSong, err = GetMetaData(pathArg, p)
		if err != nil {
			return err
		}
		key, _ := EncodeTitle(aSong.Title)
		aSong.titleH = key
		aSong.FixupOutputPath()
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
			fmt.Printf("#possible dup Song for %s %s == %s %s %s\n", aSong.inPath, aSong.Title,
				v.Title, v.Artist, v.inPath)
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
	uniqueArtists := make(map[string]Song)

	for _, aSong := range m {
		switch {
		case GetFlags().DoRenameFilename || GetFlags().DoRenameMetadata:
			outputRenameCommand(&aSong)
		case GetFlags().JustList:
			continue
		case GetFlags().ShowArtistNotInMap && !aSong.artistKnown:
			if aSong.Artist == "" {
				continue
			}
			prim, sec := EncodeArtist(aSong.Artist)
			_, ok := Gptree.Get(prim)
			if ok {
				fmt.Printf("primary found %s %s\n", prim, aSong.Artist)
				continue
			}
			if len(sec) > 0 {
				_, ok = Gptree.Get(sec)
				if ok {
					fmt.Printf("sec found %s %s\n", prim, aSong.Artist)
					continue
				}
			}
			fmt.Printf("not found %s %s\n%s\n", prim, aSong.Artist, aSong.inPath)
			uniqueArtists[prim] = aSong
			if len(sec) > 0 {
				uniqueArtists[sec] = aSong
			}
		case GetFlags().NoGroup:
			if aSong.Artist == "" {
				fmt.Printf("nogroup %s\n", aSong.inPath)
			}
		default:
		}
	}
	if GetFlags().ShowArtistNotInMap {
		for k, v := range uniqueArtists {
			fmt.Printf("addto map k: %s v: %s %s\n", k, v.Artist, v.inPath)
		}
	}
	return m
}
func outputRenameCommand(aSong *Song) {
	cmd := "mv"
	aSong.FixupOutputPath()
	if runtime.GOOS == "windows" {
		cmd = "ren "
	}
	if aSong.outPath == aSong.inPath {
		if GetFlags().Debug {
			fmt.Printf("#parseP no change for %s\n", aSong.inPath)
		}
		return
	}
	switch {
	case aSong.alreadyNew:
		return
		//fmt.Printf("pM aNew %s \"%s\" \"%s/%s; %s%s\"\n", cmd, aSong.inPath,
		//	aSong.title, aSong.artist, aSong.ext)
	case aSong.Artist == "":
		fmt.Printf("#rename artist is blank %s\n", aSong.inPath)
		cmd = "#" + cmd
		return
	case aSong.artistInDirectory:
		cmd = "#" + cmd
		fmt.Printf("%s \"%s\" \"%s\"\n", cmd, aSong.inPath, aSong.outPath)
		return
	}
	fmt.Printf("%s \"%s\" \"%s\"\n", cmd, aSong.inPath, aSong.outPath)
}

// specialized function, dumps the Artist map
func DumpGptree() {
	if GetFlags().ZDumpArtist {
		Gptree.Each(func(key string, v string) {
			fmt.Printf("\"%s\", \n", v)
		})
	}
}
