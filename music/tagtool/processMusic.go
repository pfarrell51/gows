// parseMusic
// does the main work of decending down the file structure, looking at music files
//
// this is not multi-processing safe

// bugs
// unicode search not implemented

package tagtool

import (
	"encoding/csv"
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

	if g.Flags().CSV {
		g.csvWrtr = csv.NewWriter(os.Stdout)
	}
	g.WalkFiles(pathArg)

	if g.Flags().JsonOutput {
		PrintJson(g.invSongs)
	}

	g.writeStoredData()
	g.doShutdown()
}

func (g *GlobalVars) writeStoredData() {
	g.PrintTrackSortedSongs()
}

// final close, flush, etc.
func (g *GlobalVars) doShutdown() {
	if g.Flags().CSV && g.csvWrtr != nil {
		g.csvWrtr.Flush()
		g.csvWrtr = nil
	}
	g.doSummary()
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
	if d == nil || strings.HasPrefix(p, ".") {
		return nil
	}
	if d.IsDir() {
		panic(fmt.Sprintf("impossible, pf sees directory %s\n", p))
	}
	extR := ExtRegex.FindStringIndex(p)
	if extR == nil {
		return nil // not interesting extension
	}
	if g.Flags().Debug {
		fmt.Printf("procssFile for %s\n", p)
	}
	rSong, _ := g.processSong(p)
	if rSong == nil {
		fmt.Printf("#processfile srong %s resulted in nil Song", p)
		return nil
	}
	g.songCount++
	rSong.FixupOutputPath(g)
	switch {
	case g.Flags().CopyAlbumInTrackOrder:
		g.AddSongForTrackSort(*rSong)
	case g.Flags().DoInventory:
		switch {
		case g.Flags().CSV:
			if g.Flags().Debug {
				fmt.Printf("PF will print song as CSV %s ! %s : %s\n",
					rSong.Artist, rSong.Album, rSong.Title)
			}
			g.PrintSongToCSV(rSong)
		case g.Flags().JsonOutput:
			// handles elsewhere
		default:
			fmt.Printf("Inv with no spec, %s,%s,%s\n", rSong.Artist, rSong.Album, rSong.Title)
		}
	}
	return nil
}

// isolate the actual work here, so we can cleanly
// set it off as a goroutine.
func (g *GlobalVars) processSong(p string) (*Song, error) {
	var err error
	var aSong *Song
	if g.Flags().Debug {
		fmt.Printf("procssSong for %s\n", p)
	}
	aSong, err = g.GetMetaData(p)
	if err != nil {
		return nil, err
	}
	if aSong == nil {
		panic(fmt.Sprintf("song %s resulted in nil Song", p))
	}
	if g.Flags().JsonOutput {
		if g.invSongs == nil {
			g.invSongs = make([]Song, 0)
		}
		g.invSongs = append(g.invSongs, *aSong)
	}
	return aSong, nil
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
			if p == "." {
				return nil
			}
			dir, _ := filepath.Split(p)
			switch {
			case dir == "":
				g.artistCount++
			case len(dir) > 0:
				g.albumCount++
			}
			return nil
		}

		extR := ExtRegex.FindStringIndex(p)
		if extR == nil {
			return nil // not interesting extension
		}

		var notOld = filepath.Dir(p)
		if oldDir != notOld && notOld != "." {
			oldDir = notOld
			g.numDirs++
		}
		if g.Flags().DoInventory && g.Flags().SuppressTitles && g.Flags().CSV {
			fmt.Println("skipping further processing")
			return nil // skip music files
		}

		g.processFile(fsys, p, d, err)
		return nil
	})
}

// prepopulate song structure with what we can know from the little we get from the user entered pathArg
// and the p lower path/filename.ext that we are walking through
func (s *Song) BasicPathSetup(g *GlobalVars, p string) {
	joined := filepath.FromSlash(path.Join(g.pathArg, p))
	s.InPath = joined
	s.inPathDescent, _ = path.Split(p) // ignore file name for now
	s.outPathBase = g.pathArg
	s.outPath = filepath.FromSlash(path.Join(g.pathArg, s.inPathDescent)) // start to build from here
	s.ext = path.Ext(p)
	if g.Flags().Debug {
		fmt.Printf("inpath %s\n", s.InPath)
		if s.inPathDescent != "" {
			fmt.Printf("inpath.descent %s\n", s.inPathDescent)
		}
		fmt.Printf("outpath %s\n", s.outPath)
	}
}

// build up output path in case we want to rename the file
func (s *Song) FixupOutputPath(g *GlobalVars) {
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
	if s.outPath == s.InPath {
		if g.Flags().Debug {
			fmt.Printf("#structs/fxop: no change for %s\n", s.InPath)
		}
	}
}
func (g *GlobalVars) doCompareTagsToTitle(aSong Song) {
	dir, fn := path.Split(aSong.InPath)
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
func (g *GlobalVars) doSummary() {
	fmt.Printf("#Songs: %d, Artists: %d, albums: %d, numDir: %d\n",
		g.songCount, g.artistCount, g.albumCount, g.numDirs)
}

// prints out a suitable rename/mv/ren command to put the file name
// in the format I like
func (g *GlobalVars) outputRenameCommand(aSong *Song) {
	cmd := "mv"
	if runtime.GOOS == "windows" {
		cmd = "ren "
	}
	if aSong.outPath == aSong.InPath {
		if g.Flags().Debug {
			fmt.Printf("#parseP no change for %s\n", aSong.InPath)
		}
		return
	}
	switch {
	case aSong.alreadyNew:
		fmt.Printf("#oRC  aNew %s \"%s\" \"%s - /%s; %s\"\n", cmd, aSong.InPath,
			aSong.Title, aSong.Artist, aSong.ext)
		return
	case aSong.Artist == "":
		fmt.Printf("#rename artist is blank %s\n", aSong.InPath)
		return
	case aSong.artistInDirectory:
		cmd = "#" + cmd
		fmt.Printf("%s \"%s\" \"%s\"\n", cmd, aSong.InPath, aSong.outPath)
		return
	}
	fmt.Printf("%s \"%s\" \"%s\"\n", cmd, aSong.InPath, aSong.outPath)
}
