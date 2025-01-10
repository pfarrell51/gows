// package to pull flac encoded music from a website
// first it pulls the HTML page that constains the individual album urls


package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const getURLNAME = "PULLFLAC_URL"

var getURL = "http://mfarrelltuba.hopto.org/music/music/"

func main() {
	// Check if an environment variable exists
	if val, ok := os.LookupEnv(getURLNAME); ok {
		fmt.Println("PULLFLAC_URL:", val)
	} else {
		fmt.Println("PULLFLAC_URL not found")
	}
	if len(os.Args) != 1 {
		fmt.Printf("Usage: %s \ntarget URL is entered thru envoriment variables", os.Args[0])
		return
	}

	// Get request
	var URL = getURL
	resp, err := http.Get(URL)
	if err != nil {
		fmt.Println("No response from request")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body) // response body is []byte

	result := strings.Split((string(body)), "\n")
	//	fmt.Println(result)

	for i, r := range result {
		fmt.Printf("%d %s\n", i, r)
	}
}
