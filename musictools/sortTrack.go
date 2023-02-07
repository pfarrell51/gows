// sortTrack
// sort a song by the track field, so we can replicate the order of the original album

package musictools

import (
	"fmt"
	"path"
	"sort"
)

type TrackSong struct {
	Track int
	Path  string
	Ext   string
}

func (p TrackSong) String() string {
	return fmt.Sprintf("%s: %d", p.Path, p.Track)
}

// ByTrack implements sort.Interface for []TrackSong based on
// the Age field.
type ByTrack []TrackSong

func (a ByTrack) Len() int           { return len(a) }
func (a ByTrack) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTrack) Less(i, j int) bool { return a[i].Track < a[j].Track }

var songs []TrackSong = make([]TrackSong, 0)

func AddSongForSort(a Song) error {
	//fmt.Println(a.inPath)
	songs = append(songs, TrackSong{a.Track, a.inPath, a.ext})
	return nil
}
func PrintTrackSortedSongs() {
	cmd := "cp"
	if runtime.GOOS == "windows" {
		cmd = "copy "
	}
	// First, one can define a set of methods for the slice type, as with ByTrack, and
	// call sort.Sort.
	sort.Sort(ByTrack(songs))
	for _, s := range songs {
		sp, sf := path.Split(s.Path)
		fmt.Printf("%s \"%s\"  \"%s%03d-%s\"\n", cmd, s.Path, sp, s.Track, sf)
	}
}
