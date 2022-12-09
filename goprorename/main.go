// shell program to rename files created by a GoPro into a single, sensible
// ordering of files so that the order is obvious for easy processing
// by other utulities such as RaceRender
//
// goPro naming conventions: https://community.gopro.com/s/article/GoPro-Camera-File-Naming-Convention?language=en_US

package main

import (
	"fmt"
	"github.com/pfarrell51/gows/goprorename"
	"os"
)

// get argument from shell, if any or default

func main() {
	pathArg := "."
	if len(os.Args) > 1 {
		pathArg = os.Args[1]
	}
	fs := os.DirFS(pathArg)
	fmt.Println("fs from arg: ", fs)
	goprorename.ProcessFiles(fs, pathArg)
}
