package colorune

import (
	"fmt"
	"strconv"
	"strings"
)

// Colorune is Color + Rune
type Colorune struct {
	r        rune
	fg       *Color
	bg       *Color
	hasReset bool
}

// String returns uncolored string
func (c Colorune) String() string {
	return string(c.r)
}

// Rune returns raw rune
func (c Colorune) Rune() rune {
	return c.r
}

// Fg is foreground=text color
func (c Colorune) Fg() *Color {
	return c.fg
}

// Bg is background color
func (c Colorune) Bg() *Color {
	return c.bg
}

// NewRune as constructor of Colorune
func NewRune(r rune, fg *Color, bg *Color, hasReset bool) Colorune {
	fg.Adjust()
	bg.Adjust()
	return Colorune{r: r, fg: fg, bg: bg, hasReset: hasReset}
}

// Colorunes as []Colorune, like a string
type Colorunes []Colorune

func (c Colorunes) String() string {
	var b strings.Builder
	for _, x := range c {
		b.WriteRune(x.r)
	}
	return b.String()
}

// NewString as constructor of Colorunes
func NewString(str string, fg *Color, bg *Color) Colorunes {
	l := len(str)
	c := make([]Colorune, l)

	for i, r := range str {
		if i == 0 {
			c[i] = NewRune(r, fg, bg, false)
			continue
		}
		if i == l-1 && (fg != nil || bg != nil) {
			c[i] = NewRune(r, nil, nil, true)
			continue
		}
		c[i] = NewRune(r, nil, nil, false)
	}
	return Colorunes(c)
}

// ParseANSIString make Colorunes from string
func ParseANSIString(str string) Colorunes {
	l := len(str)
	c := make([]Colorune, 0, l)
	skip := 0

	var fg *Color
	var bg *Color
	hasReset := false

	for i, char := range str {
		if skip > 0 {
			skip--
			continue
		}

		// current char = ESC, next char may have color.
		if char == esc {
			if l > i+1 && str[i+1] == '[' {
				w := 0
				for n := 2; n < 30; n++ {
					w = i + n
					if l > w && str[w] == 'm' {
						skip = n
						break
					}
				}

				codes := strings.Split(str[i+2:w], ";")
				if len(codes) != 5 || (codes[0] != fgCode && codes[0] != bgCode) {
					continue
				}

				switch codes[0] {
				case fgCode:
					fg = &Color{toInt(codes[2]), toInt(codes[3]), toInt(codes[4])}
				case bgCode:
					bg = &Color{toInt(codes[2]), toInt(codes[3]), toInt(codes[4])}
				}
			}
			continue
		}

		// next char = ESC, current char may have reset.
		if l > i+4 && rune(str[i+1]) == esc {
			if str[i+2:i+5] == "[0m" {
				hasReset = true
			}
		}

		c = append(c, NewRune(char, fg, bg, hasReset))
		fg = nil
		bg = nil
		hasReset = false
	}
	return Colorunes(c)
}

func toInt(s string) int {
	x, _ := strconv.ParseInt(s, 10, 32)
	return int(x)
}

// Color has Red Green Blue range 0-255
type Color struct {
	R int
	G int
	B int
}

// Adjust color range of 0-255
func (c *Color) Adjust() {
	if c == nil {
		return
	}

	c.R = adjust(c.R, 255, 0)
	c.G = adjust(c.G, 255, 0)
	c.B = adjust(c.B, 255, 0)
}

// FromHSV to RGB Color struct
func FromHSV(h, s, v int) *Color {
	hh := float32(adjust(h, 360, 0))
	ss := float32(adjust(s, 255, 0))
	vv := float32(adjust(v, 255, 0))

	max := vv
	min := max - ((ss / 255) * max)

	c := &Color{
		R: int(max),
		G: int(max),
		B: int(max),
	}

	x := (max - min) + min

	switch {
	case hh < 60:
		c.G = int((hh / 60) * x)
		c.B = int(min)
	case h < 120:
		c.R = int(((120 - hh) / 60) * x)
		c.B = int(min)
	case h < 180:
		c.R = int(min)
		c.B = int(((hh - 120) / 60) * x)
	case h < 240:
		c.R = int(min)
		c.G = int(((240 - hh) / 60) * x)
	case h < 300:
		c.R = int(((hh - 240) / 60) * x)
		c.G = int(min)
	case h <= 360:
		c.G = int(min)
		c.B = int(((360 - hh) / 60) * x)
	}
	return c
}

func adjust(v, max, min int) int {
	switch {
	case v > max:
		return max
	case v < min:
		return min
	}
	return v
}

const (
	fgCode    = "38"
	bgCode    = "48"
	esc       = rune(0x1b)
	resetCode = string(esc) + "[0m"
)

// FormatANSI generate string with ANSI color from Colorune
func FormatANSI(c Colorune) string {
	if c.r == 0 {
		return ""
	}

	// https://gist.github.com/XVilka/8346728
	var b strings.Builder

	if c.fg != nil {
		b.WriteRune(esc)
		b.WriteString("[")
		b.WriteString(fgCode)
		b.WriteString(";2;")
		b.WriteString(fmt.Sprintf("%d;%d;%dm", c.fg.R, c.fg.G, c.fg.B))
	}
	if c.bg != nil {
		b.WriteRune(esc)
		b.WriteString("[")
		b.WriteString(bgCode)
		b.WriteString(";2;")
		b.WriteString(fmt.Sprintf("%d;%d;%dm", c.bg.R, c.bg.G, c.bg.B))
	}

	b.WriteString(string(c.r))

	if c.hasReset {
		b.WriteString(resetCode)
	}

	return b.String()
}

// FormatANSIString generate string with ANSI color from Colorunes
func FormatANSIString(col Colorunes) string {
	var b strings.Builder
	for _, c := range col {
		b.WriteString(FormatANSI(c))
	}
	return b.String()
}
