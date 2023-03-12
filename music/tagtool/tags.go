// tags  this module deals the meta tags in mp3, flac, ogg and other neat mustic formates.
// nearly all of the work is done by  David Howden wonderful library

package tagtool

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/dhowden/tag"
	"github.com/dhowden/tag/mbz"
)

// pull meta data from music file. Full file path is pathArg + p
func (g *GlobalVars) GetMetaData(p string) (*Song, error) {
	if g.Flags().Debug {
		fmt.Printf("in GMD %s\n", p)
	}
	rval := new(Song)
	rval.BasicPathSetup(g, p)
	if foundExt := ExtRegex.FindString(p); len(foundExt) > 0 { // redundant check to prevent Bozo programmers
		if g.Flags().Debug {
			fmt.Printf("gmd:foundExt %s\n", foundExt)
		}
	}
	if rval.inPath == "" {
		panic("PIB, input path empty")
	}
	file, err := os.Open(rval.inPath)
	if err != nil {
		fmt.Printf("err : %v %s\n", err, rval.inPath)
		return nil, err
	}
	defer file.Close()
	m, err := tag.ReadFrom(file)
	if err != nil {
		fmt.Printf("%v %s", err, rval.Title)
		return nil, err
	}
	if m == nil {
		fmt.Printf("tag.ReadFrom (file) turned nil but no error for %s\n", p)
	}
	rval.Title = StandardizeTitle(m.Title()) // The title of the track (see Metadata interface for more details).
	if strings.Contains(rval.Title, "/") {
		rval.Title = strings.ReplaceAll(rval.Title, "/", " ")
	}
	//fmt.Printf("in GMD, title is %s for %s\n", rval.Title, rval.inPath)
	rval.titleH, _ = EncodeTitle(rval.Title)
	rval.Artist = StandardizeArtist(m.Artist())
	rval.Album = m.Album()
	rval.albumH, _ = EncodeAlbum(m.Album())
	rval.Genre = m.Genre()
	rval.Year = m.Year()
	rval.Track, _ = m.Track()
	disc, discCount := m.Disc()
	rval.Disc = disc
	rval.DiscCount = discCount
	info := mbz.Extract(m)
	rval.MBID = info["musicbrainz_trackid"]
	rval.AcoustID = info["acoustid_id"]
	if rval.Year == 0 {
		rw := m.Raw() // look in raw map.
		yy := rw["TORY"]
		if yy != nil {
			y, err := strconv.ParseInt(yy.(string), 10, 64)
			if err == nil {
				rval.Year = int(y)
			}
		}
	}
	if g.Flags().Debug {
		for k := range info { // loop thru extra meta data, musicbrainz, etc
			found := g.knownIds[k]
			if !found {
				g.knownIds[k] = true
			}
		}
	}
	if g.Flags().Debug {
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
	return rval, nil
}
func (g *GlobalVars) DumpKnowIDnames() {
	for k := range g.knownIds {
		fmt.Printf("%s\n", k)
	}
}
