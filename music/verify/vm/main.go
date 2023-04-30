package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/pfarrell51/gows/music/verify"
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

	var flags = new(verify.FlagST)
	flag.BoolVar(&flags.Debug, "de", false, "debug on")
	flag.BoolVar(&helpflag, "h", false, "help")
	flag.Parse()
	if helpflag {
		flag.Usage()
		return
	}

	pathArg := path.Clean(flag.Arg(0))
	verify.WalkFiles(flags, pathArg)
	duration := time.Since(start)
	fmt.Printf("# %v\n", duration)
}
