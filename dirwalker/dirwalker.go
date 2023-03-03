// dirwalker
// utility to walk a directory tree and output cool commands

package dirwalker

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
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
	if !(verb == "ffmpeg" || verb == "sox") {
		fmt.Printf("unsupported verb: %s\n", verb)
		return 0
	}
	if newExt == "" {
		fmt.Printf("empty extension not supported\n")
		return 0
	}
	switch verb {
	case "ffmpeg":
		if outpath == "" {
			outpath = "mp3"
		}
		if outpath == "mp3" {
			outpath = strings.Replace(inpath, "flac", "mp3", -1)
		}
	case "sox":
		//fmt.Printf("PIB, figure it out unsupported verb: %s\n", verb)
		//fmt.Println("sox asz.wav asz-car.wav compand 0.3,1 6:-70,-60,-20 -5 -90 0.2")
	default:
		fmt.Printf("unsupported verb: %s\n", verb)
		return 0
	}
	_, inL := filepath.Split(inpath)
	_, outL := filepath.Split(outpath)

	if inL != outL {
		fmt.Printf("input and output paths not parallel,\n%s != \n%s\n", inL, outL)
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
			return nil
		}
		newP := ExtRegex.ReplaceAllString(p, newExt)
		count++
		dir, fn := filepath.Split(filepath.Clean(filepath.Join(outpath, newP)))
		err = os.MkdirAll(dir, 0777)
		if err != nil {
			panic(fmt.Sprintf("falled to make directory %s", dir))
		}

		useIn := filepath.Join(inpath, p)
		useOut := filepath.Join(dir, fn)
		switch verb {
		case "ffmpeg":
			fmt.Printf("%s -loglevel error -y -i \"%s\" -q:a 0 \"%s\"\n", verb, useIn, useOut)
		case "sox":
			fmt.Printf("%s \"%s\" \"%s\"  compand 0.3,1 6:-70,-60,-20 -5 -90 0.2\n", verb, useIn, useOut)
		default:
			fmt.Printf("unsupported verb: %s\n", verb)
		}
		return nil
	})
	return count
}
