// shell program to rename files created by a GoPro into a single, sensible
// ordering of files so that the order is obvious for easy processing
// by other utulities such as RaceRender

package main

import (
	"github.com/pfarrell51/gows/goprorename"
	"os"
)

// get argument from shell, if any or default
func main() {
	pathArg := "."
	if len(os.Args) > 1 {
		pathArg = os.Args[1]
	}
	goprorename.ProcessFiles(pathArg)
}
