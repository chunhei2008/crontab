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
var ch chan []job = make(chan []job, 10)

func makeJobs(ctr chan bool) {

	for {
		select {
		case <-ctr:
			break
		case <-time.Tick(time.Second):
			t := time.Now()
			if t.Second() == 5 {

				fmt.Println(t.Minute(), t.Hour(), t.Day(), int(t.Month()), int(t.Weekday()))

				fmt.Println("do task#############")
				tjobs := make([]job, 0)
				for _, j := range jobs {
					tjobs = append(tjobs, j)
				}
				ch <- tjobs
			}
		}
	}
}

func doJob() {
	for {
		select {
		case <-time.Tick(time.Second):
			if time.Now().Second() == 0 {
				fmt.Println("do job##############")
				if len(ch) > 0 {
					tjobs := <-ch
					for _, tjob := range tjobs {
						go runJob(tjob)
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
