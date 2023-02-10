// this pathgames:  path and file name manipulations

package musictagtool

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

const slashPathSep = "\\" + string(os.PathSeparator)

var slashFind = regexp.MustCompile(slashPathSep)

// try to extract artist/group from directory structure
func (s *Song) processInPathDirs() {
	if len(s.inPathDescent) == 0 {
		return
	}
	var p1, p2 string
	if GetFlags().Debug {
		fmt.Printf("s.inPD: %s, all: %s\n", s.inPathDescent, s.inPath)
	}
	parts := slashFind.FindAllStringIndex(s.inPath, -1)

	switch len(parts) {
	case 1:
		panic(fmt.Sprintf("PIB, not enough slashes in %s", s.inPath))
	case 2:
		p1 = s.inPath[parts[0][1]:parts[1][0]]
		p2 = s.inPath[parts[1][1]:]
	default:
		lp := len(parts) - 1
		for i := 0; i < 3; i++ {
			j := lp - (i + 1)
			k := lp - i
			p1 = s.inPath[parts[j][1]:parts[k][0]]
			j--
			k--
			if j < 0 || k < 0 {
				fmt.Printf("Panic coming soon i: %d, j: %d, k: %d, %s\n", i, j, k, s.inPath)
				panic(fmt.Sprintf("negative j %d or k %d ", j, k))
			}
			p2 = s.inPath[parts[j][1]:parts[k][0]]
			p3 := s.inPath[parts[lp][1]:]
			var prim string
			s.Title = StandardizeTitle(p3)
			s.titleH, _ = EncodeTitle(p3)
			prim, t2 := EncodeArtist(p1)
			_, OKa := Gptree.Get(prim)
			_, OKb := Gptree.Get(t2)
			if OKa || OKb {
				s.Artist = StandardizeArtist(p1)
				s.artistH = prim
				s.Album = p2
				s.artistInDirectory = true
				s.artistKnown = true
				break
			} else {
				prim, t2 = EncodeArtist(p2)
				_, OKa = Gptree.Get(prim)
				_, OKb = Gptree.Get(t2)
				if OKa || OKb {
					s.Artist = StandardizeArtist(p2)
					s.artistH = prim
					s.Album = p1
					s.artistInDirectory = true
					s.artistKnown = true
					break
				}
			}
		}
		if GetFlags().Debug {
			fmt.Printf("art: %s, album: %s, t: %s %t %t\n", s.Artist, s.Album, s.Title,
				s.artistInDirectory, s.artistKnown)
		}
	}
}

// prepopulate song structure with what we can know from the little we get from the user entered pathArg
// and the p lower path/filename.ext that we are walking through
func (s *Song) BasicPathSetup(pathArg, p string) {
	joined := filepath.FromSlash(path.Join(pathArg, p))
	s.inPath = joined
	s.inPathDescent, _ = path.Split(p) // ignore file name for now
	s.outPathBase = pathArg
	s.outPath = filepath.FromSlash(path.Join(pathArg, s.inPathDescent)) // start to build from here
	s.processInPathDirs()
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
	if s.outPath == s.inPath {
		return
	}
	if s.ext == "" {
		panic(fmt.Sprintf("PIB, extension is empty %s\n", s.outPath))
	}
	if !strings.Contains(s.outPath, s.Title) {
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
