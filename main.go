package main

import (
	"github.com/jericho-yu/aid/reflection"
	"log"
)

type Os struct{}

func (*Os) TableName() string { return "t_os" }
func main() {
	val := reflection.New(&Os{}).CallMethodByName("TableName")
	log.Printf("%#v\n", val[0].String())
}
