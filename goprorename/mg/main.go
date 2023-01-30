/*
shell program to rename files created by a GoPro into a single, sensible
ordering of files so that the order is obvious for easy processing
by other utulities such as RaceRender.

accepts a directory spec as the first argument.
If no first argument is entered by the user, the current directory (i.e. ".") is used.

Outputs a serious of lines, each a shell command to preform the renaming.
This output can be directed to a file, as any standard output.

goPro naming conventions: [https://community.gopro.com/s/article/GoPro-Camera-File-Naming-Convention?language=en_US]
*/
package main

import (
	"fmt"
	"github.com/pfarrell51/gows/goprorename"
	"os"
)

// get argument from shell, if any or default to "."
func main() {
	pathArg := "."
	if len(os.Args) > 1 {
		pathArg = os.Args[1]
	} else {
		fmt.Printf("usage %s <directory>\n", os.Args[0])
	}
	goprorename.ProcessFiles(pathArg)
}
