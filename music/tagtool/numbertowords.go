// convert int to words
// copied from https://siongui.github.io/2018/01/10/go-convert-number-to-word-from-1-to-1000/
package tagtool

var NumberToWord = map[int]string{
	1:  "one",
	2:  "two",
	3:  "three",
	4:  "four",
	5:  "five",
	6:  "six",
	7:  "seven",
	8:  "eight",
	9:  "nine",
	10: "ten",
	11: "eleven",
	12: "twelve",
	13: "thirteen",
	14: "fourteen",
	15: "fifteen",
	16: "sixteen",
	17: "seventeen",
	18: "eighteen",
	19: "nineteen",
	20: "twenty",
	30: "thirty",
	40: "forty",
	50: "fifty",
	60: "sixty",
	70: "seventy",
	80: "eighty",
	90: "ninety",
}

func convert1to99(n int) (w string) {
	if n < 20 {
		w = NumberToWord[n]
		return
	}

	r := n % 10
	if r == 0 {
		w = NumberToWord[n]
	} else {
		w = NumberToWord[n-r] + "-" + NumberToWord[r]
	}
	return
}

func convert100to999(n int) (w string) {
	q := n / 100
	r := n % 100
	w = NumberToWord[q] + " " + "hundred"
	if r == 0 {
		return
	} else {
		w = w + " and " + convert1to99(r)
	}
	return
}

func Convert1to1000(n int) (w string) {
	if n > 1000 || n < 1 {
		panic("func Convert1to1000: n > 1000 or n < 1")
	}

	if n < 100 {
		w = convert1to99(n)
		return
	}
	if n == 1000 {
		w = "one thousand"
		return
	}
	w = convert100to999(n)
	return
}
