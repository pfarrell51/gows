// general filter for gedcom files
package gedfilter

import (
	"bytes"
	"fmt"
	"os"

	"github.com/iand/gedcom"
)

type FlagST struct {
	Basic bool
	CSV   bool
	Debug bool
	Type  bool
}
type GlobalVars struct {
	pathArg        string
	localFlags     *FlagST
	songsProcessed int
}
type et struct {
	t string // tag of event
	v int    // trust of the date of this value
}

var knownEvents map[string]int

func init() {
	var knownEventVals = []et{
		et{"ADOP", 1},
		et{"BAPL", 8},  // baptism (LDS)
		et{"BAPM", 8},  // baptism
		et{"BARM", 7},  // bar Mitzvah
		et{"BASM", 7},  // bar/bat Mitzvah
		et{"BIRT", 10}, // birth, trust it
		et{"BURI", 18}, // burial
		et{"CENS", 0},
		et{"CHR", 7},  // christening date
		et{"CHRA", 2}, // adult christening
		et{"CREM", 17},
		et{"DEAT", 20}, // death, trust it
		et{"EMIG", 1},  // emigrate?
		et{"EVEN", 0},
		et{"GRAD", 6},
		et{"IMMI", 1}, // immigration
		et{"ORDL", 5}, // ordination
		et{"ORDN", 5}, // ordination
		et{"PROB", 0},
		et{"RETI", 0}, // retirement
		et{"WILL", 0},
	}
	knownEvents = make(map[string]int)
	for _, v := range knownEventVals {
		knownEvents[v.t] = v.v
	}
}
func AllocateData() *GlobalVars {
	rval := new(GlobalVars)
	rval.localFlags = new(FlagST)
	rval.localFlags = new(FlagST)
	if rval.localFlags == nil {
		fmt.Println("PIB in allocate Data, localflags is nil")
	}
	return rval
}
func (g *GlobalVars) Flags() *FlagST {
	return g.localFlags
}

func (g *GlobalVars) ProcessFile(f string) {
	g.pathArg = f
	data, err := os.ReadFile(f)
	d := gedcom.NewDecoder(bytes.NewReader(data))
	gd, err := d.Decode()

	if err != nil {
		fmt.Println(err)
		return
	}
	var rname, bdate, ddate string
	var bdT, ddT int
	for _, rec := range gd.Individual {
		if !g.Flags().Type {
			if len(rec.Name) > 0 {
				rname = rec.Name[0].Name
			} else {
				rname = "*missing*"
			}
		}
		if len(rec.Event) > 0 {
			for i := 0; i < len(rec.Event); i++ {
				var re *gedcom.EventRecord = rec.Event[i]
				t, found := knownEvents[re.Tag]
				if !found {
					panic(fmt.Sprintf("value not found for %s", re.Tag))
				}
				switch {
				case t < 10:
					if t > bdT {
						bdT = t
						bdate = re.Date
						if g.Flags().Debug {
							fmt.Printf(" $$ for %s using %s, %s wt: %d\n", rname, re.Tag, bdate, t)
						}
					}
				case t == 10:
					bdate = re.Date
					bdT = t
				case t > 10 && t < 20:
					if t > ddT {
						ddT = t
						ddate = re.Date
						if g.Flags().Debug {
							fmt.Printf(" && for %s using %s, %s with t: %d\n", rname, re.Tag, ddate, t)
						}
					}
				case t == 20:
					ddate = re.Date
					ddT = t
				}
				if g.Flags().Type {
					fmt.Printf("%s %s\n", re.Tag, re.Date)
				}
			}
			if g.Flags().Basic {
				fmt.Printf("%s, %s, %s\n", rname, bdate, ddate)
			}
			rname = "" // clear out fields for this person
			bdate = ""
			ddate = ""
			bdT = 0
			ddT = 0
		}
	}

}
