// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

package treesort

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"
	"time"
)

const numInt = 14

func TestTreeCanned(t *testing.T) {
	data := []int{7, 2, 9, 1, 11, 15}
	var root *tree //  &tree{value: value}
	fmt.Println("fresh root ", root)
	for i := range data {
		fmt.Printf("i: %d d:%d\n", i, data[i])
		root = root.add(data[i])
		fmt.Println("root: ", root)
	}
	fmt.Println("final root ", root)
}

func TestSortRan(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())
	data := make([]int, numInt)
	for i := range data {
		data[i] = rand.Int() % numInt
	}
	Sort(data)
	if !sort.IntsAreSorted(data) {
		t.Errorf("not sorted: %v", data)
	}
}
func TestSortiRan2(t *testing.T) {
	data := make([]int, numInt)
	for i := range data {
		data[i] = rand.Int() % numInt
	}
	Sort(data)
	d := data[1]
	data[1] = data[numInt-1]
	data[numInt-1] = d
	if sort.IntsAreSorted(data) {
		t.Errorf("sorted: %v", data)
	}
}
