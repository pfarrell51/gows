// main cmd program for gedfilter package
package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/pfarrell51/gows/gedstuff/gedfilter"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s [flags] gedcom-filespec\n", os.Args[0])
		os.Exit(1)
	}
	var helpflag bool

	gv := gedfilter.AllocateData()
	var flags = gv.Flags()
	flag.BoolVar(&flags.Basic, "basic", false, "name, birth, death")
	flag.BoolVar(&flags.CSV, "csv", false, "output CSV format")
	flag.BoolVar(&flags.Debug, "debug", false, "debug on")
	flag.BoolVar(&helpflag, "h", false, "help")
	flag.BoolVar(&flags.Type, "type", false, "display Type fields")
	flag.Usage = func() {
		w := flag.CommandLine.Output() // may be os.Stderr - but not necessarily
		fmt.Fprintf(w, "Usage of %s: [flags] gedcom-spec \n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	if helpflag {
		flag.Usage()
		return
	}
	start := time.Now()
	gv.ProcessFile(flag.Arg(0))
	duration := time.Since(start)
	fmt.Printf("# %v\n", duration)
}
