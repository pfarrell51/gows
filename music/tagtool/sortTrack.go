// sortTrack
// sort a song by the track field, so we can replicate the order of the original album
//
// not multi-thread safe, uses package data store

package tagtool

import (
	"fmt"
	"path"
	"runtime"
	"sort"
)

type TrackSong struct {
	Track int
	Path  string
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

func (g *GlobalVars) AddSongForTrackSort(a Song) error {
	//fmt.Println(a.inPath)
	g.tracksongs = append(g.tracksongs, TrackSong{a.Track, a.inPath})
	return nil
}

// print out shell commands to copy the sorted and renamed fies
func (g GlobalVars) PrintTrackSortedSongs() {
	cmd := "cp"
	if runtime.GOOS == "windows" {
		cmd = "copy "
	}
	// First, one can define a set of methods for the slice type, as with ByTrack, and
	// call sort.Sort.
	sort.Sort(ByTrack(g.tracksongs))
	for _, s := range g.tracksongs {
		sp, sf := path.Split(s.Path)
		fmt.Printf("%s \"%s\"  \"%s%03d-%s\"\n", cmd, s.Path, sp, s.Track, sf)
	}
}
