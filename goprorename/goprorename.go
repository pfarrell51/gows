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
	"strings"
)

var extRegex = regexp.MustCompile(".(M|m)(p|P)4")
var nameRegex = regexp.MustCompile("(?s)(GX|H)(\\d{2})(\\d{4})")

func ProcessFiles(pathArg string) {

	// TODO: list of skips or regex for skip or function call to check skip
	//subDirToSkip := "skip"

	fmt.Println(Files(pathArg))
}

func Files(path string) (count int) {
	fsys := os.DirFS(path)
	fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		// fmt.Println("Embeded func ", d)
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
			count++
			prefix := nameParts[0][1]
			chapter := nameParts[0][2]
			clip := nameParts[0][3]
			fmt.Printf("mv %s   %2s%4s%2s.mp4\n", dName, prefix, clip, chapter)
		}
		return nil
	})
	return count
}
