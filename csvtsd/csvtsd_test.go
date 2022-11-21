// utilities (math) for latlong

package main

import (
	"fmt"
	"testing"
)

var (
	pat      = latlong{40.064587, -75.37118}
	pocono   = latlong{41.064993, -75.371267}
	doug     = latlong{35.151035, -81.0029}
	pittrace = latlong{40.844715, -80.32963}
	lyric    = latlong{37.230427, -80.41568}
)

func TestDistance(t *testing.T) {
	fmt.Println("home to Pocono")
	rval := pat.distance(pocono)
	fmt.Printf("to pocono %g\n ", rval)
	if rval < 60 || rval > 70.0 {
		t.Errorf("pocono %g", rval)
	}

	fmt.Println("home to Doug")
	rval = pat.distance(doug)
	fmt.Printf("to doug %g\n ", rval)
	if rval < 450 || rval > 500.0 {
		t.Errorf("doug %g", rval)
	}
	fmt.Println("home to Pittrace")
	rval = pat.distance(pittrace)
	fmt.Printf("to Pitt %g\n", rval)
	if rval < 200.0 || rval > 300.0 {
		t.Errorf("pitts %g", rval)
	}
	fmt.Println("home to Lyric")
	rval = pat.distance(lyric)
	fmt.Printf("to Lyric %g\n", rval)
	if rval < 300.0 || rval > 400.0 {
		t.Errorf("Lyric %g", rval)
	}
	fmt.Println("Doug to Lyric")
	rval = doug.distance(lyric)
	fmt.Printf("doug to Lyric %g\n", rval)
	if rval < 140.0 || rval > 200.0 {
		t.Errorf("Doug Lyric %g", rval)
	}
}
func TestGetVals(t *testing.T) {
	fmt.Println("testing get vals/parsing")
	record := []string{"2019-04-14T15:08:11.917Z", " 40.8505242", " -80.3463992", "334"}
	time, latlng, elev := getVals(record)
	printCSV(time, latlng, elev, 0)
	if latlng.lat < 40.0 || latlng.lat > 41.0 {
		t.Errorf("lat %g", latlng.lat)
	}
	if latlng.lng < -81.0 || latlng.lng > -80.0 {
		t.Errorf("long: %g", latlng.lng)
	}
}
