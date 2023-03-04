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
	"time"
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

// compand attack1,decay1{,attack2,decay2}
//
//	[soft-knee-dB:]in-dB1[,out-dB1]{,in-dB2,out-dB2}
//	[gain [initial-volume-dB [delay]]]
const verbosity = " -V2"
const norm = "-v 0.98  --norm -G"
const attackDelay = "0.3,1"
const softKnee = "6" // from https://sox.sourceforge.net/sox.html
const transferFun = "-70,-60,-20"
const makeupGain = "-4"
const initialVolume = "-90"
const delay = "0.2"
const soxParams = verbosity + norm + attackDelay + softKnee + ":" + transferFun + " " + makeupGain + initialVolume + delay

var currentTime = time.Now().String()
var bytestamp = []byte(soxParams + " " + currentTime + "\n")

const pathSep = string(os.PathSeparator)

func arePathsParallel(in, out string) bool {
	if in == out {
		return true
	}
	partsIn := strings.Split(in, string(filepath.Separator))
	partsOut := strings.Split(out, string(filepath.Separator))
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
func makeDirAndInfoFile(dir string) {
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		panic(fmt.Sprintf("falled to make directory %s", dir))
	}
	fpath := filepath.Join(dir, "soxdata.txt")
	file, err := os.Create(fpath)
	file.Write(bytestamp)
	file.Close()
}
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
		//fmt.Println("sox asz.wav asz-car.wav compand 0.3,1 6:-70,-60,-20 -5 -90 0.2")
	default:
		fmt.Printf("unsupported verb: %s\n", verb)
		return 0
	}
	if !arePathsParallel(inpath, outpath) {
		fmt.Printf("input and output paths not parallel,\n%s != \n%s\n", inpath, outpath)
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
		makeDirAndInfoFile(dir)
		useIn := filepath.Join(inpath, p)
		useOut := filepath.Join(dir, fn)
		switch verb {
		case "ffmpeg":
			fmt.Printf("%s -loglevel error -y -i \"%s\" -q:a 0 \"%s\"\n", verb, useIn, useOut)
		case "sox":
			fmt.Printf("%s %s %s \"%s\" \"%s\" compand %s %s:%s %s %s %s\n",
				verb, verbosity, norm, useIn, useOut, attackDelay, softKnee, transferFun, makeupGain, initialVolume, delay)
		default:
			fmt.Printf("unsupported verb: %s\n", verb)
		}
		return nil
	})
	return count
}
