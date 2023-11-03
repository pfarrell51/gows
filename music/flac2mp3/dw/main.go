package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/pfarrell51/gows/music/flac2mp3"
)

func usagePrint() {
	fmt.Printf("Usage: %s [flags] in-directory-spec out-direcctory-spec\n", os.Args[0])
	fmt.Printf("if the out-directory-spec is simply 'mp3' then the output spec will be built from the in-directory-spec,\n")
	fmt.Printf("replacing the word 'flac'with 'mp3' in the path\n")
	fmt.Printf("i.e. mumble/flac/fratz will create an output path of mumble/mp3/fratz\n")
	fmt.Printf("\nSpecial shortcut command:\n")
	fmt.Printf("dw/dw indirectory-spec mp3\n")
	fmt.Printf("will cause sox to be run on the flac directory, outputing to a 'mp3' directory\n")

}

func main() {
	if len(os.Args) < 2 {
		usagePrint()
		os.Exit(1)
	}
	start := time.Now()
	flag.Usage = func() {
		w := flag.CommandLine.Output() // may be os.Stderr - but not necessarily
		usagePrint()
		flag.PrintDefaults()
		fmt.Fprintf(w, "default is to apply verb to in-directory-spec yeilding out-directory-spec.\n")
	}

	var globals = flac2mp3.AllocateData()

	var helpflag bool
	var flags = new(flac2mp3.FlagST)
	flag.BoolVar(&flags.Debug, "de", false, "debug on")
	flag.BoolVar(&helpflag, "h", false, "help")
	flag.BoolVar(&flags.SkipIfBreadcrumbExists, "skip", false, "skip if breadcrumb file exists")
	flag.Parse()
	if helpflag {
		flag.Usage()
		return
	}
	numArgs := len(flag.Args())
	globals.SetFlagArgs(*flags)

	if globals.Flags().Debug {
		fmt.Printf("%d flag args %v \n", numArgs, flag.Args())
	}
	verb := "sox"
	var inpathArg, outpathArg string
	inpathArg = path.Clean(flag.Arg(0))
	outpathArg = path.Clean(flag.Arg(1))
	if strings.HasPrefix(inpathArg, "-") || strings.HasPrefix(outpathArg, "-") {
		fmt.Printf("WARNING, switches must be before the verb, %s i%s ignored\n", inpathArg, outpathArg)
		flag.Usage()
		return
	}
	if globals.Flags().Debug {
		fmt.Printf("#v: %s i: %s  o: %s\n", verb, inpathArg, outpathArg)
	}
	numDone := globals.Files(verb, inpathArg, outpathArg)
	fmt.Printf("#%d\n", numDone)
	duration := time.Since(start)
	fmt.Printf("# %v\n", duration)
}
