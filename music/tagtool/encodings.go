// encodings
package tagtool

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

// generate json from song slice, output to stdout
func PrintJson(m []Song) {
	PrintJsontoWriter(os.Stdout, m)
}

// generate json from song slice, output to writer
func PrintJsontoWriter(w io.Writer, m []Song) {
	data, err := json.Marshal(m)
	if err != nil {
		log.Fatalf("JSON marshalling failed: %s", err)
	}
	result := strings.ReplaceAll(string(data), "}", string("}\n"))
	b := []byte(result)
	w.Write(b)
}

// generate CSV from song slice, output to stdout
func PrintCSV(m map[string]Song) {
	PrintCSVtoWriter(os.Stdout, m)
}

// print one song as CSV, using the global csv writer
func (g *GlobalVars) PrintSongToCSV(s *Song) {
	if g.csvWrtr == nil {
		panic("csv writer is nill")
	}
	var aSong []string
	aSong = append(aSong, s.Artist)
	aSong = append(aSong, s.Album)
	if ! g.Flags().JustAlbumArtist {
		aSong = append(aSong, s.Title)
		aSong = append(aSong, s.Genre)
		aSong = append(aSong, strconv.Itoa(s.Track))
		aSong = append(aSong, strconv.Itoa(s.Year))
		aSong = append(aSong, s.MBID)
	}

	if err := g.csvWrtr.Write(aSong); err != nil {
		log.Fatalln("error writing record to csv:", err)
	}
}

// generate CSV from song slice, output to writer
func PrintCSVtoWriter(w io.Writer, m map[string]Song) {
	var songs [][]string
	for _, v := range m {
		var aSong []string
		aSong = append(aSong, v.Artist)
		aSong = append(aSong, v.Album)
		aSong = append(aSong, v.Title)
		aSong = append(aSong, v.Genre)
		aSong = append(aSong, strconv.Itoa(v.Track))
		aSong = append(aSong, strconv.Itoa(v.Year))
		aSong = append(aSong, v.MBID)

		songs = append(songs, aSong)
	}
	cw := csv.NewWriter(w)

	for _, song := range songs {
		if err := cw.Write(song); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}
	cw.Flush() // Write any buffered data to the underlying writer (standard output).
	if err := cw.Error(); err != nil {
		log.Fatal(err)
	}
}
