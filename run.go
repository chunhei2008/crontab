package main

import (
	"bufio"
	"os"
	"os/exec"
	"time"
)

/*
* 任务执行
* 开始 结束 日志
 */
var rungings map[int]job = make(map[int]job)

func runJobs(ctr chan bool) {

	for {
		select {
		case <-ctr:
			break
		case <-time.Tick(time.Second):
			t := time.Now()
			if t.Second() == 0 {
				minute := t.Minute()
				hour := t.Hour()
				dom := t.Day()
				month := int(t.Month())
				dow := int(t.Weekday())

				for _, j := range jobs {
					if inArray(j.minute, minute) &&
						inArray(j.hour, hour) &&
						inArray(j.dom, dom) &&
						inArray(j.month, month) &&
						inArray(j.dow, dow) {
						go runJob(j)
					}
				}

			}
		}
	}
}

func runJob(j job) {
	cmd := exec.Command(j.Cmd, j.Args...)
	outpipe, outErr := cmd.StdoutPipe()
	//errpipe, errErr := cmd.StderrPipe()
	if outErr != nil {
		// write into log
	}
	startErr := cmd.Start()
	if startErr != nil {
		runLog.Printf("[Err] %s %s %s %s\n", j.Cmd, j.Args, j.Out, startErr)
		return
	}
	pid := cmd.Process.Pid
	rungings[pid] = j
	defer func() {
		delete(rungings, pid)
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
