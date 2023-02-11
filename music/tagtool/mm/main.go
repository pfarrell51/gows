// this program helps clean up the mp3 and flac files in my hits playlist
// its main task is to normalize the file names to relect the artist and song title
//
// this is not multi-processing safe

// Bugs

package main

import (
	"flag"
	"fmt"
	"github.com/pfarrell51/gows/music/tagtool"
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
	var helpflag bool
	var flags = new(tagtool.FlagST)
	flag.BoolVar(&flags.ShowArtistNotInMap, "a", false, "artist map -  list artist not in source code (gpmap)")
	flag.BoolVar(&flags.CopyAlbumInTrackOrder, "c", false, "Album track order - output cp command in track order")
	flag.BoolVar(&flags.Debug, "de", false, "debug on")
	flag.BoolVar(&flags.DuplicateDetect, "dup", false, "duplicate song attempts on")
	flag.BoolVar(&helpflag, "h", false, "help")
	flag.BoolVar(&flags.DoInventory, "i", false, "inventory - basic inventory")
	flag.BoolVar(&flags.JsonOutput, "j", false, "output metadata as json")
	flag.BoolVar(&flags.JustList, "l", false, "list - list files")
	flag.BoolVar(&flags.NoGroup, "ng", false, "nogroup - list files that do not have an artist/group in the title")
	flag.BoolVar(&flags.NoTags, "nt", false, "notags - list files that do not have any meta tags")
	flag.BoolVar(&flags.DoRename, "r", false, "rename - output rename from internal metadata")
	flag.BoolVar(&flags.DoSummary, "s", false, "summary - print summary statistics")
	flag.BoolVar(&flags.ZDumpArtist, "z", false, "list artist names one per line")
	flag.Parse()
	if helpflag {
		flag.Usage()
		return
	}

	tagtool.SetFlagArgs(*flags)
	tagtool.LoadArtistMap()

	if false {
		ch := make(chan tagtool.Song)
		v := <-ch
		fmt.Println(v)
	}
	if flags.ZDumpArtist {
		tagtool.DumpGptree()
		return
	}
	pathArg := path.Clean(flag.Arg(0))
	ProcessFiles(pathArg)
	if flags.JsonOutput && flags.Debug {
		fmt.Printf("\n\n\n Dumping known ID names\n\n")
		tagtool.DumpKnowIDnames()
	}
	duration := time.Since(start)
	fmt.Printf("# %v\n", duration)
}

func ProcessFiles(pathArg string) {
	if tagtool.GetFlags().ZDumpArtist {
		tagtool.DumpGptree()
		return
	}
	rmap := tagtool.WalkFiles(pathArg)
	tagtool.ProcessMap(pathArg, rmap)
}
