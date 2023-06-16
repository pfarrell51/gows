package randnames

import (
	"fmt"
	"io/fs"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type FlagST struct {
	NoTags    bool
	TwoLetter bool
}

var localFlags = new(FlagST)
var numDirs int

// twoltr acts as a set of strings
var twoltr map[string]int

func init() {
	twoltr = make(map[string]int)
	mumble := [17]string{"As", "Do", "El", "Go", "He", "If", "In", "It", "La", "My", "No", "Oh", "On", "So", "To", "Up", "We"}
	for _, m := range mumble {
		twoltr[m] = 0
	}
}
func SetFlags(arg FlagST) {
	localFlags = &arg
}
func GetFlag() *FlagST {
	return localFlags
}

func foo() {
	s := rand.NewSource(1234567)
	r := rand.New(s)

	fmt.Println(r.Perm(20))
}
func listTwoLetterWords(arg string) {
}
func renameSongs(arg string) {
}
func countFiles(pathArg string) int {
	f, err := os.Open(pathArg)
	if err != nil {
		panic(err)
	}
	list, err := f.Readdirnames(-1)
	f.Close()
	if err != nil {
		panic(err)
	}
	return len(list)
}

// walk all files, looking for nice music files.
// fill in a map keyed by the desired new name order
func WalkFiles(pathArg string) {
	fmt.Printf("num files %d\n", countFiles(pathArg))

	var oldDir string
	fsys := os.DirFS(pathArg)
	fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("Oh No, can't process because %s\n", err)
			return err
		}
		if d.IsDir() {
			return nil
		}
		var notOld = filepath.Dir(p)
		if oldDir != notOld && notOld != "." {
			oldDir = notOld
			numDirs++
		}
		processFile(fsys, p, d, err)
		return nil
	})
}

var ExtRegex = regexp.MustCompile(`[Mm][Pp][34]|[Ff][Ll][Aa][Cc]$`)
var twoletterRegex = regexp.MustCompile(`^[[:alpha:]]{2} `)

// this is the local  process  called by WalkDir for each file
// p is the current path/name
func processFile(fsys fs.FS, p string, d fs.DirEntry, err error) error {
	if err != nil {
		fmt.Println("Error processing", p, " in ", d)
		fmt.Println("error is ", err)
		return nil
	}
	if d == nil || d.IsDir() || strings.HasPrefix(p, ".") {
		return nil
	}
	extR := ExtRegex.FindStringIndex(p)
	if extR == nil {
		return nil // not interesting extension
	}
	twoL := twoletterRegex.FindStringIndex(p)
	if twoL != nil {
		f := p[twoL[0] : twoL[1]-1]
		s := twoltr[f]
		fmt.Printf("%s ? %s\n", f, s)
	}

	return nil
}
