// msync
// utility to quickly copy a tree of music files
//

package msync

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type FlagST struct {
	Debug      bool
	ProfileCPU bool
	ProfileMEM bool
	Verify     bool
}
type GlobalVars struct {
	inPath     string
	outPath    string
	walkedDir  []string
	localFlags *FlagST
	location   *time.Location
	cacheFN    string
	cacheRval  bool
}

// holds one line with the p path to of the file and h hash value
type HashLine struct {
	p string
	h string
}
type HashedDir struct {
	p       string
	t       time.Time
	n       int        // number of music files
	lines   []HashLine // one for each music file
	dirHash string     // hash of contained Hashlines
}

var ExtRegex = regexp.MustCompile("((M|m)(p|P)(3|4))|((F|f)(L|l)(A|a)(C|c))$")

const PermBitsFile = 0766 // unix style permissions for txt file
const PermBitsDir = 0777  // unix style permissions for directory
const SumFN = "msync.sum"

var UTCloc, _ = time.LoadLocation("UTC")
var LastCentury = time.Date(1999, time.December, 31, 0, 0, 0, 0, time.UTC)

func (hl *HashLine) String() string {
	return hl.p + " " + hl.h + "\n"
}
func newHashedDir() *HashedDir {
	var rval = new(HashedDir)
	rval.lines = make([]HashLine, 0, 15)
	return rval
}
func (g *GlobalVars) Flags() *FlagST {
	return g.localFlags
}
func (hd *HashedDir) addHash(p, h string) {
	hl := HashLine{p: p, h: h}
	hd.lines = append(hd.lines, hl)
	hd.n++
}

func (hd *HashedDir) String() string {
	var rval strings.Builder
	if hd.t.Before(LastCentury) {
		hd.t = time.Now()
	}

	rval.WriteString(fmt.Sprintf("%s\n", hd.t.In(UTCloc).Format(time.RFC822)))
	rval.WriteString(fmt.Sprintf("%s %d\n", hd.p, hd.n))
	for i := 0; i < hd.n; i++ {
		rval.WriteString(hd.lines[i].String())
	}
	rval.WriteString(fmt.Sprintf("%s\n", hd.dirHash))
	return rval.String()
}

// copy user set flags to a local store
func (g *GlobalVars) SetFlagArgs(f FlagST) {
	g.localFlags = new(FlagST)
	g.localFlags.Debug = f.Debug
	g.localFlags.ProfileCPU = f.ProfileCPU
	g.localFlags.ProfileMEM = f.ProfileMEM
	g.localFlags.Verify = f.Verify
}
func AllocateData() *GlobalVars {
	rval := new(GlobalVars)
	rval.localFlags = new(FlagST)
	rval.walkedDir = make([]string, 100)
	rval.location, _ = time.LoadLocation("UTC")
	return rval
}

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

func (g GlobalVars) ifExistsBreadcrumbfile(dir string) bool {
	if dir == g.cacheFN {
		return g.cacheRval
	}
	g.cacheFN = dir
	fpath := filepath.Join(dir, SumFN)

	info, err := os.Stat(fpath) // warning for possible race condition
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			g.cacheRval = false
			return false // fmt.Println("error is ErrNotExist")
		} else {
			fmt.Printf("can't Stat file %s got %v\n", fpath, err)
			g.cacheRval = false
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
		g.cacheRval = true
		return true // breadcrumb exists
	}
}

func (g GlobalVars) makeDirAndInfoFile(dir string) {
	err := os.MkdirAll(dir, PermBitsDir)
	if err != nil {
		panic(fmt.Sprintf("falled to make directory %s", dir))
	}
	fpath := filepath.Join(dir, SumFN)
	if g.ifExistsBreadcrumbfile(dir) {
		fmt.Println("mDIG found it") // breadcrumb exists
	} else {
		fmt.Printf("mDIF will touch %s\n", fpath)
		// breadcrumb file does *not* exist
		file, err := os.OpenFile(fpath, os.O_RDWR|os.O_CREATE|os.O_EXCL, PermBitsFile)
		if err != nil {
			fmt.Printf("create error: %s", err)
		}
		file.Write(g.getBytestamp())
		file.Close()
	}
}
func writeDirSumFile(g *GlobalVars, path string, sumLines []string, hs string) {
	info, err := os.Stat(path) // warning for possible race condition
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return // fmt.Println("error is ErrNotExist")
		} else {
			fmt.Printf("can't Stat file %s got %v\n", path, err)
			return
		}
	}
	if !info.IsDir() {
		return
	}
	fpath := filepath.Join(path, SumFN)
	file, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE, PermBitsFile)
	if err != nil {
		fmt.Printf("create error: %s", err)
		file.Close()
		return
	}
	defer file.Close()
	file.Write(g.getBytestamp())
	if len(sumLines) == 0 {
		return
	}
	t := fmt.Sprintf("%s %d\n", path, len(sumLines))
	fmt.Printf("%s", t)
	fmt.Fprintf(file, "%s", t)
	for _, l := range sumLines {
		t = fmt.Sprintf("wDSD %s\n", l)
		fmt.Printf("%s", t)
		fmt.Fprintf(file, "%s", t)
	}
	t = fmt.Sprintf("hDir %s\n", hs) //base64.StdEncoding.EncodeToString(h.Sum(nil)))
	fmt.Printf("%s", t)
	fmt.Fprintf(file, "%s", t)
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
	var count int
	var oldInDir string
	var sumLines = make([]string, 0, 10)
	hsDir := newHashedDir()
	hashedDirEntries := sha256.New()
	fsys := os.DirFS(inp)
	fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		joined := filepath.Clean(filepath.Join(inp, p))
		if d.IsDir() {
			g.walkedDir = append(g.walkedDir, joined)
			if joined != oldInDir {
				if oldInDir != "" {
					dH := base64.StdEncoding.EncodeToString(hashedDirEntries.Sum(nil))
					writeDirSumFile(g, oldInDir, sumLines, dH)
					fmt.Printf(hsDir.String())
				}
				oldInDir = joined
				sumLines = make([]string, 0, 10)
				hashedDirEntries = sha256.New()
				hsDir = newHashedDir()
				hsDir.p = joined
				hsDir.t = time.Now()
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
		fileHash := base64.StdEncoding.EncodeToString(result)
		fmt.Fprintf(hashedDirEntries, "%x  %s\n", result, joined) // add in line to hashedDirEntries (whole directory hash)
		hsDir.addHash(joined, fileHash)
		nameSum := fmt.Sprintf("nameSum: %s %s", joined, fileHash)
		sumLines = append(sumLines, nameSum)
		return nil
	})
	fmt.Println("#falling out of dir walk, last line")
	dirH := base64.StdEncoding.EncodeToString(hashedDirEntries.Sum(nil))
	writeDirSumFile(g, oldInDir, sumLines, dirH)
	fmt.Printf("final hashedDir\n%s\n", hsDir.String())
	return
}
