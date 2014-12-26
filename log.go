package main

import (
	"log"
	"os"
)

var sysLog *log.Logger
var runLog *log.Logger

func initLog() {
	sysLogFile, _ := os.OpenFile(*logs+"sys.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	runLogFile, _ := os.OpenFile(*logs+"run.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)

	sysLog = log.New(sysLogFile, "", log.LstdFlags)
	runLog = log.New(runLogFile, "", log.LstdFlags)
}
