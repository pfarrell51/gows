// verify that the filename and the internal meta data are essentially the same

package verify

import (
	"fmt"
	"io/fs"
	"os"
	pl "path"
	"path/filepath"
	"strings"

	"github.com/dhowden/tag"
	"github.com/texttheater/golang-levenshtein/levenshtein"
)

type FlagST struct {
	Debug bool
}

func WalkFiles(flags *FlagST, path string) {
	fsys := os.DirFS(path)
	fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {

		if !strings.EqualFold(filepath.Ext(p), ".mp3") {
			return nil
		}

		joined := filepath.FromSlash(pl.Join(path, p))

		file, err := os.Open(joined)
		if err != nil {
			fmt.Printf("err : %v %s\n", err, joined)
			return err
		}
		defer file.Close()
		m, err := tag.ReadFrom(file)
		if err != nil {
			fmt.Printf("%v %s", err, p)
			return err
		}
		if m == nil {
			fmt.Printf("tag.ReadFrom (file) turned nil but no error for %s\n", p)
		}
		target := p[0 : len(p)-4]
		parts := strings.Split(p, ";")
		if len(parts) > 1 {
			target = parts[0]
		}

		distance := levenshtein.DistanceForStrings([]rune(m.Title()), []rune(target), levenshtein.DefaultOptions)
		if distance > 4 {
			fmt.Printf("%s and %s distance %d\n", target, m.Title(), distance)
		}
		return nil
	})

}
