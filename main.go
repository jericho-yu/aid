package main

import (
	"github.com/jericho-yu/aid/str"
)

func main() {
	str.NewTerminalLog("Hello, %s").Default("aaa")
	str.NewTerminalLog("Hello, %s").Info("aaa")
	str.NewTerminalLog("Hello, %s").Success("aaa")
	str.NewTerminalLog("Hello, %s").Wrong("aaa")
	str.NewTerminalLog("Hello, %s").Error("aaa")
}
