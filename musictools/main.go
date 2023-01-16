// this program helps clean up the mp3 and flac files in my hits playlist
// its main task is to normalize the file names to relect the artist and song title

package main

import (
	"bytes"
	"fmt"
	"github.com/dlclark/metaphone3"
	g "github.com/zyedidia/generic"
	"github.com/zyedidia/generic/btree"
	"io/fs"
	"os"
	"path"
	"regexp"
	"strings"
	"unicode"
)

var tree = btree.New[string, string](g.Less[string])
var enc metaphone3.Encoder
var extRegex = regexp.MustCompile(".((M|m)(p|P)(3|4))|((F|f)(L|l)(A|a)(C|c))")

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s DIRNAME", os.Args[0])
		os.Exit(1)
	}
	loadMetaPhone()
	pathArg := path.Clean(os.Args[1])
	ProcessFiles(pathArg)
}

func ProcessFiles(pathArg string) {
	rmap := walkFiles(pathArg)
	processMap(rmap)
}
func justLetter(a string) string {
	buff := bytes.Buffer{}
	for _, c := range a {
		if unicode.IsLetter(c) {
			buff.WriteRune(c)
		} else if c == '_' {
			// ignore it
		} else if c == '-' {
			break
		}
	}
	return buff.String()
}
func loadMetaPhone() {
	groupNames := [...]string{"ABBA", "Alison_Krauss", "AllmanBrothers", "Almanac_Singers", "Animals",
		"Arlo_Guthrie", "Band", "Basia", "BeachBoys", "Beatles", "BlindFaith", "BloodSweatTears", "Boston",
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
	for _, n := range groupNames {
		prim, sec := enc.Encode(justLetter(n))
		tree.Put(prim, n)
		if len(sec) > 0 {
			tree.Put(sec, n)
		}
	}
}

func processFile(pathArg string, fsys fs.FS, root string, d fs.DirEntry, err error) error {
	fmt.Printf("pa: %s  p: %s, d: %v, \n", pathArg, fsys, root)
	return nil
}

// walk all files, looking for nice GoPro created video files.
// fill in a map keyed by the desired new name order
func walkFiles(pathArg string) map[string]string {
	theMap := make(map[string]string)
	fsys := os.DirFS(pathArg)
	fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		err = processFile(pathArg, fsys, p, d, err)
		if err != nil {
			fmt.Println("Error processing", p, " in ", d)
			fmt.Println("error is ", err)
			return nil
		}
		if d == nil {
			fmt.Println("d is nil")
			return nil
		}
		if strings.HasPrefix(p, ".") {
			return nil // hidden, not interesting
		}

		ext := extRegex.FindString(p)
		if len(ext) == 0 {
			return nil // not interesting extension
		}
		var addThe, cmd, group string
		ps, pn := path.Split(p)
		if strings.HasPrefix(pn, "The ") {
			addThe = "The "
			pn = pn[4:]
		}
		prim, sec := enc.Encode(justLetter(pn))
		_, ok := tree.Get(prim)
		if ok {
			group = pn
			theMap[prim] = pn
			cmd = fmt.Sprintf("1mv '%s/%s'\n  -> '%s/%s%s%s", pathArg, p,
				pathArg, ps, addThe, group)
		} else {
			_, ok = tree.Get(sec)
			if !ok {
				group = findGroup(pn)
				if len(group) == 0 {
					fmt.Println("no group found for song ", p)
					return nil
				}
				cmd = fmt.Sprintf("2mv '%s/%s'\n  -> '%s/+%s%s%s: - %s", pathArg, p,
					pathArg, ps, addThe, group, pn)
			}
		}
		if len(cmd) > 0 {
			fmt.Println(cmd)
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
	return m
}
