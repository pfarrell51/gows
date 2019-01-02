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

		if strings.HasPrefix(line, "<trkpt") {
			fmt.Printf("C %d %s\n", ln, line)
		} else if strings.HasPrefix(line, "<speed") {
			fmt.Printf("B %d %s\n", ln, line)
			s := strings.Index(line, ">")+1
			sp := strings.Index(line[s:], "<")
			val := line[s : s+sp]
			fmt.Printf("C %d %d %d %d %s\n", ln, s, sp, s+sp, val)
			fmt.Printf("%s\n", val)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
