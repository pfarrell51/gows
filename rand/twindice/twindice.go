// this program helps clean up the mp3 and flac files in my hits playlist
// its main task is to normalize the file names to relect the artist and song title
//
// this is not multi-processing safe

// Bugs

package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	start := time.Now()
	var values [7][7]int

	const samples = 100* 1000 * 1000
	const sides = 6
	for is := 0; is < samples; is++ {
		nl := rand.Intn(sides) + 1
		nr := rand.Intn(sides) + 1
		values[nl][nr]++

	}
	for k := 1; k < 70; k++ {
		i := k / 10
		j := k % 10
		if i < 7 && j < 7 {
			fmt.Printf("%d,%d\n", i*10+j, values[i][j])
		} else {
			fmt.Printf("%d,0\n", i*10+j)
		}
	}

	duration := time.Since(start)
	fmt.Printf("# %v\n", duration)
}
