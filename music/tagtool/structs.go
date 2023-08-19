// structs  this  file contains structures and other "global" data stores

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
	"unicode"

	"github.com/dlclark/metaphone3"
	"github.com/zyedidia/generic"
	"github.com/zyedidia/generic/avl"
	"github.com/zyedidia/generic/btree"
)

type Song struct {
	Artist            string
	artistH           string
	Album             string
	albumH            string
	Title             string
	titleH            string
	smapKey           string
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
type PairSongs struct {
	a *Song
	b *Song
}
type FlagST struct {
	CompareTagsToTitle    bool
	CopyAlbumInTrackOrder bool
	CpuProfile            string
	CSV                   bool
	Debug                 bool
	DoInventory           bool
	DoRename              bool
	DoSummary             bool
	DupJustTitle          bool
	DupTitleAlbumArtist   bool
	JsonOutput            bool
	JustList              bool
	MemProfile            string
	NoGroup               bool
	NoTags                bool
	ShowArtistNotInMap    bool
	ShowNoSongs           bool
	UnicodePunct          bool
	ZDumpArtist           bool
}

type GlobalVars struct {
	pathArg                              string
	localFlags                           *FlagST
	songsProcessed                       int
	numNoAcoustId, numNoTitle, numNoMBID int
	artistCountTree                      *avl.Tree[string, int]
	albumCountTree                       *avl.Tree[string, int]
	songCountTree                        *avl.Tree[string, int]
	gptree                               *btree.Tree[string, string]
	songTree                             map[string]Song
	tracksongs                           []TrackSong // for sorting by track number within an album
	invTriples                           []InventorySong
	dupSongs                             []PairSongs
	knownIds                             map[string]bool
	numDirs                              int
}

// copy user set flags to a local store
func (g *GlobalVars) SetFlagArgs(f FlagST) {
	g.localFlags = new(FlagST)
	g.localFlags.CompareTagsToTitle = f.CompareTagsToTitle
	g.localFlags.CopyAlbumInTrackOrder = f.CopyAlbumInTrackOrder
	g.localFlags.CpuProfile = f.CpuProfile
	g.localFlags.CSV = f.CSV
	g.localFlags.Debug = f.Debug
	g.localFlags.DupJustTitle = f.DupJustTitle
	g.localFlags.DupTitleAlbumArtist = f.DupTitleAlbumArtist
	g.localFlags.DoRename = f.DoRename
	g.localFlags.DoInventory = f.DoInventory
	g.localFlags.DoSummary = f.DoSummary
	g.localFlags.NoTags = f.NoTags
	g.localFlags.JsonOutput = f.JsonOutput
	g.localFlags.JustList = f.JustList
	g.localFlags.MemProfile = f.MemProfile
	g.localFlags.NoGroup = f.NoGroup
	g.localFlags.NoTags = f.NoTags
	g.localFlags.ShowArtistNotInMap = f.ShowArtistNotInMap
	g.localFlags.ShowNoSongs = f.ShowNoSongs
	g.localFlags.UnicodePunct = f.UnicodePunct
	g.localFlags.ZDumpArtist = f.ZDumpArtist
}
func (g *GlobalVars) Flags() *FlagST {
	return g.localFlags
}
func (g *GlobalVars) GetSongTree() map[string]Song {
	return g.songTree
}

var enc metaphone3.Encoder

const maxEncode = 20

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
	rval = strings.TrimPrefix(rval, "The ")
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
	//fmt.Printf("Encoding title of %s\n", s)
	prim, sec := enc.Encode(justLetter(StandardizeTitle(s)))
	return prim, sec
}
func EncodeArtist(s string) (string, string) {
	prim, sec := enc.Encode(justLetter(StandardizeArtist(s)))
	return prim, sec
}
func EncodeAlbum(s string) (string, string) {
	prim, sec := enc.Encode(justLetter(StandardizeTitle(s)))
	return prim, sec
}

func AllocateData() *GlobalVars {
	enc.Encode("ignore this")
	enc.MaxLength = maxEncode
	rval := new(GlobalVars)
	rval.localFlags = new(FlagST)
	rval.artistCountTree = avl.New[string, int](generic.Less[string])
	rval.albumCountTree = avl.New[string, int](generic.Less[string])
	rval.songCountTree = avl.New[string, int](generic.Less[string])
	rval.gptree = btree.New[string, string](generic.Less[string])
	rval.songTree = make(map[string]Song, 1000)
	rval.tracksongs = make([]TrackSong, 0, 1500)
	rval.dupSongs = make([]PairSongs, 0, 1500)
	rval.invTriples = make([]InventorySong, 0, 1000)
	rval.knownIds = make(map[string]bool, 1000)
	if rval.localFlags == nil {
		fmt.Println("PIB in allocate Data, localflags is nil")
	}
	if rval.songTree == nil {
		panic("Allocate Data post setup songTree is nil")
	}
	rval.loadArtistMap()
	return rval
}

func (g *GlobalVars) loadArtistMap() {
	if g.gptree.Size() > 0 {
		panic("PIB group/artist array already populated")
	}
	artists := make([]string, 0, 600)
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
			lit = strings.TrimPrefix(lit, "\"")
			lit = strings.TrimSuffix(lit, "\"")
			artists = append(artists, lit)
		}
	}
	sort.Strings(artists)
	for _, n := range artists {
		prim, sec := EncodeArtist(n)
		g.gptree.Put(prim, n)
		if len(sec) > 0 {
			g.gptree.Put(sec, n)
		}
		if g.Flags().Debug {
			fmt.Printf("%s, %s, %s\n", prim, sec, n)
		}
	}
}
