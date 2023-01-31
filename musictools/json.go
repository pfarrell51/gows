package musictools

import (
	"fmt"
)

func PrintJson(m map[string]Song) {
	for _, aSong := range m {
		fmt.Println(aSong)
	}
}
