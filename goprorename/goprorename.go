// rename files created by a GoPro into a single, sensible
// ordering of files so that the order is obvious for easy processing
// by other utulities such as RaceRender
//
// goPro naming conventions: https://community.gopro.com/s/article/GoPro-Camera-File-Naming-Convention?language=en_US

package goprorename

import (
	"fmt"
	"github.com/emirpasic/gods/maps/treemap"
	"io/fs"
	"os"
	"regexp"
	"strings"
)

type gpfile struct {
	oldname string
	newname string
}

var extRegex = regexp.MustCompile(".(M|m)(p|P)4")
var nameRegex = regexp.MustCompile("(?s)(GX|H)(\\d{2})(\\d{4})")

func ProcessFiles(pathArg string) {
	theMap := treemap.NewWithStringComparator()
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
			entry := gpfile{dName, key}
			fmt.Println(" dname ", dName, " key ", key, " entry ", entry)
			theMap.Put(key, entry)
			val, _ := theMap.Get(key)
			fmt.Println(key, val)
		}
		return nil
	})
	//	var gt interface{} = &gpfile{}
	keyList := theMap.Keys()
	for k, v := range keyList {
		fmt.Println("k of keylist ", k)
		g, ok := theMap.Get(k)
		gv, ok := g.(*gpfile)
		if ok {
			//	v, _ := theMap.Get(k)
			fmt.Printf("mv %s %s\n", gv.oldname, gv.newname)
			fmt.Printf("k,v=%s %s\n", k, v)
		} else {
			fmt.Println("error, bad type")
		}
	}
}
