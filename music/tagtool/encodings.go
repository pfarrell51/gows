package tagtool

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

func PrintJson(m map[string]Song) {
	var songs []Song
	for _, v := range m {
		songs = append(songs, v)
	}
	data, err := json.Marshal(songs)
	if err != nil {
		log.Fatalf("JSON marshalling failed: %s", err)
	}
	fmt.Println(strings.ReplaceAll(string(data), "}", string("}\n")))
}

func PrintCSV(m map[string]Song) {
	var songs [][]string
	for _, v := range m {
		var aSong []string
		aSong = append(aSong, v.Artist)
		aSong = append(aSong, v.Album)
		aSong = append(aSong, v.Title)

		songs = append(songs, aSong)
	}
	w := csv.NewWriter(os.Stdout)

	for _, song := range songs {
		if err := w.Write(song); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}

	// Write any buffered data to the underlying writer (standard output).
	w.Flush()

	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
}
