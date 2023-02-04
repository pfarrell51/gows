// this pathgames:  path and file name manipulations

package musictools

import (
	"fmt"
	"path"
	"regexp"
	"strings"
)

var slashFind = regexp.MustCompile("/")

// try to extract artist/group from diretory structure
func (s *Song) processInPathDirs() {
	if len(s.inPathDescent) == 0 {
		return
	}
	var p1, p2 string
	parts := slashFind.FindAllStringIndex(s.inPathDescent, -1)
	switch {
	case len(parts) == 1:
		p1 = s.inPathDescent[:parts[0][0]]
		p2 = ""
	case len(parts) >= 2:
		p1 = s.inPathDescent[:parts[0][0]]
		p2 = s.inPathDescent[parts[0][1]:parts[1][0]]
	}

	t1, t2 := EncodeArtist(p1)
	_, OKa := Gptree.Get(t1)
	_, OKb := Gptree.Get(t2)
	if OKa || OKb {
		s.Artist = p1
		s.Album = p2
		s.artistInDirectory = true
		s.artistKnown = true
	} else {
		t1, t2 = EncodeArtist(p2)
		_, OKa = Gptree.Get(t1)
		_, OKb = Gptree.Get(t2)
		if OKa || OKb {
			s.Artist = p2
			s.Album = p1
			s.artistInDirectory = true
			s.artistKnown = true
		}
	}
}

// prepopulate song structure with what we can know from the little we get from the user entered pathArg
// and the p lower path/filename.ext that we are walking through
func (s *Song) BasicPathSetup(pathArg, p string) {
	joined := path.Join(pathArg, p)
	s.inPath = joined
	s.inPathDescent, _ = path.Split(p) // ignore file name for now
	s.outPathBase = pathArg
	s.outPath = path.Join(pathArg, s.inPathDescent) // start to build from here
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
	if GetFlags().Debug {
		fmt.Printf("leaving FOP %s\n", s.outPath)
	}
	if s.outPath == s.inPath {
		if GetFlags().Debug {
			fmt.Printf("#structs/fxop: no change for %s\n", s.inPath)
		}
	}
}
