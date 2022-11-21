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
const makeRadians = math.Pi / 180.0


type latlong struct {
	lat, lng float64
}

func main() {
	var oldTime time.Time
	var oldlatlng latlong
	var oldelev, speed int32
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

		// do stuff
		deltaTime := time.Sub(oldTime).Seconds()
		dist := latlng.distance(oldlatlng)
		speed = int32(dist / (deltaTime / 3600.0))
		printCSV(oldTime, oldlatlng, oldelev, speed)
		oldTime = time
		oldlatlng = latlng
		oldelev = elev
	}
	printCSV(oldTime, oldlatlng, oldelev, speed)
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
	relev = (relev * meterTo1kFeet) / 1000
	return rtime, rlatlong, int32(relev)
}
func printCSV(tval time.Time, latlng latlong, ele, spd int32) {
	fmt.Printf("%v, %g, %g, %d, %d\n", tval, latlng.lat, latlng.lng, ele, spd)
}

// distance using law of cosines
//    d = acos( sin φ1 * sin φ2 + cos φ1 * cos φ2 * cos Δλ ) * R
// note, lat long incoming is in degrees, we want radians
func (l1 latlong) distance(l2 latlong) float64 {
	lr1 := latlong{l1.lat * makeRadians, l1.lng * makeRadians}		// convert degrees to radians
	lr2 := latlong{l2.lat * makeRadians, l2.lng * makeRadians}
	partial := math.Sin(lr1.lat)*math.Sin(lr2.lat) + math.Cos(lr1.lat)*math.Cos(lr2.lat)*math.Cos(lr1.lng-lr2.lng)
	var rval = math.Acos(partial) * earthRadiusMiles
	return rval
}
