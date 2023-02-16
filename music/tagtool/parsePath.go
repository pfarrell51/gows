// parsePath
//
// this is not multi-processing safe

// bugs
// do not remove level if album name is same as song title
// count physical albums

package tagtool

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"

	g "github.com/zyedidia/generic"
	"github.com/zyedidia/generic/avl"
)

func init() {
	LoadArtistMap()
}

var songsProcessed int
var numNoAcoustId, numNoTitle, numNoMBID int
var numAlbums, numArtists int
var artistTree = avl.New[string, int](g.Less[string])
var albumTree = avl.New[string, int](g.Less[string])
var songTree = avl.New[string, int](g.Less[string])

// top level entry point, takes the path to the directory to walk/process
func ProcessFiles(pathArg string) {
	if GetFlags().ZDumpArtist {
		DumpGptree()
		return
	}
	rmap := WalkFiles(pathArg)
	ProcessMap(pathArg, rmap)
}

var ExtRegex = regexp.MustCompile(`[Mm][Pp][34]|[Ff][Ll][Aa][Cc]$`) //"((M|m)(p|P)(3|4))|((F|f)(L|l)(A|a)(C|c))$")

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
	songsProcessed++
	key, _ := EncodeTitle(aSong.Title)
	aSong.titleH = key
	aSong.FixupOutputPath()
	sMap[key] = *aSong
	if GetFlags().NoTags {
		if aSong.AcoustID == "" {
			numNoAcoustId++
		}
		if aSong.Title == "" {
			numNoTitle++
		}
		if aSong.MBID == "" {
			numNoMBID++
		}
		if aSong.AcoustID == "" && aSong.Title == "" && aSong.MBID == "" {
			fmt.Printf("#No tags found for %s\n", aSong.inPath)
		}
	}

	if GetFlags().DoSummary {
		updateUniqueCounts(*aSong)
	}
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

// prepopulate song structure with what we can know from the little we get from the user entered pathArg
// and the p lower path/filename.ext that we are walking through
func (s *Song) BasicPathSetup(pathArg, p string) {
	joined := filepath.FromSlash(path.Join(pathArg, p))
	s.inPath = joined
	s.inPathDescent, _ = path.Split(p) // ignore file name for now
	s.outPathBase = pathArg
	s.outPath = filepath.FromSlash(path.Join(pathArg, s.inPathDescent)) // start to build from here
	s.ext = path.Ext(p)
	if GetFlags().Debug {
		fmt.Printf("inpath %s\n", s.inPath)
		if s.inPathDescent != "" {
			fmt.Printf("inpath.descent %s\n", s.inPathDescent)
		}
		fmt.Printf("outpath %s\n", s.outPath)
		fmt.Printf("outpath.base %s\n", s.outPathBase)
		fmt.Printf("ext: %s\n", s.ext)
	}
}

// build up output path in case we want to rename the file
func (s *Song) FixupOutputPath() {
	if GetFlags().Debug {
		fmt.Printf("FOP %s\n", s.outPath)
	}
	if s.ext == "" {
		panic(fmt.Sprintf("PIB, extension is empty %s\n", s.outPath))
	}
	if s.Album == s.Title {
		s.outPath = path.Join(s.outPath, s.Title)
	} else if !strings.Contains(s.outPath, s.Title) {
		s.outPath = path.Join(s.outPath, s.Title)
	}
	if !s.artistInDirectory && s.Artist != "" && !strings.Contains(s.outPath, s.Artist) {
		s.outPath += "; " + s.Artist
	}
	if !strings.Contains(s.outPath, s.ext) {
		s.outPath = s.outPath + s.ext
	}
	s.outPath = filepath.FromSlash(s.outPath)
	if GetFlags().Debug {
		fmt.Printf("leaving FOP %s\n", s.outPath)
	}
	if s.outPath == s.inPath {
		if GetFlags().Debug {
			fmt.Printf("#structs/fxop: no change for %s\n", s.inPath)
		}
	}
}

// called for each song, we see if we have this artist/album/song in the appropriate
// tree, and either increment it or insert it with 1 (seen once)
func updateUniqueCounts(s Song) {
	v, ok := artistTree.Get(s.Artist)
	if ok {
		v++
	} else {
		v = 1
	}
	artistTree.Put(s.Artist, v)
	v, ok = artistTree.Get(s.Artist)

	key := s.Artist + " " + s.Album
	v, ok = albumTree.Get(key)
	if ok {
		v++
	} else {
		v = 1
	}
	albumTree.Put(key, v)

	key = s.Title + " " + s.Artist
	v, ok = songTree.Get(key)
	if ok {
		v++
	} else {
		v = 1
	}
	songTree.Put(key, v)
}

// this is the output routine. it goes thru the map and produces output
// appropriate for the specified flag, unless -r was specified, there
// is a separate routine to output the rename command

func ProcessMap(pathArg string, m map[string]Song) {
	if GetFlags().JsonOutput {
		PrintJson(m)
		return
	}
	if GetFlags().CopyAlbumInTrackOrder {
		PrintTrackSortedSongs()
		return
	}
	uniqueArtists := make(map[string]Song)
	var countSongs, countNoGroup int
	for _, aSong := range m {
		countSongs++
		switch {
		case GetFlags().DoRename:
			outputRenameCommand(&aSong)
		case GetFlags().JustList:
			fmt.Printf("%s by %s\n", aSong.Title, aSong.Artist)
			continue
		case GetFlags().ShowArtistNotInMap && !aSong.artistKnown:
			if aSong.Artist == "" {
				countNoGroup++
				continue
			}
			prim, sec := EncodeArtist(aSong.Artist)
			_, ok := Gptree.Get(prim)
			if ok {
				// fmt.Printf("primary found %s %s\n", prim, aSong.Artist)
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
				countNoGroup++
				fmt.Printf("nogroup %s\n", aSong.inPath)
			}
		default:
		}
	}

	if GetFlags().NoGroup {
		fmt.Printf("#scanned %d songs, %d had no artist/group\n", countSongs, countNoGroup)
	}
	if GetFlags().NoTags {
		fmt.Printf("#scanned %d songs, %d had no artist, %d no AcoustId, %d no title, %d no MBID\n",
			countSongs, countNoGroup, numNoAcoustId, numNoTitle, numNoMBID)
	}
	if GetFlags().ShowArtistNotInMap {
		for k, v := range uniqueArtists {
			fmt.Printf("addto map k: %s v: %s %s\n", k, v.Artist, v.inPath)
		}
	}
	if GetFlags().DoSummary {
		if GetFlags().Debug {
			fmt.Println("artists. Count is number of songs across all albums for this artist")
			artistTree.Each(func(k string, v int) {
				fmt.Printf("%d %s\n", v, k)
			})
			fmt.Printf(">#2 found %d artists, %d albums and %d songs\n",
				artistTree.Size(), albumTree.Size(), songTree.Size())
		}
		if GetFlags().Debug {
			fmt.Println("albums. Count is number of songs in the given artist/album")
			albumTree.Each(func(k string, v int) {
				fmt.Printf("%d %s\n", v, k)
			})
			fmt.Printf(">#3 found %d artists, %d albums and %d songs\n",
				artistTree.Size(), albumTree.Size(), songTree.Size())
		}
		if GetFlags().Debug {
			fmt.Println("songs")
			songTree.Each(func(k string, v int) {
				fmt.Printf("%d %s\n", v, k)
			})
		}
		fmt.Printf("found %d artists, %d albums and %d songs\n", artistTree.Size(), albumTree.Size(), songTree.Size())
	}
	return
}

// prints out a suitable rename/mv/ren command to put the file name
// in the format I like
func outputRenameCommand(aSong *Song) {
	aSong.FixupOutputPath()
	cmd := "mv"
	if runtime.GOOS == "windows" {
		cmd = "ren "
	}
	// fmt.Printf("#oRC start  %s \"%s\" \"%s-/%s; %s\"\n", cmd, aSong.inPath,
	// 	aSong.Title, aSong.Artist, aSong.ext)
	if aSong.outPath == aSong.inPath {
		if GetFlags().Debug {
			fmt.Printf("#parseP no change for %s\n", aSong.inPath)
		}
		return
	}
	switch {
	case aSong.alreadyNew:
		fmt.Printf("#oRC  aNew %s \"%s\" \"%s - /%s; %s\"\n", cmd, aSong.inPath,
			aSong.Title, aSong.Artist, aSong.ext)
		return
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
