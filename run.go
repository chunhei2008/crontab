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
		runLog.Println("start err")
		return
	}
	//fmt.Println(cmd.Process.Pid) 获取进程ID
	runLog.Printf("%s %s %s start .\n", j.Cmd, j.Args, j.Out)
	//errrd := bufio.NewReader(errpipe)
	//errrd.WriteTo("run log")
	if j.Out != "" {
		of, ofErr := os.OpenFile(j.Out, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if ofErr != nil {
			runLog.Printf("%s %s %s %s", j.Cmd, j.Args, j.Out, ofErr)
		} else {
			defer of.Close()
			outrd := bufio.NewReader(outpipe)
			outrd.WriteTo(of)
		}
	}
	cmd.Wait()
	runLog.Printf("%s %s %s end .\n", j.Cmd, j.Args, j.Out)
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
