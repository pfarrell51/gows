// convert an GPX file into a simple csv file
// reads stdin

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	var tag, tval, lat, lng, ele, spd string
	var pos int
	active := false
	fmt.Printf("#time, lat, long, ele, spd\n")
	scanner := bufio.NewScanner(os.Stdin)
	ln := 0
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), " \t")
		ln++
		if strings.Compare(line, "</gpx>") == 0 {
			printCSV(tval, lat, lng, ele, spd)
			return
		}
		tag, pos = getTag(line)
		switch tag {
		case "trkpt":
			if active {
				printCSV(tval, lat, lng, ele, spd)
				tval = ""
				lat = ""
				lng = ""
				ele = ""
				spd = ""
			}
			active = true
			i := strings.Index(line[pos:], "lat=")
			i += (pos + 4) // length of lat="
			lat, pos = getQuoted(line[i:])
			i += pos
			lns := strings.Index(line[i:], "lon=")
			i += (lns + 4)
			lng, pos = getQuoted(line[i:])

		case "speed":
			spd = getXmlVal(line)
		case "time":
			tval = getXmlVal(line)
		case "ele":
			ele = getXmlVal(line)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

// get the initial XML tag from the current line, returns tag and offset to end of tag
func getTag(line string) (tag string, pos int) {
	var i, j int
	i = strings.Index(line, "<") + 1
	if i < 0 {
		return "", 0
	}
	j = strings.Index(line[i:], " ")
	if j > 1 {
		j = j + i
	} else {

		j = strings.Index(line[i:], ">")
		if j < 0 {
			fmt.Printf("GT(bail) negative %d\n", j)
		} else {
			j++

			if j == len(line) {
				fmt.Printf("GT(bail) %d %s\n", j, line[j:])
				return "", 0
			}
		}
	}
	return line[i:j], i + j
}

// getXmlVal returns the value between a pair of matching XML tags, such as <foo>bar</foo>
// returns bar
func getXmlVal(line string) string {
	if len(line) == 0 {
		return ""
	}
	i := strings.Index(line, "<")
	if i < 0 {
		return ""
	}
	j := strings.Index(line, ">")
	if j < 0 {
		return ""
	}
	k := strings.Index(line[1:], "<")
	if k < 0 {
		return ""
	}
	j++
	k++
	return line[j:k]
}

// gets the value between double quotes, sush as <foo bar="mumble">
// returns mumble
func getQuoted(line string) (val string, pos int) {
	i := strings.Index(line, "\"")
	if i < 0 {
		return "", 0
	}
	i++
	j := strings.Index(line[i:], "\"")
	if j < 0 {
		return "", 0
	}
	j++ // skip over tailing quote
	return line[i:j], i + j
}
func printCSV(tval, lat, lng, ele, spd string) {
	fmt.Printf("%s, %s, %s, %s, %s\n", tval, lat, lng, ele, spd)
}
