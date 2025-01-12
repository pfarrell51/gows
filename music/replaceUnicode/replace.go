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

package replaceUnicode

import (
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"unicode"
)

var repuni = map[string]string{
	"à":      "a",
	"á":      "a",
	"ã":      "a",
	"ç":      "c",
	"é":      "e",
	"è":      "e",
	"ë":      "e",
	"ì":      "i",
	"í":      "i",
	"ñ":      "n",
	"ò":      "o",
	"ó":      "o",
	"ô":      "o",
	"ö":      "o",
	"ś":      "s",
	"ù":      "u",
	"ú":      "u",
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

const divP = " -+" // want space for names like Led Zeppelin - Bron-Yr-Aur

var dashRegex = regexp.MustCompile(divP)

var extRegex = regexp.MustCompile(".((M|m)(p|P)(3|4))|((F|f)(L|l)(A|a)(C|c))$")

func Fratz() {
	fmt.Println("Hello World")
}
func DoReplacement(s string, c *bool) (r string) {
	var sb, ub strings.Builder
	for _, runeValue := range s {
		if runeValue > unicode.MaxASCII { // Check if the rune is not an ASCII character
			slog.Warn("%c (U+%04X)\n", runeValue, runeValue)
			ub.WriteRune(runeValue)
		} else {
			if ub.Len() == 0 {
				sb.WriteRune(runeValue)
			} else {
				replace(&sb, &ub)
				*c = true
			}
		}
		if ub.Len() > 0 {
			replace(&sb, &ub)
			*c = true
		}
	}

	return sb.String()
}
func replace(sb, ub *strings.Builder) string {
	// lookup
	k := ub.String()
	v, ok := repuni[k]
	if !ok {
		slog.Error("**** error: lookup failed for %s\n", k)
	}
	//fmt.Printf("k: %s f?:%t v: %s (U+%04X)\n", k, ok, v, v)
	sb.WriteString(v)
	ub.Reset()
	return v
}
