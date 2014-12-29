package main

import (
	"log"
	"os"
	"sync"
	"time"
)

var sysLog *log.Logger
var runLog *wyLogger

const DATEFORMAT = `20060102`

func initLog() {
	sysLogFile, _ := os.OpenFile(*logs+"sys.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	//	runLogFile, _ := os.OpenFile(*logs+"run.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)

	sysLog = log.New(sysLogFile, "", log.LstdFlags)
	//	runLog = log.New(runLogFile, "", log.LstdFlags)
	runLog = newWyLogger(*logs, RUN_LOG_POSTFIX)
}

type wyLogger struct {
	dir      string
	filename string
	_date    *time.Time
	mu       *sync.RWMutex
	logfile  *os.File
	lg       *log.Logger
}

// 启动时判断 是否需要rename  判断文件是否存在 存在open 不存在创建（一样的），将文件挂到log上，

func newWyLogger(dir string, filename string) *wyLogger {

	logger := &wyLogger{dir: dir, filename: filename}
	logger._date = new(time.Time)
	logger.mu = new(sync.RWMutex)
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logger.rename()

	go fileMonitor(logger)
	return logger
}

//当前时间跟记录的日志时间比较是否稍后

func (l *wyLogger) isMustRename() bool {
	t, _ := time.Parse(DATEFORMAT, time.Now().Format(DATEFORMAT))
	if t.After(*l._date) {
		return true
	}
	return false
}

//关闭之前的log日志文件，用当前日期创建文件作为日志记录文件

func (l *wyLogger) rename() {

	if l.isMustRename() {
		if l.logfile != nil {
			l.logfile.Close()
		}
		tf := time.Now().Format(DATEFORMAT)
		t, _ := time.Parse(DATEFORMAT, tf)
		l._date = &t
		fn := l.dir + tf + "_" + l.filename
		l.logfile, _ = os.OpenFile(fn, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		l.lg = log.New(l.logfile, "", log.LstdFlags)
	}
}

func fileMonitor(l *wyLogger) {
	for {
		select {
		case <-time.Tick(time.Second):
			go fileCheck(l)
		}
	}
}

func fileCheck(l *wyLogger) {
	if l != nil && l.isMustRename() {
		l.mu.Lock()
		defer l.mu.Unlock()
		l.rename()
	}
}
