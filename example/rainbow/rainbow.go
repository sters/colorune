package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/sters/colorune"
)

func main() {
	in, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	ins := colorune.ParseANSIString(strings.TrimSpace(string(in)))

	// if you try randomizartion, enable this line
	// rand.Seed(time.Now().UnixNano() + int64(10*os.Getpid()))
	// step := rand.Int()
	step := 0

	len := len(ins)
	col := make([]colorune.Colorune, len)
	for i, inr := range ins {
		col[i] = inr

		if inr.Fg() == nil && inr.Bg() == nil {
			c := colorune.FromHSV(int((float32((i+step)%len)/float32(len))*360), 200, 255)
			col[i] = colorune.NewRune(inr.Rune(), c, nil, true)
		}
	}

	fmt.Println("Rainbow! " + colorune.FormatANSIString(colorune.Colorunes(col)))
}
