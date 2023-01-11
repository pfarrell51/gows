package main

import (
	"fmt"
	"io/fs"
	"os"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s DIRNAME", os.Args[0])
		os.Exit(1)
	}
	pathArg := os.Args[1]
	fmt.Println(pathArg)
	fmt.Println(runtime.GOOS)
}

var extRegex = regexp.MustCompile(".(M|m)(p|P)4")
var nameRegex = regexp.MustCompile("(?s)(GX|H)(\\d{2})(\\d{4})")

func ProcessFiles(pathArg string) {
	rmap := walkFiles(pathArg)
	processMap(rmap)

}

// walk all files, looking for nice GoPro created video files.
// fill in a map keyed by the desired new name order
func walkFiles(pathArg string) map[string]string {
	theMap := make(map[string]string)
	fsys := os.DirFS(pathArg)
	fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Println("Error processing", p, " in ", d)
			fmt.Println("error is ", err)
			return nil
		}
		if d == nil {
			fmt.Println("d is nil")
			return nil
		}
		dName := d.Name()
		if strings.HasPrefix(p, ".") {
			return nil
		}
		ext := extRegex.FindString(p)
		if len(ext) == 0 {
			return nil
		}
		nameParts := nameRegex.FindAllStringSubmatch(p, -1)
		if len(nameParts) > 0 && len(nameParts[0]) > 3 {
			prefix := nameParts[0][1]
			chapter := nameParts[0][2]
			clip := nameParts[0][3]
			key := prefix + clip + chapter
			theMap[key] = dName
		}
		return nil
	})
	return theMap
}

// go thru the map, sort by key
// then create new ordering that makes sense to human
func processMap(m map[string]string) map[string]string {
	var cmd = "mv "
	if runtime.GOOS == "windows" {
		cmd = "ren "
	}
	keys := make([]string, 0, len(m)) // copy to a slice to sort
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var firstChpt string
	var cnum int
	for _, k := range keys {
		if cnum < 1 {
			firstChpt = k[0:6]
			c, _ := strconv.Atoi(k[6:8])
			cnum = c
		} else {
			cnum++
		}
		source, _ := m[k]
		delete(m, k)
		m[k] = source
		fmt.Printf("%s%s %2s%02d%4s.mp4\n", cmd, source, firstChpt[0:2], cnum, firstChpt[2:6])
	}
	return m
}
