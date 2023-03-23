// encodings_test

package tagtool

import (
	"bufio"
	"bytes"
	"fmt"
	"sort"
	"strings"
	"testing"
	"unicode/utf8"
)

var tsongs []Song
var m map[string]Song

func init() {
	tsongs = []Song{
		Song{Artist: "Heart", Album: "Dreamboat Annie", Title: "Crazy On You", Genre: "Rock", Disc: 1, DiscCount: 1, Track: 3, Year: 1976},
		Song{Artist: "Elvis Presley", Album: "Forest Gump", Title: "Hound Dog", Genre: "Rock", Disc: 1, DiscCount: 2, Year: 1957},
		Song{Artist: "Santana", Album: "All That I Am", Title: "Brown Skin Girl", Year: 2005, Genre: "Rock", Disc: 1, DiscCount: 1, Track: 11},
		Song{Artist: "Heart", Album: "Dog & Butterfly", Title: "Baracuda", Genre: "Rock", Year: 1980, Track: 6, Disc: 1, DiscCount: 1},
		Song{Artist: "Animals", Album: "Best of 60s", Title: "House of Rising Sun", Year: 1965, Genre: "Rock", Disc: 1, DiscCount: 1},
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
func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
func compareLine(a, b string) int {
	wa := strings.TrimSpace(a)
	la := utf8.RuneCountInString(wa)
	wb := strings.TrimSpace(b)
	lb := utf8.RuneCountInString(wb)
	if la != lb {
		if la < lb {
			return -100
		}
		return 100
	}
	ln := min(la, lb)
	for i := 0; i < ln; i++ {
		ra, sa := utf8.DecodeRuneInString(a)
		rb, sb := utf8.DecodeRuneInString(b)
		if sa != sb {
			panic("rune sizes not same")
		}
		switch {
		case ra == rb:
			continue
		case ra < rb:
			return -i
		case ra > rb:
			return i
		}
	}
	return 0
}
func compareLineArrays(a, b []string) bool {
	var la, lb, matched int
	for i := 0; i < len(a); i++ {
		if len(a[i]) > 2 {
			la++
			a[i] = strings.TrimPrefix(a[i], "[")
			a[i] = strings.TrimPrefix(a[i], ",")
			a[i] = strings.TrimSuffix(a[i], "]")
			a[i] = strings.TrimSuffix(a[i], "]")
		} else {
			continue
		}
		//fmt.Printf("cLAa) %d %s\n", i, a[i])
		for j := 0; j < len(b); j++ {
			if len(b[j]) > 2 {
				lb++
				b[j] = strings.TrimPrefix(b[j], "[")
				b[j] = strings.TrimPrefix(b[j], ",")
				b[j] = strings.TrimSuffix(b[j], "]")
				b[j] = strings.TrimSuffix(b[j], "]")
				if rv := compareLine(a[i], b[j]); rv == 0 {
					//	fmt.Printf("matched! %d  %s\n", j, b[j])
					matched++
					break
				} else {
					// fmt.Printf("%d cLAb)  %d >%s< != \n%d >%s<\n", rv, j, b[j], i, a[i])
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
	expected := `[{"Artist":"Crosby, Still \u0026 Nash","Album":"Deja Vu","Title":"Suite, Judy Blue Eyes","Genre":"Rock","Disc":1,"DiscCount":1,"Track":1,"Year":1969}
,{"Artist":"Santana","Album":"All That I Am","Title":"Brown Skin Girl","Year":2005,"Genre":"Rock","Disc":1,"DiscCount":1,"Track":11}
,{"Artist":"Animals","Album":"Best of 60s","Title":"House of Rising Sun","Genre":"Rock","Disc":1,"DiscCount":1,"Year":1965}
,{"Artist":"Heart","Album":"Dog \u0026 Butterfly","Title":"Baracuda","Genre":"Rock","Disc":1,"DiscCount":1,"Track":6,"Year":1980}
,{"Artist":"Heart","Album":"Dreamboat Annie","Title":"Crazy On You","Genre":"Rock","Disc":1,"DiscCount":1,"Track":3,"Year":1976}
,{"Artist":"Elvis Presley","Album":"Forest Gump","Title":"Hound Dog","Genre":"Rock","Disc":1,"DiscCount":2,"Year":1957} ]`
	elines := splitLines(expected)
	if !compareLineArrays(jlines, elines) {
		t.Errorf("did not match expeccted JSON  lines")
	}
}
func TestPrintCSV(t *testing.T) {
	fmt.Printf("testing CSV\n")
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	PrintCSVtoWriter(w, m)
	cs := b.String()
	clines := splitLines(cs)
	expected := `"Crosby, Still & Nash",Deja Vu,"Suite, Judy Blue Eyes",Rock,1,1969,
Animals,Best of 60s,House of Rising Sun,Rock,0,1965,
Heart,Dog & Butterfly,Baracuda,Rock,6,1980,
Heart,Dreamboat Annie,Crazy On You,Rock,3,1976,
Elvis Presley,Forest Gump,Hound Dog,Rock,0,1957,
Santana,All That I Am,Brown Skin Girl,Rock,11,2005,
`
	elines := splitLines(expected)
	if !compareLineArrays(clines, elines) {
		t.Errorf("did not match expected CSV")
		fmt.Printf("cl: %s\n", clines)
		fmt.Printf("ex: %s\n", elines)
	}

}
func TestPrintSortedCSV(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	var triParts []InventorySong
	for _, s := range tsongs {
		tP := InventorySong{artist: s.Artist,
			title: s.Title}
		fmt.Printf("tp: %v\n", tP)
		triParts = append(triParts, tP)
		fmt.Printf("triP: %v\n", triParts)
	}
	fmt.Println()
	sort.Sort(ByThree(triParts))
	PrintSortedCSVtoWriter(w, triParts, true)
	cs := b.String()
	clines := splitLines(cs)
	expected := `"Animals,Best of 60s,House of Rising Sun,
Crosby, Still & Nash",Deja Vu,"Suite, Judy Blue Eyes",
Elvis Presley,Forest Gump,Hound Dog,
Heart,Dog & Butterfly,Baracuda,
Heart,Dreamboat Annie,Crazy On You,
Santana, All That I Am,
`
	elines := splitLines(expected)
	if !compareLineArrays(clines, elines) {
		t.Errorf("did not match expected Sorted CSV")
		fmt.Printf("Sorted CSV lines: %s\n", clines)
		fmt.Printf("expected %s\n", elines)
	}
}
