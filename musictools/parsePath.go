// parsePath
//
// this is not multi-processing safe

package musictools

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"regexp"
	"runtime"
	"strings"

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

var ExtRegex = regexp.MustCompile("((M|m)(p|P)(3|4))|((F|f)(L|l)(A|a)(C|c))$")
var underToSpace = regexp.MustCompile("_")
var dReg = regexp.MustCompile("-\\s")
var commaExp = regexp.MustCompile(",\\s")
var slashReg = regexp.MustCompile("/")
var sortKeyDashExp = regexp.MustCompile("^[A-Z] - ")
var sortKeyUnderExp = regexp.MustCompile("^[A-Z]_")

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
	return StandardizeArtist(artist), album
}
func identifyArtistFromPair(a, b string, songpath string) (artist, title string) {
	prim, sec := EncodeArtist(a)
	_, ok := Gptree.Get(prim)
	if ok {
		return StandardizeArtist(a), b
	}
	_, ok = Gptree.Get(sec)
	if ok {
		return StandardizeArtist(a), b
	}
	prim, sec = EncodeArtist(b)
	_, ok = Gptree.Get(prim)
	if ok {
		return StandardizeArtist(b), a
	}
	_, ok = Gptree.Get(sec)
	if ok {
		return StandardizeArtist(b), a
	}
	fmt.Printf("did not match group name with either argument, %s, %s for %s\n", a, b, songpath)
	switch {
	case a == "" && b != "":
		return a, b // assume b is title with no group
	case a != "" && b == "":
		return b, a // assume a is title with no group
	case a == "" && b == "":
		return a, b // makes no difference
	case a != "" && b != "":
		return a, b // just a guess
	}
	panic("PIB, in identify.ArtistFromPair")
}

const divP = " -+" // want space for names like Led Zeppelin - Bron-Yr-Aur
var dashRegex = regexp.MustCompile(divP)

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
	if sortKeyDashExp.Match(nameB) {
		nameB = nameB[4:]
	}
	if sortKeyUnderExp.Match(nameB) {
		nameB = nameB[2:]
	}
	nameB = underToSpace.ReplaceAll(nameB, []byte(" "))
	extR := ExtRegex.FindIndex(nameB)
	if extR == nil || len(extR) == 0 {
		return rval
	}
	ext := path.Ext(p)
	rval.ext = ext
	if rval.artistInDirectory && rval.artistKnown && rval.Title != "" {
		return rval
	}
	nameB = nameB[0 : extR[0]-1]
	var groupN, SongN string
	groupN = rval.Artist // default to using the name from the subdirectory

	ps, _ := path.Split(string(nameB))
	if len(ps) > 0 {
		rval.outPath = path.Join(rval.outPath, ps)
		if !rval.artistInDirectory && !rval.artistKnown {
			rval.artistInDirectory = true
			tArtist, tAlbum := pathLastTwo(ps)
			if tArtist == rval.Artist {
				rval.Artist = tArtist
				rval.Album = tAlbum
			} else {
				fmt.Printf("artist names not same two ways %s != %s in %s\n", tArtist, rval.Artist, p)
			}
		}
		nameB = nameB[len(ps):]
		SongN = StandardizeTitle(string(nameB))
	}
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
			groupN, SongN = identifyArtistFromPair(sa, sb, p)
		} else {
			commaLoc := commaExp.FindIndex(nameB)
			if len(commaLoc) > 0 {
				sa := string(nameB[:commaLoc[0]]) // have a comma, so split by it
				sb := string(nameB[commaLoc[1]:])
				groupN, SongN = identifyArtistFromPair(sa, sb, p)
			} else {
				groupN = ""
				SongN = string(nameB)
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
	if GetFlags().DuplicateDetect {
		v := sMap[aSong.titleH]
		if len(v.titleH) > 0 {
			if aSong.artistH == v.artistH {
				fmt.Printf("#existing duplicate Song for %s %s == %s\n", aSong.inPath, v.inPath)
			} else {
				fmt.Printf("#possible dup Song or cover for %s %s == %s %s %s\n", aSong.inPath, v.inPath)
				aSong.titleH += "1"
			}
			return nil
		}
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
			fmt.Printf("%s by %s\n", aSong.Title, aSong.Artist)
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
