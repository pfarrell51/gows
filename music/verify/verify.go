// verify that the filename and the internal meta data are essentially the same

package verify

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type FlagST struct {
	Debug bool
}

func WalkFiles(flags *FlagST, path string) (count int) {
	fsys := os.DirFS(path)
	fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		if strings.EqualFold(filepath.Ext(p), ".mp4") {
			count++
		}
		return nil
	})
	return count
}
