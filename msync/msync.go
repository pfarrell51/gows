// msync
// utility to quickly copy a tree of music files
//

package msync

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

type FlagST struct {
	Debug      bool
	ProfileCPU bool
	ProfileMEM bool
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
	g.localFlags.ProfileCPU = f.ProfileCPU
	g.localFlags.ProfileMEM = f.ProfileMEM
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
func writeDirSumFile(g *GlobalVars, path string, sumLines []string) {
	fmt.Printf("wDSF %s, %d\n", path, len(sumLines))
}

// process files, walking all of 'inpath' and creating the proper command
// and arguments to execute the verb with the processed files going
// to the parallel 'outpath' directory with the extension of 'newExt'
// I expect that newExt will always be 'mp3' but lets see over time
func (g *GlobalVars) WalkDirectories(inpath, outpath string) {
	var count int
	g.inPath = inpath
	g.outPath = outpath
	var oldInDir string
	var sumLines = make([]string, 0, 10)
	hDir := sha256.New()
	fsys := os.DirFS(inpath)
	fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		t := filepath.Clean(filepath.Join(inpath, p))
		dir, _ := filepath.Split(t)
		fmt.Printf("t: %s  dir: %s\n", t, dir)
		if d.IsDir() {
			g.walkedDir = append(g.walkedDir, dir)
		}
		if dir != oldInDir {
			fmt.Printf("# oldIn: %s, current: %s\n", oldInDir, dir)
			writeDirSumFile(g, t, sumLines)
			oldInDir = dir
			sumLines = make([]string, 0, 10)
			hDir = sha256.New()
		}
		hFile := sha256.New()
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
		//	useOut := filepath.Join(dir, fn)
		f, err := os.Open(useIn)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		if _, err := io.Copy(hFile, f); err != nil {
			log.Fatal(err)
		}
		result := hFile.Sum(nil)
		fmt.Fprintf(hDir, "%x  %s\n", result, useIn) // add in line to hDir (whole directory hash)
		nameSum := fmt.Sprintf("%s %s", useIn, base64.StdEncoding.EncodeToString(result))
		sumLines = append(sumLines, nameSum)
		return nil
	})
	fmt.Printf("%d\n", len(sumLines))
	for _, l := range sumLines {
		fmt.Printf("%s\n", l)
	}
	fmt.Printf("hDir: %s\n", base64.StdEncoding.EncodeToString(hDir.Sum(nil)))
	return
}
