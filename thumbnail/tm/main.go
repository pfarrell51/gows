// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// The thumbnail command produces thumbnails of JPEG files
// whose names are provided on each line of the standard input.
package main

import (
	"bufio"
	"fmt"
	"github.com/pfarrell51/gows/thumbnail"
	"log"
	"os"
)

func main() {
	if false {
		if len(os.Args) == 1 {
			fmt.Printf("usage %s <directory>\n", os.Args[0])
			return
		}
	}
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		fmt.Println(input.Text())
		thumb, err := thumbnail.ImageFile(input.Text())
		if err != nil {
			log.Print(err)
			continue
		}
		fmt.Println(thumb)
	}
	if err := input.Err(); err != nil {
		log.Fatal(err)
	}
}
