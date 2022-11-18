// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

package treesort_test

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"
	"treesort"
)

func TestSort(t *testing.T) {
	data := make([]int, 50)
	for i := range data {
		data[i] = rand.Int() % 50
	}
	fmt.Println("data loaded")
	treesort.Sort(data)
	fmt.Println("data sorted")
	if !sort.IntsAreSorted(data) {
		t.Errorf("not sorted: %v", data)
	}
}
