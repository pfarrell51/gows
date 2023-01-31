// this program helps clean up the mp3 and flac files in my hits playlist
// its main task is to normalize the file names to relect the artist and Song title
//
// this is not multi-processing safe

// Bugs

package musictools

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path"
	"regexp"
	"runtime"
	"strings"
	"unicode"

	"github.com/dhowden/tag"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func init() {
	GetEncoder().MaxLength = 24
	LoadArtistMap()
}

func ProcessFiles(pathArg string) {
	if !GetFlags().ZDumpArtist {
		rmap := WalkFiles(pathArg)
		ProcessMap(pathArg, rmap)
	} else {
		DumpGptree()
	}
}

var regAnd = regexp.MustCompile("(?i) (and|the) ")

func JustLetter(a string) string {
	buff := bytes.Buffer{}
	loc := []int{0, 0}
	for j := 0; j < 4; j++ { // 4 allows to and and two the, but it will nearly alwaysbreak before that
		loc = regAnd.FindStringIndex(a)
		if len(loc) < 1 {
			break
		}
		a = a[:loc[0]] + a[loc[1]-1:]
	}
	for _, c := range a {
		if unicode.IsLetter(c) {
			buff.WriteRune(c)
		} else if c == '_' || c == '&' || unicode.IsSpace(c) {
			// ignore it
		} else if c == '-' {
			break
		}
	}
	return buff.String()
}

var sortKeyExp = regexp.MustCompile("^[A-Z](-|_)")
var extRegex = regexp.MustCompile("((M|m)(p|P)(3|4))|((F|f)(L|l)(A|a)(C|c))$")
var underToSpace = regexp.MustCompile("_")

var cReg = regexp.MustCompile(",\\s")
var dReg = regexp.MustCompile("-\\s")
var commaExp = regexp.MustCompile(",\\s")

// parse the file info to find artist and Song title
// most of my music files have file names with the artist name, a hyphen and then the track title
// so this pulls out the information and fills in the "Song" object.
func parseFilename(pathArg, p string) *Song {
	// fmt.Printf("sf: %s\n", p)
	var rval = new(Song)
	rval.inPath = path.Join(pathArg, p)
	rval.outPath = pathArg
	nameB := []byte(strings.TrimSpace(p))
	if sortKeyExp.Match(nameB) {
		nameB = nameB[2:]
	}
	nameB = underToSpace.ReplaceAll(nameB, []byte(" "))
	extR := extRegex.FindIndex(nameB)
	if extR == nil || len(extR) == 0 {
		return rval
	}
	ext := path.Ext(p)
	rval.ext = ext
	nameB = nameB[0 : extR[0]-1]

	var groupN, SongN string
	ps, _ := path.Split(string(nameB))
	if len(ps) > 0 {
		rval.artistInDirectory = true
		rval.outPath = path.Join(rval.outPath, ps)
		nameB = nameB[len(ps):]
		SongN = string(nameB)
		groupN = ps[0 : len(ps)-1] // cut off trailing slash
	}
	words := cReg.FindAllIndex(nameB, -1)
	dash := dReg.FindIndex(nameB)

	var semiExp = regexp.MustCompile("; ")
	if semiExp.Match(nameB) {
		semiLoc := semiExp.FindIndex(nameB)
		SongN = strings.TrimSpace(string(nameB[:semiLoc[0]]))
		groupN = strings.TrimSpace(string(nameB[semiLoc[1]:]))
		rval.alreadyNew = true
	} else {
		//  no ;
		if dash != nil {
			sa := string(nameB[:dash[0]])
			sb := string(nameB[dash[1]:])
			if len(words) == 0 {
				groupN = sa
				SongN = sb
			} else {
				sa := sa
				sb := sb
				if strings.HasPrefix(sa, "The ") {
					sa = sa[4:]
				}
				if strings.HasPrefix(sb, "The ") {
					sb = sb[4:]
				}
				ta, _ := GetEncoder().Encode(JustLetter(sa))
				tb, _ := GetEncoder().Encode(JustLetter(sb))
				_, OKa := gptree.Get(ta)
				if OKa {
					SongN = sb
					groupN = sa
				}
				_, OKb := gptree.Get(tb)
				if OKb {
					groupN = sb
					SongN = sa
				}
			}
		} else {
			// fmt.Println("no dash and no ; try , ")
			commasS := commaExp.FindIndex(nameB)
			if commasS == nil || len(commasS) == 0 {
				SongN = string(nameB)
			} else {
				SongN = string(nameB[:commasS[0]])
				groupN = string(nameB[commasS[1]:])
			}
		}
	}
	groupN = cases.Title(language.English, cases.NoLower).String(strings.TrimSpace(groupN))
	rval.title = strings.TrimSpace(SongN)
	rval.titleH, _ = GetEncoder().Encode(JustLetter(SongN))
	rval.artist = strings.TrimSpace(groupN)
	//fmt.Printf("main %s by %s\n", rval.title, rval.artist)
	if strings.HasPrefix(rval.artist, "The ") {
		rval.artistHasThe = true
		rval.artist = rval.artist[4:]
	}
	rval.artistH, _ = GetEncoder().Encode(JustLetter(rval.artist))
	_, ok := gptree.Get(rval.artistH)
	rval.artistKnown = ok
	return rval
}

const divP = " -+" // want space for names like Led Zeppelin - Bron-Yr-Aur
var dashRegex = regexp.MustCompile(divP)

func GetMetaData(pathArg, p string) *Song {
	rval := new(Song)
	rval.inPath = path.Join(pathArg, p)
	rval.outPath = pathArg
	file, err := os.Open(rval.inPath)
	defer file.Close()
	if err != nil {
		fmt.Println(err)
		return rval
	}
	m, err := tag.ReadFrom(file)
	if err != nil {
		fmt.Printf("%v %s", err, rval.title)
		return rval
	}
	rval.title = m.Title() // The title of the track (see Metadata interface for more details).
	if rval.title == "" {
		_, filename := path.Split(rval.inPath)
		punchIdx := dashRegex.FindStringIndex(filename)
		if punchIdx != nil {
			rval.title = strings.TrimSpace(filename[punchIdx[1]:])
		}
	}
	rval.titleH, _ = GetEncoder().Encode(JustLetter(rval.title))
	rval.artist = m.Artist()
	rval.album = m.Album()

	fmt.Printf("Format() %#v\n`", m.Format())
	fmt.Printf("FileType() %#v\n", m.FileType())
	fmt.Printf("Title() %#v\n", m.Title())
	fmt.Printf("Album() %#v\n", m.Album())
	fmt.Printf("Artist() %#v\n", m.Artist())
	fmt.Printf("AlbumArtist() %#v\n", m.AlbumArtist())
	fmt.Printf("Composer() %#v\n", m.Composer())
	fmt.Printf("Genre() %#v\n", m.Genre())
	fmt.Printf("Year() %#v\n", m.Year())
	t, tt := m.Track()
	fmt.Printf("Track() (int, int) %d of %d v\n", t, tt) // (int, int) // Number,Total
	t, tt = m.Disc()
	fmt.Printf("Disc() (int, int) %d of %d\n", t, tt)      // (int, int) // Number,Total
	fmt.Printf("Picture() *Picture // %#v\n", m.Picture()) // *Picture //Artwork
	fmt.Printf("Lyrics() %#v\n", m.Lyrics())
	fmt.Printf("Comment() %#v\n", m.Comment())

	return rval
}

// this is the local  WalkDirFunc called by WalkDir for each file
// pathArg is the path to the base of our walk
// p is the current path/name
func processFile(pathArg string, sMap map[string]Song, fsys fs.FS, p string, d fs.DirEntry, err error) error {
	if err != nil {
		fmt.Println("Error processing", p, " in ", d)
		fmt.Println("error is ", err)
		return nil
	}
	if d == nil || d.IsDir() || strings.HasPrefix(p, ".") {
		return nil
	}
	extR := extRegex.FindStringIndex(p)
	if extR == nil {
		return nil // not interesting extension
	}
	aSong := new(Song)
	if GetFlags().JsonOutput {
		aSong = GetMetaData(pathArg, p)
	} else {
		aSong = parseFilename(pathArg, p)
	}
	if aSong == nil {
		return nil
	}
	v := sMap[aSong.titleH]
	if len(v.titleH) > 0 {
		if aSong.artistH == v.artistH {
			fmt.Printf("#existing duplicate Song for %s %s == %s\n", aSong.inPath, aSong.title, v.title)
		} else {
			fmt.Printf("#possible dup Song for %s %s == %s %s\n", aSong.inPath, aSong.title,
				v.title, v.artist)
			aSong.titleH += "1"
		}
		return nil
	}
	sMap[aSong.titleH] = *aSong
	return nil
}

// walk all files, looking for nice GoPro created video files.
// fill in a map keyed by the desired new name order
func WalkFiles(pathArg string) map[string]Song {
	songMap := make(map[string]Song)
	fsys := os.DirFS(pathArg)
	fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		err = processFile(pathArg, songMap, fsys, p, d, err)
		return nil
	})
	return songMap
}

// go thru the map, sort by key
// then create new ordering that makes sense to human
func ProcessMap(pathArg string, m map[string]Song) map[string]Song {
	if GetFlags().JsonOutput {
		PrintJson(m)
		return m
	}
	uniqueArtists := make(map[string]string) // we just need a set, but use a map

	for _, aSong := range m {
		switch {
		case GetFlags().DoRename:
			cmd := "mv"
			if runtime.GOOS == "windows" {
				cmd = "ren "
			}
			switch {
			case aSong.alreadyNew:
				continue
				//fmt.Printf("pM aNew %s \"%s\" \"%s/%s; %s%s\"\n", cmd, aSong.inPath,
				//	pathArg, aSong.title, aSong.artist, aSong.ext)
			case aSong.artist == "":
				fmt.Printf("#rename artist is blank %s\n", aSong.inPath)
				cmd = "#" + cmd
				continue
			case aSong.artistInDirectory:
				fmt.Printf("%s \"%s\" \"%s/%s%s\"\n", cmd, aSong.inPath,
					aSong.outPath, aSong.title, aSong.ext)
				continue
			}
			if aSong.artist == "" {
				fmt.Printf("%s \"%s\" \"%s/%s%s\"\n", cmd, aSong.inPath,
					aSong.outPath, aSong.title, aSong.ext)
			} else {
				fmt.Printf("%s \"%s\" \"%s/%s; %s%s\"\n", cmd, aSong.inPath,
					aSong.outPath, aSong.title, aSong.artist, aSong.ext)
			}
		case GetFlags().JustList:
			the := ""
			if aSong.artistHasThe {
				the = "The "
			}
			fmt.Printf("%s by %s%s\n", aSong.title, the, aSong.artist)
		case GetFlags().ShowArtistNotInMap && !aSong.artistKnown:
			prim, _ := GetEncoder().Encode(JustLetter(aSong.artist))
			if len(prim) == 0 && len(aSong.artist) == 0 {
				// fmt.Printf("prim: %s, a: %v\n", prim, aSong)
				continue
			}
			uniqueArtists[prim] = aSong.artist
		case GetFlags().NoGroup:
			if aSong.artist == "" {
				fmt.Printf("nogroup %s\n", aSong.inPath)
			}
		default:
		}
	}
	if GetFlags().ShowArtistNotInMap {
		for k, v := range uniqueArtists {
			fmt.Printf("addto map k: %s v: %s\n", k, v)
		}
	}
	return m
}
func DumpGptree() {
	if GetFlags().ZDumpArtist {
		gptree.Each(func(key string, v string) {
			fmt.Printf("\"%s\", \n", v)
		})
	}
}
