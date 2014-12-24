package main

import (
	"os/exec"
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
	cmd         *exec.Cmd //命令行
	out         string
}

var lock *sync.RWMutex
var tasks map[string]task = make(map[string]task)

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
