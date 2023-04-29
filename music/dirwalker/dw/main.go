package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/pfarrell51/gows/music/dirwalker"
)

func usagePrint() {
	fmt.Printf("Usage: %s [flags] verb in-directory-spec out-direcctory-spec\n", os.Args[0])
	fmt.Printf("For now, only 'ffmpeg' 'sox' and 'both' are allowed as a verb\n")
	fmt.Printf("if the out-directory-spec is simply 'mp3' then the output spec will be built from the in-directory-spec,\n")
	fmt.Printf("replacing the word 'flac'with 'mp3' in the path\n")
	fmt.Printf("i.e. mumble/flac/fratz will create an output path of mumble/mp3/fratz\n")
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

	var globals = dirwalker.AllocateData()

	var helpflag bool
	var flags = new(dirwalker.FlagST)
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

	verb := flag.Arg(0)
	if !(verb == "ffmpeg" || verb == "sox" || verb == "both") {
		fmt.Printf("3 Verb must be either ffmpeg sox or both, not %s\n", verb)
		flag.Usage()
		return
	}
	var inpathArg, outpathArg string
	inpathArg = path.Clean(flag.Arg(1))
	outpathArg = path.Clean(flag.Arg(2))
	ext := "mp3"
	switch {
	case numArgs == 1:
		flag.Usage()
		return
	case numArgs == 2:
		switch verb {
		case "ffmpeg", "sox", "both":
			inpathArg = path.Clean(flag.Arg(1))
			outpathArg = path.Clean(flag.Arg(2))
			outpathArg = "mp3"

		default:
			flag.Usage()
			return
		}
	case numArgs > 2:
		if strings.HasPrefix(inpathArg, "-") || strings.HasPrefix(outpathArg, "-") {
			fmt.Printf("WARNING, switches must be before the verb, %s i%s ignored\n", inpathArg, outpathArg)
			flag.Usage()
			return
		}
	}
	if globals.Flags().Debug {
		fmt.Printf("#v: %s i: %s  o: %s\n", verb, inpathArg, outpathArg)
	}
	numDone := globals.Files(verb, inpathArg, outpathArg, ext)
	fmt.Printf("#%d\n", numDone)
	duration := time.Since(start)
	fmt.Printf("# %v\n", duration)
}
