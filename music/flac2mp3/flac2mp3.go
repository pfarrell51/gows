// flac2mp3
// utility to walk a directory tree and output cool commands
//

package flac2mp3

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

// parameters for ffmpeg and sox
// compand attack1,decay1{,attack2,decay2}
//
//	[soft-knee-dB:]in-dB1[,out-dB1]{,in-dB2,out-dB2}
//	[gain [initial-volume-dB [delay]]]
const audioOutputParams = " -b:a 350k -q:a 0 " // for ffmpeg
const verbosity = " -V2"
const norm = " --norm " // = "-v 0.98  --norm -G"
const betweenSwitch = "-C 320 "
const attackDelay = "0.3,1"
const softKnee = "6" // from https://sox.sourceforge.net/sox.html
const transferFun = "-70,-60,-20"
const makeupGain = "-8"
const initialVolume = "-90"
const delay = "0.2"
const soxParams = verbosity + norm + attackDelay + softKnee + ":" + transferFun + " " + makeupGain + initialVolume + delay

//  sox -V2 --norm infile.flac -C 320 outfile.mp3 compand 0.3,1 6:-70,-60,-20 -8 -90 0.2

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

// process files, walking all of 'inpath' and creating the proper command
// and arguments to execute the verb with the processed files going
// to the parallel 'outpath' directory
func (g *GlobalVars) Files(verb, inpath, outpath string) (count int) {
	// if !(verb == "ffmpeg" || verb == "sox" || verb == "both") {
	if verb != "sox" {
		fmt.Printf("unsupported verb: %s\n", verb)
		return 0
	}
	g.verb = verb
	var songCount int

	if outpath == "mp3" {
		outpath = strings.Replace(inpath, "flac", "mp3", -1)
	}
	if !arePathsParallel(inpath, outpath) {
		fmt.Printf("input and output paths not parallel,\n%s != \n%s\n", inpath, outpath)
		return 0
	}

	fsys := os.DirFS(inpath)
	fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Println("error walking dir")
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

		useIn := filepath.Join(inpath, p)
		useOut := filepath.Join(dir, fn)
		fmt.Printf("sox %s %s \"%s\" %s \"%s\" compand %s %s:%s %s %s %s\n", verbosity, norm,
			useIn, betweenSwitch, useOut, attackDelay, softKnee, transferFun, makeupGain, initialVolume, delay)
		return nil
	})
	return songCount
}
