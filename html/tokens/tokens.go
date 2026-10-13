package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"os"
	"strings"
)

func main() {
	var htmlIn string

	if len(os.Args) == 1 {
		htmlIn = sampleHtml
	} else {
		// Reads the entire file content into memory
		content, _ := os.ReadFile(os.Args[1])
		htmlIn = string(content)
	}
	tokenizer := html.NewTokenizer(strings.NewReader(htmlIn))
	for {
		tokenType := tokenizer.Next()
		token := tokenizer.Token()
		if tokenType == html.ErrorToken {
			if tokenizer.Err() == io.EOF {
				return
			}
			fmt.Printf("Error: %v", tokenizer.Err())
			return
		}
		fmt.Printf("Token: %v\n", html.UnescapeString(token.String()))
	}
}

const sampleHtml = `<!DOCTYPE html><html><head><style> body {background-color: powderblue;} h1 {color: red;} p {color: orange;}</style><title>Sample HTML Code</title><script src="my-script.js">abc</script></head><body><h1>Main title</h1><p id="demo"></p><a href="https://medium.com/">Medium Blogs</a><script>document.getElementById("demo").innerHTML = "Hello JavaScript!";</script></body></html>`
