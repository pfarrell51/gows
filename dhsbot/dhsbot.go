// package to call DHS's trusted traveler appointment website
// a good GET is  https://ttp.cbp.dhs.gov/schedulerapi/slots?orderBy=soonest&limit=100&locationId=5445&minimum=1
//
// this link will return JSON of all the locations
// https://ttp.cbp.dhs.gov/schedulerapi/slots/asLocations?limit=100

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Response struct {
	LocationID     int    `json:"locationId"`
	StartTimestamp string `json:"startTimestamp"`
	EndTimestamp   string `json:"endTimestamp"`
	Active         bool   `json:"active"`
	Duration       int    `json:"duration"`
	RemoteInd      bool   `json:"remoteInd"`
}

const getURL = "https://ttp.cbp.dhs.gov/schedulerapi/slots?orderBy=soonest&limit=1&locationId=5004&minimum=1"

const dhsResponse = `{
  "locationId" : 5004,
  "startTimestamp" : "2023-07-14T09:45",
  "endTimestamp" : "2023-07-14T10:00",
  "active" : true,
  "duration" : 15,
  "remoteInd" : false
} 
`

func main() {
	var body []byte
	if true {
		// Get request
		resp, err := http.Get(getURL)
		if err != nil {
			fmt.Println("No response from request")
		}
		defer resp.Body.Close()
		body, err = ioutil.ReadAll(resp.Body) // response body is []byte
	} else {
		body = []byte(dhsResponse)
	}
	fmt.Println(string(body))
	var result Response
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to the go struct pointer
		fmt.Printf("Can not unmarshal JSON %v\n", err)
	}
	fmt.Println(result)
}
