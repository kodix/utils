package strings

import "testing"

func TestStringLess(t *testing.T) {
	s1 := "a2"
	s2 := "a11"
	if !StringLess(s1, s2) {
		t.Fatalf("%s must be less, then %s", s1, s2)
	}
}
