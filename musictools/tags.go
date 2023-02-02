// this module deals the meta tags in mp3, flac, ogg and other neat mustic formates.
// nearly all of the work is done by  David Howden wonderful library

package musictools

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/dhowden/tag"
	"github.com/dhowden/tag/mbz"
)

var knownIds map[string]bool

func init() {
	knownIds = make(map[string]bool)
}
func GetMetaData(pathArg, p string) *Song {
	if GetFlags().Debug {
		fmt.Printf("in GMD %s\n", p)
	}
	rval := new(Song)
	rval.inPath = path.Join(pathArg, p)
	rval.outPath = pathArg
	file, err := os.Open(rval.inPath)
	defer file.Close()
	if err != nil {
		fmt.Println(err)
		return rval
	}
	m, err := tag.ReadFrom(file)
	if err != nil {
		fmt.Printf("%v %s", err, rval.Title)
		return rval
	}
	rval.Title = m.Title() // The title of the track (see Metadata interface for more details).
	if rval.Title == "" {
		_, filename := path.Split(rval.inPath)
		punchIdx := dashRegex.FindStringIndex(filename)
		if punchIdx != nil {
			rval.Title = strings.TrimSpace(filename[punchIdx[1]:])
		}
	}
	rval.titleH, _ = GetEncoder().Encode(JustLetter(rval.Title))
	rval.Artist = m.Artist()
	rval.Album = m.Album()
	rval.Year = m.Year()
	rval.Track, _ = m.Track()
	info := mbz.Extract(m)
	rval.MBID, _ = info["musicbrainz_trackid"]
	rval.AcoustID, _ = info["acoustid_id"]

	if GetFlags().Debug {
		for k, _ := range info {
			found := knownIds[k]
			if !found {
				knownIds[k] = true
			}
		}
	}

	if GetFlags().Debug {
		fmt.Printf("Format %s Type %s\n", m.Format(), m.FileType())
		if m.Title() != "" {
			fmt.Printf("Title() %v\n", m.Title())
		}
		if m.Album() != "" {
			fmt.Printf("Album() %v\n", m.Album())
		}
		fmt.Printf("Artist() %v\n", m.Artist())
		if m.AlbumArtist() != "" && m.Artist() != m.AlbumArtist() {
			fmt.Printf("AlbumArtist() %v\n", m.AlbumArtist())
		}
		if m.Composer() != "" {
			fmt.Printf("Composer() %v\n", m.Composer())
		}
		//fmt.Printf("Genre() %v\n", m.Genre())
		if m.Year() > 0 {
			fmt.Printf("Year() %#v\n", m.Year())
		}
		if t, _ := m.Track(); t > 0 {
			fmt.Printf("Track %d \n", t)
		}
		if rval.MBID != "" {
			fmt.Printf("Musicbrainz track id %s\n\n", rval.MBID)
		}
	}
	return rval
}
func DumpKnowIDnames() {
	for k, _ := range knownIds {
		fmt.Printf("%s\n", k)
	}
}
