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

	"github.com/pfarrell51/gows/music/tagtool"
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
	var globals = tagtool.AllocateData()

	var helpflag bool

	var flags = new(tagtool.FlagST)
	flag.BoolVar(&flags.CSV, "csv", false, "output CSV format")
	flag.BoolVar(&flags.CompareTagsToTitle, "ctt", false, "compare Tags to Title")
	flag.BoolVar(&flags.CopyAlbumInTrackOrder, "c", false, "Album track order - output cp command in track order")
	flag.BoolVar(&flags.Debug, "de", false, "debug on")
	flag.BoolVar(&flags.DoInventory, "i", false, "inventory - basic inventory")
	flag.BoolVar(&flags.DoRename, "r", false, "rename - output rename from internal metadata")
	flag.BoolVar(&flags.DoSummary, "s", false, "summary - print summary statistics")
	flag.BoolVar(&helpflag, "h", false, "help")
	flag.BoolVar(&flags.JsonOutput, "json", false, "output metadata as json")
	flag.BoolVar(&flags.JustArtistAlbum, "jaa", false, "output Artist & Album (no song title)")
	flag.BoolVar(&flags.NoGroup, "ng", false, "nogroup - list files that do not have an artist/group in the title")
	flag.BoolVar(&flags.NoTags, "nt", false, "notags - list files that do not have any meta tags")
	flag.BoolVar(&flags.UnicodePunct, "u", false, "show songs with Unicode punct")
	flag.Parse()
	if helpflag {
		flag.Usage()
		return
	}
	globals.SetFlagArgs(*flags)

	pathArg := path.Clean(flag.Arg(0))
	globals.ProcessFiles(pathArg)

	duration := time.Since(start)
	fmt.Printf("# %v\n", duration)
}
