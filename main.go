package main

import (
	"github.com/jericho-yu/aid/log"
)

func main() {

	log.NewZapProvider("logs", false).Info("info")
	log.NewZapProvider("logs", true).Debug("debug")
	log.NewZapProvider("logs", true).Error("error")

}
