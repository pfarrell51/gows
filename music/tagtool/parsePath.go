// parsePath
//
// this is not multi-processing safe

package tagtool

import (
	"fmt"
	"io/fs"
	"os"
	"regexp"
	"runtime"
	"sort"
	"strings"
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
var slashReg = regexp.MustCompile("\\" + string(os.PathSeparator))
var sortKeyDashExp = regexp.MustCompile("^[A-Z] - ")
var sortKeyUnderExp = regexp.MustCompile("^[A-Z]_")

// parse the last two directories of a song's path, trying to find which is the artist
func pathLastTwo(s string) (artist, album string) {
	panic("pathLastTwo called")
}

// identify the artist from two strings, a & b which were split from the file name
// looks up possilbe artists from the GroupTree hash map
func identifyArtistFromPair(a, b string, songpath string) (artist, title string) {
	panic("identifyArtistFromParts called")
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
	panic("parseFilename called")
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
	aSong, err = GetMetaData(pathArg, p)
	if err != nil {
		return err
	}
	key, _ := EncodeTitle(aSong.Title)
	aSong.titleH = key
	aSong.FixupOutputPath()
	sMap[key] = *aSong
	if GetFlags().CopyAlbumInTrackOrder {
		AddSongForSort(*aSong)
	}
	if GetFlags().DuplicateDetect {
		fmt.Printf("#dupDetect Not Yet Implemented\n")
	}
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
	if GetFlags().CopyAlbumInTrackOrder {
		PrintTrackSortedSongs()
		return m
	}
	uniqueArtists := make(map[string]Song)

	for _, aSong := range m {
		switch {
		case GetFlags().DoRename:
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

// prints out a suitable rename/mv/ren command to put the file name
// in the format I like
func outputRenameCommand(aSong *Song) {
	aSong.FixupOutputPath()
	cmd := "mv"
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
	if !GetFlags().ZDumpArtist {
		return
	}
	var arts []string
	Gptree.Each(func(key string, v string) {
		var t string
		if GetFlags().Debug {
			t = fmt.Sprintf("%s  \"%s\"", v, key)
		} else {
			t = fmt.Sprintf("%s", v)
		}
		arts = append(arts, t)
	})
	sort.Strings(arts)
	for _, v := range arts {
		fmt.Printf("%s\n", v)
	}
}
