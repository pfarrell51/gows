// this  file contains structures and other "global" data stores

package tagtool

import (
	"bytes"
	_ "embed"
	"fmt"
	"go/scanner"
	"go/token"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"unicode"

	"github.com/dlclark/metaphone3"
	g "github.com/zyedidia/generic"
	"github.com/zyedidia/generic/btree"
)

type Song struct {
	Artist            string
	artistH           string
	Album             string
	albumH            string
	Title             string
	titleH            string
	Genre             string
	Disc, DiscCount   int
	Track             int    `json:",omitempty"`
	Year              int    `json:",omitempty"`
	ISRC              string `json:",omitempty"` // International Standard Recording Code
	MBID              string `json:",omitempty"` // musicbrainz ID
	AcoustID          string `json:",omitempty"` // Acoust ID
	alreadyNew        bool
	artistInDirectory bool
	artistKnown       bool
	inPath            string
	inPathDescent     string // any descent below the pathArg aka outPathBase
	outPath           string
	outPathBase       string // copied from pathArg entered by the user
	ext               string
}
type FlagST struct {
	ShowArtistNotInMap    bool
	DoRename              bool
	DoInventory           bool
	DoSummary             bool
	JustList              bool
	NoGroup               bool
	NoTags                bool
	ZDumpArtist           bool
	JsonOutput            bool
	Debug                 bool
	DuplicateDetect       bool
	CopyAlbumInTrackOrder bool
}

var localFlags = new(FlagST)

// copy user set flags to a local store
func SetFlagArgs(f FlagST) {
	localFlags.ShowArtistNotInMap = f.ShowArtistNotInMap
	localFlags.DoRename = f.DoRename
	localFlags.DoInventory = f.DoInventory
	localFlags.DoSummary = f.DoSummary
	localFlags.JustList = f.JustList
	localFlags.NoGroup = f.NoGroup
	localFlags.NoTags = f.NoTags
	localFlags.ZDumpArtist = f.ZDumpArtist
	localFlags.JsonOutput = f.JsonOutput
	localFlags.Debug = f.Debug
	localFlags.DuplicateDetect = f.DuplicateDetect
	localFlags.CopyAlbumInTrackOrder = f.CopyAlbumInTrackOrder
}
func GetFlags() *FlagST {
	return localFlags
}

var enc metaphone3.Encoder

const maxEncode = 20

func init() {
	enc.Encode("ignore this")
	enc.MaxLength = maxEncode
}

var (
	//go:embed data/artist.txt
	artistnames []byte
)

// takes a string and returns just the letters.
func justLetter(a string) string {
	buff := bytes.Buffer{}
	for _, c := range a {
		if unicode.IsLetter(c) {
			buff.WriteRune(c)
		} else if c == '_' || c == '&' || unicode.IsSpace(c) {
			// ignore it
		} else if c == '-' {
			break
		}
	}
	return buff.String()
}

var regAndThe = regexp.MustCompile("(?i) (and|the) ")

func StandardizeArtist(art string) string {
	if len(art) == 0 {
		return art
	}
	rval := strings.TrimSpace(art)
	if rval == "The The" { // UK post-punk band
		return rval
	}
	if matched, _ := regexp.MatchString("^[A-Z]( |-|_)", rval); matched {
		rval = rval[2:]
	}
	if strings.HasPrefix(rval, "The ") {
		rval = rval[4:]
	}
	for j := 0; j < 4; j++ { // 4 allows the space before the keyword (and/the), as we back up
		loc := regAndThe.FindStringIndex(rval)
		if len(loc) < 1 {
			break
		}
		rval = rval[:loc[0]] + rval[loc[1]-1:]
	}
	return rval
}

var digLetters = []byte{'B', 'C', 'D', 'F', 'G', 'J', 'K', 'L', 'M', 'R'}

func StandardizeTitle(title string) string {
	if len(title) == 0 {
		return title
	}
	rval := strings.TrimSpace(title)
	if matched, _ := regexp.MatchString("^[A-Z](-|_)", rval); matched {
		rval = rval[2:]
	}
	rval = strings.ReplaceAll(rval, "/", " ")
	rval = strings.ReplaceAll(rval, "_", " ")
	if strings.HasPrefix(rval, "Track") {
		if n, err := strconv.Atoi(rval[5:]); err == nil {
			tens := n / 10
			ones := n % 10
			rval = fmt.Sprintf("%s%s%c%c", rval[:5], Convert1to1000(n), digLetters[tens], digLetters[ones])
		}
	}
	return rval
}
func EncodeTitle(s string) (string, string) {
	prim, sec := enc.Encode(justLetter(StandardizeTitle(s)))
	return prim, sec
}
func EncodeArtist(s string) (string, string) {
	prim, sec := enc.Encode(justLetter(StandardizeArtist(s)))
	return prim, sec
}

var Gptree = btree.New[string, string](g.Less[string])

func GetArtistMap() *btree.Tree[string, string] {
	return Gptree
}

var onlyOnce sync.Once

func LoadArtistMap() {
	onlyOnce.Do(func() {
		enc.Encode("ignore this")
		enc.MaxLength = maxEncode
		var artists []string
		// Initialize the scanner.
		var s scanner.Scanner
		fset := token.NewFileSet()                              // positions are relative to fset
		file := fset.AddFile("", fset.Base(), len(artistnames)) // register input "file"
		s.Init(file, artistnames, nil /* no error handler */, scanner.ScanComments)

		// Repeated calls to Scan yield the token sequence found in the input.
		for {
			_, tok, lit := s.Scan()
			if tok == token.EOF {
				break
			}
			if tok == token.STRING {
				if strings.HasPrefix(lit, "\"") {
					lit = lit[1:]
				}
				if strings.HasSuffix(lit, "\"") {
					lit = lit[:len(lit)-1]
				}
				artists = append(artists, lit)
			}
		}
		sort.Strings(artists)
		for _, n := range artists {
			prim, sec := EncodeArtist(n)
			Gptree.Put(prim, n)
			if len(sec) > 0 {
				Gptree.Put(sec, n)
			}
			if GetFlags().Debug {
				fmt.Printf("%s, %s, %s\n", prim, sec, n)
			}
		}
	})
}
