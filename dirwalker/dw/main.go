package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/pfarrell51/gows/music/dirwalker"
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
	/*
		var globals = tagtool.AllocateData()
		if globals.GetSongTree() == nil {
			panic("main.go global songtree is nil")
		}
	*/

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

	pathArg := os.Args[1]

	fmt.Println(Files(pathArg))
	duration := time.Since(start)
	fmt.Printf("# %v\n", duration)

}
