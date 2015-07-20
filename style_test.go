package tps

import (
	"strings"
	"testing"
)

type AlignTest struct {
	alignment int
	expected  string
}

func TestConvertAlignment(t *testing.T) {
	var value string
	s := Style{}
	s.Alignment = AlignLeft
	// key, value = alignment test, expected value. order does not matter
	tests := []AlignTest{
		{AlignLeft, "L"},
		{AlignLeft | AlignTop, "LT"},
		{AlignBottom | AlignRight, "RB"},
		{AlignMiddle, "M"},
	}
	for _, test := range tests {
		s.Alignment = test.alignment
		value = s.convertAlignment()
		for _, c := range test.expected {
			if !strings.ContainsRune(value, c) {
				t.Error(
					"Style.convertAlignment failed. Expected \"%s\" got \"%s\"",
					test.expected,
					value,
				)
				continue
			}
		}
	}
}
