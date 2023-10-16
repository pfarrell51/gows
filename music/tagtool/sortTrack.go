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

type InventorySong struct {
	artist string
	album  string
	title  string
}

func (s InventorySong) String() string {
	return fmt.Sprintf("%s: %s -->  %s", s.artist, s.album, s.title)
}

// ByTrack implements sort.Interface for []TrackSong based on
// the track number field.
type ByTrack []TrackSong

func (a ByTrack) Len() int           { return len(a) }
func (a ByTrack) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTrack) Less(i, j int) bool { return a[i].Track < a[j].Track }

// ByTwo implements sort.Interface for []InventorySong based on artist & album fields.
type ByTwo []InventorySong

func (a ByTwo) Len() int      { return len(a) }
func (a ByTwo) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByTwo) Less(i, j int) bool {
	p, q := a[i].artist, a[j].artist
	switch {
	case p < q:
		return true
	case q < p:
		return false
	}
	return a[i].album < a[j].album
}

// ByThree implements sort.Interface for []InventorySong based on the three fields.
type ByThree []InventorySong

func (a ByThree) Len() int      { return len(a) }
func (a ByThree) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByThree) Less(i, j int) bool {
	p, q := a[i].artist, a[j].artist
	switch {
	case p < q:
		return true
	case q < p:
		return false
	}
	// p == q, so try next one
	p, q = a[i].album, a[j].album
	switch {
	case p < q:
		return true
	case q < p:
		return false
	}
	// album and artist same
	return a[i].title < a[j].title
}

type TrackSong struct {
	Track int
	Path  string
}

func (p TrackSong) String() string {
	return fmt.Sprintf("%s: %d", p.Path, p.Track)
}

func (g *GlobalVars) AddSongForTrackSort(a Song) {
	g.tracksongs = append(g.tracksongs, TrackSong{a.Track, a.inPath})
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
