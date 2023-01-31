package musictools

import (
	"encoding/json"
	"fmt"
	"log"
)

func PrintJson(m map[string]Song) {
	fmt.Printf("map size: %d\n", len(m))
	var songs [500]Song
	i := 0
	for k, v := range m {
		fmt.Printf("i: %d k: %#v v %#v\n", i, k, v)
		songs[i] = v
		i++
	}
	data, err := json.Marshal(songs)
	if err != nil {
		log.Fatalf("JSON marshalling failed: %s", err)
	}
	fmt.Printf("len: %d, %s\n", len(data), string(data))
	fmt.Println("begin Json dump")
	fmt.Println(data)
	fmt.Println("end json dump")
}
