// this pathgames:  path and file name manipulations

package musictools

import (
	"fmt"
	"path"
	"strings"
)

// prepopulate song structure with what we can know from the little we get from the user entered pathArg
// and the p lower path/filename.ext that we are walking through
func (s *Song) BasicPathSetup(pathArg, p string) {
	joined := path.Join(pathArg, p)
	s.inPath = joined
	s.inPathDescent, _ = path.Split(p) // ignore file name for now
	s.outPathBase = pathArg
	s.outPath = path.Join(pathArg, s.inPathDescent) // start to build from here
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
func (s Song) FixupOutputPath() {
	if s.outPath == s.inPath {
		return
	}
	if strings.Contains(s.outPath, s.Title) {
		panic(fmt.Sprintf("PIB, title aleady in output path %s\n%s\n", s.outPath, s.inPath))
	}
	if s.ext == "" {
		panic(fmt.Sprintf("PIB, extension is empty %s\n", s.outPath))
	}
	if s.ext != "" && strings.Contains(s.outPath, s.ext) {
		panic(fmt.Sprintf("PIB, extension >%s< aleady in output path %s\n", s.ext, s.outPath))
	}
	var outputPath = path.Join(s.outPath, s.Title)
	if !s.artistInDirectory {
		if s.Artist != "" {
			outputPath += "; " + s.Artist
		}
	}
	s.outPath = outputPath + s.ext
	if s.outPath == s.inPath {
		if GetFlags().Debug {
			fmt.Printf("#structs/fxop: no change for %s\n", s.inPath)
		}
	}
}
