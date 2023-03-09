// msync
// utility to quickly copy a tree of music files
//

package msync

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type FlagST struct {
	Debug bool
}
type GlobalVars struct {
	inPath     string
	outPath    string
	walkedDir  []string
	localFlags *FlagST
}

func (g *GlobalVars) Flags() *FlagST {
	return g.localFlags
}

// copy user set flags to a local store
func (g *GlobalVars) SetFlagArgs(f FlagST) {
	g.localFlags = new(FlagST)
	g.localFlags.Debug = f.Debug
}
func AllocateData() *GlobalVars {
	rval := new(GlobalVars)
	rval.localFlags = new(FlagST)
	rval.walkedDir = make([]string, 100)
	return rval
}

var ExtRegex = regexp.MustCompile("((M|m)(p|P)(3|4))|((F|f)(L|l)(A|a)(C|c))$")

var currentTime = time.Now().String()
var bytestamp = []byte(" " + currentTime + "\n")

const SumFN = "msync.txt"
const pathSep = string(os.PathSeparator)

func arePathsParallel(in, out string) bool {
	if in == out {
		return true
	}
	partsIn := strings.Split(in, pathSep)
	partsOut := strings.Split(out, pathSep)
	if len(partsIn) != len(partsOut) {
		return false
	}
	var cDiff int
	for i := 0; i < len(partsIn); i++ {
		if partsIn[i] != partsOut[i] {
			cDiff++
		}
		if cDiff > 1 {
			return false
		}
	}
	if cDiff == 1 {
		return true
	}
	return false
}

func (g GlobalVars) ifExistsBreadcrumbfile(dir string) bool {
	fpath := filepath.Join(dir, SumFN)
	var err error
	if _, err = os.Stat(fpath); err == nil {
		if g.Flags().Debug {
			fmt.Printf("found breadcrumb for %s\n", dir)
		}
		return true // breadcrumb exists
	}
	return false
}
func (g GlobalVars) makeDirAndInfoFile(dir string) {
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		panic(fmt.Sprintf("falled to make directory %s", dir))
	}
	fpath := filepath.Join(dir, SumFN)
	if g.ifExistsBreadcrumbfile(dir) {
		// breadcrumb exists
	} else {
		file, err := os.Create(fpath) // breadcrumb file does *not* exist
		if err != nil {
			fmt.Printf("create error: %s", err)
		}
		file.Write(bytestamp)
		file.Close()
	}
}
func writeDirSumFile(g *GlobalVars, sumLines []string) {
}

// process files, walking all of 'inpath' and creating the proper command
// and arguments to execute the verb with the processed files going
// to the parallel 'outpath' directory with the extension of 'newExt'
// I expect that newExt will always be 'mp3' but lets see over time
func (g *GlobalVars) WalkDirectories(inpath, outpath string) {
	var count int
	if !arePathsParallel(inpath, outpath) {
		fmt.Printf("input and output paths not parallel,\n%s != \n%s\n", inpath, outpath)
		return
	}
	g.inPath = inpath
	g.outPath = outpath
	var oldOutDir string
	var sumLines = make([]string, 10)
	fsys := os.DirFS(inpath)
	fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}

		t := filepath.Clean(filepath.Join(outpath, p))
		dir, fn := filepath.Split(t)
		if d.IsDir() {
			g.walkedDir = append(g.walkedDir, dir)
		}
		if dir != oldOutDir {
			writeDirSumFile(g, sumLines)
			oldOutDir = dir
			sumLines = make([]string, 10)
		}
		//g.makeDirAndInfoFile(dir)
		if g.ifExistsBreadcrumbfile(dir) {
			if g.Flags().Debug {
				fmt.Printf("#breadcrumb found, skipping directory %s\n", dir)
			}
		}
		if !ExtRegex.MatchString(p) {
			return nil
		}

		if count++; count%500 == 0 {
			fmt.Printf("echo \"processing %d\"\n", count)
		}

		if g.ifExistsBreadcrumbfile(dir) {
			if g.Flags().Debug {
				fmt.Printf("#breadcrumb found, skipping directory %s\n", dir)
			}
		}

		count++
		useIn := filepath.Join(inpath, p)
		useOut := filepath.Join(dir, fn)
		fmt.Printf("i: %s, o: %s\n", useIn, useOut)
		return nil
	})
	return
}
