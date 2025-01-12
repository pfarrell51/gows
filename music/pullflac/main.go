// package to pull flac encoded music from a website
// first it pulls the HTML page that constains the individual album urls

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	//	"github.com/pfarrell51/gows/music/replaceUnicode"
)

const getURLNAME = "PULLFLAC_URL"

var hrefStart = regexp.MustCompile(`(\<a href=\")(.*?)("\>)`)

const hrefTarget = 2 // good link is the second group in above regex
var musicEnd = regexp.MustCompile(`.((M|m)(p|P)(3|4))|((F|f)(L|l)(A|a)(C|c))$`)

var hrefBackup = regexp.MustCompile(`^\?C=[DMN];`) // mike's javascript creates these
var hrefQuestionC = regexp.MustCompile(`^\?C=`)    // catch any unknown cases
var skipfiles = regexp.MustCompile(`^(/music/)|(/music/music/)|(desktop.ini)|(.*\.jpg)$`)
var debugFlag bool
var helpFlag bool

type filenames struct {
	baseURL       string
	outputPath    string
	artistInPath  string
	artistOutPath string
	albumInPath   string
	albumOutPath  string
	songInPath    string
	songOutPath   string
}

var fnames = new(filenames)

func main() {
	// Check if an environment variable exists
	if val, ok := os.LookupEnv(getURLNAME); ok {
		if debugFlag {
			fmt.Printf("found %s  %s:\n", getURLNAME, val)
		}
		fnames.baseURL = val
	} else {
		slog.Error("%s not found\n", getURLNAME)
	}
	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "Usage of %s: [flags]  [dest]\n", os.Args[0])
		fmt.Fprintf(w, " if no destination directory given, will use .\n")
		flag.PrintDefaults()
	}
	flag.BoolVar(&debugFlag, "d", false, "debug on")
	flag.BoolVar(&helpFlag, "h", false, "print help")

	flag.Parse()

	if helpFlag || (len(flag.Args()) != 1 && len(flag.Args()) != 0) {
		flag.Usage()
		return
	}
	if len(flag.Args()) == 1 {
		fnames.outputPath = flag.Arg(0)
	} else {
		fnames.outputPath, _ = os.Getwd()
	}
	if stat, err := os.Stat(fnames.outputPath); err == nil && stat.IsDir() {
		if debugFlag {
			fmt.Println(stat)
		}
	} else {
		slog.Error("No such directory as %s\n", fnames.outputPath)
		return
	}
	getTopIndex()
}
func getTopIndex() {
	result, err := getURLstrings(fnames.baseURL)
	if err != nil {
		slog.Error("Error or No response %v from request: %s\n", err, fnames.baseURL)
		return
	}

	for i, row := range result {
		if i < 9 {
			continue
		}
		loc := hrefStart.FindStringSubmatch(row)
		if len(loc) > 0 {
			fnames.artistInPath = loc[hrefTarget]
			fnames.artistOutPath = cleanupWriteDirectory(fnames.artistInPath)
			if skipfiles.MatchString(fnames.artistInPath) { // special case parent directory
				continue
			}
			backup := hrefBackup.FindString(fnames.artistInPath)
			if backup != "" {
				continue
			}
			backup = hrefQuestionC.FindString(fnames.artistInPath)
			if backup != "" {
				slog.Error("getTopIdx bad, found ?= $s\n", i, backup)
				continue
			}
			getArtistDirectory(fnames.artistInPath)
		}
		if i > 12 {
			break
		}
	}
}
func cleanupWriteDirectory(s string) string {
	decodedString, err := url.QueryUnescape(s)
	if err != nil {
		fmt.Println("Error decoding:", err)
	} else if debugFlag {
		slog.Info("Decoded string:", decodedString)
	}
	//var changed bool
	rval := decodedString //replaceUnicode.doReplacement(decodedString, &changed)
	return rval
}
func doGetURL(u string) ([]byte, error) {
	fmt.Printf("doing HTTP get %s\n", u)
	resp, err := http.Get(u)
	if err != nil {
		fmt.Printf("No response from request: %s\n", fnames.baseURL)
	}
	body, err := ioutil.ReadAll(resp.Body) // response body is []byte
	resp.Body.Close()

	return body, err
}
func getURLstrings(u string) (s []string, e error) {
	body, err := doGetURL(u)
	if err != nil {
		fmt.Printf("error or No response (%v) from request: %s\n", err, u)
		return nil, err
	}
	result := strings.Split((string(body)), "\n")
	return result, err
}

func getArtistDirectory(p string) {
	fmt.Printf("looking at Artist directory %s\n", p)

	if debugFlag {
		fmt.Printf("Art in: %s\nout: %s\n", fnames.artistInPath, fnames.artistOutPath)
	}
	// Get request
	target := fnames.baseURL + p

	result, err := getURLstrings(target)
	if err != nil {
		fmt.Printf("Error or No response %v from request: %s\n", err, target)
		return
	}

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

			fnames.albumOutPath = cleanupWriteDirectory(path)
			albumPath := target + path
			fnames.albumInPath = albumPath
			getAlbumDirectory(albumPath)
		}
		if i > 16 {
			break
		}
	}
}
func getAlbumDirectory(sd string) {
	fmt.Printf("\n\nreading album directory %s\n", sd)
	if debugFlag {
		fmt.Printf("Art in: %s\nout: %s\n", fnames.artistInPath, fnames.artistOutPath)
		fmt.Printf("Alb in: %s\nout: %s\n", fnames.albumInPath, fnames.albumOutPath)
	}
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
			fnames.songOutPath = cleanupWriteDirectory(path)
			songPath := sd + path
			fnames.songInPath = songPath
			getSong(songPath)
		}
		if i > 18 {
			break
		}
	}
}
func getSong(sp string) {
	fmt.Printf("handling song %s\n", sp)
	if debugFlag {
		fmt.Printf("Art in: %s\nout: %s\n", fnames.artistInPath, fnames.artistOutPath)
		fmt.Printf("Alb in: %s\nout: %s\n", fnames.albumInPath, fnames.albumOutPath)
		fmt.Printf("song in: %s\nout: %s\n", fnames.songInPath, fnames.songOutPath)
	}
}
