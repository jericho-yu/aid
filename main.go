package main

import (
	"github.com/jericho-yu/aid/reflection"
)

type Example struct{}

func main() {
	var exa *Example
	println(reflection.New(&exa).GetReflectionType())
}
