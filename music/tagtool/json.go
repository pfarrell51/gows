package tagtool

import (
	"encoding/json"
	"fmt"
	"log"
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
