package main

import (
//"os/exec"
//"sync"
)

/*
*  任务列表管理（添加，删除，更新）
 */

type job struct {
	Time    string   `json:"time"`
	Cmd     string   `json:"cmd"`
	Args    []string `json:"args"`
	Out     string   `json:"out"`
	Comment string   `json:"comment"`
}

func getTask() {

}

func setTask() {

}

func delTask() {

}

func nxtTask() {

}
