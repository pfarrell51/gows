// utilities (math) for latlong

package main

import (
	"fmt"
	"testing"
)

var (
	pat      = Latlong{40.064587, -75.37118}
	doug     = Latlong{35.151035, -81.0029}
	pittrace = Latlong{40.844715, -80.32963}
	lyric    = Latlong{37.230427, -80.41568}
)

func TestDistance(t *testing.T) {
	rval := pat.distance(doug)
	fmt.Printf("to doug %g\n ", rval)
	if rval < 500 || rval > 1000.0 {
		t.Errorf("doug %g", rval)
	}
	rval = pat.distance(pittrace)
	fmt.Printf("to Pitt %g\n", rval)
	if rval < 200.0 || rval > 300.0 {
		t.Errorf("pitts %g", rval)
	}
	rval = pat.distance(lyric)
	fmt.Printf("to Lyric %g\n", rval)
	if rval < 400.0 || rval > 500.0 {
		t.Errorf("Lyric %g", rval)
	}
}
