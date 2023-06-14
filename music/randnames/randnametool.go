package randnames

import (
	"fmt"
	"math/rand"
)

func foo() {
	s := rand.NewSource(1234567)
	r := rand.New(s)

	fmt.Println(r.Perm(20))
}
