// shell program to rename files created by a GoPro into a single, sensible
// ordering of files so that the order is obvious for easy processing
// by other utulities such as RaceRender

package main

import (
	"fmt"
	"io/fs"
	"os"
	"regexp"
	"strings"
)

func Files(fsys fs.FS, patharg string) {
	comRegex, err := regexp.Compile(".(Mm)r(pP)4")
	if err != nil {
		fmt.Println("PIB")
		return
	}
	fs.WalkDir(fsys, patharg, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Println("Error processing", d)
			fmt.Println("error is ", err)
			return nil
		}
		if strings.HasPrefix(path, ".") {
			return nil
		}
		pathAsBytes := []byte(path)
		fmt.Println(pathAsBytes)
		ext := comRegex.Find(pathAsBytes)
		if ext == nil {
			fmt.Println("regex failed ", path)
			return nil
		} else {
			fmt.Println("success ", ext)
		}
		fmt.Println(path)

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
	Files(fs, pathArg)
}
