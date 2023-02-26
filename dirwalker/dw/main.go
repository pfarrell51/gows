package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/pfarrell51/gows/dirwalker"
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

	var globals = dirwalker.AllocateData()

	var helpflag bool
	var flags = new(dirwalker.FlagST)
	flag.BoolVar(&flags.Debug, "de", false, "debug on")
	flag.BoolVar(&helpflag, "h", false, "help")
	flag.Parse()
	if helpflag {
		flag.Usage()
		return
	}

	globals.SetFlagArgs(*flags)

	pathArg := path.Clean(flag.Arg(0))

	fmt.Println(dirwalker.Files(pathArg))
	duration := time.Since(start)
	fmt.Printf("# %v\n", duration)

}
