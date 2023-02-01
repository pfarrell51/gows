// this module deals the meta tags in mp3, flac, ogg and other neat mustic formates.
// nearly all of the work is done by  David Howden wonderful library

package musictools

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/dhowden/tag"
)

func GetMetaData(pathArg, p string) *Song {
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
	}
	return rval
}
