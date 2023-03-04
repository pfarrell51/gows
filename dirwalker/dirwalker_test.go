package dirwalker

import (
	"testing"
)

func TestArePathsParallel(t *testing.T) {
	type Trip struct {
		ip  string
		op  string
		exp bool
	}
	tp := []Trip{
		Trip{"/home/pfarrell/whome/Music/flac/Animals", "/home/pfarrell/whome/flac/Music/", false},
		Trip{"/home/pfarrell/whome/Music/flac/Animals", "/home/pfarrell/whome/Music/flac/Animals", true},
		Trip{"/home/pfarrell/whome/Music/flac/Animals", "/home/pfarrell/whome/flac/Music/Animals", false},
		Trip{"/home/pfarrell/whome/Music/flac/Animals", "/home/pfarrell/whome/Music/mp3/Animals", true},
		Trip{"/home/pfarrell/whome/Music/flac/Animals/Best of The60s", "/home/pfarrell/whome/Music/mp3/Animals", false},
		Trip{"/home/pfarrell/whome/Music/flac/Animals", "/home/pfarrell/whome/Music/mp3/Animals/Best Of The60s/Roadrunner.mp3", false},
	}
	for i, s := range tp {
		rval := arePathsParallel(s.ip, s.op)
		if rval != s.exp {
			t.Errorf("%d incorrect for %v", i, s)
		}
	}
}
