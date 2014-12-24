package main

import (
	"fmt"
	"time"
)

/*
* 任务执行
* 开始 结束 日志
 */

func doTasks(ctr chan bool) {
	for {
		select {
		case <-ctr:
			break

		case <-time.Tick(1 * time.Minute):
			fmt.Println("tick ############  do tasks ")
			//	go runCmd(t)
		}
	}
}

func runCmd(t task) {
	fmt.Println("start")
	t.cmd.Start()

	t.cmd.Wait()
	fmt.Println("end")
}
