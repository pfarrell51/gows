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

// generate json from song map, output to stdout
func PrintJson(m map[string]Song) {
	PrintJsontoWriter(os.Stdout, m)
}

// generate json from song map, output to writer
func PrintJsontoWriter(w io.Writer, m map[string]Song) {
	var songs []Song
	for _, v := range m {
		songs = append(songs, v)
	}
	data, err := json.Marshal(songs)
	if err != nil {
		log.Fatalf("JSON marshalling failed: %s", err)
	}
	result := strings.ReplaceAll(string(data), "}", string("}\n"))
	b := []byte(result)
	w.Write(b)
}

// generate CSV from song map, output to stdout
func PrintCSV(m map[string]Song) {
	PrintCSVtoWriter(os.Stdout, m)
}

// generate CSV from song map, output to writer
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