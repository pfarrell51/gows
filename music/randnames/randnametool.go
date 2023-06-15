package randnames

import (
	"fmt"
	"math/rand"
)

type FlagST struct {
	NoTags    bool
	TwoLetter bool
}

var localFlags = new(FlagST)

func foo() {
	s := rand.NewSource(1234567)
	r := rand.New(s)

	fmt.Println(r.Perm(20))
}
