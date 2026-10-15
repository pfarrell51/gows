// use HTML parse to read an HTML file and process the DOM
// usage, if no file name is specified, use ../index.html
// if specified, read it.
package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func main() {
	var inStr string

	if len(os.Args) == 1 {
		inStr = "../index.html"
	} else {
		inStr = os.Args[1]
	}

	// Reads the entire file content into memory
	content, err := os.ReadFile(inStr)
	if err != nil {
		log.Fatalf("failed to read  html: %v", err)
	}
	htmlIn := bytes.NewBuffer(content)

	// Parse the HTML content from the file (which implements io.Reader)
	doc, err := html.Parse(htmlIn)
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

func dump(n *html.Node) {
	fmt.Printf("\nDump\nType %s  Element tag: %s n: %s\n", n.Type.String(), n.Data, n.DataAtom.String())
	fmt.Printf("Element tag: %s n: %s\n", n.Data, n.DataAtom.String())
	fmt.Printf("Parent %v\n", n.Parent)
	fmt.Printf("FirstChild %v\n", n.FirstChild)
	fmt.Printf("LastChild %v\n", n.LastChild)
	fmt.Printf("PrevSib %v\n", n.PrevSibling)
	fmt.Printf("NextSib %v\n\n", n.NextSibling)
}
