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
	for n := range doc.Descendants() {
		checkParent := false
		fmt.Printf("n %p nType: %s atom: %s\n", &n, n.Type, n.DataAtom)
		switch n.Type {
		case html.ElementNode:
			if n.DataAtom == atom.A {
				fmt.Printf("got an A, children??\n")
				for c := range n.Descendants() {
					checkParent = true
					fmt.Printf("c %p cType: %s atom: %s AD: %s\n", &c, c.Type, c.DataAtom, c.Data)
				}
				fmt.Println()
			}
		case html.TextNode:
			text := n.Data
			fmt.Printf("TextNode: :%s; len %d no& %p CP: %t parent %p\n", text, len(text), n.Parent, checkParent, &n.Parent)
			checkParent = false
		default:
			fmt.Printf("defaulted t: %s: %s\n", n.DataAtom.String(), n.Data)
		}
	}
}
