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

	"github.com/pfarrell51/gows/music/filenametool"
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
	var helpflag bool
	var flags = new(filenametool.FlagST)
	flag.BoolVar(&flags.ShowArtistNotInMap, "a", false, "artist map -  list artist not in source code (data/artists.txt)")
	flag.BoolVar(&flags.Debug, "de", false, "debug on")
	flag.BoolVar(&flags.DuplicateDetect, "dup", false, "duplicate song -"+
		" attempts to identify duplicate songs, very buggy")
	flag.BoolVar(&helpflag, "h", false, "help")
	flag.BoolVar(&flags.JustList, "l", false, "list - list files")
	flag.BoolVar(&flags.NoGroup, "ng", false, "nogroup - list files that do not have an artist/group in the title")
	flag.BoolVar(&flags.DoRename, "re", false, "rename - output rename from path/filename")
	flag.BoolVar(&flags.ZDumpArtist, "z", false, "list artist names one per line")
	flag.Parse()
	if helpflag {
		flag.Usage()
		return
	}

	filenametool.SetFlagArgs(*flags)
	filenametool.LoadArtistMap()

	if false {
		ch := make(chan filenametool.Song)
		v := <-ch
		fmt.Println(v)
	}
	if flags.ZDumpArtist {
		filenametool.DumpGptree()
		return
	}
	pathArg := path.Clean(flag.Arg(0))
	ProcessFiles(pathArg)
	duration := time.Since(start)
	fmt.Printf("# %v\n", duration)
}

func ProcessFiles(pathArg string) {
	if filenametool.GetFlags().ZDumpArtist {
		filenametool.DumpGptree()
		return
	}
	rmap := filenametool.WalkFiles(pathArg)
	filenametool.ProcessMap(pathArg, rmap)
}
