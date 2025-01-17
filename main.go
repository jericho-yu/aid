package main

import (
	"github.com/jericho-yu/aid/str"
)

type Os struct{}

func (*Os) TableName() string { return "t_os" }
func main() {
	str.NewTerminalLog("哈哈", "%s").Info("呵呵")
}
