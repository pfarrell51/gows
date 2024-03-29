/*
create word list for john the ripper
*/

package main

import (
	"fmt"
	"os"
)

const max = 20 * 1000

// get argument from shell, if any or default to "xyz"
func main() {
	prefix := "xyz"
	if len(os.Args) > 1 {
		prefix = os.Args[1]
	} else {
		fmt.Printf("usage %s <prefix>\n", os.Args[0])
		return
	}
	for i := 0; i < max; i++ {
		fmt.Printf("%s%04d\n", prefix, i)
	}
}
