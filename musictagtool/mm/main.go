// this program helps clean up the mp3 and flac files in my hits playlist
// its main task is to normalize the file names to relect the artist and song title
//
// this is not multi-processing safe

// Bugs

package main

import (
	"flag"
	"fmt"
	"github.com/pfarrell51/gows/musictagtool`"
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
	var flags = new(musictagtool.FlagST)
	flag.BoolVar(&flags.ShowArtistNotInMap, "a", false, "artist map -  list artist not in source code (gpmap)")
	flag.BoolVar(&flags.CopyAlbumInTrackOrder, "c", false, "Album track order - output cp command in track order")
	flag.BoolVar(&flags.Debug, "de", false, "debug on")
	flag.BoolVar(&flags.DuplicateDetect, "dup", false, "duplicate song attempts on")
	flag.BoolVar(&helpflag, "h", false, "help")
	flag.BoolVar(&flags.JsonOutput, "j", false, "output metadata as json")
	flag.BoolVar(&flags.JustList, "l", false, "list - list files")
	flag.BoolVar(&flags.NoGroup, "n", false, "nogroup - list files that do not have an artist/group in the title")
	flag.BoolVar(&flags.DoRenameMetadata, "mr", false, "rename - output rename from internal metadata")
	flag.BoolVar(&flags.ZDumpArtist, "z", false, "list artist names one per line")
	flag.Parse()
	if helpflag {
		flag.Usage()
		return
	}

	musictagtool.SetFlagArgs(*flags)
	musictagtool.LoadArtistMap()

	if false {
		ch := make(chan musictagtool.Song)
		v := <-ch
		fmt.Println(v)
	}
	if flags.ZDumpArtist {
		musictagtool.DumpGptree()
		return
	}
	pathArg := path.Clean(flag.Arg(0))
	ProcessFiles(pathArg)
	if flags.JsonOutput && flags.Debug {
		fmt.Printf("\n\n\n Dumping known ID names\n\n")
		musictagtool.DumpKnowIDnames()
	}
	duration := time.Since(start)
	fmt.Printf("# %v\n", duration)
}

func ProcessFiles(pathArg string) {
	if !musictagtool.GetFlags().ZDumpArtist {
		rmap := musictagtool.WalkFiles(pathArg)
		musictagtool.ProcessMap(pathArg, rmap)
	} else {
		musictagtool.DumpGptree()
	}
}
