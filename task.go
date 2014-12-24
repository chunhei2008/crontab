package main

import (
	//"os/exec"
	"sync"
)

/*
*  任务列表管理（添加，删除，更新）
 */
type task struct {
	hash        string //hash值
	prevRuntime int    //上次执行时间
	commit      string //备注
	time        string //crontab时间 * * * * *
	crond       string
	cmd         string
	args        []string
	out         string
}

var lock *sync.RWMutex
var tasks map[string]task

func init() {
	tasks = make(map[string]task)
	tasks["1"] = task{cmd: "php", args: []string{"-v"}}
	tasks["2"] = task{cmd: "php", args: []string{"-v"}}

}

func getTask() string {
	lock.RLock()
	defer lock.RUnlock()
	taskList := ""
	for _, t := range tasks {
		taskList += t.crond + "\n"
	}
	return taskList
}

func setTask(h string, t task) {
	lock.Lock()
	defer lock.Unlock()
	tasks[h] = t
}

func delTask(h string) {
	lock.Lock()
	defer lock.Unlock()
	delete(tasks, h)
}

func nxtTask() {

}
