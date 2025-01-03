// structs  this  file contains structures and other "global" data stores

package tagtool

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

type FlagST struct {
	CompareTagsToTitle    bool
	CopyAlbumInTrackOrder bool
	CSV                   bool
	Debug                 bool
	DoInventory           bool
	DoRename              bool
	DoSummary             bool
	SuppressTitles        bool
	JsonOutput            bool
	NoGroup               bool
	NoTags                bool
	UnicodePunct          bool
}

type GlobalVars struct {
	pathArg                              string
	SecondArg                            string
	localFlags                           *FlagST
	songsProcessed                       int
	numNoAcoustId, numNoTitle, numNoMBID int
	artistCount                          int
	oldArtist                            string
	albumCount                           int
	numDirs                              int
	oldAlbum                             string
	songCount                            int
	tracksongs                           []TrackSong // for sorting by track number within an album
	invSongs                             []Song
	csvWrtr                              *csv.Writer
}

func (f FlagST) String() string {
	var r strings.Builder

	r.WriteString(fmt.Sprintf("CompareTagsToTitle: %t\n", f.CompareTagsToTitle))
	r.WriteString(fmt.Sprintf("CopyAlbumInTrackOrder: %t\n", f.CopyAlbumInTrackOrder))
	r.WriteString(fmt.Sprintf("CSV: %t\n", f.CSV))
	r.WriteString(fmt.Sprintf("Debug: %t\n", f.Debug))
	r.WriteString(fmt.Sprintf("DoInventory: %t\n", f.DoInventory))
	r.WriteString(fmt.Sprintf("DoRename: %t\n", f.DoRename))
	r.WriteString(fmt.Sprintf("DoSummary: %t\n", f.DoSummary))
	r.WriteString(fmt.Sprintf("SuppressTitles: %t\n", f.SuppressTitles))
	r.WriteString(fmt.Sprintf("JsonOutput: %t\n", f.JsonOutput))
	r.WriteString(fmt.Sprintf("NoGroup: %t\n", f.NoGroup))
	r.WriteString(fmt.Sprintf("NoTags: %t\n", f.NoTags))
	r.WriteString(fmt.Sprintf("UnicodePunct: %t\n", f.UnicodePunct))
	return r.String()
}

// copy user set flags to a local store
func (g *GlobalVars) SetFlagArgs(f FlagST) {
	g.localFlags = new(FlagST)
	g.localFlags.CompareTagsToTitle = f.CompareTagsToTitle
	g.localFlags.CopyAlbumInTrackOrder = f.CopyAlbumInTrackOrder
	g.localFlags.CSV = f.CSV
	g.localFlags.Debug = f.Debug
	g.localFlags.DoRename = f.DoRename
	g.localFlags.DoInventory = f.DoInventory
	g.localFlags.DoSummary = f.DoSummary
	g.localFlags.NoTags = f.NoTags
	g.localFlags.JsonOutput = f.JsonOutput
	g.localFlags.SuppressTitles = f.SuppressTitles
	g.localFlags.NoGroup = f.NoGroup
	g.localFlags.NoTags = f.NoTags
	g.localFlags.UnicodePunct = f.UnicodePunct
}
func (g *GlobalVars) Flags() *FlagST {
	return g.localFlags
}

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

func AllocateData() *GlobalVars {
	rval := new(GlobalVars)
	rval.localFlags = new(FlagST)
	if rval.localFlags == nil {
		fmt.Println("PIB in allocate Data, localflags is nil")
	}
	return rval
}
