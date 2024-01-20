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
	t string
	v int
}

var knownEvents map[string]int

func init() {
	var knownEventVals = []et{
		et{"ADOP", 0},
		et{"BAPM", 0},
		et{"BARM", 0},
		et{"BASM", 0},
		et{"BIRT", 10},
		et{"BURI", 0},
		et{"CENS", 0},
		et{"CHR", 0},
		et{"CHRA", 0},
		et{"CREM", 0},
		et{"DEAT", 20},
		et{"EMIG", 0}, // emigrate?
		et{"EVEN", 0},
		et{"GRAD", 0},
		et{"IMMI", 0},
		et{"ORDN", 0},
		et{"PROB", 0},
		et{"RETI", 0},
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
	rval.localFlags.Basic = true
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
	for _, rec := range gd.Individual {
		if !g.Flags().Type {
			//fmt.Println(reflect.TypeOf(rec))
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
				switch t {
				case 10:
					bdate = re.Date
				case 20:
					ddate = re.Date
				}
				if g.Flags().Type {
					fmt.Printf("%s %s\n", re.Tag, re.Date)
				}
			}
			if g.Flags().Basic {
				fmt.Printf("%s, %s, %s\n", rname, bdate, ddate)
			}
			rname = ""
			bdate = ""
			ddate = ""
		}
	}

}
