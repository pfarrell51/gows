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
	for n := range doc.Descendants() {
		checkParent := false
		fmt.Printf("n %p nType: %s atom: %s AD: %s\n", &n, n.Type, n.DataAtom, n.Data)
		switch n.Type {
		case html.ElementNode:
			for i, a := range n.Attr {
				fmt.Printf("%d %v k: %v V: %v\n", i, a.Namespace, a.Key, a.Val)
			}
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
			// Print text contents if not empty/whitespace
			switch len(text) {
			case 0:
			// do nothing
			case 1, 2:
				byteSlice := []byte(text)
				fmt.Println(byteSlice)
			default:
				fmt.Printf("Text: :%s; len %d\n", text, len(text))
			}
			checkParent = false
		default:
			fmt.Printf("defaulted t: %s: %s\n", n.DataAtom.String(), n.Data)
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
