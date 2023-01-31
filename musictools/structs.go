// this program helps clean up the mp3 and flac files in my hits playlist
// its main task is to normalize the file names to relect the artist and song title
//
// this is not multi-processing safe

// Bugs

package musictools

import (
	"github.com/dlclark/metaphone3"
	g "github.com/zyedidia/generic"
	"github.com/zyedidia/generic/btree"
)

type Song struct {
	alreadyNew        bool
	artist            string
	artistH           string
	artistHasThe      bool
	artistInDirectory bool
	artistKnown       bool
	album             string
	albumH            string
	title             string
	titleH            string
	track             int
	year              int
	inPath            string
	outPath           string
	ext               string
}

type FlagST struct {
	ShowArtistNotInMap bool
	DoRename           bool
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
	localFlags.DoRename = f.DoRename
	localFlags.JustList = f.JustList
	localFlags.NoGroup = f.NoGroup
	localFlags.ZDumpArtist = f.ZDumpArtist
	localFlags.JsonOutput = f.JsonOutput
	localFlags.Debug = f.Debug
	if localFlags.Debug {
		localFlags.JsonOutput = true
	}
}
func GetFlags() *FlagST {
	return localFlags
}

var enc metaphone3.Encoder

func GetEncoder() *metaphone3.Encoder {
	return &enc
}

var gptree = btree.New[string, string](g.Less[string])

func GetArtistMap() *btree.Tree[string, string] {
	return gptree
}
func LoadArtistMap() {
	groupNames := [...]string{
		"5th_Dimension", "ABBA", "Alice Cooper", "Alison_Krauss", "AllmanBrothers", "Almanac_Singers",
		"Animals", "Aquarius", "Aretha Franklin", "Arlo_Guthrie", "Association", "Average White Band",
		"Band", "Basia", "BeachBoys", "Beatles", "Bee Gees", "Billy Joel", "BlindFaith",
		"BloodSweatTears", "Blue Oyster Cult", "Blues Brothers", "Bob Dylan", "Boston", "Box Tops", "Bread",
		"Brewer and Shipley", "Brewer & Shipley", "BuffaloSpringfield", "Byrds",
		"Carole King", "Carpenters", "Cheap Trick", "Chesapeake", "Cream", "Crosby & Nash",
		"Crosby and Nash", "Crosby Stills & Nash", "Crosby Stills And Nash", "CSN&Y",
		"Crosby Stills Nash Young", "Crosby Stills Nash & Young", "David Allan Coe",
		"David Bowie", "David_Bromberg", "Deep Purple", "Derek and the Dominos",
		"Derek_Dominos", "Detroit Wheels",
		"Dire_Straits", "Doc Watson", "Don McLean", "Doobie_Brothers", "Doors", "Dylan",
		"Elton_John", "Emerson, Lake & Palmer", "Emmylou_Harris", "Fifth Dimension",
		"Fleetwood_Mac", "Genesis",
		"George Harrison", "Graham Nash", "Gram Parsons", "Hall and Oates", "Hall & Oates",
		"Heart", "Isley Brothers", "Jackie Wilson", "Jackson Browne",
		"James_Taylor", "Jefferson_Airplane", "Jethro_Tull", "Jimmy Buffett", "John_Denver",
		"John_Hartford", "John_Starling", "Joni_Mitchell", "Judy_Collins", "Kansas",
		"KC The Sunshine Band", "Kingston_Trio", "Led_Zepplin", "Linda_Ronstadt",
		"Lovin Spoonful", "Lynyrd_Skynyrd",
		"Mamas And Papas", "Mamas & The Papas", "Maria Muldaur",
		"Meatloaf", "Mike_Auldridge", "Mith Ryder & Detroit Wheels", "Moody Blues",
		"Neal Young", "Neil Diamond",
		"New Riders of the Purple Sage", "New_Riders_Purple_Sage",
		"Nitty Gritty Dirt Band", "Oates", "Otis Redding", "Pablo_Cruise", "Palmer",
		"Paul_Simon", "Peter_Paul_Mary", "Rascals", "Ringo Starr", "Roberta Flack", "Rolling_Stones",
		"Roy_Orbison", "Sam And Dave", "Santana", "Seals and Crofts", "Seals_Croft", "Seldom_Scene",
		"Shadows Of Knight", "Simon and Garfunkel", "Simon_Garfunkel", "Sonny And Cher",
		"Spoonful", "Seals & Crofts", "Steely_Dan", "Steppenwolf", "Steven_Stills",
		"Stevie Ray Vaughan and Double Trouble", "Sting", "Sunshine Band",
		"Three Dog Night", "TonyRice", "Traveling_Wilburys", "Turtles", "Warren Zevon",
		"Who", "Wilson Pickett", "Yes",
	}
	for _, n := range groupNames {
		prim, sec := enc.Encode(JustLetter(n))
		gptree.Put(prim, n)
		if len(sec) > 0 {
			gptree.Put(sec, n)
		}
	}
}
