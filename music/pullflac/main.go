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

var hrefStart = regexp.MustCompile(`(\<a href=\")(.*?)("\>)`)

const hrefTarget = 2 // good link is the second group in above regex
var musicEnd = regexp.MustCompile(`.((M|m)(p|P)(3|4))|((F|f)(L|l)(A|a)(C|c))$`)

var hrefBackup = regexp.MustCompile(`^\?C=[DMN];`)
var hrefQuestionC = regexp.MustCompile(`^\?C=`)
var skipfiles = regexp.MustCompile(`^(/music/)|(/music/music/)(desktop.ini)$`)
var getURL string

func main() {
	// Check if an environment variable exists
	if val, ok := os.LookupEnv(getURLNAME); ok {
		fmt.Printf("found %s  %s:\n", getURLNAME, val)
		getURL = val
	} else {
		fmt.Printf("%s not found\n, getURLNAME")
	}
	if len(os.Args) != 1 {
		fmt.Printf("Usage: %s \ntarget URL is entered thru envoriment variables", os.Args[0])
		return
	}

	result, err := getURLstrings(getURL)
	if err != nil {
		fmt.Printf("Error or No response %v from request: %s\n", err, getURL)
		return
	}

	for i, row := range result {
		if i < 9 {
			continue
		}
		loc := hrefStart.FindStringSubmatch(row)
		if len(loc) > 0 {
			artistPath := loc[hrefTarget]
			if skipfiles.MatchString(artistPath) { // special case parent directory
				continue
			}
			backup := hrefBackup.FindString(artistPath)
			if backup != "" {
				continue
			}
			backup = hrefQuestionC.FindString(artistPath)
			if backup != "" {
				fmt.Printf("main bad, found ?= $s\n", i, backup)
				continue
			}

			getArtistDirectory(artistPath)
		}
		if i > 17 {
			break
		}
	}
}
func doGetURL(u string) ([]byte, error) {
	fmt.Printf("doing HTTP get %s\n", u)
	resp, err := http.Get(u)
	if err != nil {
		fmt.Printf("No response from request: %s\n", getURL)
	}
	body, err := ioutil.ReadAll(resp.Body) // response body is []byte
	resp.Body.Close()

	return body, err
}
func getURLstrings(u string) (s []string, e error) {
	body, err := doGetURL(u)
	if err != nil {
		fmt.Printf("No response from request: %s\n", u)
		return nil, err
	}
	result := strings.Split((string(body)), "\n")
	return result, err
}

func getArtistDirectory(p string) {
	fmt.Printf("looking at Artist directory %s\n", p)

	// Get request
	target := getURL + p

	result, err := getURLstrings(target)
	if err != nil {
		fmt.Printf("Error or No response %v from request: %s\n", err, target)
		return
	}
	//fmt.Println(result)

	for i, row := range result {
		if i < 8 {
			continue
		}
		loc := hrefStart.FindStringSubmatch(row)
		if len(loc) > 0 {
			path := loc[hrefTarget]
			if skipfiles.MatchString(path) { // special case parent directory
				continue
			}
			backup := hrefBackup.FindString(path)
			if backup != "" {
				continue
			}
			backup = hrefQuestionC.FindString(path)
			if backup != "" {
				fmt.Printf("gAD bad, found ?= $s\n", i, backup)
				continue
			}

			albumPath := target + path
			getAlbumDirectory(albumPath)
		}
		if i > 18 {
			break
		}
	}
}
func getAlbumDirectory(sd string) {
	fmt.Printf("\n\nreading album directory %s\n", sd)
	// Get request
	result, err := getURLstrings(sd)
	if err != nil {
		fmt.Printf("Error or No response %v from request: %s\n", err, sd)
		return
	}
	for i, row := range result {
		loc := hrefStart.FindStringSubmatch(row)
		if len(loc) > 0 {
			path := loc[hrefTarget]
			if skipfiles.MatchString(path) { // special case parent directory
				continue
			}
			backup := hrefBackup.FindString(path)
			if backup != "" {
				continue
			}
			backup = hrefQuestionC.FindString(path)
			if backup != "" {
				fmt.Printf("gAlbD bad, found ?= $s\n", i, backup)
				continue
			}
			res := musicEnd.FindString(path)
			if len(res) == 0 {
				continue
			}
			songPath := sd + path
			getSong(songPath)
		}
		if i > 18 {
			break
		}
	}
}
func getSong(sp string) {
	fmt.Printf("handling song %s\n", sp)
}
