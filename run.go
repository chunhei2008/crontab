package main

import (
	"fmt"
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
	fmt.Println("start")
	o, _ := exec.Command(j.Cmd, j.Args...).Output()
	fmt.Printf("%q\n", o)
	fmt.Println("end")
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
