package util

import (
	"fmt"
	"testing"
)

const bvSize = 26

var turnOn = []int{3, 8, 9, 13, 14, 15, 16, 20}

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
func makeTestBV() BitVector {
	//const	turnOn := []int{3, 8, 9, 13, 14, 15, 16, 20}
	var bv = NewBitVector(bvSize)
	for i := 0; i < len(turnOn); i++ {
		bv.Set(turnOn[i])
	}
	return bv
}
func TestSetClear(t *testing.T) {
	var bv = NewBitVector(bvSize)
	bv.Set(3)
	bv.Set(8)
	bv.Set(9)
	fmt.Printf("c: %b %x\n", bv.store, bv.store)
	bv.Clear(8)
	fmt.Printf("c: %b %x\n", bv.store, bv.store)
}
func TestAllTrue(t *testing.T) {
	var banana = NewBitVector(6)
	banana.Set(1)
	banana.Set(3)
	banana.Set(5)
	ib := banana.LogicalInvert()
	fmt.Printf("b: %x  ib: %x\n", banana.store, ib.store)

	br := banana.AllTrue()
	fmt.Printf("br: %v\n", br)

	var fHB = makeTestBV()
	fr := fHB.AllTrue()
	fmt.Printf("fr: %v\n", fr)
	if len(fr) != 4 {
		t.Errorf("wrong length returned: %d", len(fr))
	}
	if fr[1][0] != 8 || fr[1][1] != 10 {
		t.Errorf("wrong run1 position returned: %v", fr)
	}
}
func TestRun(t *testing.T) {
	var fHB = makeTestBV()
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
func TestOn(t *testing.T) {
	var fHB = makeTestBV()
	var any = fHB.AnyOn()
	if !any {
		t.Errorf("wrong AnyOn: %v", any)
	}
}
func TestInvert(t *testing.T) {
	var fHB = makeTestBV()
	var inv = fHB.LogicalInvert()
	if fHB.Get(0) {
		t.Errorf("inv wrong for 0")
	}
	if inv.Get(0) {
		t.Errorf("wrong inv0")
	}
}
