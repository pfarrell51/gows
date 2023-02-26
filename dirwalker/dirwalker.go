// dirwalker
// utility to walk a directory tree and output cool commands

package dirwalker

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"text/template"
)

type FlagST struct {
	CopyAlbumInTrackOrder bool
	CSV                   bool
	Debug                 bool
}
type GlobalVars struct {
	pathArg    string
	localFlags *FlagST
}
type DirEntry struct {
	InPath  string
	OutPath string
	Cmd     string
}

// copy user set flags to a local store
func (g *GlobalVars) SetFlagArgs(f FlagST) {
	g.localFlags = new(FlagST)
	g.localFlags.CopyAlbumInTrackOrder = f.CopyAlbumInTrackOrder
	g.localFlags.CSV = f.CSV
	g.localFlags.Debug = f.Debug
}
func AllocateData() *GlobalVars {
	rval := new(GlobalVars)
	rval.localFlags = new(FlagST)
	if true {
		std1 := DirEntry{"/home/Music/flac", "/home/music/mp3", "rename "}
		tmpl, err := template.New("test").Parse("{{.Cmd}} {{.InPath}} of {{.OutPath}}")
		if err != nil {
			panic(err)
		}
		err = tmpl.Execute(os.Stdout, std1)
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(os.Stdout, "\n")
	}
	return rval
}
func Files(path string) (count int) {
	fsys := os.DirFS(path)
	fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		if filepath.Ext(p) == ".flac" {
			count++
		}
		return nil
	})
	return count
}
