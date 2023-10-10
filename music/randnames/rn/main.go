// this program deals with creatging random prefixes for music files
// and managing them
// its main task is to normalize the file names to relect the artist and song title
//
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

	var helpflag bool

	var flags = new(randnames.FlagST)
	flag.BoolVar(&flags.AddTag, "add", false, "add prefix")
	flag.BoolVar(&helpflag, "h", false, "help")
	flag.BoolVar(&flags.Debug, "d", false, "debug")
	flag.BoolVar(&flags.RmTag, "rm", false, "rm (remove) prefix")
	flag.BoolVar(&flags.TwoLetter, "tl", false, "two letter - list leading two letter words in titles")
	flag.Parse()
	if helpflag {
		flag.Usage()
		return
	}
	randnames.SetFlags(*flags)

	pathArg := path.Clean(flag.Arg(0))
	randnames.WalkFiles(pathArg)

	duration := time.Since(start)
	fmt.Printf("# %v\n", duration)
}
