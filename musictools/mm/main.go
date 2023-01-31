// this program helps clean up the mp3 and flac files in my hits playlist
// its main task is to normalize the file names to relect the artist and song title
//
// this is not multi-processing safe

// Bugs

package main

import (
	"flag"
	"fmt"
	"github.com/pfarrell51/gows/musictools"
	"os"
	"path"
	"time"
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
	var flags = new(musictools.FlagST)
	flag.BoolVar(&flags.ShowArtistNotInMap, "a", false, "artist map -  list artist not in source code (gpmap)")
	flag.BoolVar(&flags.JustList, "l", false, "list - list files")
	flag.BoolVar(&flags.NoGroup, "n", false, "nogroup - list files that do not have an artist/group in the title")
	flag.BoolVar(&flags.DoRename, "r", false, "rename - output command to perform rename function on needed files")
	flag.BoolVar(&flags.ZDumpArtist, "z", false, "list artist names one per line")
	flag.BoolVar(&flags.JsonOutput, "j", false, "output metadata as json")
	flag.BoolVar(&flags.Debug, "d", false, "debug on")
	flag.Parse()
	musictools.SetFlagArgs(*flags)

	if false {
		ch := make(chan musictools.Song)
		v := <-ch
		fmt.Println(v)
	}

	pathArg := path.Clean(flag.Arg(0))
	ProcessFiles(pathArg)
	duration := time.Since(start)
	fmt.Printf("# %v\n", duration)
}

func ProcessFiles(pathArg string) {
	if !musictools.GetFlags().ZDumpArtist {
		rmap := musictools.WalkFiles(pathArg)
		musictools.ProcessMap(pathArg, rmap)
	} else {
		musictools.DumpGptree()
	}
}
