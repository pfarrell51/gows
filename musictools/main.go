package main

import (
	"fmt"
	"github.com/dlclark/metaphone3"
	g "github.com/zyedidia/generic"
	"github.com/zyedidia/generic/btree"
	"io/fs"
	"os"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

var tree = btree.New[string, string](g.Less[string])
var enc metaphone3.Encoder
var extRegex = regexp.MustCompile(".((M|m)(p|P)3)|((M|m)(p|P)4)|((F|f)(L|l)(A|a)(C|c))")

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s DIRNAME", os.Args[0])
		os.Exit(1)
	}
	pathArg := os.Args[1]
	ProcessFiles(pathArg)
}

func ProcessFiles(pathArg string) {
	loadMetaPhone()
	rmap := walkFiles(pathArg)
	processMap(rmap)
}
func loadMetaPhone() {
	albumNames := [...]string{"ABBA", "Alison_Krauss", "AllmanBrothers", "Almanac_Singers", "Animals",
		"Arlo_Guthrie", "Band_The", "Basia", "BeachBoys", "Beatles", "BlindFaith", "BloodSweatTears", "Boston",
		"BrewerAndShipley", "BuffaloSpringfield", "Byrds", "CensorBeep.mp4", "Chesapeake",
		"Cream", "Crosby_Stills_Nash",
		"David_Bromberg", "Derek_Dominos", "Dire_Straits", "Doobie_Brothers", "Doors", "Dylan", "Elton_John",
		"Emmylou_Harris", "Fleetwood_Mac", "Heart", "James_Taylor", "Jefferson_Airplane", "Jethro_Tull",
		"John_Denver", "John_Hartford", "John_Starling", "Joni_Mitchell", "Judy_Collins", "Kingston_Trio",
		"Led_Zepplin", "Linda_Ronstadt", "Lynyrd_Skynyrd", "Mamas_Popas", "Meatloaf", "Mike_Auldridge",
		"New_Riders_Purple_Sage", "Pablo_Cruise", "Paul_Simon", "Peter_Paul_Mary", "Rolling_Stones",
		"Roy_Orbison", "Santana", "Seals_Croft", "Seldom_Scene", "Simon_Garfunkel", "Steely_Dan",
		"5th_Dimension", "TonyRice", "Traveling_Wilburys", "Who", "Yes",
	}
	for _, n := range albumNames {
		prim, sec := enc.Encode(n)
		tree.Put(prim, n)
		if len(sec) > 0 {
			fmt.Printf("found a secondary: %s for %s\n", sec, n)
			tree.Put(sec, n)
		}
	}
}

// walk all files, looking for nice GoPro created video files.
// fill in a map keyed by the desired new name order
func walkFiles(pathArg string) map[string]string {
	theMap := make(map[string]string)
	fsys := os.DirFS(pathArg)
	fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Println("Error processing", p, " in ", d)
			fmt.Println("error is ", err)
			return nil
		}
		if d == nil {
			fmt.Println("d is nil")
			return nil
		}
		dName := d.Name()
		if strings.HasPrefix(p, ".") {
			return nil
		}
		fmt.Println(p, d, dName)
		ext := extRegex.FindString(p)
		if len(ext) == 0 {
			fmt.Println("no extension for ", p)
			return nil
		}

		prim, sec := enc.Encode(dName)
		_, ok := tree.Get(prim)
		if !ok {
			fmt.Printf("Song %s did not find primary match for %s\n", dName, prim)
			_, ok := tree.Get(sec)
			if !ok {
				fmt.Println("no match for either primary or secondary ", dName)
				group := findGroup(dName)
				if len(group) == 0 {
					fmt.Println("no group found for song ", dName)
					return nil
				}
				fmt.Printf("found group %s for song %s\n", group, dName)
			}
		}
		return nil
	})
	return theMap
}
func findGroup(s string) string {
	nameRegex := regexp.MustCompile("\\S?(\\w*)\\S")
	group := nameRegex.FindString(s)
	fmt.Printf(">%s<\n", group)
	prim, _ := enc.Encode(group)
	group, ok := tree.Get(prim)
	if !ok {
		fmt.Println("very bad")
		return ""
	}
	return group
}

// go thru the map, sort by key
// then create new ordering that makes sense to human
func processMap(m map[string]string) map[string]string {
	var cmd = "mv "
	if runtime.GOOS == "windows" {
		cmd = "ren "
	}
	keys := make([]string, 0, len(m)) // copy to a slice to sort
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var firstChpt string
	var cnum int
	for _, k := range keys {
		if cnum < 1 {
			firstChpt = k[0:6]
			c, _ := strconv.Atoi(k[6:8])
			cnum = c
		} else {
			cnum++
		}
		source, _ := m[k]
		delete(m, k)
		m[k] = source
		fmt.Printf("%s%s %2s%02d%4s.mp4\n", cmd, source, firstChpt[0:2], cnum, firstChpt[2:6])
	}
	return m
}
