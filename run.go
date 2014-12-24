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
var ch chan map[string]task = make(chan map[string]task, 10)

func doTasks(ctr chan bool) {

	for {
		select {
		case <-ctr:
			break

		case <-time.Tick(time.Second):
			if time.Now().Second() == 5 {
				fmt.Println("do task#############")
				ch <- tasks
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
					ts := <-ch
					for _, t := range ts {
						go runCmd(t)
					}
				}

			}
		}
	}
}

func runCmd(t task) {
	fmt.Println("start")
	o, _ := exec.Command(t.cmd, t.args...).Output()
	fmt.Printf("%q\n", o)
	fmt.Println("end")
}
