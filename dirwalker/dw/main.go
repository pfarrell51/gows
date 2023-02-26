package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/pfarrell51/gows/dirwalker"
)

func usagePrint() {
	fmt.Printf("Usage: %s [flags] verb in-directory-spec out-direcctory-spec extension\nFor now, only 'ffmpeg' is allowed as a verb", os.Args[0])
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("arglen %d\n", len(os.Args))
		usagePrint()
		os.Exit(1)
	}
	start := time.Now()
	flag.Usage = func() {
		w := flag.CommandLine.Output() // may be os.Stderr - but not necessarily
		usagePrint()
		flag.PrintDefaults()
		fmt.Fprintf(w, "default is to apply verb to in-directory-slec yeilding out-directory-spec.\n")
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
	if len(flag.Args()) < 4 {
		flag.Usage()
		return
	}

	globals.SetFlagArgs(*flags)
	verb := flag.Arg(0)
	inpathArg := path.Clean(flag.Arg(1))
	outpathArg := path.Clean(flag.Arg(2))
	ext := strings.TrimSpace(flag.Arg(3))

	fmt.Printf("#%d\n", dirwalker.Files(verb, inpathArg, outpathArg, ext))
	duration := time.Since(start)
	fmt.Printf("# %v\n", duration)

}
