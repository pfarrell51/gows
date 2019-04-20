// calculate speed and bearing from the simple lat/long CSV file
// extracts timestamp, lat, long, elevation and speed
// reads stdin

package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

const meterTo1kFeet = 3280
const earthRadiusMiles = 3959

type Latlong struct {
	lat, lng float64
}

func main() {
	var oldTime time.Time
	var oldlatlng Latlong
	var oldelev int32
	fmt.Printf("#time, lat, long, ele, spd, bearing\n")
	ln := 0
	r := csv.NewReader(os.Stdin)
	for {
		record, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
		}
		ln++
		time, latlng, elev := getVals(record)
		printCSV(time, latlng, elev, 0)
		// do stuff
		deltaTime := time.Sub(oldTime)
		fmt.Println(deltaTime)

		oldTime = time
		oldlatlng = latlng
		oldelev = elev
	}
	printCSV(oldTime, oldlatlng, oldelev, 0)
}

// parse out current line's values, converting to useful values
func getVals(line []string) (rtime time.Time, rlatlng Latlong, rele int32) {
	var rlatlong = Latlong{0.0, 0.0}
	rtime, err := time.Parse(time.RFC3339, strings.Trim(line[0], " "))
	if err != nil {
		log.Fatal(err)
		return rtime, rlatlng, 0
	}
	rlatlong.lat, err = strconv.ParseFloat(strings.Trim(line[1], " 	"), 64)
	rlatlong.lat *= (math.Pi / 180.0)
	rlatlong.lng, err = strconv.ParseFloat(strings.Trim(line[2], " 	"), 64)
	rlatlong.lng *= (math.Pi / 180.0)
	relev, err := strconv.ParseInt(strings.Trim(line[3], " 	"), 0, 32)
	relev *= meterTo1kFeet
	relev /= 1000 // remove 1000 back to feet
	fmt.Printf("%v %g %g %d\n", rtime, rlatlong.lat, rlatlong.lng, relev)
	return rtime, rlatlong, rele
}
func printCSV(tval time.Time, latlng Latlong, ele, spd int32) {
	fmt.Printf("%v, %g, %g, %d, %d\n", tval, latlng.lat, latlng.lng, ele, spd)
}

// distance using law of cosines
//    d = acos( sin φ1 * sin φ2 + cos φ1 * cos φ2 * cos Δλ ) * R
func (l1 Latlong) distance(l2 Latlong) float64 {
	var deltaLong = l1.lng - l2.lng
	deltaLat := math.Abs(l1.lat - l2.lat)
	fmt.Printf("d lat: %g  dLong: %g\n", deltaLat, deltaLong)
	partial := math.Sin(l1.lat)*math.Sin(l2.lat) + math.Cos(l1.lat)*math.Cos(l2.lat)*math.Cos(deltaLong)
	fmt.Printf("partial: %g aCos(P) %g\n", partial, math.Acos(partial))
	var rval = math.Acos(partial) * earthRadiusMiles
	return rval
}
