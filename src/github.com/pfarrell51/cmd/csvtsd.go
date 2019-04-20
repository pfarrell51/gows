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

type latlong struct {
	lat, lng float64
}

func main() {
	var oldTime time.Time
//	var oldlatlng latlong
//	var oldelev int32
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
		fmt.Printf("delta time: %d\n", deltaTime)

		oldTime = time
//		oldlatlng = latlng
//		oldelev = elev
	}
//	printCSV(oldTime, oldlatlng, oldelev, 0)
}

// parse out current line's values, converting to useful values
func getVals(line []string) (rtime time.Time, rlatlng latlong, rele int32) {
	var rlatlong = latlong{0.0, 0.0}
	rtime, err := time.Parse(time.RFC3339, strings.Trim(line[0], " "))
	if err != nil {
		log.Fatal(err)
		return rtime, rlatlng, 0
	}
	rlatlong.lat, err = strconv.ParseFloat(strings.Trim(line[1], " 	"), 64)
	rlatlong.lng, err = strconv.ParseFloat(strings.Trim(line[2], " 	"), 64)
	relev, err := strconv.ParseInt(strings.Trim(line[3], " 	"), 0, 32)
	relev *= meterTo1kFeet
	relev /= 1000 // remove 1000 back to feet
	fmt.Printf("%v %g %g %d\n", rtime, rlatlong.lat, rlatlong.lng, relev)
	return rtime, rlatlong, rele
}
func printCSV(tval time.Time, latlng latlong, ele, spd int32) {
	fmt.Printf("%v, %g, %g, %d, %d\n", tval, latlng.lat, latlng.lng, ele, spd)
}

// distance using law of cosines
//    d = acos( sin φ1 * sin φ2 + cos φ1 * cos φ2 * cos Δλ ) * R
// note, lat long incoming is in degrees, we want radians
func (l1 latlong) distance(l2 latlong) float64 {
	lr1 := latlong{l1.lat * (math.Pi / 180.0), l1.lng * (math.Pi / 180.0)}
	lr2 := latlong{l2.lat * (math.Pi / 180.0), l2.lng * (math.Pi / 180.0)}
	fmt.Printf("1 deg (%g, %g) rads (%g, %g)\n", l1.lat, l1.lng, lr1.lat, lr1.lng)
	fmt.Printf("2 deg (%g, %g) rads (%g, %g)\n", l2.lat, l2.lng, lr2.lat, lr2.lng)
	deltaLong := (lr1.lng - lr2.lng)
	deltaLat := math.Abs(lr1.lat - lr2.lat)
	fmt.Printf("d lat: %g  dLong: %g\n", deltaLat, deltaLong)
	partial := math.Sin(lr1.lat)*math.Sin(lr2.lat) + math.Cos(lr1.lat)*math.Cos(lr2.lat)*math.Cos(deltaLong)
	fmt.Printf("partial: %g aCos(P) %g\n", partial, math.Acos(partial))
	var rval = math.Acos(partial) * earthRadiusMiles
	return rval
}
