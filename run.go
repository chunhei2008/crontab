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

var configJobs *Jobs = NewJobs()
var runningJobs *Jobs = NewJobs()

func jobHandle() {
	tick := time.NewTicker(time.Second)
	for {
		select {
		case <-stopCh:
			tick.Stop()
			sysLog.Println("Stop crontab")
		case <-startCh:
			tick = time.NewTicker(time.Second)
			sysLog.Println("Start crontab")
		case <-tick.C:
			go configJobs.runJobs()
		}
	}
}

func runJob(j job) {
	cmd := exec.Command(j.Cmd, j.Args...)
	outpipe, outErr := cmd.StdoutPipe()
	if outErr != nil {
		runLog.lg.Printf("[Err] %s %s %s %s\n", j.Cmd, j.Args, j.Out, outErr)
	}
	startErr := cmd.Start()
	if startErr != nil {
		runLog.lg.Printf("[Err] %s %s %s %s\n", j.Cmd, j.Args, j.Out, startErr)
		return
	}
	pid := cmd.Process.Pid
	spid := strconv.Itoa(pid)
	j.Start = time.Now().Format(TIMEFORMAT)
	runningJobs.add(spid, &j)
	defer func() {
		runningJobs.del(spid)
		runLog.lg.Printf("[End] pid.%d %s %s %s\n", pid, j.Cmd, j.Args, j.Out)
	}()
	runLog.lg.Printf("[Start] pid.%d %s %s %s\n", pid, j.Cmd, j.Args, j.Out)
	if j.Out != "" {
		of, ofErr := os.OpenFile(j.Out, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if ofErr != nil {
			runLog.lg.Printf("[Err] pid.%d %s %s %s %s", pid, j.Cmd, j.Args, j.Out, ofErr)
		} else {
			defer of.Close()
			outrd := bufio.NewReader(outpipe)
			outrd.WriteTo(of)
		}
	}
	waitErr := cmd.Wait()
	if waitErr != nil {
		runLog.lg.Printf("[Err] pid.%d %s %s %s %s\n", pid, j.Cmd, j.Args, j.Out, waitErr)
	}
}

func inArray(array []int, item int) bool {
	if len(array) < 1 {
		return false
	}
	if array[0] == -1 {
		return true
	}
	for _, v := range array {
		if item == v {
			return true
		}
	}
	return false
}
