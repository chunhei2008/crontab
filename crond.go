package main

import (
	"flag"
	"net/http"
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
	flag.Parse()

	loadConf()

	go runJobs(ctr)

	http.HandleFunc("/set", set)
	http.HandleFunc("/get", get)
	http.HandleFunc("/del", del)
	http.HandleFunc("/log", log)
	http.HandleFunc("/load", load)
	http.HandleFunc("/stop", stop)

	http.ListenAndServe(*port, nil)
}
