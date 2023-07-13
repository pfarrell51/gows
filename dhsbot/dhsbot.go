// package to call DHS's trusted traveler appointment website
// a good GET is  https://ttp.cbp.dhs.gov/schedulerapi/slots?orderBy=soonest&limit=100&locationId=5445&minimum=1
//
// this link will return JSON of some of the locations
// https://ttp.cbp.dhs.gov/schedulerapi/slots/asLocations?limit=1000

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Response []struct {
	LocationID     int    `json:"locationId"`
	StartTimestamp string `json:"startTimestamp"`
	EndTimestamp   string `json:"endTimestamp"`
	Active         bool   `json:"active"`
	Duration       int    `json:"duration"`
	RemoteInd      bool   `json:"remoteInd"`
}

const getURL = "https://ttp.cbp.dhs.gov/schedulerapi/slots?orderBy=soonest&limit=1&minimum=1&locationId="

func main() {
	if len(os.Args) != 2 {
		fmt.Println("must specify locationId as argument")
		fmt.Println("PHL = 5445, Laredo = 5004")
		return
	}
	// Get request
	var URL = getURL + os.Args[1]
	resp, err := http.Get(URL)
	if err != nil {
		fmt.Println("No response from request")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body) // response body is []byte

	fmt.Println(string(body))
	var result Response
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to the go struct pointer
		fmt.Printf("Can not unmarshal JSON %v\n", err)
	}
	fmt.Println(result)
}
