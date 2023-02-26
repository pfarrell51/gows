// dirwalker
// utility to walk a directory tree and output cool commands

package dirwalker

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
)

type FlagST struct {
	CopyAlbumInTrackOrder bool
	CSV                   bool
	Debug                 bool
}
type GlobalVars struct {
	pathArg    string
	outPath    string
	verb       string
	localFlags *FlagST
}

func (g *GlobalVars) Flags() *FlagST {
	return g.localFlags
}

// copy user set flags to a local store
func (g *GlobalVars) SetFlagArgs(f FlagST) {
	g.localFlags = new(FlagST)
	g.localFlags.CopyAlbumInTrackOrder = f.CopyAlbumInTrackOrder
	g.localFlags.CSV = f.CSV
	g.localFlags.Debug = f.Debug
}
func AllocateData() *GlobalVars {
	rval := new(GlobalVars)
	rval.localFlags = new(FlagST)
	return rval
}

var ExtRegex = regexp.MustCompile("((M|m)(p|P)(3|4))|((F|f)(L|l)(A|a)(C|c))$")

func Files(verb, inpath, outpath, newExt string) (count int) {
	if verb != "ffmpeg" {
		fmt.Printf("unsupported verb: %s\n", verb)
		return 0
	}
	fsys := os.DirFS(inpath)
	fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Println(err)
			return nil
		}
		if d.IsDir() {
			return nil
		}
		if !ExtRegex.MatchString(p) {
			fmt.Printf("#not interesting: %s\n", p)
			return nil
		}
		newP := ExtRegex.ReplaceAllString(p, newExt)
		count++
		dir, fn := filepath.Split(filepath.Clean(filepath.Join(outpath, newP)))
		err = os.MkdirAll(dir, 0777)
		if err != nil {
			panic(fmt.Sprintf("falled to make directory %s", dir))
		}
		fmt.Printf("%s -loglevel error -y -i \"%s\" -q:a 0 \"%s\"\n", verb, filepath.Join(inpath, p), filepath.Join(dir, fn))
		return nil
	})
	return count
}
