package strings

import (
	"math"
	"unicode/utf8"
)

// StringLess compares two alphanumeric strings correctly.
func StringLess(s1, s2 string) (less bool) {
	// uint64 = max 19 digits
	n1, n2 := make([]rune, 0, 18), make([]rune, 0, 18)

	for i, j := 0, 0; i < len(s1) || j < len(s2); {
		var r1, r2 rune
		var w1, w2 int
		var d1, d2 bool

		// read rune from former string available
		if i < len(s1) {
			r1, w1 = utf8.DecodeRuneInString(s1[i:])
			i += w1

			// if digit, accumulate
			if d1 = ('0' <= r1 && r1 <= '9'); d1 {
				n1 = append(n1, r1)
			}
		}

		// read rune from latter string if available
		if j < len(s2) {
			r2, w2 = utf8.DecodeRuneInString(s2[j:])
			j += w2

			// if digit, accumulate
			if d2 = ('0' <= r2 && r2 <= '9'); d2 {
				n2 = append(n2, r2)
			}
		}

		// if have rune and other non-digit rune
		if (!d1 || !d2) && r1 > 0 && r2 > 0 {
			// and accumulators have digits
			if len(n1) > 0 && len(n2) > 0 {
				// make numbers from digit group
				in1 := digitsToNum(n1)
				in2 := digitsToNum(n2)
				// and compare
				if in1 != in2 {
					return in1 < in2
				}
				// if equal, empty accumulators and continue
				n1, n2 = n1[0:0], n2[0:0]
			}
			// detect if non-digit rune from former or latter
			if r1 != r2 {
				return r1 < r2
			}
		}
	}

	// if reached end of both strings and accumulators
	// have some digits
	if len(n1) > 0 || len(n2) > 0 {
		in1 := digitsToNum(n1)
		in2 := digitsToNum(n2)
		if in1 != in2 {
			return in1 < in2
		}
	}

	// last hope
	return len(s1) < len(s2)
}

// Convert a set of runes (digits 0-9) to uint64 number
func digitsToNum(d []rune) (n uint64) {
	if l := len(d); l > 0 {
		n += uint64(d[l-1] - 48)
		k := uint64(l - 1)
		for _, r := range d[:l-1] {
			n, k = n+uint64(r-48)*uint64(math.Pow10(int(k))), k-1
		}
	}
	return
}
