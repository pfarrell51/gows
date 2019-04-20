// utilities (math) for latlong

package main

import (
	"fmt"
	"testing"
)

var (
	pat      = Latlong{0.6992589566, -1.315475252} // 40.064587, -75.37118
	pocono   = Latlong{0.7167193352, -1.315476771} // 41.064993, -75.371267}
	doug     = Latlong{0.6135012962, -1.413767309} //35.151035, -81.0029}
	pittrace = Latlong{0.7128747588, -1.40201653}  // 40.844715, -80.32963}
	lyric    = Latlong{0.6497935331, -1.403518386} //37.230427, -80.41568}
)

func TestDistance(t *testing.T) {
	rval := pat.distance(pocono)
	fmt.Printf("to pocono %g\n ", rval)
	if rval < 60 || rval > 70.0 {
		t.Errorf("pocono %g", rval)
	}
	rval = pat.distance(doug)
	fmt.Printf("to doug %g\n ", rval)
	if rval < 450 || rval > 500.0 {
		t.Errorf("doug %g", rval)
	}
	rval = pat.distance(pittrace)
	fmt.Printf("to Pitt %g\n", rval)
	if rval < 200.0 || rval > 300.0 {
		t.Errorf("pitts %g", rval)
	}
	rval = pat.distance(lyric)
	fmt.Printf("to Lyric %g\n", rval)
	if rval < 300.0 || rval > 400.0 {
		t.Errorf("Lyric %g", rval)
	}
}
