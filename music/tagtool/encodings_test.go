package tagtool

import (
	"fmt"
	"testing"
)

var tsongs []Song
var m map[string]Song

func init() {
	tsongs = []Song{
		Song{Artist: "Animals", Album: "Best of 60s", Title: "House of Rising Sun"},
		Song{Artist: "Heart", Album: "Dog & Butterfly", Title: "Baracuda"},
		Song{Artist: "Elvis Presley", Album: "Forest Gump", Title: "Hound Dog"},
	}
	m = make(map[string]Song)
	for i, s := range tsongs {
		key := fmt.Sprintf("A%3d", i)
		m[key] = s
	}
}
func TestPrintJson(t *testing.T) {
	PrintJson(m)
}
func TestPrintCSV(t *testing.T) {
	PrintCSV(m)
}
