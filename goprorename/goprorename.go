// rename files created by a GoPro into a single, sensible
// ordering of files so that the order is obvious for easy processing
// by other utulities such as RaceRender
//
// goPro naming conventions: https://community.gopro.com/s/article/GoPro-Camera-File-Naming-Convention?language=en_US

package goprorename

import (
	"fmt"
	"io/fs"
	"os"
	"regexp"
	"sort"
	"strings"
)

var extRegex = regexp.MustCompile(".(M|m)(p|P)4")
var nameRegex = regexp.MustCompile("(?s)(GX|H)(\\d{2})(\\d{4})")

func ProcessFiles(pathArg string) {
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

	keys := make([]string, 0, len(theMap))
	for k := range theMap {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		fmt.Printf("mv %s.mp4 %s\n", theMap[k], k)
	}
}
