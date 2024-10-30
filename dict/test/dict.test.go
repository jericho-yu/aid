package main

import (
	"fmt"

	"github.com/jericho-yu/aid/dict"
)

type A struct{}

func main() {
	var (
		t = map[string]A{
			"a": A{},
			"b": A{},
			"c": A{},
			"d": A{},
			"e": A{},
		}
	)

	b := dict.GetValues[string, A](t)
	fmt.Printf("%+v\n", b)
}
