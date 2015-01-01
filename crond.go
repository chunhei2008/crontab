package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"runtime"
)

var port *string = flag.String("port", ":8080", "web port")
var logs *string = flag.String("logs", "logs/", "log path")
var conf *string = flag.String("conf", "crontab.conf", "crontab config")
var stopCh chan bool = make(chan bool)
var startCh chan bool = make(chan bool)

const (
	RUN_LOG_POSTFIX = `run.log`
	SVR_LOG         = `svr.log`
	DATEFORMAT      = `20060102`
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.Parse()

	initLog()

	loaded, loadErr := loadConf()
	if !loaded {
		sysLog.Printf("Err %s exit.", loadErr)
		os.Exit(1)
	}

	go runJobs()

	http.HandleFunc("/set", set)
	http.HandleFunc("/get", get)
	http.HandleFunc("/del", del)
	http.HandleFunc("/log", loger)
	http.HandleFunc("/load", load)
	http.HandleFunc("/stop", stop)
	http.HandleFunc("/start", start)
	http.HandleFunc("/status", status)

	startErr := http.ListenAndServe(*port, nil)
	if startErr != nil {
		fmt.Println("Start server failed.", startErr)
		os.Exit(1)
	}
}
