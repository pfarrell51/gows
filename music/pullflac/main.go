// package to pull flac encoded music from a website
// first it pulls the HTML page that constains the individual album urls

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const getURLNAME = "PULLFLAC_URL"

var getURL string

func main() {
	// Check if an environment variable exists
	if val, ok := os.LookupEnv(getURLNAME); ok {
		fmt.Println("PULLFLAC_URL:", val)
		getURL = val
	} else {
		fmt.Println("PULLFLAC_URL not found")
	}
	if len(os.Args) != 1 {
		fmt.Printf("Usage: %s \ntarget URL is entered thru envoriment variables", os.Args[0])
		return
	}

	// Get request
	resp, err := http.Get(getURL)
	if err != nil {
		fmt.Println("No response from request")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body) // response body is []byte

	result := strings.Split((string(body)), "\n")

	var hrefStart = regexp.MustCompile(`\<a href="`)
	//var hrefEnd = regexp.MustCompile(`\"\>`)
	for i, row := range result {
		if i < 9 {
			continue
		}
		fmt.Printf("%d %s\n", i, row)
		//560 <tr><td valign="top"><img src="/icons/folder.gif" alt="[DIR]"></td><td><a href="Youngbloods/">Youngbloods/</a></td><td align="right">2024-09-26 18:28  </td><td align="right">  - </td><td>&nbsp;</td></tr>
		loc := hrefStart.FindStringIndex(row)
		if len(loc) > 0 {
			fmt.Printf("loc %v\n", loc)
			//loce := hrefEnd.FindStringIndex(row[loc[0])
			//fmt.Printf("loce %v\n", loce)
		}
		if i > 15 {
			break
		}
	}
}
