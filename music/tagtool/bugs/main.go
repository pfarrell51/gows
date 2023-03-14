package main

import (
	"fmt"
	"os"

	"github.com/dhowden/tag"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Usage, bugs [filespec]")
		return
	}
	arg := os.Args[1]
	if arg == "" {
		fmt.Println("PIB, input path empty")
		return
	}
	file, err := os.Open(arg)
	if err != nil {
		fmt.Printf("err : %v %s\n", err, arg)
		return
	}
	defer file.Close()

	m, err := tag.ReadFrom(file)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	if m == nil {
		fmt.Printf("tag.ReadFrom (file) turned nil but no error for %s\n", arg)
	}
}
