package util

import (
	"fmt"
	"testing"
)

const bvSize = 26

func TestSize(t *testing.T) {
	var fiveHundredBits = NewBitVector(bvSize)
	l := len(fiveHundredBits.store)
	if l != (bvSize+63)/64 {
		t.Errorf("not equal %d, got %d", (bvSize+63)/64, l)
	}
	m := fiveHundredBits.max
	if m != bvSize {
		t.Errorf("not equal %d, got %d", bvSize, m)
	}
}

func TestAllTrue(t *testing.T) {
	var fHB = NewBitVector(bvSize)
	fHB.Set(3)
	fHB.Set(8)
	fHB.Set(9)
	fHB.Set(13)
	fHB.Set(14)
	fHB.Set(15)
	fHB.Set(16)
	fHB.Set(20)
	fr := fHB.AllTrue()
	fmt.Printf("fr: %v\n", fr)
	if len(fr) != 4 {
		t.Errorf("wrong length returned: %d", len(fr))
	}
	if fr[1][0] != 8 || fr[1][1] != 9 {
		t.Errorf("wrong run1 position returned: %v", fr)
	}
}
func TestRun(t *testing.T) {
	var fHB = NewBitVector(bvSize)
	fHB.Set(3)
	fHB.Set(8)
	fHB.Set(9)
	fHB.Set(13)
	fHB.Set(14)
	fHB.Set(15)
	fHB.Set(16)
	fr := fHB.FirstRun(-1)
	if len(fr) != 2 {
		t.Errorf("wrong length returned: %d", len(fr))
	}
	if fr[0] != 8 || fr[1] != 9 {
		t.Errorf("wrong run1 position returned: %v", fr)
	}

	skip := fr[1]
	fr = fHB.FirstRun(bvSize * 2)
	if len(fr) > 1 {
		t.Errorf("for bad index, wrong length returned: %d", len(fr))
	}

	fr = fHB.FirstRun(skip)

	if len(fr) != 2 {
		t.Errorf("wrong length returned: %d", len(fr))
		return
	}
	if fr[0] != 13 || fr[1] != 16 {
		t.Errorf("wrong run2 position returned: %v", fr)
	}
}
