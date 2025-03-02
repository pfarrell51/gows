// this package  fix(remove or replace) unicode characters a string
//
//
// this is not multi-processing safe
// this is not general, it uses a specific look up table rather
// than the official "PRECIS"
// (Preparation, Enforcement, and Comparison of Internationalized Strings in Application Protocols)
// and is documented in RFC7564.
//
// for example Joshua Bell's music often uses Cyrillic or Polish, which this does not handle.
//
// handy reference https://www.i18nqa.com/debug/utf8-debug.html
//
// bugs:
//   does not handle multiple bad unicode characters properly.
//   for example, if there are two pairs of unicode points that are bad
//    such as  {"AntonÃ­n DvoÅ™Ã¡k", "Antonin Dvorak", true},
//   it should do the Å™ and then the Ã¡

package util

import (
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

var repuni = map[string]string{
	"à":  "a",
	"á":  "a",
	"ã":  "a",
	"ç":  "c",
	"é":  "e",
	"è":  "e",
	"ë":  "e",
	"ì":  "i",
	"í":  "i",
	"Ã±": "n", // ñ
	"ñ":  "n",
	"ò":  "o",
	"ó":  "o",
	"ô":  "o",
	"ö":  "o",
	"ř":  "r",
	"řá": "ra",
	"ś":  "s",
	"ù":  "u",
	"ú":  "u",
	// AntonÃ­n DvoÅ™Ã¡k   // maybe with final;
	"Ã":      "A",
	"Å":      "A",
	"Å™":     "r",
	"Ã­":     "i", // í
	"Ã³":     "o", // really ó
	"ã©":     "e", // really é
	"Ã©":     "e", // really é
	"â€™":    "'", // ’ closing single
	"â€˜":    "'", // ‘ opening single
	"Ã™":     "U", // really Ù
	"Å¡":     "s", // really š
	"È":      "E",
	"É":      "E",
	"Ù":      "U",
	"Ú":      "U",
	"’":      "'",
	"´":      "'",
	"`":      "'",
	"“":      "\"",
	"”":      "\"",
	"«":      "\"",
	"»":      "\"",
	"…":      "",
	"⁄":      " ",
	"\u00A0": " ", // non-breaking space
	"\u2010": "-", // hyphen
	"\u2013": "-", //En dash
	"\u2014": "-", //Em dash
	"\u2015": "―", // Horizontal bar
}
var longestUnicodeString int
var noTheRegex = regexp.MustCompile("^((T|t)(H|h)(E|e)) ")

var extRegex = regexp.MustCompile(".((M|m)(p|P)(3|4))|((F|f)(L|l)(A|a)(C|c))$")

func init() {
	for key := range repuni {
		runeCount := utf8.RuneCountInString(key)
		if runeCount > longestUnicodeString {
			longestUnicodeString = runeCount
		}
	}
}
func Fratz() {
	fmt.Println("Hello World")
}

// slog.Warn("found high rune", "runeValue",
//
//	fmt.Sprintf("runeval: %c as hex (U+%04X)", runeValue, runeValue))
func CleanUni(s string, c *bool) (r string) {
	wellformed := utf8.FullRuneInString(s)
	if !wellformed {
		slog.Error("input stream not well formed UTF-8", "S", s)
		return
	}
	var sb strings.Builder
	runes := []rune(s)
	highRunes := NewBitVector(len(runes))
	// loop thru input, checking each  rune is not an ASCII character
	for i := 0; i < len(runes); i++ {
		runeValue := runes[i]
		if runeValue > unicode.MaxASCII {
			highRunes.Set(i)
		}
	}
	bitsOn := highRunes.TruePositions()

	for _, onRange := range bitsOn {
		//fmt.Printf("%d, %v\n", i, onRange)
		for i := 0; i < onRange[0]; i++ {
			//	fmt.Printf("copying %c   %s\n", runes[i], string(runes[i]))
			sb.WriteRune(runes[i])
		}
		//lr := (onRange[1] - onRange[0])
		aRune := runes[onRange[0]:onRange[1]]
		k := string(aRune)
		//fmt.Printf("%d, lr: %d >%s< %v %c (U+%04X)\n",
		//		i, lr, k, onRange, aRune, aRune)
		v, ok := repuni[k]
		if ok {
			//fmt.Printf("will replace %s with %s\n", k, v)
			sb.WriteString(v)
			*c = true
		} else {
			fmt.Printf("**** failed %s (U+%04X)\n", k, aRune)
			//			slog.Error("**** error: lookup failed ", "key", k)
		}
	}

	return sb.String()
}

var puncts = regexp.MustCompile(`[\.,'"’]`)

func RemovePunct(s string) (string, bool) {
	rval := puncts.ReplaceAllString(s, "")
	return rval, rval != s
}
