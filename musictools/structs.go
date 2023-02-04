// this  file contains structures and other "global" data stores

package musictools

import (
	"fmt"
	"regexp"
	"strings"

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
	ShowArtistNotInMap bool
	DoRenameFilename   bool
	DoRenameMetadata   bool
	JustList           bool
	NoGroup            bool
	ZDumpArtist        bool
	JsonOutput         bool
	Debug              bool
}

var localFlags = new(FlagST)

// copy user set flags to a local store
func SetFlagArgs(f FlagST) {
	localFlags.ShowArtistNotInMap = f.ShowArtistNotInMap
	localFlags.DoRenameFilename = f.DoRenameFilename
	localFlags.DoRenameMetadata = f.DoRenameMetadata
	localFlags.JustList = f.JustList
	localFlags.NoGroup = f.NoGroup
	localFlags.ZDumpArtist = f.ZDumpArtist
	localFlags.JsonOutput = f.JsonOutput
	localFlags.Debug = f.Debug
}
func GetFlags() *FlagST {
	return localFlags
}

var enc metaphone3.Encoder

func init() {
	enc.Encode("ignore this")
	enc.MaxLength = 20
}
func StandardizeArtist(art string) string {
	if len(art) == 0 {
		return art
	}
	rval := strings.TrimSpace(art)
	if rval == "The The" {
		return rval
	}
	if matched, _ := regexp.MatchString("^[A-Z]( |-|_)", rval); matched {
		rval = rval[2:]
	}
	if strings.HasPrefix(rval, "The ") {
		rval = rval[4:]
	}
	return rval
}
func StandardizeTitle(title string) string {
	if len(title) == 0 {
		return title
	}
	rval := strings.TrimSpace(title)
	if matched, _ := regexp.MatchString("^[A-Z]( |-|_)", rval); matched {
		rval = rval[2:]
	}
	if strings.HasPrefix(rval, "The ") {
		rval = rval[4:]
	}
	rval = strings.ReplaceAll(rval, "/", " ")
	rval = strings.ReplaceAll(rval, "_", " ")
	return rval
}
func EncodeTitle(s string) (string, string) {
	prim, sec := enc.Encode(StandardizeTitle(s))
	return prim, sec
}
func EncodeArtist(s string) (string, string) {
	prim, sec := enc.Encode(StandardizeArtist(s))
	return prim, sec
}

var Gptree = btree.New[string, string](g.Less[string])

func GetArtistMap() *btree.Tree[string, string] {
	return Gptree
}
func LoadArtistMap() {
	groupNames := [...]string{
		"5th_Dimension", "ABBA", "Alice Cooper", "Alison Krauss", "Alison Krauss Union Station",
		"Allman Brothers", "Allman Brothers Band", "Almanac Singers",
		"Animals", "Aquarius", "Aretha Franklin", "Arlo Guthrie", "Association", "Average White Band",
		"Band", "Basia", "Beach Boys", "Beatles", "Bee Gees", "Billy Joel", "Blind Faith",
		"Blood Sweat Tears", "Blue Oyster Cult", "Blues Brothers", "Bob Dylan", "Boston", "Box Tops", "Bread",
		"Brewer and Shipley", "Brewer & Shipley", "Buffalo Springfield", "Byrds",
		"Carole King", "Carpenters", "Cheap Trick", "Chesapeake", "Cream", "Crosby & Nash",
		"Crosby and Nash", "Crosby Stills & Nash", "Crosby Stills And Nash", "CSN&Y",
		"Crosby Stills Nash Young", "Crosby Stills Nash & Young", "David Allan Coe",
		"David Bowie", "David Bromberg", "Deep Purple", "Derek and the Dominos",
		"Derek_Dominos", "Detroit Wheels",
		"Dire_Straits", "Doc Watson", "Don McLean", "Doobie Brothers", "Doors", "Dylan",
		"Elton_John", "Emerson, Lake & Palmer", "Emmylou_Harris", "Fifth Dimension",
		"Fleetwood_Mac", "Genesis",
		"George Harrison", "Gillian Welch", "Gillian Welch Alison Krauss", "Graham Nash",
		"Gram Parsons", "Hall and Oates", "Hall & Oates",
		"Heart", "Isley Brothers", "Jackie Wilson", "Jackson Browne",
		"James_Taylor", "Jefferson_Airplane", "Jethro_Tull", "Jimmy Buffett", "John_Denver",
		"John_Hartford", "John_Starling", "Joni_Mitchell", "Judy_Collins", "Kansas",
		"KC The Sunshine Band", "Kingston_Trio", "Led_Zeppelin", "Linda_Ronstadt",
		"Lovin Spoonful", "Lynyrd_Skynyrd",
		"Mamas And Papas", "Mamas & The Papas", "Maria Muldaur",
		"Meatloaf", "Mike_Auldridge", "Mith Ryder & Detroit Wheels", "Moody Blues",
		"Neal Young", "Neil Diamond", "New Riders of the Purple Sage", 
		"Nitty Gritty Dirt Band", "Oates", "Original Soundtrack", "Otis Redding", "Pablo Cruise",
		"Paul_Simon", "Pete Seeger", "Peter Paul Mary", "Rascals", "Ringo Starr",
		"Robert Plant Alison Kraus", "Roberta Flack", "Rolling_Stones",
		"Roy _Orbison", "Sam And Dave", "Santana", "Seals and Crofts", "Seals Croft", "Seldom Scene",
		"Shadows Of Knight", "Simon and Garfunkel", "Simon Garfunkel", "Soggy Mountain Boys", "Sonny And Cher",
		"Spoonful", "Seals Crofts", "Steely_Dan", "Steppenwolf", "Steven_Stills",
		"Stevie Ray Vaughan and Double Trouble", "Sting", "Sunshine Band",
		"Three Dog Night", "TonyRice", "Traveling Wilburys", "Turtles", "Warren Zevon",
		"Who", "Wilson Pickett", "Yes",
	}
	for _, n := range groupNames {
		prim, sec := EncodeArtist(n)
		Gptree.Put(prim, n)
		if len(sec) > 0 {
			Gptree.Put(sec, n)
		}
		if GetFlags().Debug {
			fmt.Printf("%s, %s, %s\n", prim, sec, n)
		}
	}
}
