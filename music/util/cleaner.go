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
var noTheRegex = regexp.MustCompile("^((T|t)(H|h)(E|e)) ")

var extRegex = regexp.MustCompile(".((M|m)(p|P)(3|4))|((F|f)(L|l)(A|a)(C|c))$")

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
	var sb, ub strings.Builder
	var inUni bool
	runes := []rune(s)
	highRunes := NewBitVector(len(runes))
	for i := 0; i < len(runes); i++ {
		runeValue := runes[i]
		if runeValue > unicode.MaxASCII { // Check if the rune is not an ASCII character
			highRunes.Set(i)
		}
	}
	for i := 0; i < len(runes); i++ {
		runeValue := runes[i]
		if runeValue > unicode.MaxASCII { // Check if the rune is not an ASCII character
			inUni = true
			//fmt.Printf("high rune %c\n", runeValue)
			ub.WriteRune(runeValue)
			k := ub.String()
			_, ok := repuni[k]
			if ok {
				replace(&sb, &ub)
				inUni = false
				*c = true
			} else {
				slog.Error("**** error: lookup failed ", "key", k)
			}
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
	k := ub.String()
	v, ok := repuni[k] // lookup
	if !ok {
		slog.Error("**** error: lookup failed ", "key", k)
	}
	sb.WriteString(v)
	ub.Reset()
	return v
}

var puncts = regexp.MustCompile(`[\.,'"’]`)

func RemovePunct(s string) (r string, changed bool) {
	var sb strings.Builder
	var cf = false
	locA := puncts.FindAllStringIndex(s, -1)
	if len(locA) > 0 {
		var pos int = 0
		for _, loc := range locA {
			sb.WriteString(s[pos:loc[0]])
			sb.WriteString(s[loc[1]:])
			pos += loc[1] - loc[0]
			cf = true
		}
	}
	return sb.String(), cf
}
