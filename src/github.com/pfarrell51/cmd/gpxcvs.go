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

	var line string
	scanner := bufio.NewScanner(os.Stdin)
	ln := 0
	for scanner.Scan() {
		ln++
		line = strings.Trim(scanner.Text(), " \t")
		if ln < 20 {
			fmt.Printf("B %d %s\n", ln, line)
		}
		part := strings.Split(line, "<")
		for i := 0; i < len(part); i++ {
			fmt.Printf("%d %s\n", i, part[i])
		}
		switch {

		case strings.HasPrefix(line, "<trkpt"):
			fmt.Printf("C %d %s\n", ln, line)

		case strings.HasPrefix(line, "<speed"):
			fmt.Printf("B %d %s\n", ln, line)
			s := strings.Index(line, ">") + 1
			sp := strings.Index(line[s:], "<")
			val := line[s : s+sp]
			fmt.Printf("%s\n", val)
		case strings.HasPrefix(line, "<time"):
			fmt.Printf("C %d %s\n", ln, line)

		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

// parse xml tag and return parts
func getXmlParts(line []string) int {
	return 0
}
