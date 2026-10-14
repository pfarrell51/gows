// test of parent/child relationship
// of nested HTML elements.

package main

import (
	"fmt"
	"log"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func main() {
	const inStr = "<html><a href=\"foo\"><img src=baz></a>fratz</html>"
	aReader :=
		strings.NewReader(inStr)

	// Parse the HTML content from the file (which implements io.Reader)
	doc, err := html.Parse(aReader)
	if err != nil {
		log.Fatalf("failed to parse html: %v", err)
	}
	checkParent := false
	for n := range doc.Descendants() {
		fmt.Printf("n %p nType: %s atom: %s %t\n", n, n.Type, n.DataAtom, checkParent)
		switch n.Type {
		case html.ElementNode:
			checkParentType(n, checkParent)
			if n.DataAtom == atom.A {
				fmt.Printf("got an A, children??\n")
				for c := range n.Descendants() {
					checkParent = true
					fmt.Printf("c %p cType: %s atom: %s AD: %s\n", c, c.Type, c.DataAtom, c.Data)
				}
				fmt.Println()
			}
		case html.TextNode:
			fmt.Printf("TextNode: :%s; len %d no& %p CP: %t parent %p\n", n.Data, len(n.Data), n.Parent, checkParent, n.Parent)
			checkParentType(n, checkParent)
			checkParent = false
		default:
			fmt.Printf("defaulted t: %s: %s\n", n.DataAtom.String(), n.Data)
		}
	}
}
func checkParentType(n *html.Node, c bool) {
	if c {
		if n.Parent.DataAtom == atom.A {
			fmt.Printf("Yay, parent is an A href %p\n", n.Parent)
		} else {
			fmt.Printf("bad, %p a: %s\n", n.Parent, n.Parent.DataAtom)
		}
	}
}
