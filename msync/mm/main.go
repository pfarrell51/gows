package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/pfarrell51/gows/msync"
	"github.com/pkg/profile"
)

func usagePrint() {
	fmt.Printf("Usage: %s [flags] in-directory-spec out-direcctory-spec \n", os.Args[0])
	fmt.Printf(" a quick copy, like rsync, but tuned for writing to flash USB drives\n")
}

func main() {
	if len(os.Args) < 2 {
		usagePrint()
		os.Exit(1)
	}
	start := time.Now()
	flag.Usage = func() {
		flag.CommandLine.Output() // may be os.Stderr - but not necessarily
		usagePrint()
		flag.PrintDefaults()
	}

	var globals = msync.AllocateData()

	var helpflag bool
	var flags = new(msync.FlagST)
	flag.BoolVar(&flags.Debug, "de", false, "debug on")
	flag.BoolVar(&helpflag, "h", false, "help")
	flag.BoolVar(&flags.ProfileCPU, "profcpu", false, "turn on profiling for CPU")
	flag.BoolVar(&flags.ProfileMEM, "profmem", false, "turn on profiling for memory usage")
	flag.Parse()
	if helpflag {
		flag.Usage()
		return
	}
	numArgs := len(flag.Args())
	globals.SetFlagArgs(*flags)

	var inpathArg, outpathArg string
	switch {
	case numArgs < 2:
		flag.Usage()
		return
	case numArgs >= 2:
		inpathArg = path.Clean(flag.Arg(0))
		outpathArg = path.Clean(flag.Arg(1))
		if strings.HasPrefix(inpathArg, "-") || strings.HasPrefix(outpathArg, "-") {
			fmt.Printf("WARNING, switches must be before the verb, %s i%s ignored\n", inpathArg, outpathArg)
			flag.Usage()
			return
		}
	}
	if globals.Flags().Debug {
		fmt.Printf("i: %s  o: %s\n", inpathArg, outpathArg)
	}
	if globals.Flags().ProfileCPU {
		defer profile.Start(profile.CPUProfile, // profile.ProfilePath(os.Getenv("HOME")),
			profile.NoShutdownHook)
	}
	globals.WalkDirectories(inpathArg, outpathArg)
	duration := time.Since(start)
	fmt.Printf("# %v\n", duration)
}
