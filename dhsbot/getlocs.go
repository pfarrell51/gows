// pulls data from DNS lising locations
//
// this link will return JSON of some of the locations
// https://ttp.cbp.dhs.gov/schedulerapi/slots/asLocations?limit=1000

package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type Generated []struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	ShortName         string `json:"shortName"`
	LocationType      string `json:"locationType"`
	LocationCode      string `json:"locationCode"`
	Address           string `json:"address"`
	AddressAdditional string `json:"addressAdditional"`
	City              string `json:"city"`
	State             string `json:"state"`
	PostalCode        string `json:"postalCode"`
	CountryCode       string `json:"countryCode"`
	TzData            string `json:"tzData"`
}

func main() {
	// Get request
	var URL = `https://ttp.cbp.dhs.gov/schedulerapi/slots/asLocations?limit=100`
	resp, err := http.Get(URL)
	if err != nil {
		fmt.Println("No response from request")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body) // response body is []byte

	var result Generated
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to the go struct pointer
		fmt.Printf("Can not unmarshal JSON %v\n", err)
		return
	}
	cw := csv.NewWriter(os.Stdout)
	for _, r := range result {
		outline := make([]string, 5)
		outline[0] = strconv.Itoa(r.ID)
		outline[1] = r.Name
		outline[2] = r.Address
		outline[3] = r.City
		outline[4] = r.State
		cw.Write(outline)
	}
	cw.Flush()
}
