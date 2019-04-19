// calculate speed and bearing from the simple lat/long CSV file
// extracts timestamp, lat, long, elevation and speed
// reads stdin

package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	fmt.Printf("#time, lat, long, ele, spd, bearing\n")
	ln := 0
	r := csv.NewReader(os.Stdin)
	for {
		record, err := r.Read()
		if err != nil {
			log.Fatal(err)
		}
		ln++
		timeStr := strings.Trim(record[0], " ")
		latStr := strings.Trim(record[1], " 	")
		lngStr := strings.Trim(record[2], " 	")
		eleStr := strings.Trim(record[3], " 	")
		fmt.Printf("%d %d %s\n", ln, len(record), record)
		time, err := time.Parse(time.RFC3339, timeStr)
		var lat float64
		if lat, err = strconv.ParseFloat(latStr, 64); err == nil {
			fmt.Printf("str: %s val: %g\n", latStr, lat)
		} else {
			fmt.Println(err)
		}
		lng, err := strconv.ParseFloat(lngStr, 64)
		elev, err := strconv.ParseInt(eleStr, 0, 32)
		fmt.Printf("%v %g %g %d\n", time, lat, lng, elev)
	}

}

func printCSV(tval, lat, lng, ele, spd string) {
	fmt.Printf("%s, %s, %s, %s, %s\n", tval, lat, lng, ele, spd)
}
