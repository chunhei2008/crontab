package main

import (
	"bufio"
	"os"
	"os/exec"
	"strconv"
	"time"
)

/*
* 任务执行
* 开始 结束 日志
 */
var runnings map[string]job = make(map[string]job, 20)
var tick *time.Ticker

func runJobs() {
	tick = time.NewTicker(time.Second)
	for {
		select {
		case <-stopCh:
			tick.Stop()
			sysLog.Println("Stop crontab")
		case <-startCh:
			tick = time.NewTicker(time.Second)
			sysLog.Println("Start crontab")
		case <-tick.C:
			t := time.Now()
			if t.Second() == 0 {
				minute := t.Minute()
				hour := t.Hour()
				dom := t.Day()
				month := int(t.Month())
				dow := int(t.Weekday())
				lock.RLock()
				for _, j := range jobs {
					if inArray(j.minute, minute) &&
						inArray(j.hour, hour) &&
						inArray(j.dom, dom) &&
						inArray(j.month, month) &&
						inArray(j.dow, dow) {
						go runJob(j)
					}
				}
				lock.RUnlock()
			}
		}
	}
}

func runJob(j job) {
	cmd := exec.Command(j.Cmd, j.Args...)
	outpipe, outErr := cmd.StdoutPipe()
	if outErr != nil {
		runLog.Printf("[Err] %s %s %s %s\n", j.Cmd, j.Args, j.Out, outErr)
	}
	startErr := cmd.Start()
	if startErr != nil {
		runLog.Printf("[Err] %s %s %s %s\n", j.Cmd, j.Args, j.Out, startErr)
		return
	}
	pid := cmd.Process.Pid
	spid := strconv.Itoa(pid)
	j.Start = time.Now().Format(`2006-01-02 15:04:05`)
	runnings[spid] = j
	defer func() {
		delete(runnings, spid)
		runLog.Printf("[End] pid.%d %s %s %s\n", pid, j.Cmd, j.Args, j.Out)
	}()
	runLog.Printf("[Start] pid.%d %s %s %s\n", pid, j.Cmd, j.Args, j.Out)
	if j.Out != "" {
		of, ofErr := os.OpenFile(j.Out, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if ofErr != nil {
			runLog.Printf("[Err] pid.%d %s %s %s %s", pid, j.Cmd, j.Args, j.Out, ofErr)
		} else {
			defer of.Close()
			outrd := bufio.NewReader(outpipe)
			outrd.WriteTo(of)
		}
	}
	waitErr := cmd.Wait()
	if waitErr != nil {
		runLog.Printf("[Err] pid.%d %s %s %s %s\n", pid, j.Cmd, j.Args, j.Out, waitErr)
	}
}

func inArray(array []int, item int) bool {
	if len(array) < 1 {
		return false
	}
	for _, v := range array {
		if item == v {
			return true
		}
	}
	return false
}
