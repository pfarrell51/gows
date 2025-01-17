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

package util

import (
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"unicode"
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
	"Ã":   "A",
	"Å":   "A",
	"Å™":  "r",
	"Ã­":  "i",
	"Ã³":  "o", // really ó
	"ã©":  "e", // really é
	"Ã©":  "e", // really é
	"â€™": "'", // really ’
	"Ã™":  "U", // really Ù
	"Å¡":  "s", // really š
	// řá  Antonín Dvořák
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
var noTheRegex = regexp.MustCompile("^((T|t)(H|h)(E|e)) ")

var extRegex = regexp.MustCompile(".((M|m)(p|P)(3|4))|((F|f)(L|l)(A|a)(C|c))$")

func Fratz() {
	fmt.Println("Hello World")
}

// slog.Warn("found high rune", "runeValue",
//
//	fmt.Sprintf("runeval: %c as hex (U+%04X)", runeValue, runeValue))
func CleanUni(s string, c *bool) (r string) {
	fmt.Printf("CU %s\n", s)
	var sb, ub strings.Builder
	var inUni bool
	for i, runeValue := range s {
		fmt.Printf("0i: %d rune: %c (U+%04X) sb: %s ub: %s\n",
			i, runeValue, runeValue, sb.String(), ub.String())
		if runeValue > unicode.MaxASCII { // Check if the rune is not an ASCII character
			inUni = true
			fmt.Printf("high rune %c\n", runeValue)
			ub.WriteRune(runeValue)
			fmt.Printf("1i: %d rune: %c sb: %s ub: %s\n", i, runeValue, sb.String(), ub.String())
		} else {
			// boring ASCII
			if ub.Len() == 0 {
				if inUni {
					fmt.Printf("non-Uni %c (U+%04X) while inUni  %s to %s\n",
						runeValue, runeValue, sb.String(), ub.String())
				}
			} else {
				if !inUni {
					panic("PIB not inUni")
				}
				//fmt.Printf("2i do replace: %d rune: %c sb: %s ub: %s\n", i, runeValue, sb.String(), ub.String())
				replace(&sb, &ub)
				inUni = false
				*c = true
			}
			sb.WriteRune(runeValue)
		}
	}
	if ub.Len() > 0 {
		replace(&sb, &ub)
		*c = true
	}

	return sb.String()
}
func replace(sb, ub *strings.Builder) string {
	//fmt.Printf("4x sb: %s ub: %s\n", sb.String(), ub.String())
	// lookup
	k := ub.String()
	v, ok := repuni[k]
	if !ok {
		slog.Error("**** error: lookup failed ", "key", k)
	}
	//fmt.Printf("k: %s f?:%t v: %s (U+%04X)\n", k, ok, v, v)
	sb.WriteString(v)
	ub.Reset()
	//fmt.Printf("3 sb: %s ub: %s\n", sb.String(), ub.String())
	return v
}

var puncts = regexp.MustCompile(`[\.,'"’]`)

func RemovePunct(s string) (r string, changed bool) {
	fmt.Println(s)
	var sb strings.Builder
	var cf = false
	locA := puncts.FindAllStringIndex(s, -1)
	if len(locA) > 0 {
		fmt.Printf("loc: %v len(loc) %d\n", locA, len(locA))
		var pos int = 0
		for i, loc := range locA {
			part := s[loc[0]:loc[1]]
			fmt.Printf("in for %d,%v [%d:%d], %s\n", i, loc, loc[0], loc[1], part)
			sb.WriteString(s[pos:loc[0]])
			fmt.Printf("first part %s\n", sb.String())
			sb.WriteString(s[loc[1]:])
			fmt.Printf("second part %s\n", sb.String())
			pos += loc[1] - loc[0]
			cf = true
		}
	}
	fmt.Printf("will return >%s< and %t\n", sb.String(), cf)
	return sb.String(), cf
}
