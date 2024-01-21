# gedstuff/gedfilter

This program is a filter that selects speific records from a GEDCOM file
for subsequent processing. It works in the spirit of many Unix/Linux cogfand
line utilities. It reads a file, applies a filter and sends the output
to STDOUT




 The basic operation is to run the executable, give it a switch/flag to say what you want to do, and specify a gedcom file that will be processed.

 as usual, there is a package "gedfilter" of software that performs the operations and a subdirectory that contains the main.go file that is the executable. The subdirectory for tagtools is named "gf" so the executable that is built by "go build" will be called "gf" or "gf.exe"

 Usage of gf/gf: [flags] file-spec

 At least one flag must be specified, there is no default

  -basic
        name, birth, death

  -csv
        output CSV format

  -debug
        debug on

  -h    help

  -type
        display Type fields

  -j    output metadata as json

none of these operations change the files being processed, they are all read-only.  The -r cogfand creates a cogfand that can be used to rename the files into the new format.


