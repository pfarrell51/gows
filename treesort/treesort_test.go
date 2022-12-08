// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

package treesort

import (
	"math/rand"
	"sort"
	"testing"
)

const numInt = 50

func TestSort(t *testing.T) {
	data := make([]int, numInt)
	for i := range data {
		data[i] = rand.Int() % numInt
	}
	Sort(data)
	if !sort.IntsAreSorted(data) {
		t.Errorf("not sorted: %v", data)
	}
}
func TestSort2(t *testing.T) {
	data := make([]int, numInt)
	for i := range data {
		data[i] = rand.Int() % numInt
	}
	Sort(data)
	d := data[1]
	data[1] = data[numInt-1]
	data[numInt-1] = d
	if sort.IntsAreSorted(data) {
		t.Errorf("not sorted: %v", data)
	}
}
