package main

import (
	"github.com/jericho-yu/aid/str"
)

func main() {
	str.NewTerminalLog("Hello, %s").Error("aaa")
}
