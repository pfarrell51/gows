package musictools

import (
	"encoding/json"
	"fmt"
	"log"
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
	fmt.Printf("%s\n", string(data))
}
