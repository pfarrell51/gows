// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 101.

// Package treesort provides insertion sort using an unbalanced binary tree.
package treesort

import (
	"fmt"
)

type tree struct {
	left  *tree
	value int
	right *tree
}

// Sort sorts values in place. Not a tree sort
func Sort(values []int) {
	var root *tree
	for _, v := range values {
		root = root.add(v)
	}
	fmt.Print(root)
	appendValues(values[:0], root)
}

// appendValues appends the elements of t to values in order
// and returns the resulting slice.
func appendValues(values []int, t *tree) []int {
	if t != nil {
		values = appendValues(values, t.left)
		values = append(values, t.value)
		values = appendValues(values, t.right)
	}
	return values
}

func (t *tree) add(value int) *tree {
	if t == nil {
		fmt.Printf("tree was nil, creating val: %d\n", value)
		// Equivalent to return &tree{value: value}.
		t = new(tree)
		t.value = value
		return t
	}
	if value < t.value {
		fmt.Printf("value: %d less adding left\n", value)
		t.left = t.left.add(value)
	} else {
		fmt.Printf("value %d >= adding right\n", value)
		t.right = t.right.add(value)
	}
	return t
}
func (t *tree) String() string {
	if t == nil {
		return "()"
	}
	s := ""
	if t.left != nil {
		s += t.left.String() + " "
	}
	s += fmt.Sprint(t.value)
	if t.right != nil {
		s += " " + t.right.String()
	}
	return "(" + s + ")"
}
