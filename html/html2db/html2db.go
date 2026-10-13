package main

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/net/html"
)

func main() {
	filePath := "../index.html"

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer file.Close()

	// Parse the HTML content from the file (which implements io.Reader)
	doc, err := html.Parse(file)
	if err != nil {
		log.Fatalf("failed to parse html: %v", err)
	}

	// Traverse the DOM tree and process nodes
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode {
			fmt.Printf("Element tag: <%s>\n", n.Data)
		} else if n.Type == html.TextNode {
			// Print text contents if not empty/whitespace
			if text := n.Data; text != "" {
				fmt.Printf("Text: %s\n", text)
			}
		}

		// Recursively process child nodes
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(doc)
}
