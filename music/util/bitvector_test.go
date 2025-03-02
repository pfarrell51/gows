package util

import (
	//	"fmt"
	"testing"
)

const bvSize = 69

var turnOn = []int{3, 8, 9, 13, 14, 15, 16, 20, 63, 64, 67}

func makeTestBV() BitVector {
	var bv = NewBitVector(bvSize)
	for i := 0; i < len(turnOn); i++ {
		bv.Set(turnOn[i])
	}
	//fmt.Printf("make: %b %x\n", bv.store, bv.store)
	return bv
}
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
func TestSetClear(t *testing.T) {
	var bv = NewBitVector(bvSize)
	bv.Set(3)
	bv.Set(8)
	bv.Set(9)
	if !bv.Get(9) {
		t.Errorf("set not got")
	}
	const cb = 8
	bv.Clear(cb)
	if bv.Get(cb) {
		t.Errorf("setclear not got")
	}
}
func TestTruePositions(t *testing.T) {
	var fHB = makeTestBV()
	fr := fHB.TruePositions()
	if len(fr) != 6 {
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
	var banana = NewBitVector(6)
	banana.Set(1)
	banana.Set(3)
	banana.Set(5)
	ib := banana.LogicalInvert()
	if ib.Get(1) {
		t.Errorf("inv banana inv wrong for 1")
	}

	var fHB = makeTestBV()
	var inv = fHB.LogicalInvert()
	if fHB.Get(0) {
		t.Errorf("inv wrong for 0")
	}
	if !inv.Get(0) {
		t.Errorf("wrong inv0")
	}
}
