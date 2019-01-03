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
	var tag string
	var pos int
	scanner := bufio.NewScanner(os.Stdin)
	ln := 0
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), " \t")
		ln++
		if ln < 20 {
			fmt.Printf("A %d %s\n", ln, line)
		}
		tag, pos = getTag(line)
		switch tag {
		case "trkpt":
			fmt.Printf("C %d %d %s\n", ln, pos, tag, line[pos])
			i := strings.Index(line, "lat=\"")
			i = i + 5 // length of lat="
			fmt.Printf("A %d %d %s\n", i, len(line[i:]), line[i:])
			j := strings.Index(line[i+1:], "\"")
			p := line[i:j]

			fmt.Printf("D %d %d %s\n", i, j, p)

		case "speed":
			fmt.Printf("B %d %d %s\n", ln, pos, tag, line[pos])
			val := getXmlVal(line)
			fmt.Printf("B2 %s\n", val)
		case "time":
			val := getXmlVal(line)
			fmt.Printf("C %s\n", val)

		case "ele":
			val := getXmlVal(line)
			fmt.Printf("D %s\n", val)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

// parse xml tag and return parts
func getXmlParts(line string) string {
	if len(line) == 0 {
		return ""
	}
	return line[1:len(line)]
}
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
	fmt.Printf("Va %d:%d %s\n", i, j, line[i:j])
	k := strings.Index(line[1:], "<")
	if k < 0 {
		return ""
	}
	j++
	k++

	fmt.Printf("V3 %d %d %s\n", j, k, line[j:k])
	return line[j:k]
}
func getTag(line string) (string, int) {
	fmt.Printf("GT1 with %s\n", line)
	var tag string
	var i, j int
	i = strings.Index(line, "<") + 1
	fmt.Printf("GT2 i=%d %s\n", i, line[i:])

	j = strings.Index(line[i:], " ")
	fmt.Printf("GT3 %d %d\n", j, len(line))
	if j > 1 {
		j = j + i
	} else {

		j = strings.Index(line[i:], ">")
		if j < 0 {
			fmt.Printf("GT(bail) negative %d\n", j)
		} else {
			j++

			fmt.Printf("GT4 %d %s\n", j, line[i:j])
			if j == len(line) {
				fmt.Printf("GT(bail) %d %s\n", j, line[j:])
				return "", 0
			}
		}
	}
	tag = line[i:j]
	fmt.Printf("GT5 %d : %d %s\n", i, j, tag)
	return tag, j
}
