package tagtool

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"testing"
)

var tsongs []Song
var m map[string]Song

func init() {
	tsongs = []Song{
		Song{Artist: "Animals", Album: "Best of 60s", Title: "House of Rising Sun", Year: 1965, Genre: "Rock", Disc: 1, DiscCount: 1},
		Song{Artist: "Heart", Album: "Dog & Butterfly", Title: "Baracuda", Genre: "Rock", Year: 1980, Track: 6, Disc: 1, DiscCount: 1},
		Song{Artist: "Elvis Presley", Album: "Forest Gump", Title: "Hound Dog", Genre: "Rock", Year: 1957, Disc: 1, DiscCount: 2},
		Song{Artist: "Crosby, Still & Nash", Album: "Deja Vu", Title: "Suite, Judy Blue Eyes", Year: 1969,
			Genre: "Rock", Track: 1, Disc: 1, DiscCount: 1},
	}
	m = make(map[string]Song)
	for i, s := range tsongs {
		key := fmt.Sprintf("A%3d", i)
		m[key] = s
	}
}
func splitLines(in string) []string {
	var rval []string
	sc := bufio.NewScanner(strings.NewReader(in))
	for sc.Scan() {
		rval = append(rval, sc.Text())
	}
	return rval
}
func compareLineArrays(a, b []string) bool {
	var la, lb, matched int
	for i := 0; i < len(a); i++ {
		if len(a[i]) > 2 {
			la++
			a[i] = strings.TrimPrefix(a[i], "[")
			a[i] = strings.TrimPrefix(a[i], ",")
		} else {
			continue
		}
		for j := 0; j < len(b); j++ {
			if len(b[j]) > 2 {
				lb++
				b[j] = strings.TrimPrefix(b[j], "[")
				b[j] = strings.TrimPrefix(b[j], ",")
				if strings.Compare(a[i], b[j]) == 0 {
					matched++
				}
			} else {
				// ignore short line  fmt.Printf("b[j] >%s<\n", b[j])
			}
		}
	}

	return la == matched
}
func TestPrintJson(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	PrintJsontoWriter(w, m)
	w.Flush()
	js := b.String()
	jlines := splitLines(js)

	expected := `[{"Artist":"Animals","Album":"Best of 60s","Title":"House of Rising Sun","Genre":"Rock","Disc":1,"DiscCount":1,"Year":1965}
,{"Artist":"Heart","Album":"Dog \u0026 Butterfly","Title":"Baracuda","Genre":"Rock","Disc":1,"DiscCount":1,"Track":6,"Year":1980}
,{"Artist":"Elvis Presley","Album":"Forest Gump","Title":"Hound Dog","Genre":"Rock","Disc":1,"DiscCount":2,"Year":1957}
,{"Artist":"Crosby, Still \u0026 Nash","Album":"Deja Vu","Title":"Suite, Judy Blue Eyes","Genre":"Rock","Disc":1,"DiscCount":1,"Track":1,"Year":1969}
]`
	elines := splitLines(expected)

	if !compareLineArrays(jlines, elines) {
		t.Errorf("did not match expeccted JSON  lines")
	}
}
func TestPrintCSV(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	PrintCSVtoWriter(w, m)
	cs := b.String()
	clines := splitLines(cs)
	expected := `"Crosby, Still & Nash",Deja Vu,"Suite, Judy Blue Eyes",Rock,1,1969,
Animals,Best of 60s,House of Rising Sun,Rock,0,1965,
Heart,Dog & Butterfly,Baracuda,Rock,6,1980,
Elvis Presley,Forest Gump,Hound Dog,Rock,0,1957,
`
	elines := splitLines(expected)
	if !compareLineArrays(clines, elines) {
		t.Errorf("did not match expected CSV")
	}

}
