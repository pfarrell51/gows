// shell program to rename files created by a GoPro into a single, sensible
// ordering of files so that the order is obvious for easy processing
// by other utulities such as RaceRender

package main

import (
	"fmt"
	"github.com/pfarrell51/gows/goprorename"
	"io/fs"
	"os"
)

// get argument from shell, if any or default

func main() {
	pathArg := "."
	if false {
		if len(os.Args) > 1 {
			pathArg = os.Args[1]
		}
	}
	pathArg = "/usr/lib"
	ok := fs.ValidPath(pathArg)
	fmt.Println("is OK? ", ok)
	fsr := os.DirFS(pathArg)
	fmt.Println("fs from arg: ", fsr)
	fmt.Println(fsr)
	goprorename.ProcessFiles(fsr, pathArg)
}
