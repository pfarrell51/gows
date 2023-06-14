// this program helps clean up the mp3 and flac files in my hits playlist
// its main task is to normalize the file names to relect the artist and song title
//
// this is not multi-processing safe

// Bugs

package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/pfarrell51/gows/music/randnames"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s [flags] directory-spec\n", os.Args[0])
		os.Exit(1)
	}
	start := time.Now()
	flag.Usage = func() {
		w := flag.CommandLine.Output() // may be os.Stderr - but not necessarily
		fmt.Fprintf(w, "Usage of %s: [flags] directory-spec \n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(w, "default is to list files that need love.\n")

	}
	var globals = randnames.AllocateData()
	if globals.GetSongTree() == nil {
		panic("main.go global songtree is nil")
	}

	var helpflag bool

	var flags = new(tagtool.FlagST)
	flag.StringVar(&flags.CpuProfile, "cpuprofile", "", "write cpu profile to `file`")
	flag.BoolVar(&helpflag, "h", false, "help")
	flag.BoolVar(&flags.JustList, "l", false, "list - list files")
	flag.BoolVar(&flags.NoTags, "nt", false, "notags - list files that do not have any meta tags")
	flag.BoolVar(&flags.TwoLetter, "tl", false, "two letter - list leading two letter words in titles")
	flag.BoolVar(&flags.ZDumpArtist, "z", false, "list artist names one per line")
	flag.Parse()
	if helpflag {
		flag.Usage()
		return
	}

	globals.SetFlagArgs(*flags)

	if flags.ZDumpArtist {
		globals.DumpGptree()
		return
	}
	pathArg := path.Clean(flag.Arg(0))
	//globals.ProcessFiles(pathArg)
	 
	duration := time.Since(start)
	fmt.Printf("# %v\n", duration)
}
