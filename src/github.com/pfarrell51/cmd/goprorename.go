// shell program to rename files created by a GoPro into a single, sensible
// ordering of files so that the order is obvious for easy processing
// by other utulities such as RaceRender

package main

import (
	"fmt"
	"io/fs"
	"os"
)

func Files(fsys fs.FS, patharg string) {

	fs.WalkDir(fsys, patharg, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Println(d)
			fmt.Println("error is ")
			fmt.Println(err)
		}
		fmt.Println(path)
		return nil
	})

}

// get argument from shell, if any or default
func main() {
	pathArg := "./*.[mMpP]4"
	if len(os.Args) > 1 {
		pathArg = os.Args[1]
	}
	fs := os.DirFS(pathArg)
	fmt.Println("fs from arg: ", fs)
	Files(fs, pathArg)

}
