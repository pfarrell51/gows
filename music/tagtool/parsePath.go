// parsePath
//
// this is not multi-processing safe

// bugs
// unicode search not implemented

package tagtool

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/texttheater/golang-levenshtein/levenshtein"
)

// top level entry point, takes the path to the directory to walk/process
func (g *GlobalVars) ProcessFiles(pathArg string) {
	if g.Flags() == nil {
		fmt.Println("PIB, g.Flags() is nil")
	}
	if len(g.pathArg) == 0 {
		g.pathArg = pathArg
	}
	g.WalkFiles(pathArg)
	g.ProcessMap()
}

var ExtRegex = regexp.MustCompile(`[Mm][Pp][34]|[Ff][Ll][Aa][Cc]$`)

// this is the local  WalkDirFunc called by WalkDir for each file
// p is the current path/name
func (g *GlobalVars) processFile(fsys fs.FS, p string, d fs.DirEntry, err error) error {
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
	rSong, _ := g.processSong(p)
	if rSong == nil {
		fmt.Printf("#processfile srong %s resulted in nil Song", p)
		return nil
	}
	g.getSong(rSong)
	return nil
}

// isolate the actual work here, so we can cleanly
// set it off as a goroutine.
func (g *GlobalVars) processSong(p string) (*Song, error) {
	var err error
	var aSong *Song
	aSong, err = g.GetMetaData(p)
	if err != nil {
		return nil, err
	}
	if aSong == nil {
		panic(fmt.Sprintf("song %s resulted in nil Song", p))
	}
	return aSong, nil
}
func (g *GlobalVars) getSong(aSong *Song) {
	if aSong == nil {
		panic(fmt.Sprintf("getsong has nil Song"))
	}
	g.songsProcessed++
	key, _ := EncodeTitle(aSong.Title)
	combKey := key
	if aSong.Artist == "" {
		combKey += "#"
	} else {
		if aSong.artistH == "" {
			tmp, _ := EncodeArtist(aSong.Artist)
			aSong.artistH = tmp
			combKey += aSong.artistH
		} else {
			combKey += aSong.artistH
		}
	}
	if aSong.Album != "" {
		combKey += aSong.albumH
	} else {
		combKey += "+"
	}
	aSong.smapKey = combKey
	aSong.FixupOutputPath(g)
	switch {
	case g.Flags().NoTags:
		if aSong.AcoustID == "" {
			g.numNoAcoustId++
		}
		if aSong.Title == "" {
			g.numNoTitle++
		}
		if aSong.MBID == "" {
			g.numNoMBID++
		}
		if aSong.AcoustID == "" && aSong.Title == "" && aSong.MBID == "" {
			fmt.Printf("#No tags found for %s\n", aSong.inPath)
		}
	case g.Flags().DoSummary:
		fmt.Println("no imp for DoSummery")
	case g.Flags().CopyAlbumInTrackOrder:
		g.AddSongForTrackSort(*aSong)
	case g.Flags().DoInventory:
		g.AddSongForTripleSort(*aSong)
	}
}

// walk all files, looking for nice music files.
func (g *GlobalVars) WalkFiles(pathArg string) {
	g.pathArg = pathArg
	var oldDir string
	fsys := os.DirFS(pathArg)
	fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("Oh No, can't process because %s\n", err)
			return err
		}
		if d.IsDir() {
			fmt.Println(p)
			return nil
		}
		var notOld = filepath.Dir(p)
		if oldDir != notOld && notOld != "." {
			oldDir = notOld
			g.numDirs++
		}
		// g.processFile(fsys, p, d, err)
		return nil
	})
}

// prepopulate song structure with what we can know from the little we get from the user entered pathArg
// and the p lower path/filename.ext that we are walking through
func (s *Song) BasicPathSetup(g *GlobalVars, p string) {
	joined := filepath.FromSlash(path.Join(g.pathArg, p))
	s.inPath = joined
	s.inPathDescent, _ = path.Split(p) // ignore file name for now
	s.outPathBase = g.pathArg
	s.outPath = filepath.FromSlash(path.Join(g.pathArg, s.inPathDescent)) // start to build from here
	s.ext = path.Ext(p)
	if g.Flags().Debug {
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
func (s *Song) FixupOutputPath(g *GlobalVars) {
	if g.Flags().Debug {
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
	if g.Flags().Debug {
		fmt.Printf("leaving FOP %s\n", s.outPath)
	}
	if s.outPath == s.inPath {
		if g.Flags().Debug {
			fmt.Printf("#structs/fxop: no change for %s\n", s.inPath)
		}
	}
}


// this is the output routine. it goes thru the map and produces output
// appropriate for the specified flag, unless -r was specified, there
// is a separate routine to output the rename command

func (g *GlobalVars) ProcessMap() {
	switch {
	case g.Flags().JsonOutput:
	fmt.Println("no impl for JsonOutput")
		return
	case g.Flags().CopyAlbumInTrackOrder:
		g.PrintTrackSortedSongs()
		return
	case g.Flags().DoInventory:
		g.doInventory()
		return
	case g.Flags().DupTitleAlbumArtist:
		for i, d := range g.dupSongs {
			fmt.Printf("%2d d.a %s\n   d.b %s\n", i, d.a.inPath, d.b.inPath)
		}
		return
	}

	if g.Flags().NoTags {
		fmt.Printf("#scanned  %d no AcoustId, %d no title, %d no MBID\n",
			g.numNoAcoustId, g.numNoTitle, g.numNoMBID)
	}
	if g.Flags().DoSummary {
		g.doSummary()
	}
}
func (g *GlobalVars) doSearchForUnicodePunct(aSong Song) {
	dir, fn := path.Split(aSong.inPath)
	fname := strings.TrimSuffix(fn, filepath.Ext(fn))
	fmt.Println("not yet implemented %s %s and %s\n", dir, fn, fname)
}
func (g *GlobalVars) doCompareTagsToTitle(aSong Song) {
	dir, fn := path.Split(aSong.inPath)
	fname := strings.TrimSuffix(fn, filepath.Ext(fn))
	var ttl, art string
	if strings.Contains(fname, ";") {
		parts := strings.Split(fname, ";")
		if parts == nil || len(parts) == 0 {
			panic("file name has no parts but has a ;")
		}
		ttl = parts[0]
		art = parts[1]
	} else {
		ttl = fname
	}
	var disTTL, disART int
	disTTL = levenshtein.DistanceForStrings([]rune(ttl), []rune(aSong.Title), levenshtein.DefaultOptions)
	if art != "" {
		disART = levenshtein.DistanceForStrings([]rune(art), []rune(aSong.Artist), levenshtein.DefaultOptions)
	}
	if disTTL > 4 || disART > 4 {
		fmt.Printf(" dT: %d, dA: %d, d: %s  fn:%s\n", disTTL, disART, dir, fn)
	}
}
func (g *GlobalVars) doInventory() {
	g.PrintTripleSortedSongs()
}
func (g *GlobalVars) doSummary() {
	if g.Flags().Debug {
		fmt.Println("artists. Count is number of songs across all albums for this artist")
		g.artistCountTree.Each(func(k string, v int) {
			fmt.Printf("%d %s\n", v, k)
		})

		fmt.Printf(">#2 found %d artists, %d albums (%d directories) and %d songs\n",
			g.artistCountTree.Size(), g.albumCountTree.Size(), g.numDirs, g.songCountTree.Size())
		fmt.Println("albums. Count is number of songs in the given artist/album")
		g.albumCountTree.Each(func(k string, v int) {
			fmt.Printf("%d %s\n", v, k)
		})
		fmt.Printf(">#3 found %d artists, %d albums and %d songs\n",
			g.artistCountTree.Size(), g.albumCountTree.Size(), g.songCountTree.Size())

		if g.Flags().Debug {
			fmt.Println("songs")
			g.songCountTree.Each(func(k string, v int) {
				fmt.Printf("%d %s\n", v, k)
			})
		}
	}
	fmt.Printf("found %d artists, %d albums (%d directories) and %d songs or sP %d\n",
		g.artistCountTree.Size(), g.albumCountTree.Size(), g.numDirs, g.songCountTree.Size(), g.songsProcessed)
}

// prints out a suitable rename/mv/ren command to put the file name
// in the format I like
func (g *GlobalVars) outputRenameCommand(aSong *Song) {
	cmd := "mv"
	if runtime.GOOS == "windows" {
		cmd = "ren "
	}
	if aSong.outPath == aSong.inPath {
		if g.Flags().Debug {
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
		return
	case aSong.artistInDirectory:
		cmd = "#" + cmd
		fmt.Printf("%s \"%s\" \"%s\"\n", cmd, aSong.inPath, aSong.outPath)
		return
	}
	fmt.Printf("%s \"%s\" \"%s\"\n", cmd, aSong.inPath, aSong.outPath)
}
