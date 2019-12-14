package colorune

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFormatANSIString(t *testing.T) {
	tests := []struct {
		c    Colorunes
		want string
	}{
		{NewString("", nil, nil), ""},
		{
			NewString("No Color Hello World", nil, nil),
			"No Color Hello World",
		},
		{
			NewString(
				"Hello World With Color",
				&Color{255, 128, 64},
				&Color{128, 64, 255},
			),
			"\x1b[38;2;255;128;64m\x1b[48;2;128;64;255mHello World With Color\x1b[0m",
		},
	}
	for _, test := range tests {
		t.Run(test.c.String(), func(t *testing.T) {
			got := FormatANSIString(test.c)
			assert.Equal(t, test.want, got)

			pgot := ParseANSIString(got)
			assert.Equal(t, test.c.String(), pgot.String())
			assert.Equal(t, got, FormatANSIString(pgot))
		})
	}
}
