package main

import (
	"encoding/json"
	"os"
)

type OneRec struct {
	Leadin     string `json:"leadin,omitempty"`
	Anchor     string `json:"anchor,omitempty"`
	AnchorText string `json:"anchortext,omitempty"`
	EmbedImg   string `json:"embedimage,omitempty"`
	TwoColumn  bool   `json:"twocolumn,omitempty"`
}

func main() {
	data := []OneRec{
		{Leadin: "Bob", Anchor: "href mumble goes here", AnchorText: "Bob Dahlquist"},
		{Leadin: "Harry", Anchor: "href Photo of Harry Scogin", AnchorText:  "Harry Scogin Dahlquist"},
		{Leadin: "Hank", Anchor: "href Photo of Hank Farrell",  AnchorText: "Hank Farrell", TwoColumn: true},
	}	

	// 1. Create or open the destination file
	file, err := os.Create("output_stream.json")
	if err != nil {
		panic(err)
	}
	defer file.Close() // Ensure the file closes when done

	// 2. Create an encoder pointing to the file and stream the data
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Optional: Remove this line for minified JSON

	err = encoder.Encode(data)
	if err != nil {
		panic(err)
	}
}
