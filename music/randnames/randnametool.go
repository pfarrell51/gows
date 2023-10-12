// this manages randomized two latter prefixes for song nmes to that
// my silly Mazda's infotainment system will play all of my songs and not
// restart the order every time I turn off the car.
//
// not re-entrant, thread safe, etc.

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
	AddTag    bool
	NoTags    bool
	RmTag     bool
	TwoLetter bool
	Debug     bool
}
type FileDoer interface {
	FileDo(fsys fs.FS, p string, d fs.DirEntry, err error) error
}

var localFlags = new(FlagST)
var numDirs int

const maxsongs = 800
const lowers = "abcdefghighlmnopqrstuvwxyz"
const alphas = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghighlmnopqrstuvwxyz"

var tmpprefix, songprefix [maxsongs]string

// stopwords acts as a set of strings
var stopwords map[string]int

func init() {
	stopwords = make(map[string]int)
	ktl := [17]string{"As", "Do", "El", "Go", "He", "If", "In", "It", "La", "My", "No", "Oh", "On", "So", "To", "Up", "We"}
	for _, m := range ktl {
		stopwords[m] = 0
	}
	s := rand.NewSource(1234567)
	r := rand.New(s)

	ranIdx := r.Perm(maxsongs)
	i := 0
Outerloop:
	for j := 0; j < len(lowers); j++ {
		for k := 0; k < len(alphas); k++ {
			val := string(lowers[j]) + string(alphas[k])
			_, ok := stopwords[val]
			if !ok {
				tmpprefix[i] = val
				i++
				if i >= maxsongs {
					break Outerloop
				}

			} else {
				fmt.Printf("found i: %d j: %d k: %d %s\n", i, j, k, tmpprefix[i])
			}
		}
	}
	for i = 0; i < maxsongs; i++ {
		songprefix[i] = tmpprefix[ranIdx[i]]
	}
}
func SetFlags(arg FlagST) {
	localFlags = &arg
}
func GetFlag() *FlagST {
	return localFlags
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
	if len(list) > maxsongs {
		panic(fmt.Sprintf("too many songs %d in directory, rebuild with bigger maxsongs %d\n", len(list), maxsongs))
	}

	return len(list)
}

var WalkFileNum int

// walk all files, looking for nice music files.
// fill in a map keyed by the desired new name order
func WalkFiles(pathArg string) {
	fmt.Printf("num files %d\n", countFiles(pathArg))

	WalkFileNum++
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

		perFile(fsys, p, d, err)
		return nil
	})
}

var ExtRegex = regexp.MustCompile(`[Mm][Pp][34]|[Ff][Ll][Aa][Cc]$`)
var twoletterRegex = regexp.MustCompile(`^[[:alpha:]]{2}[[:space:]]`)
var goodPrefixRegex = regexp.MustCompile(`^[a-z][[:alpha:]]_`)
var brokenPrefixRegex = regexp.MustCompile(`^[a-z][[:alpha:]][A-Z][[:alpha:]]`)

// this is the local  process  called by WalkDir for each file
// p is the current path/name
func perFile(fsys fs.FS, p string, d fs.DirEntry, err error) error {
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

	switch {
	case localFlags.AddTag:
		err = processAddTag(fsys, p, d, err)
		if err != nil {
			return err
		}
	case localFlags.RmTag:
		err = processRmTag(fsys, p, d, err)
		if err != nil {
			return err
		}
	case localFlags.TwoLetter:
		err = processTwoLetter(fsys, p, d, err)
		if err != nil {
			return err
		}
	default:
		fmt.Println("default triggered")
	}
	return nil
}

// check format of file name (p) and return new name
// with random prefix
func pAndNewP(p string) (newP string) {
	newP = p
	twoPunctL := goodPrefixRegex.FindStringIndex(p)
	if twoPunctL != nil {
		fmt.Printf("found leading sort in %s will use %s\n", p, p[3:])
		newP = p[3:]
	}

	twoL := twoletterRegex.FindStringIndex(p)
	if twoL == nil {
		newP = songprefix[WalkFileNum] + "_" + p
	} else {
		f := p[twoL[0] : twoL[1]-1]
		_, ok := stopwords[f]
		if ok {
			stopwords[f]++
			newP = songprefix[WalkFileNum] + "_" + p
			if GetFlag().Debug {
				v := stopwords[f]
				fmt.Printf("%s '%s' ? %d\n", f, p, v)
			}
		} else {
			panic(fmt.Sprintf("unknown two letter word %s in song %s\n", f, p))
		}
	}
	return newP
}
func processAddTag(fsys fs.FS, p string, d fs.DirEntry, err error) error {
	newP := pAndNewP(p)
	fmt.Printf("mv \"%s\" \"%s\"\n", p, newP)
	WalkFileNum++
	return nil
}
func processRmTag(fsys fs.FS, p string, d fs.DirEntry, err error) error {
	goodIdx := goodPrefixRegex.FindStringIndex(p)
	if goodIdx != nil {
		fmt.Printf("mv \"%s\" \"%s\"\n", p, p[3:])
	}
	return nil
}
func processTwoLetter(fsys fs.FS, p string, d fs.DirEntry, err error) error {
	goodIdx := goodPrefixRegex.FindStringIndex(p)
	if goodIdx != nil {
		fmt.Printf("good %s\n", p)
		return nil
	}
	badIdx := brokenPrefixRegex.FindStringIndex(p)
	if badIdx != nil {
		fmt.Printf("found bad index in %s\n", p)
	} else {
		fmt.Printf("no prefix in %s\n", p)
	}
	return nil
}
