package main

import (
	"github.com/jericho-yu/aid/reflection"
)

func main() {
	var exa int
	println(reflection.New(exa).GetReflectionType())

}
