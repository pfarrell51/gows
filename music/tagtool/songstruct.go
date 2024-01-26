// struct for Song

package tagtool

type Song struct {
	Artist            string
	Album             string
	Title             string
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
	InPath            string
	inPathDescent     string // any descent below the pathArg aka outPathBase
	outPath           string
	outPathBase       string // copied from pathArg entered by the user
	ext               string
}

const NAllocateSongs = 10 * 1000
