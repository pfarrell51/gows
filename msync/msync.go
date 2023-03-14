// msync
// utility to quickly copy a tree of music files
//

package msync

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"hash"
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
	location   *time.Location
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
	rval.location, _ = time.LoadLocation("UTC")
	return rval
}

var ExtRegex = regexp.MustCompile("((M|m)(p|P)(3|4))|((F|f)(L|l)(A|a)(C|c))$")

const PermBits = 0766 // unix style permissions
const SumFN = "msync.txt"

func (g *GlobalVars) getBytestamp() []byte {
	s := g.toUTC(time.Now())
	return []byte(s + "\n")
}

func dummy() {
	now := time.Now()
	loc, err := time.LoadLocation("UTC")
	if err != nil {
		panic(err)
	}
	fmt.Println(now)
	fmt.Println(now.In(loc))
	s := now.In(loc).Format(time.RFC822)
	fmt.Println(s)
	t, e := time.Parse(time.RFC822, s)
	fmt.Printf("time: %v, e: %v\n", t, e)
}
func (g GlobalVars) toUTC(t time.Time) string {
	rval := t.In(g.location).Format(time.RFC822)
	return rval
}

// these should be in the GV struct, this is a bug farm
var cacheFN string
var cacheRval bool

func (g GlobalVars) ifExistsBreadcrumbfile(dir string) bool {
	if dir == cacheFN {
		return cacheRval
	}
	cacheFN = dir
	fpath := filepath.Join(dir, SumFN)

	info, err := os.Stat(fpath) // warning for possible race condition
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			cacheRval = false
			return false // fmt.Println("error is ErrNotExist")
		} else {
			fmt.Printf("can't Stat file %s got %v\n", fpath, err)
			cacheRval = false
			return false
		}
	} else {
		if true || g.Flags().Debug {
			fmt.Printf("found breadcrumb for %s\n", fpath)
		}
		if info == nil {
			fmt.Printf("most strange, nil info for %s\n", fpath)
		}
		if info.IsDir() {
			fmt.Printf("%s is a directory\n", fpath)
		}
		modTime := info.ModTime()
		fmt.Printf("%s mod time %s\n", fpath, g.toUTC(modTime))
		cacheRval = true
		return true // breadcrumb exists
	}
}

func (g GlobalVars) makeDirAndInfoFile(dir string) {
	fmt.Printf("mDIF %s\n", dir)
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		panic(fmt.Sprintf("falled to make directory %s", dir))
	}
	fpath := filepath.Join(dir, SumFN)
	if g.ifExistsBreadcrumbfile(dir) {
		fmt.Println("mDIG found it") // breadcrumb exists
	} else {
		fmt.Printf("mDIF will touch %s\n", fpath)
		// breadcrumb file does *not* exist
		file, err := os.OpenFile(fpath, os.O_RDWR|os.O_CREATE|os.O_EXCL, PermBits)
		if err != nil {
			fmt.Printf("create error: %s", err)
		}
		file.Write(g.getBytestamp())
		file.Close()
	}
}
func writeDirSumFile(g *GlobalVars, path string, sumLines []string, h hash.Hash) {
	info, err := os.Stat(path) // warning for possible race condition
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			cacheRval = false
			return // fmt.Println("error is ErrNotExist")
		} else {
			fmt.Printf("can't Stat file %s got %v\n", path, err)
			cacheRval = false
			return
		}
	}
	if !info.IsDir() {
		return
	}
	fpath := filepath.Join(path, SumFN)
	file, err := os.OpenFile(fpath, os.O_RDWR|os.O_CREATE|os.O_EXCL, PermBits)
	if err != nil {
		fmt.Printf("create error: %s", err)
	}
	defer file.Close()
	file.Write(g.getBytestamp())
	if len(sumLines) == 0 {
		return
	}
	t := fmt.Sprintf("%s %d\n", path, len(sumLines))
	fmt.Println(t)
	fmt.Fprintf(file, t)
	for _, l := range sumLines {
		t = fmt.Sprintf("%s\n", l)
		fmt.Println(t)
		fmt.Fprintf(file, t)
	}
	t = fmt.Sprintf("%s\n", base64.StdEncoding.EncodeToString(h.Sum(nil)))
	fmt.Println(t)
	fmt.Fprintf(file, t)
}

// Called from main with both input and output paths.
// actual processing is split.
// do input processing first
func (g *GlobalVars) WalkDirectories(inpath, outpath string) {
	g.inPath = inpath
	g.WalkInputDirectories(inpath)
	return
}

// process files, walking all of 'inpath'
func (g *GlobalVars) WalkInputDirectories(inp string) {
	fmt.Printf("\nWID with %s\n", inp)
	var count int
	var oldInDir string
	var sumLines = make([]string, 0, 10)
	hDir := sha256.New()
	fsys := os.DirFS(inp)
	fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		joined := filepath.Clean(filepath.Join(inp, p))
		if d.IsDir() {
			fmt.Printf("in loop, have directory  %s\n", joined)
			g.walkedDir = append(g.walkedDir, joined)

			if joined != oldInDir {
				fmt.Printf("#  oldIn: %s, joined: %s\n", oldInDir, joined)
				if oldInDir != "" {
					fmt.Printf("#1 oldIn: %s, current: %s\n", oldInDir, joined)
					writeDirSumFile(g, oldInDir, sumLines, hDir)
				}
				oldInDir = joined
				sumLines = make([]string, 0, 10)
				hDir = sha256.New()
			}
		}
		if !ExtRegex.MatchString(p) {
			return nil // only sum on flac and mp3 files
		}
		if count++; count%500 == 0 {
			fmt.Printf("echo \"processing %d\"\n", count)
		}
		count++
		hFile := sha256.New()
		f, err := os.Open(joined)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		if _, err := io.Copy(hFile, f); err != nil {
			log.Fatal(err)
		}
		result := hFile.Sum(nil)
		fmt.Fprintf(hDir, "%x  %s\n", result, joined) // add in line to hDir (whole directory hash)
		nameSum := fmt.Sprintf("%s %s", joined, base64.StdEncoding.EncodeToString(result))
		sumLines = append(sumLines, nameSum)
		return nil
	})
	fmt.Println("#falling out of dir walk, last line")
	writeDirSumFile(g, oldInDir, sumLines, hDir)
	return
}
