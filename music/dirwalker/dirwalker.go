// dirwalker
// utility to walk a directory tree and output cool commands
//
// todo:
// detect the "soxdata.txt" file and don't re-compress using 'sox' command unless you really want to

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
	CopyAlbumInTrackOrder  bool
	CSV                    bool
	Debug                  bool
	SkipIfBreadcrumbExists bool
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
	g.localFlags.SkipIfBreadcrumbExists = f.SkipIfBreadcrumbExists
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
const norm = " --norm " // = "-v 0.98  --norm -G"
const attackDelay = "0.3,1"
const softKnee = "6" // from https://sox.sourceforge.net/sox.html
const transferFun = "-70,-60,-20"
const makeupGain = "-8"
const initialVolume = "-90"
const delay = "0.2"
const soxParams = verbosity + norm + attackDelay + softKnee + ":" + transferFun + " " + makeupGain + initialVolume + delay

var currentTime = time.Now().String()
var bytestamp = []byte(soxParams + " " + currentTime + "\n")

const pathSep = string(os.PathSeparator)

const interPathPart = "mp3u"
const BreadcrumbFN = "soxdata.txt"

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
	fpath := filepath.Join(dir, BreadcrumbFN)
	var err error
	if _, err = os.Stat(fpath); err == nil {
		if g.Flags().Debug {
			fmt.Printf("found breadcrumb for %s\n", dir)
		}
		return true // breadcrumb exists
	}
	return false
}
func (g GlobalVars) makeDirAndBreadcrumbFile(dir string) {
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		panic(fmt.Sprintf("falled to make directory %s", dir))
	}
	if g.verb == "ffmpeg" {
		return
	}
	fpath := filepath.Join(dir, BreadcrumbFN)
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

// process files, walking all of 'inpath' and creating the proper command
// and arguments to execute the verb with the processed files going
// to the parallel 'outpath' directory
func (g *GlobalVars) Files(verb, inpath, outpath string) (count int) {
	if !(verb == "ffmpeg" || verb == "sox" || verb == "both") {
		fmt.Printf("unsupported verb: %s\n", verb)
		return 0
	}
	g.verb = verb
	var songCount int
	var tmpPath = inpath
	switch verb {
	case "ffmpeg":
		if outpath == "mp3" {
			outpath = strings.Replace(inpath, "flac", "mp3", -1)
		}

	case "both":
		tmpPath = strings.Replace(inpath, "flac", interPathPart, -1)
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
		if d.IsDir() || !ExtRegex.MatchString(p) {
			return nil
		}
		if songCount++; songCount%500 == 0 {
			fmt.Printf("echo \"processing %d\"\n", songCount)
		}

		newP := ExtRegex.ReplaceAllString(p, "mp3")

		dir, fn := filepath.Split(filepath.Clean(filepath.Join(outpath, newP)))
		if g.ifExistsBreadcrumbfile(dir) {
			if g.localFlags.SkipIfBreadcrumbExists {
				if g.Flags().Debug {
					fmt.Printf("#breadcrumb found, skipping directory %s\n", dir)
				}
				return nil
			}
		}

		g.makeDirAndBreadcrumbFile(dir)
		count++
		useIn := filepath.Join(inpath, p)
		useFN := filepath.Join(tmpPath, newP)
		useOut := filepath.Join(dir, fn)
		switch verb {
		case "ffmpeg":
			fmt.Printf("%s -loglevel error -y -i \"%s\" -q:a 0 \"%s\"\n", verb, useIn, useOut)
		case "sox":
			fmt.Printf("%s %s %s \"%s\" \"%s\" compand %s %s:%s %s %s %s\n",
				verb, verbosity, norm, useIn, useOut, attackDelay, softKnee, transferFun, makeupGain, initialVolume, delay)
		case "both":
			fmt.Printf("%s -loglevel error -y -i \"%s\" -q:a 0 \"%s\"\n", "ffmpeg", useIn, useFN)
			fmt.Printf("%s %s %s \"%s\" \"%s\" compand %s %s:%s %s %s %s\n",
				"sox", verbosity, norm, useFN, useOut, attackDelay, softKnee, transferFun, makeupGain, initialVolume, delay)
		default:
			fmt.Printf("unsupported verb: %s\n", verb)
		}
		return nil
	})
	return count
}
