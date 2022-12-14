package main

import (
	"fmt"
	"math/rand"
)

type Tree[T any] struct {
	cmp  func(a, b string) int
	root *node
}
type nodei[T any] struct {
	left, right *node
	data        T
}

func (bt *Tree[T]) find(val T) **node[T] {
	pl := &bt.root
	for *pl != nil {
		switch cmp := bt.cmp(val, (*pl).data); {
		case cmp < 0:
			pl = &(*pl).left
		case cmp > 0:
			pl = &(*pl).right
		default:
			return pl
		}
	}
	return pl
}
func insert(n node) *node {
	if root == nil {

		return &node{nil, v, nil}
	}
	if it.cmp(v, t.data) < 0 {
		t.left = insert(t.left, v)
		return t
	}
	t.right = insert(t.right, v)
	return t
}
func cmp(a, b string) {
	return strings.compare(a, b)
}
func printTree(t *Tree) {
	if t == nil {
		fmt.Println("tree is empty")
		return
	}
	fmt.println(t.root)
}
func populateTree(t *Tree) {
	t := new(Tree)
	for i := 0; i < 1; i++ {
		n := new(node)
		d := fmt.Sprint("dt %i", i)
		t = t.insert(n, d)
	}
	return t
}
func main() {
	theTree := populateTree()
	printTree(theTree)

}
