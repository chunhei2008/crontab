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
var ctr chan bool

const (
	RUN_LOG_POSTFIX = `_run.log`
	SVR_LOG         = `svr.log`
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.Parse()

	if runtime.GOOS == "linux" && *d {
		daemon()
	}

	initLog()

	loadConf()

	go runJobs(ctr)

	http.HandleFunc("/set", set)
	http.HandleFunc("/get", get)
	http.HandleFunc("/del", del)
	http.HandleFunc("/log", loger)
	http.HandleFunc("/load", load)
	http.HandleFunc("/stop", stop)
	http.HandleFunc("/status", status)

	startErr := http.ListenAndServe(*port, nil)
	if startErr != nil {
		fmt.Println("start server failed.", startErr)
		os.Exit(1)
	}
	sysLog.Println("start server success.")
}
