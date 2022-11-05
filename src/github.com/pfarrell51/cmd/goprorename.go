// shell program to rename files created by a GoPro into a single, sensible
// ordering of files so that the order is obvious for easy processing
// by other utulities such as RaceRender
//
// goPro naming conventions: https://community.gopro.com/s/article/GoPro-Camera-File-Naming-Convention?language=en_US

package main

import (
	"fmt"
	"io/fs"
	"os"
	"regexp"
	"strings"
)

// the new testing/fstest package provides a TestFS function that checks for and
// reports common mistakes. It also provides a simple in-memory file systemi
// implementation, MapFS, which can be useful for testing code that accepts fs.FS implementations."

func ProcessFiles(fsys fs.FS, patharg string) {
	extRegex := regexp.MustCompile(".(M|m)(p|P)4")
	nameRegex := regexp.MustCompile("(?s)(GX|H)(\\d{2})(\\d{4})")
	fs.WalkDir(fsys, patharg, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Println("Error processing", d, " in ", path)
			fmt.Println("error is ", err)
			return nil
		}
		if strings.HasPrefix(path, ".") {
			return nil
		}
		ext := extRegex.FindString(path)
		if len(ext) == 0 {
			return nil
		}
		nameParts := nameRegex.FindAllStringSubmatch(path, -1)
		if len(nameParts) > 0 && len(nameParts[0]) > 3 {
			prefix := nameParts[0][1]
			chapter := nameParts[0][2]
			clip := nameParts[0][3]
			fmt.Println("p: ", prefix, " chpt ", chapter, " clip ", clip)
		}
		return nil
	})
}

// get argument from shell, if any or default

func main() {
	pathArg := "."
	if len(os.Args) > 1 {
		pathArg = os.Args[1]
	}
	fs := os.DirFS(pathArg)
	fmt.Println("fs from arg: ", fs)
	ProcessFiles(fs, pathArg)
}
